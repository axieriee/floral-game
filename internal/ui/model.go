package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/floral-game/floral-realms/internal/game"
	"github.com/floral-game/floral-realms/internal/save"
)

type tab int

const (
	tabGarden tab = iota
	tabUpgrades
	tabFlowers
	tabPrestige
)

var tabNames = []string{"Garden", "Upgrades", "Flowers", "Prestige"}

type tickMsg time.Time

type Model struct {
	state       *game.GameState
	activeTab   tab
	cursor      int
	selectedPlot int // for planting flowers
	message     string
	messageTime time.Time
	width       int
	height      int
}

func NewModel() Model {
	loaded, err := save.Load()
	if err != nil {
		loaded = nil
	}

	var state *game.GameState
	if loaded != nil {
		state = loaded
	} else {
		state = game.NewGameState()
	}

	return Model{
		state:     state,
		activeTab: tabGarden,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	// Calculate offline progress
	petals, harvests := m.state.CalculateOfflineProgress()
	if harvests > 0 {
		m.message = fmt.Sprintf("Welcome back! Earned %.0f petals from %d offline harvests", petals, harvests)
		m.messageTime = time.Now()
	}
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.state.Tick()
		// Auto-save every 30 seconds
		if time.Since(m.state.LastTick) > 30*time.Second {
			m.state.LastTick = time.Now()
			_ = save.Save(m.state)
		}
		// Clear old messages
		if !m.messageTime.IsZero() && time.Since(m.messageTime) > 3*time.Second {
			m.message = ""
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			_ = save.Save(m.state)
			return m, tea.Quit
		case "tab", "right", "l":
			m.activeTab = (m.activeTab + 1) % tab(len(tabNames))
			m.cursor = 0
			return m, nil
		case "shift+tab", "left", "h":
			m.activeTab = (m.activeTab - 1 + tab(len(tabNames))) % tab(len(tabNames))
			m.cursor = 0
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "j":
			m.cursor++
			return m, nil
		case "s":
			_ = save.Save(m.state)
			m.message = "Game saved!"
			m.messageTime = time.Now()
			return m, nil
		}

		// Tab-specific key handling
		switch m.activeTab {
		case tabGarden:
			return m.updateGarden(msg)
		case tabUpgrades:
			return m.updateUpgrades(msg)
		case tabFlowers:
			return m.updateFlowers(msg)
		case tabPrestige:
			return m.updatePrestige(msg)
		}
	}
	return m, nil
}

func (m Model) updateGarden(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Clamp cursor
	maxCursor := len(m.state.Plots) - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	switch msg.String() {
	case "enter", " ":
		if m.cursor < len(m.state.Plots) {
			if m.state.IsReady(m.cursor) {
				petals, seeds, doubled := m.state.Harvest(m.cursor)
				msg := fmt.Sprintf("+%.0f petals", petals)
				if seeds > 0 {
					msg += fmt.Sprintf(", +%.1f seeds", seeds)
				}
				if doubled {
					msg += " (DOUBLE!)"
				}
				m.message = msg
				m.messageTime = time.Now()
			} else {
				m.message = "Not ready yet..."
				m.messageTime = time.Now()
			}
		}
	case "p":
		// Switch to planting mode — go to flowers tab with selected plot
		m.selectedPlot = m.cursor
		m.activeTab = tabFlowers
		m.cursor = 0
	}
	return m, nil
}

func (m Model) updateUpgrades(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxCursor := len(game.Upgrades) - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	switch msg.String() {
	case "enter", " ":
		if m.state.BuyUpgrade(m.cursor) {
			u := game.Upgrades[m.cursor]
			m.message = fmt.Sprintf("Purchased %s!", u.Name)
			m.messageTime = time.Now()
		} else {
			m.message = "Can't afford that upgrade"
			m.messageTime = time.Now()
		}
	}
	return m, nil
}

func (m Model) updateFlowers(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxCursor := len(game.FlowerTypes) - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	switch msg.String() {
	case "enter", " ":
		ft := game.FlowerTypes[m.cursor]
		if !m.state.Unlocked[m.cursor] {
			// Try to unlock
			if m.state.UnlockFlower(m.cursor) {
				m.message = fmt.Sprintf("Unlocked %s %s!", ft.Emoji, ft.Name)
				m.messageTime = time.Now()
			} else {
				m.message = fmt.Sprintf("Need %.0f petals to unlock %s", ft.UnlockCost, ft.Name)
				m.messageTime = time.Now()
			}
		} else if m.selectedPlot >= 0 && m.selectedPlot < len(m.state.Plots) {
			// Plant in selected plot
			if m.state.PlantFlower(m.selectedPlot, m.cursor) {
				m.message = fmt.Sprintf("Planted %s %s in plot %d", ft.Emoji, ft.Name, m.selectedPlot+1)
				m.messageTime = time.Now()
				m.activeTab = tabGarden
				m.cursor = m.selectedPlot
				m.selectedPlot = -1
			} else {
				if ft.SeedCost > 0 {
					m.message = fmt.Sprintf("Need %.0f seeds to plant %s", ft.SeedCost, ft.Name)
				} else {
					m.message = "Can't plant here"
				}
				m.messageTime = time.Now()
			}
		}
	}
	return m, nil
}

func (m Model) updatePrestige(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		nectar := m.state.Prestige()
		if nectar > 0 {
			m.message = fmt.Sprintf("Prestige! Earned %.0f nectar. All growth is now faster!", nectar)
			m.messageTime = time.Now()
			m.activeTab = tabGarden
			m.cursor = 0
		} else {
			m.message = "Need at least 10,000 lifetime petals to prestige"
			m.messageTime = time.Now()
		}
	}
	return m, nil
}
