package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/floral-game/floral-realms/internal/game"
)

var (
	// Colors
	petalColor  = lipgloss.Color("#FF69B4")
	seedColor   = lipgloss.Color("#8B4513")
	nectarColor = lipgloss.Color("#FFD700")
	readyColor  = lipgloss.Color("#00FF7F")
	growColor   = lipgloss.Color("#3CB371")
	dimColor    = lipgloss.Color("#555555")
	accentColor = lipgloss.Color("#BA55D3")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(accentColor).
			Padding(0, 2)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(dimColor)

	activeTabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(accentColor)

	resourceStyle = lipgloss.NewStyle().
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#444444"))

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor)
)

func (m Model) View() string {
	var b strings.Builder

	// Title bar
	b.WriteString(titleStyle.Render("  Floral Realms  "))
	b.WriteString("\n\n")

	// Resource bar
	b.WriteString(m.renderResources())
	b.WriteString("\n\n")

	// Tabs
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Content
	switch m.activeTab {
	case tabGarden:
		b.WriteString(m.renderGarden())
	case tabUpgrades:
		b.WriteString(m.renderUpgrades())
	case tabFlowers:
		b.WriteString(m.renderFlowers())
	case tabPrestige:
		b.WriteString(m.renderPrestige())
	}

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(messageStyle.Render(m.message))
	}

	// Help
	b.WriteString("\n\n")
	b.WriteString(m.renderHelp())

	return b.String()
}

func (m Model) renderResources() string {
	petals := lipgloss.NewStyle().Foreground(petalColor).Render(
		fmt.Sprintf("🌸 %.0f petals", m.state.Petals))
	seeds := lipgloss.NewStyle().Foreground(seedColor).Render(
		fmt.Sprintf("🌱 %.0f seeds", m.state.Seeds))

	resources := petals + "  " + seeds

	if m.state.Nectar > 0 {
		nectar := lipgloss.NewStyle().Foreground(nectarColor).Render(
			fmt.Sprintf("✨ %.0f nectar", m.state.Nectar))
		resources += "  " + nectar
	}

	pps := m.state.PetalsPerSecond()
	if pps > 0 {
		rate := lipgloss.NewStyle().Foreground(growColor).Render(
			fmt.Sprintf("(%.1f/s)", pps))
		resources += "  " + rate
	}

	return resourceStyle.Render(resources)
}

func (m Model) renderTabs() string {
	var tabs []string
	for i, name := range tabNames {
		if tab(i) == m.activeTab {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, tabStyle.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderGarden() string {
	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Your Garden"))
	lines = append(lines, "")

	for i, plot := range m.state.Plots {
		ft := game.FlowerTypes[plot.FlowerType]
		progress := m.state.PlotProgress(i)
		ready := progress >= 1.0

		// Progress bar
		barWidth := 20
		filled := int(progress * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}

		var bar string
		if ready {
			bar = lipgloss.NewStyle().Foreground(readyColor).Render(
				strings.Repeat("█", barWidth))
		} else {
			bar = lipgloss.NewStyle().Foreground(growColor).Render(
				strings.Repeat("█", filled)) +
				lipgloss.NewStyle().Foreground(dimColor).Render(
					strings.Repeat("░", barWidth-filled))
		}

		status := ""
		if ready {
			status = lipgloss.NewStyle().Foreground(readyColor).Bold(true).Render(" READY!")
		} else {
			remaining := m.state.EffectiveGrowTime(ft) - time.Since(plot.Planted)
			if remaining < 0 {
				remaining = 0
			}
			status = lipgloss.NewStyle().Foreground(dimColor).Render(
				fmt.Sprintf(" %s", formatDuration(remaining)))
		}

		line := fmt.Sprintf(" %s %s [%s]%s", ft.Emoji, padRight(ft.Name, 14), bar, status)

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}

		lines = append(lines, line)
	}

	return boxStyle.Render(strings.Join(lines, "\n"))
}

func (m Model) renderUpgrades() string {
	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Upgrades"))
	lines = append(lines, "")

	for i, u := range game.Upgrades {
		level := m.state.UpgradeLevels[u.Effect]
		cost := game.UpgradeCost(u, level)

		maxed := u.MaxLevel > 0 && level >= u.MaxLevel
		affordable := m.state.Petals >= cost

		var line string
		if maxed {
			line = fmt.Sprintf(" ✓ %s (MAX)", u.Name)
			line = lipgloss.NewStyle().Foreground(dimColor).Render(line)
		} else {
			costStr := formatNumber(cost)
			levelStr := ""
			if level > 0 {
				levelStr = fmt.Sprintf(" [Lv.%d]", level)
			}

			line = fmt.Sprintf(" %s%s - %s (%.0s petals)",
				padRight(u.Name, 18), levelStr, u.Description, costStr)

			costDisplay := fmt.Sprintf(" Cost: %s petals", costStr)
			if affordable {
				costDisplay = lipgloss.NewStyle().Foreground(readyColor).Render(costDisplay)
			} else {
				costDisplay = lipgloss.NewStyle().Foreground(dimColor).Render(costDisplay)
			}
			line += costDisplay
		}

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		lines = append(lines, line)
	}

	return boxStyle.Render(strings.Join(lines, "\n"))
}

func (m Model) renderFlowers() string {
	var lines []string

	header := "Flower Collection"
	if m.selectedPlot >= 0 && m.selectedPlot < len(m.state.Plots) {
		header = fmt.Sprintf("Choose flower for Plot %d", m.selectedPlot+1)
	}
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(header))
	lines = append(lines, "")

	for i, ft := range game.FlowerTypes {
		unlocked := m.state.Unlocked[i]

		var line string
		if unlocked {
			growTime := m.state.EffectiveGrowTime(ft)
			yield := (ft.PetalYield + m.state.FlatBonus()) * m.state.PetalMultiplier()
			line = fmt.Sprintf(" %s %s  Yield: %.0f  Time: %s",
				ft.Emoji, padRight(ft.Name, 14), yield, formatDuration(growTime))
			if ft.SeedCost > 0 {
				seedStr := lipgloss.NewStyle().Foreground(seedColor).Render(
					fmt.Sprintf("  Plant: %.0f seeds", ft.SeedCost))
				line += seedStr
			}
		} else {
			line = fmt.Sprintf(" 🔒 %s  Unlock: %.0f petals",
				padRight(ft.Name, 14), ft.UnlockCost)
			if m.state.Petals >= ft.UnlockCost {
				line = lipgloss.NewStyle().Foreground(readyColor).Render(line)
			} else {
				line = lipgloss.NewStyle().Foreground(dimColor).Render(line)
			}
		}

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		lines = append(lines, line)
	}

	return boxStyle.Render(strings.Join(lines, "\n"))
}

func (m Model) renderPrestige() string {
	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Prestige - The Eternal Garden"))
	lines = append(lines, "")

	lines = append(lines,
		fmt.Sprintf(" Lifetime petals: %s", formatNumber(m.state.TotalPetals)))
	lines = append(lines,
		fmt.Sprintf(" Current nectar:  %s", formatNumber(m.state.Nectar)))
	lines = append(lines,
		fmt.Sprintf(" Prestige count:  %d", m.state.PrestigeCount))
	lines = append(lines, "")

	nectarGain := m.state.NectarFromPrestige()
	if nectarGain > 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(nectarColor).Bold(true).Render(
			fmt.Sprintf(" Press ENTER to prestige for %.0f nectar!", nectarGain)))
		lines = append(lines, "")
		lines = append(lines, " This will reset your garden, upgrades, and flowers.")
		lines = append(lines, " You keep your nectar, which gives +10% petals per prestige.")
	} else {
		needed := 10000 - m.state.TotalPetals
		if needed < 0 {
			needed = 0
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
			fmt.Sprintf(" Earn %.0f more lifetime petals to unlock prestige", needed)))
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
		" Nectar bonuses:"))
	lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
		fmt.Sprintf("   +%d%% petal yield (from %d prestiges)",
			m.state.PrestigeCount*10, m.state.PrestigeCount)))

	return boxStyle.Render(strings.Join(lines, "\n"))
}

func (m Model) renderHelp() string {
	var help string
	switch m.activeTab {
	case tabGarden:
		help = "↑/↓ select  •  ENTER harvest  •  p plant  •  TAB switch tab  •  s save  •  q quit"
	case tabUpgrades:
		help = "↑/↓ select  •  ENTER buy  •  TAB switch tab  •  s save  •  q quit"
	case tabFlowers:
		help = "↑/↓ select  •  ENTER unlock/plant  •  TAB switch tab  •  s save  •  q quit"
	case tabPrestige:
		help = "ENTER prestige  •  TAB switch tab  •  s save  •  q quit"
	}
	return helpStyle.Render(help)
}

// Helpers

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "0s"
	}
	s := int(d.Seconds())
	if s < 60 {
		return fmt.Sprintf("%ds", s)
	}
	m := s / 60
	s = s % 60
	return fmt.Sprintf("%dm%ds", m, s)
}

func formatNumber(n float64) string {
	if n < 1000 {
		return fmt.Sprintf("%.0f", n)
	}
	suffixes := []string{"", "K", "M", "B", "T"}
	i := int(math.Log10(n) / 3)
	if i >= len(suffixes) {
		i = len(suffixes) - 1
	}
	val := n / math.Pow(10, float64(i*3))
	return fmt.Sprintf("%.1f%s", val, suffixes[i])
}
