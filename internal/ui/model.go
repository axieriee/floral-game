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
	tabBreeding
	tabPrestige
	tabAchievements
)

var tabNames = []string{"Garden", "Upgrades", "Flowers", "Breeding", "Prestige", "Achievements"}

type tickMsg time.Time

type Model struct {
	state        *game.GameState
	activeTab    tab
	cursor       int
	selectedPlot int // for planting flowers
	message      string
	messageTime  time.Time
	width        int
	height       int
	lastSeason   game.Season
}

func NewModel() Model {
	loaded, err := save.Load()
	if err != nil {
		loaded = nil
	}

	var state *game.GameState
	if loaded != nil {
		state = loaded
		// Ensure maps are initialized for old saves
		if state.EssenceUpgrades == nil {
			state.EssenceUpgrades = make(map[string]int)
		}
		if state.UpgradeLevels == nil {
			state.UpgradeLevels = make(map[string]int)
		}
		if state.CompletedAchievements == nil {
			state.CompletedAchievements = make(map[game.AchievementID]bool)
		}
		if state.EventHarvests == nil {
			state.EventHarvests = make(map[string]bool)
		}
		// Ensure unlocked slice is big enough for hybrids
		totalFlowers := len(game.FlowerTypes) + len(game.HybridFlowers)
		if len(state.Unlocked) < totalFlowers {
			extended := make([]bool, totalFlowers)
			copy(extended, state.Unlocked)
			state.Unlocked = extended
		}
		state.Log = nil
	} else {
		state = game.NewGameState()
	}

	return Model{
		state:        state,
		activeTab:    tabGarden,
		selectedPlot: -1,
		lastSeason:   state.CurrentSeason(),
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	petals, harvests := m.state.CalculateOfflineProgress()
	if harvests > 0 {
		m.message = fmt.Sprintf("Welcome back! Earned %s petals from %d offline harvests",
			formatNumber(petals), harvests)
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

		// Check for season change
		currentSeason := m.state.CurrentSeason()
		if currentSeason != m.lastSeason {
			s := game.Seasons[currentSeason]
			m.state.AddLog(fmt.Sprintf("%s %s has arrived! %s", s.Emoji, s.Name, s.Special), "#FFD700")
			m.message = fmt.Sprintf("%s %s has arrived! %s", s.Emoji, s.Name, s.Special)
			m.messageTime = time.Now()
			m.lastSeason = currentSeason
		}

		// Try to trigger random garden events
		if evt := m.state.TryTriggerEvent(); evt != nil {
			m.state.AddLog(fmt.Sprintf("%s %s — %s", evt.Emoji, evt.Name, evt.Description), "#FFD700")
			m.message = fmt.Sprintf("%s %s — %s", evt.Emoji, evt.Name, evt.Description)
			m.messageTime = time.Now()
		}

		// Check achievements
		newAch := m.state.CheckAchievements()
		for _, a := range newAch {
			m.message = fmt.Sprintf("%s Achievement Unlocked: %s! (+%.0f %s)",
				a.Emoji, a.Name, a.Reward, a.RewardType)
			m.messageTime = time.Now()
		}

		// Auto-save every 30 seconds
		if time.Since(m.state.LastSave) > 30*time.Second {
			m.state.LastSave = time.Now()
			_ = save.Save(m.state)
		}

		// Clear old messages
		if !m.messageTime.IsZero() && time.Since(m.messageTime) > 4*time.Second {
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
		case "1":
			m.activeTab = tabGarden
			m.cursor = 0
			return m, nil
		case "2":
			m.activeTab = tabUpgrades
			m.cursor = 0
			return m, nil
		case "3":
			m.activeTab = tabFlowers
			m.cursor = 0
			return m, nil
		case "4":
			m.activeTab = tabBreeding
			m.cursor = 0
			return m, nil
		case "5":
			m.activeTab = tabPrestige
			m.cursor = 0
			return m, nil
		case "6":
			m.activeTab = tabAchievements
			m.cursor = 0
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
		case tabBreeding:
			return m.updateBreeding(msg)
		case tabPrestige:
			return m.updatePrestige(msg)
		case tabAchievements:
			return m.updateAchievements(msg)
		}
	}
	return m, nil
}

func (m Model) updateGarden(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxCursor := len(m.state.Plots) - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	switch msg.String() {
	case "enter", " ":
		if m.cursor < len(m.state.Plots) {
			if m.state.IsReady(m.cursor) {
				// Register combo before harvest so multiplier applies
				m.state.RegisterManualHarvest()
				comboMult := m.state.ComboMultiplier()

				petals, seeds, doubled, hybridIdx := m.state.Harvest(m.cursor)
				// Apply combo multiplier to manual harvests
				if comboMult > 1.0 {
					bonus := petals * (comboMult - 1.0)
					petals += bonus
					m.state.Petals += bonus
					m.state.TotalPetals += bonus
				}

				msg := fmt.Sprintf("+%s petals", formatNumber(petals))
				if seeds > 0 {
					msg += fmt.Sprintf(", +%.1f seeds", seeds)
				}
				if doubled {
					msg += " (DOUBLE!)"
				}
				if m.state.ComboCount > 1 {
					msg += fmt.Sprintf(" [%dx COMBO! %.0f%%]", m.state.ComboCount, (comboMult-1)*100)
				}
				if hybridIdx >= 0 {
					h := game.HybridFlowers[hybridIdx]
					msg += fmt.Sprintf("  NEW HYBRID: %s %s!", h.Emoji, h.Name)
				}
				m.message = msg
				m.messageTime = time.Now()
				m.state.AddLog(msg, "#3CB371")
			} else {
				m.message = "Not ready yet..."
				m.messageTime = time.Now()
			}
		}
	case "p":
		m.selectedPlot = m.cursor
		m.activeTab = tabFlowers
		m.cursor = 0
	}
	return m, nil
}

func (m Model) updateUpgrades(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	totalItems := len(game.Upgrades)
	if m.state.Essence > 0 || m.state.Prestige2Count > 0 {
		totalItems += len(game.EssenceUpgrades)
	}
	maxCursor := totalItems - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	switch msg.String() {
	case "enter", " ":
		if m.cursor < len(game.Upgrades) {
			if m.state.BuyUpgrade(m.cursor) {
				u := game.Upgrades[m.cursor]
				m.message = fmt.Sprintf("Purchased %s!", u.Name)
				m.messageTime = time.Now()
				m.state.AddLog(fmt.Sprintf("⚒ Bought %s", u.Name), "#00FF7F")
			} else {
				m.message = "Can't afford that upgrade"
				m.messageTime = time.Now()
			}
		} else {
			// Essence upgrade
			essIdx := m.cursor - len(game.Upgrades)
			if m.state.BuyEssenceUpgrade(essIdx) {
				eu := game.EssenceUpgrades[essIdx]
				m.message = fmt.Sprintf("Purchased %s!", eu.Name)
				m.messageTime = time.Now()
				m.state.AddLog(fmt.Sprintf("💫 Bought %s", eu.Name), "#FF69B4")
			} else {
				m.message = "Not enough essence"
				m.messageTime = time.Now()
			}
		}
	}
	return m, nil
}

func (m Model) updateFlowers(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	totalFlowers := len(game.FlowerTypes) + len(game.HybridFlowers)
	maxCursor := totalFlowers - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	switch msg.String() {
	case "enter", " ":
		globalIdx := m.cursor
		if !m.state.Unlocked[globalIdx] {
			// Try to unlock (only base flowers)
			if globalIdx < len(game.FlowerTypes) {
				if m.state.UnlockFlower(globalIdx) {
					ft := game.FlowerTypes[globalIdx]
					m.message = fmt.Sprintf("Unlocked %s %s!", ft.Emoji, ft.Name)
					m.messageTime = time.Now()
					m.state.AddLog(fmt.Sprintf("🔓 Unlocked %s %s", ft.Emoji, ft.Name), "#FFD700")
				} else {
					ft := game.FlowerTypes[globalIdx]
					m.message = fmt.Sprintf("Need %s petals to unlock %s",
						formatNumber(ft.UnlockCost), ft.Name)
					m.messageTime = time.Now()
				}
			} else {
				m.message = "Hybrids are discovered through breeding!"
				m.messageTime = time.Now()
			}
		} else if m.selectedPlot >= 0 && m.selectedPlot < len(m.state.Plots) {
			name, emoji, _, _, seedCost, _ := game.GetFlowerInfo(globalIdx)
			if m.state.PlantFlower(m.selectedPlot, globalIdx) {
				m.message = fmt.Sprintf("Planted %s %s in plot %d", emoji, name, m.selectedPlot+1)
				m.messageTime = time.Now()
				m.activeTab = tabGarden
				m.cursor = m.selectedPlot
				m.selectedPlot = -1
			} else {
				if seedCost > 0 {
					m.message = fmt.Sprintf("Need %.0f seeds to plant %s", seedCost, name)
				} else {
					m.message = "Can't plant here"
				}
				m.messageTime = time.Now()
			}
		} else {
			m.message = "Press 'p' on a garden plot first to select where to plant"
			m.messageTime = time.Now()
		}
	}
	return m, nil
}

func (m Model) updateBreeding(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxCursor := len(game.HybridRecipes) - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}
	// Breeding tab is informational — no enter action needed
	return m, nil
}

func (m Model) updateAchievements(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxCursor := len(game.Achievements) - 1
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}
	// Read-only tab
	return m, nil
}

func (m Model) updatePrestige(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.cursor > 1 {
		m.cursor = 1
	}

	switch msg.String() {
	case "enter", " ":
		if m.cursor == 0 {
			nectar := m.state.Prestige()
			if nectar > 0 {
				m.message = fmt.Sprintf("Prestige! Earned %.0f nectar", nectar)
				m.messageTime = time.Now()
				m.activeTab = tabGarden
				m.cursor = 0
			} else {
				m.message = "Need at least 10,000 lifetime petals to prestige"
				m.messageTime = time.Now()
			}
		} else if m.cursor == 1 {
			essence := m.state.Prestige2()
			if essence > 0 {
				m.message = fmt.Sprintf("Transcendence! Earned %.0f essence", essence)
				m.messageTime = time.Now()
				m.activeTab = tabGarden
				m.cursor = 0
			} else {
				m.message = "Need at least 50 nectar to transcend"
				m.messageTime = time.Now()
			}
		}
	}
	return m, nil
}
