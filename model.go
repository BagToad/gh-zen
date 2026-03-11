package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// layout constants
const (
	borderSize    = 1
	statusBarRows = 2 // status bar + padding
)

type model struct {
	garden   *Garden
	rake     Rake
	width    int
	height   int
	won      bool
	zenQuote string
	zenErr   error
	quitting bool
	debug    bool
}

func newModel() model {
	return model{
		debug: os.Getenv("GH_ZEN_DEBUG") != "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// initGarden creates a fresh garden sized to the terminal.
func (m *model) initGarden() {
	// reserve space for border and status bar
	gw := m.width - borderSize*2
	gh := m.height - borderSize*2 - statusBarRows
	if gw < 10 {
		gw = 10
	}
	if gh < 6 {
		gh = 6
	}

	m.garden = newGarden(gw, gh)
	m.rake = newRake(gw/2, gh/2)
	for _, c := range m.rake.cells() {
		m.garden.set(c[0], c[1], CellFlattened)
	}
	m.garden.placeRocks(m.rake.X, m.rake.Y)
	m.won = false
	m.zenQuote = ""
	m.zenErr = nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initGarden()
		return m, nil

	case tea.KeyMsg:
		if m.quitting {
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "r":
			m.initGarden()
			return m, nil

		case "f", "tab":
			m.rake.Flipped = !m.rake.Flipped
			return m, nil

		case "c":
			if m.debug {
				return m.clearAllSand()
			}

		case "up", "k":
			return m.moveRake(DirUp)
		case "down", "j":
			return m.moveRake(DirDown)
		case "left", "h":
			return m.moveRake(DirLeft)
		case "right", "l":
			return m.moveRake(DirRight)
		}

	case zenMsg:
		m.zenQuote = msg.quote
		m.zenErr = msg.err
		return m, nil
	}
	return m, nil
}

func (m model) clearAllSand() (tea.Model, tea.Cmd) {
	if m.garden == nil || m.won {
		return m, nil
	}
	for y := range m.garden.Cells {
		for x := range m.garden.Cells[y] {
			c := m.garden.Cells[y][x]
			if c == CellSand || c == CellRakedH || c == CellRakedV {
				m.garden.Cells[y][x] = CellFlattened
			}
		}
	}
	m.won = true
	return m, fetchZen
}

func (m model) moveRake(dir Direction) (tea.Model, tea.Cmd) {
	if m.garden == nil || m.won {
		return m, nil
	}

	m.rake.Dir = dir
	dx, dy := dir.delta()
	newX, newY := m.rake.X+dx, m.rake.Y+dy
	newCells := perpCells(newX, newY, dir)

	// check all 3 target cells are passable
	for _, pos := range newCells {
		if !m.garden.inBounds(pos[0], pos[1]) {
			return m, nil
		}
		if m.garden.at(pos[0], pos[1]) == CellRock {
			return m, nil
		}
	}

	// apply rake effect on the 3 cells we're entering
	for _, pos := range newCells {
		cell := m.garden.at(pos[0], pos[1])
		if m.rake.Flipped {
			// flat end: smooth everything into flattened
			if cell == CellSand || cell == CellRakedH || cell == CellRakedV {
				m.garden.set(pos[0], pos[1], CellFlattened)
			}
		} else {
			// spokes end: rake everything into directional lines
			if cell == CellSand || cell == CellFlattened || cell == CellRakedH || cell == CellRakedV {
				m.garden.set(pos[0], pos[1], rakedCell(dir))
			}
		}
	}

	m.rake.X = newX
	m.rake.Y = newY

	// check win condition
	if m.garden.countUnraked() == 0 && !m.won {
		m.won = true
		return m, fetchZen
	}

	return m, nil
}

func (m model) View() string {
	if m.garden == nil {
		return "Initializing garden..."
	}

	var b strings.Builder
	g := m.garden

	// colours
	sandStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#C2B280"))
	rakedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7355"))
	flatStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B5B3A"))
	rockStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#696969")).Bold(true)
	rakeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6347")).Bold(true)
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#5C4033"))

	// top border
	b.WriteString(borderStyle.Render("╔" + strings.Repeat("═", g.Width) + "╗"))
	b.WriteRune('\n')

	for y := 0; y < g.Height; y++ {
		b.WriteString(borderStyle.Render("║"))
		for x := 0; x < g.Width; x++ {
			if m.rake.occupies(x, y) {
				b.WriteString(rakeStyle.Render(string(rakeRune(m.rake.Dir, m.rake.Flipped))))
				continue
			}
			c := g.Cells[y][x]
			t := g.Texture[y][x]
			switch c {
			case CellSand:
				var ch rune
				switch {
				case t < 15:
					ch = '▒'
				case t < 22:
					ch = '▓'
				default:
					ch = '░'
				}
				b.WriteString(sandStyle.Render(string(ch)))
			case CellRakedH:
				b.WriteString(rakedStyle.Render("─"))
			case CellRakedV:
				b.WriteString(rakedStyle.Render("│"))
			case CellFlattened:
				var ch rune
				switch {
				case t < 10:
					ch = '∙'
				case t < 18:
					ch = ' '
				default:
					ch = '·'
				}
				b.WriteString(flatStyle.Render(string(ch)))
			case CellRock:
				up := g.at(x, y-1) == CellRock
				down := g.at(x, y+1) == CellRock
				left := g.at(x-1, y) == CellRock
				right := g.at(x+1, y) == CellRock
				b.WriteString(rockStyle.Render(string(rockChar(up, down, left, right))))
			default:
				b.WriteRune(' ')
			}
		}
		b.WriteString(borderStyle.Render("║"))
		b.WriteRune('\n')
	}

	// bottom border
	b.WriteString(borderStyle.Render("╚" + strings.Repeat("═", g.Width) + "╝"))
	b.WriteRune('\n')

	// status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(m.width)

	remaining := g.countUnraked()
	progress := ""
	if m.won {
		if m.zenQuote != "" {
			zenStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7FFFD4")).
				Bold(true).
				Italic(true)
			progress = "🧘 " + zenStyle.Render(m.zenQuote)
		} else if m.zenErr != nil {
			progress = fmt.Sprintf("✨ Garden complete! (zen fetch failed: %v)", m.zenErr)
		} else {
			progress = "✨ Garden complete! Fetching wisdom..."
		}
	} else {
		total := remaining
		for y := 0; y < g.Height; y++ {
			for _, c := range g.Cells[y] {
				if c == CellRakedH || c == CellRakedV || c == CellFlattened {
					total++
				}
			}
		}
		pct := 0
		if total > 0 {
			pct = ((total - remaining) * 100) / total
		}
		progress = fmt.Sprintf("Sand remaining: %d (%d%% raked)", remaining, pct)
	}

	modeLabel := "╞ spokes"
	if m.rake.Flipped {
		modeLabel = "━ flat"
	}
	keys := "↑↓←→/hjkl: move │ f: flip rake │ r: new garden │ q: quit │ " + modeLabel
	if m.debug {
		keys += " │ c: clear sand (debug)"
	}
	b.WriteString(statusStyle.Render(progress))
	b.WriteRune('\n')
	b.WriteString(statusStyle.Render(keys))

	return b.String()
}
