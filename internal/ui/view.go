package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/floral-game/floral-realms/internal/game"
)

// Season-themed color palettes
type seasonPalette struct {
	primary   lipgloss.Color
	secondary lipgloss.Color
	accent    lipgloss.Color
	bg        lipgloss.Color
	border    lipgloss.Color
}

var seasonPalettes = [game.NumSeasons]seasonPalette{
	{ // Spring
		primary:   lipgloss.Color("#7FFF7F"),
		secondary: lipgloss.Color("#FFB7D5"),
		accent:    lipgloss.Color("#90EE90"),
		bg:        lipgloss.Color("#1A2F1A"),
		border:    lipgloss.Color("#7FFF7F"),
	},
	{ // Summer
		primary:   lipgloss.Color("#FFD700"),
		secondary: lipgloss.Color("#FF8C00"),
		accent:    lipgloss.Color("#FFA500"),
		bg:        lipgloss.Color("#2F2A1A"),
		border:    lipgloss.Color("#FFD700"),
	},
	{ // Autumn
		primary:   lipgloss.Color("#CD853F"),
		secondary: lipgloss.Color("#FF6347"),
		accent:    lipgloss.Color("#DAA520"),
		bg:        lipgloss.Color("#2F1F1A"),
		border:    lipgloss.Color("#CD853F"),
	},
	{ // Winter
		primary:   lipgloss.Color("#ADD8E6"),
		secondary: lipgloss.Color("#E0FFFF"),
		accent:    lipgloss.Color("#87CEEB"),
		bg:        lipgloss.Color("#1A1F2F"),
		border:    lipgloss.Color("#ADD8E6"),
	},
}

var (
	petalColor  = lipgloss.Color("#FF69B4")
	seedColor   = lipgloss.Color("#8B4513")
	nectarColor = lipgloss.Color("#FFD700")
	essenceColor = lipgloss.Color("#FF69B4")
	readyColor  = lipgloss.Color("#00FF7F")
	growColor   = lipgloss.Color("#3CB371")
	dimColor    = lipgloss.Color("#555555")
	hybridColor = lipgloss.Color("#DA70D6")
	warnColor   = lipgloss.Color("#FF4444")

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#444444"))
)

func (m Model) palette() seasonPalette {
	return seasonPalettes[m.state.CurrentSeason()]
}

func (m Model) titleStyle() lipgloss.Style {
	p := m.palette()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(p.primary).
		Padding(0, 2)
}

func (m Model) boxStyle() lipgloss.Style {
	p := m.palette()
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.border).
		Padding(0, 1)
}

func (m Model) activeTabStyle() lipgloss.Style {
	p := m.palette()
	return lipgloss.NewStyle().
		Padding(0, 2).
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(p.primary)
}

func (m Model) View() string {
	var b strings.Builder
	season := game.Seasons[m.state.CurrentSeason()]

	// Title bar with season
	title := fmt.Sprintf("  %s Floral Realms  %s  ", season.Emoji, season.Emoji)
	b.WriteString(m.titleStyle().Render(title))
	b.WriteString("  ")
	b.WriteString(m.renderSeasonBar())
	b.WriteString("\n\n")

	// Resource bar
	b.WriteString(m.renderResources())
	b.WriteString("\n\n")

	// Tabs
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Active event banner
	if m.state.IsEventActive() {
		b.WriteString(m.renderEventBanner())
		b.WriteString("\n")
	}

	// Combo indicator
	if m.state.ComboActive() && m.state.ComboCount > 1 {
		comboStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF4500"))
		comboMult := m.state.ComboMultiplier()
		b.WriteString(comboStyle.Render(
			fmt.Sprintf("  COMBO x%d (%.0f%% bonus)  ", m.state.ComboCount, (comboMult-1)*100)))
		b.WriteString("\n")
	}

	// Content
	switch m.activeTab {
	case tabGarden:
		b.WriteString(m.renderGarden())
	case tabUpgrades:
		b.WriteString(m.renderUpgrades())
	case tabFlowers:
		b.WriteString(m.renderFlowers())
	case tabBreeding:
		b.WriteString(m.renderBreeding())
	case tabPrestige:
		b.WriteString(m.renderPrestige())
	case tabAchievements:
		b.WriteString(m.renderAchievements())
	}

	// Event log (bottom right area)
	b.WriteString("\n")
	b.WriteString(m.renderLog())

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

func (m Model) renderSeasonBar() string {
	season := game.Seasons[m.state.CurrentSeason()]
	progress := m.state.SeasonProgress()
	p := m.palette()

	barWidth := 15
	filled := int(progress * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	bar := lipgloss.NewStyle().Foreground(p.primary).Render(
		strings.Repeat("━", filled)) +
		lipgloss.NewStyle().Foreground(dimColor).Render(
			strings.Repeat("─", barWidth-filled))

	label := lipgloss.NewStyle().Foreground(p.primary).Bold(true).Render(season.Name)
	flavor := lipgloss.NewStyle().Foreground(p.secondary).Italic(true).Render(season.Special)

	return fmt.Sprintf("%s [%s] %s", label, bar, flavor)
}

func (m Model) renderResources() string {
	petals := lipgloss.NewStyle().Foreground(petalColor).Bold(true).Render(
		fmt.Sprintf("🌸 %s", formatNumber(m.state.Petals)))
	seeds := lipgloss.NewStyle().Foreground(seedColor).Render(
		fmt.Sprintf("🌱 %s", formatNumber(m.state.Seeds)))

	parts := []string{petals, seeds}

	if m.state.Nectar > 0 || m.state.PrestigeCount > 0 {
		nectar := lipgloss.NewStyle().Foreground(nectarColor).Render(
			fmt.Sprintf("✨ %s", formatNumber(m.state.Nectar)))
		parts = append(parts, nectar)
	}

	if m.state.Essence > 0 || m.state.Prestige2Count > 0 {
		essence := lipgloss.NewStyle().Foreground(essenceColor).Render(
			fmt.Sprintf("💫 %s", formatNumber(m.state.Essence)))
		parts = append(parts, essence)
	}

	pps := m.state.PetalsPerSecond()
	if pps > 0 {
		rate := lipgloss.NewStyle().Foreground(growColor).Render(
			fmt.Sprintf("(%s/s)", formatNumber(pps)))
		parts = append(parts, rate)
	}

	return lipgloss.NewStyle().Padding(0, 1).Render(strings.Join(parts, "  "))
}

func (m Model) renderTabs() string {
	tabStyle := lipgloss.NewStyle().Padding(0, 2).Foreground(dimColor)
	var tabs []string
	for i, name := range tabNames {
		if tab(i) == m.activeTab {
			tabs = append(tabs, m.activeTabStyle().Render(name))
		} else {
			tabs = append(tabs, tabStyle.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderGarden() string {
	var lines []string
	season := game.Seasons[m.state.CurrentSeason()]
	p := m.palette()

	header := lipgloss.NewStyle().Bold(true).Foreground(p.primary).Render(
		fmt.Sprintf("%s Your Garden %s", season.Emoji, season.Emoji))
	lines = append(lines, header)
	lines = append(lines, "")

	// Visual garden bed
	lines = append(lines, m.renderGardenBed())
	lines = append(lines, "")

	// Plot details
	for i, plot := range m.state.Plots {
		name, emoji, _, _, _, tier := game.GetFlowerInfo(plot.FlowerType)
		progress := m.state.PlotProgress(i)
		ready := progress >= 1.0

		// Greenhouse indicator
		ghIcon := ""
		if plot.IsGreenhouse {
			ghIcon = "🏠"
		}

		// Progress bar with gradient
		barWidth := 20
		filled := int(progress * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}

		var bar string
		if ready {
			// Pulsing effect based on time
			pulse := math.Sin(float64(time.Now().UnixMilli())/300) * 0.5
			if pulse > 0 {
				bar = lipgloss.NewStyle().Foreground(readyColor).Bold(true).Render(
					strings.Repeat("█", barWidth))
			} else {
				bar = lipgloss.NewStyle().Foreground(lipgloss.Color("#00CC55")).Render(
					strings.Repeat("█", barWidth))
			}
		} else {
			// Color based on season/tier
			growFg := growColor
			seasonMult := m.state.SeasonGrowthMult(tier)
			if seasonMult > 1.2 && !plot.IsGreenhouse {
				growFg = p.primary // boosted = season color
			} else if seasonMult < 0.8 && !plot.IsGreenhouse {
				growFg = warnColor // penalized = red
			}
			bar = lipgloss.NewStyle().Foreground(growFg).Render(
				strings.Repeat("█", filled)) +
				lipgloss.NewStyle().Foreground(dimColor).Render(
					strings.Repeat("░", barWidth-filled))
		}

		status := ""
		if ready {
			status = lipgloss.NewStyle().Foreground(readyColor).Bold(true).Render(" ✿ READY!")
		} else {
			growTime := m.state.EffectiveGrowTimeForPlot(i)
			remaining := growTime - time.Since(plot.Planted)
			if remaining < 0 {
				remaining = 0
			}
			status = lipgloss.NewStyle().Foreground(dimColor).Render(
				fmt.Sprintf(" %s", formatDuration(remaining)))
		}

		// Season modifier indicator
		seasonInd := ""
		if !plot.IsGreenhouse {
			seasonMult := m.state.SeasonGrowthMult(tier)
			if seasonMult > 1.1 {
				seasonInd = lipgloss.NewStyle().Foreground(p.primary).Render(" ▲")
			} else if seasonMult < 0.9 {
				seasonInd = lipgloss.NewStyle().Foreground(warnColor).Render(" ▼")
			}
		}

		line := fmt.Sprintf(" %s%s %s [%s]%s%s",
			ghIcon, emoji, padRight(name, 14), bar, status, seasonInd)

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}

		lines = append(lines, line)
	}

	return m.boxStyle().Render(strings.Join(lines, "\n"))
}

// renderGardenBed draws a visual representation of the garden.
func (m Model) renderGardenBed() string {
	var rows []string

	// Top fence
	width := len(m.state.Plots)*4 + 3
	if width < 20 {
		width = 20
	}

	p := m.palette()
	fenceColor := p.secondary

	fence := lipgloss.NewStyle().Foreground(fenceColor).Render(
		"╔" + strings.Repeat("═", width) + "╗")
	rows = append(rows, " "+fence)

	// Flower display row
	var flowers []string
	for i, plot := range m.state.Plots {
		_, emoji, _, _, _, _ := game.GetFlowerInfo(plot.FlowerType)
		progress := m.state.PlotProgress(i)
		if progress >= 1.0 {
			// Full bloom with sparkle
			flowers = append(flowers, " "+emoji+" ")
		} else if progress > 0.5 {
			// Growing
			flowers = append(flowers, " "+emoji+" ")
		} else if progress > 0.1 {
			// Sprouting
			flowers = append(flowers, " 🌱 ")
		} else {
			// Just planted
			flowers = append(flowers, " · ")
		}
	}

	flowerRow := lipgloss.NewStyle().Foreground(fenceColor).Render("║") +
		" " + strings.Join(flowers, "") + " " +
		lipgloss.NewStyle().Foreground(fenceColor).Render("║")
	rows = append(rows, " "+flowerRow)

	// Ground
	season := m.state.CurrentSeason()
	var groundChar string
	switch season {
	case game.Spring:
		groundChar = "~"
	case game.Summer:
		groundChar = "."
	case game.Autumn:
		groundChar = ","
	case game.Winter:
		groundChar = "*"
	}
	groundColor := p.accent
	ground := lipgloss.NewStyle().Foreground(fenceColor).Render("╚") +
		lipgloss.NewStyle().Foreground(groundColor).Render(
			strings.Repeat(groundChar, width)) +
		lipgloss.NewStyle().Foreground(fenceColor).Render("╝")
	rows = append(rows, " "+ground)

	return strings.Join(rows, "\n")
}

func (m Model) renderUpgrades() string {
	var lines []string
	p := m.palette()
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(p.primary).Render("⚒  Upgrades"))
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
			levelStr := ""
			if level > 0 {
				levelStr = lipgloss.NewStyle().Foreground(p.accent).Render(
					fmt.Sprintf(" [%d]", level))
			}

			desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render(u.Description)
			costStr := formatNumber(cost)

			costDisplay := ""
			if affordable {
				costDisplay = lipgloss.NewStyle().Foreground(readyColor).Bold(true).Render(
					fmt.Sprintf("  %s 🌸", costStr))
			} else {
				costDisplay = lipgloss.NewStyle().Foreground(dimColor).Render(
					fmt.Sprintf("  %s 🌸", costStr))
			}

			line = fmt.Sprintf(" %s%s  %s%s", padRight(u.Name, 18), levelStr, desc, costDisplay)
		}

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		lines = append(lines, line)
	}

	// Essence upgrades section (if any essence)
	if m.state.Essence > 0 || m.state.Prestige2Count > 0 {
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(essenceColor).Render(
			"💫 Essence Upgrades (permanent)"))
		lines = append(lines, "")

		baseIdx := len(game.Upgrades)
		for i, eu := range game.EssenceUpgrades {
			level := m.state.EssenceUpgrades[eu.Effect]
			maxed := eu.MaxLevel > 0 && level >= eu.MaxLevel
			affordable := m.state.Essence >= eu.Cost

			var line string
			if maxed {
				line = fmt.Sprintf(" ✓ %s (MAX)", eu.Name)
				line = lipgloss.NewStyle().Foreground(dimColor).Render(line)
			} else {
				levelStr := ""
				if level > 0 {
					levelStr = lipgloss.NewStyle().Foreground(essenceColor).Render(
						fmt.Sprintf(" [%d]", level))
				}
				desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render(eu.Description)
				costDisplay := ""
				if affordable {
					costDisplay = lipgloss.NewStyle().Foreground(readyColor).Bold(true).Render(
						fmt.Sprintf("  %.0f 💫", eu.Cost))
				} else {
					costDisplay = lipgloss.NewStyle().Foreground(dimColor).Render(
						fmt.Sprintf("  %.0f 💫", eu.Cost))
				}
				line = fmt.Sprintf(" %s%s  %s%s", padRight(eu.Name, 20), levelStr, desc, costDisplay)
			}

			if baseIdx+i == m.cursor {
				line = selectedStyle.Render(line)
			}
			lines = append(lines, line)
		}
	}

	return m.boxStyle().Render(strings.Join(lines, "\n"))
}

func (m Model) renderFlowers() string {
	var lines []string
	p := m.palette()

	header := "🌺 Flower Collection"
	if m.selectedPlot >= 0 && m.selectedPlot < len(m.state.Plots) {
		header = fmt.Sprintf("🌺 Choose flower for Plot %d", m.selectedPlot+1)
	}
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(p.primary).Render(header))
	lines = append(lines, "")

	// Base flowers
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render(" ── Base Flowers ──"))
	for i, ft := range game.FlowerTypes {
		lines = append(lines, m.renderFlowerLine(i, ft.Name, ft.Emoji, ft.PetalYield,
			ft.SeedCost, ft.UnlockCost, ft.Tier, i == m.cursor))
	}

	// Hybrid flowers (only show discovered or hints)
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(hybridColor).Render(" ── Hybrid Flowers ──"))
	for i, h := range game.HybridFlowers {
		globalIdx := len(game.FlowerTypes) + i
		cursorIdx := len(game.FlowerTypes) + i
		if m.state.Unlocked[globalIdx] {
			lines = append(lines, m.renderFlowerLine(globalIdx, h.Name, h.Emoji, h.PetalYield,
				h.SeedCost, 0, h.Tier, cursorIdx == m.cursor))
		} else {
			// Show hint
			hint := lipgloss.NewStyle().Foreground(dimColor).Italic(true).Render(
				fmt.Sprintf(" 🔮 ???  (Breed adjacent flowers to discover)"))
			if cursorIdx == m.cursor {
				hint = selectedStyle.Render(hint)
			}
			lines = append(lines, hint)
		}
	}

	return m.boxStyle().Render(strings.Join(lines, "\n"))
}

func (m Model) renderFlowerLine(globalIdx int, name, emoji string, yield, seedCost, unlockCost float64, tier int, selected bool) string {
	unlocked := m.state.Unlocked[globalIdx]
	season := m.state.CurrentSeason()
	seasonInfo := game.Seasons[season]

	var line string
	if unlocked {
		growTime := m.state.EffectiveGrowTime(globalIdx)
		effectiveYield := (yield + m.state.FlatBonus()) * m.state.PetalMultiplier()
		yieldMult, ok := seasonInfo.YieldMult[tier]
		if ok {
			effectiveYield *= yieldMult
		}

		// Season indicator
		growMult, _ := seasonInfo.GrowthMult[tier]
		seasonInd := ""
		if growMult > 1.1 {
			seasonInd = lipgloss.NewStyle().Foreground(readyColor).Render(" ▲ in season")
		} else if growMult < 0.9 {
			seasonInd = lipgloss.NewStyle().Foreground(warnColor).Render(" ▼ off season")
		}

		yieldStr := lipgloss.NewStyle().Foreground(petalColor).Render(fmt.Sprintf("%.0f", effectiveYield))
		timeStr := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render(formatDuration(growTime))
		line = fmt.Sprintf(" %s %s  ⚘ %s  ⏱ %s", emoji, padRight(name, 16), yieldStr, timeStr)
		if seedCost > 0 {
			line += lipgloss.NewStyle().Foreground(seedColor).Render(
				fmt.Sprintf("  🌱%.0f", seedCost))
		}
		line += seasonInd
	} else if unlockCost > 0 {
		line = fmt.Sprintf(" 🔒 %s  Unlock: %s petals", padRight(name, 16), formatNumber(unlockCost))
		if m.state.Petals >= unlockCost {
			line = lipgloss.NewStyle().Foreground(readyColor).Render(line)
		} else {
			line = lipgloss.NewStyle().Foreground(dimColor).Render(line)
		}
	}

	if selected {
		line = selectedStyle.Render(line)
	}
	return line
}

func (m Model) renderBreeding() string {
	var lines []string
	p := m.palette()
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(hybridColor).Render(
		"🧬 Hybrid Breeding"))
	lines = append(lines, "")

	lines = append(lines, lipgloss.NewStyle().Foreground(p.secondary).Render(
		" Place parent flowers in adjacent plots to discover hybrids!"))
	lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
		" Hybrids are discovered on harvest when parents are neighbors."))
	lines = append(lines, "")

	// Show recipes (discovered or as hints)
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Bold(true).Render(
		" Known Recipes:"))
	lines = append(lines, "")

	for i, recipe := range game.HybridRecipes {
		globalIdx := len(game.FlowerTypes) + recipe.ResultIdx
		discovered := m.state.Unlocked[globalIdx]

		parent1 := game.FlowerTypes[recipe.Parent1]
		parent2 := game.FlowerTypes[recipe.Parent2]

		var line string
		if discovered {
			hybrid := game.HybridFlowers[recipe.ResultIdx]
			line = fmt.Sprintf(" %s %s + %s %s  →  %s %s",
				parent1.Emoji, parent1.Name, parent2.Emoji, parent2.Name,
				hybrid.Emoji, hybrid.Name)
			line = lipgloss.NewStyle().Foreground(hybridColor).Render(line)

			// Show lore
			lore := lipgloss.NewStyle().Foreground(dimColor).Italic(true).Render(
				fmt.Sprintf("   \"%s\"", hybrid.Lore))
			line += "\n" + lore
		} else {
			// Show parents as hint if both are unlocked
			if m.state.Unlocked[recipe.Parent1] && m.state.Unlocked[recipe.Parent2] {
				line = fmt.Sprintf(" %s %s + %s %s  →  🔮 ???",
					parent1.Emoji, parent1.Name, parent2.Emoji, parent2.Name)
				line = lipgloss.NewStyle().Foreground(dimColor).Render(line)
			} else {
				line = lipgloss.NewStyle().Foreground(dimColor).Render(
					" ??? + ???  →  🔮 ???")
			}
		}

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		lines = append(lines, line)
	}

	// Pollination upgrade status
	polLevel := m.state.UpgradeLevels["pollination"]
	if polLevel > 0 {
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Foreground(readyColor).Render(
			fmt.Sprintf(" 🐝 Pollination Lv.%d (+%d%% discovery chance)", polLevel, polLevel*25)))
	}

	// Discovery count
	discovered := 0
	for _, h := range game.HybridFlowers {
		_ = h
		discovered = len(m.state.DiscoveredHybrids)
		break
	}
	total := len(game.HybridFlowers)
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(p.primary).Bold(true).Render(
		fmt.Sprintf(" Discovered: %d / %d", discovered, total)))

	return m.boxStyle().Render(strings.Join(lines, "\n"))
}

func (m Model) renderPrestige() string {
	var lines []string
	p := m.palette()

	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(p.primary).Render(
		"✨ Prestige — The Eternal Garden"))
	lines = append(lines, "")

	// Layer 1: Nectar
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(nectarColor).Render(
		" ─── Layer 1: Nectar ───"))
	lines = append(lines, fmt.Sprintf(" Lifetime petals: %s",
		lipgloss.NewStyle().Foreground(petalColor).Render(formatNumber(m.state.TotalPetals))))
	lines = append(lines, fmt.Sprintf(" Current nectar:  %s",
		lipgloss.NewStyle().Foreground(nectarColor).Render(formatNumber(m.state.Nectar))))
	lines = append(lines, fmt.Sprintf(" Prestige count:  %d", m.state.PrestigeCount))
	lines = append(lines, "")

	nectarGain := m.state.NectarFromPrestige()
	if nectarGain > 0 {
		gainStr := lipgloss.NewStyle().Foreground(nectarColor).Bold(true).Render(
			fmt.Sprintf("%.0f", nectarGain))
		if m.cursor == 0 {
			lines = append(lines, selectedStyle.Render(
				fmt.Sprintf(" ► Prestige for %s nectar (+10%% yield)", gainStr)))
		} else {
			lines = append(lines, fmt.Sprintf(" ► Press ENTER on this to prestige for %s nectar", gainStr))
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
			"   Resets: garden, upgrades, flowers. Keeps: nectar, essence upgrades"))
	} else {
		needed := 10000 - m.state.TotalPetals
		if needed < 0 {
			needed = 0
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
			fmt.Sprintf(" Need %s more lifetime petals to prestige", formatNumber(needed))))
	}

	lines = append(lines, "")

	// Layer 2: Essence
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(essenceColor).Render(
		" ─── Layer 2: Transcendence ───"))
	lines = append(lines, fmt.Sprintf(" Current essence:    %s",
		lipgloss.NewStyle().Foreground(essenceColor).Render(formatNumber(m.state.Essence))))
	lines = append(lines, fmt.Sprintf(" Transcendence count: %d", m.state.Prestige2Count))
	lines = append(lines, "")

	essenceGain := m.state.EssenceFromPrestige2()
	if essenceGain > 0 {
		gainStr := lipgloss.NewStyle().Foreground(essenceColor).Bold(true).Render(
			fmt.Sprintf("%.0f", essenceGain))
		if m.cursor == 1 {
			lines = append(lines, selectedStyle.Render(
				fmt.Sprintf(" ► Transcend for %s essence", gainStr)))
		} else {
			lines = append(lines, fmt.Sprintf(" ► Transcend for %s essence (select & ENTER)", gainStr))
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
			"   Resets: EVERYTHING except essence & essence upgrades"))
	} else {
		lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render(
			" Need 50+ nectar to transcend"))
	}

	// Bonuses summary
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render(
		" Current bonuses:"))
	if m.state.PrestigeCount > 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(nectarColor).Render(
			fmt.Sprintf("   ✨ +%d%% petal yield (from %d prestiges)",
				m.state.PrestigeCount*10, m.state.PrestigeCount)))
	}
	for _, eu := range game.EssenceUpgrades {
		level := m.state.EssenceUpgrades[eu.Effect]
		if level > 0 {
			lines = append(lines, lipgloss.NewStyle().Foreground(essenceColor).Render(
				fmt.Sprintf("   💫 %s Lv.%d — %s", eu.Name, level, eu.Description)))
		}
	}

	return m.boxStyle().Render(strings.Join(lines, "\n"))
}

func (m Model) renderEventBanner() string {
	if !m.state.IsEventActive() {
		return ""
	}
	info := game.GetEventInfo(m.state.ActiveEvent.Type)
	remaining := m.state.EventTimeRemaining()

	bgColor := lipgloss.Color("#2F4F2F")
	fgColor := lipgloss.Color("#00FF7F")
	if !info.Positive {
		bgColor = lipgloss.Color("#4F2F2F")
		fgColor = lipgloss.Color("#FF6347")
	}

	bannerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(fgColor).
		Background(bgColor).
		Padding(0, 2)

	return bannerStyle.Render(fmt.Sprintf(" %s %s — %s [%s] ",
		info.Emoji, info.Name, info.Description, formatDuration(remaining)))
}

func (m Model) renderAchievements() string {
	var lines []string
	p := m.palette()

	completed, total := game.AchievementProgress(m.state)
	header := fmt.Sprintf("🏆 Achievements (%d/%d)", completed, total)
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(p.primary).Render(header))
	lines = append(lines, "")

	// Show best combo
	if m.state.BestCombo > 1 {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500")).Render(
			fmt.Sprintf(" Best combo: %dx", m.state.BestCombo)))
		lines = append(lines, "")
	}

	for i, a := range game.Achievements {
		done := m.state.CompletedAchievements[a.ID]

		var line string
		if done {
			rewardStr := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render(
				fmt.Sprintf("+%.0f %s", a.Reward, a.RewardType))
			line = fmt.Sprintf(" %s %s  %s  %s",
				a.Emoji,
				lipgloss.NewStyle().Foreground(nectarColor).Render(a.Name),
				lipgloss.NewStyle().Foreground(dimColor).Render(a.Description),
				rewardStr)
		} else if a.Hidden {
			line = lipgloss.NewStyle().Foreground(dimColor).Render(" 🔮 ???  Hidden achievement")
		} else {
			line = fmt.Sprintf(" ○ %s  %s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(a.Name),
				lipgloss.NewStyle().Foreground(dimColor).Render(a.Description))
		}

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		lines = append(lines, line)
	}

	return m.boxStyle().Render(strings.Join(lines, "\n"))
}

func (m Model) renderLog() string {
	if len(m.state.Log) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, lipgloss.NewStyle().Foreground(dimColor).Render("── Recent ──"))

	start := len(m.state.Log) - 5
	if start < 0 {
		start = 0
	}
	for _, entry := range m.state.Log[start:] {
		age := time.Since(entry.Time)
		timeStr := lipgloss.NewStyle().Foreground(dimColor).Render(
			fmt.Sprintf("[%s ago]", formatDuration(age)))
		color := lipgloss.Color(entry.Color)
		if entry.Color == "" {
			color = dimColor
		}
		msg := lipgloss.NewStyle().Foreground(color).Render(entry.Message)
		lines = append(lines, fmt.Sprintf(" %s %s", timeStr, msg))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderHelp() string {
	var help string
	switch m.activeTab {
	case tabGarden:
		help = "↑/↓ select  •  ENTER harvest  •  p plant  •  TAB next tab  •  s save  •  q quit"
	case tabUpgrades:
		help = "↑/↓ select  •  ENTER buy  •  TAB next tab  •  s save  •  q quit"
	case tabFlowers:
		help = "↑/↓ select  •  ENTER unlock/plant  •  TAB next tab  •  s save  •  q quit"
	case tabBreeding:
		help = "↑/↓ browse  •  TAB next tab  •  s save  •  q quit"
	case tabPrestige:
		help = "↑/↓ select layer  •  ENTER prestige/transcend  •  TAB next tab  •  s save  •  q quit"
	case tabAchievements:
		help = "↑/↓ browse  •  TAB next tab  •  s save  •  q quit"
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
	if m < 60 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	h := m / 60
	m = m % 60
	return fmt.Sprintf("%dh%dm", h, m)
}

func formatNumber(n float64) string {
	if n < 0 {
		return fmt.Sprintf("-%.0f", -n)
	}
	if n < 1000 {
		return fmt.Sprintf("%.0f", n)
	}
	suffixes := []string{"", "K", "M", "B", "T", "Qa", "Qi"}
	i := int(math.Log10(n) / 3)
	if i >= len(suffixes) {
		i = len(suffixes) - 1
	}
	val := n / math.Pow(10, float64(i*3))
	return fmt.Sprintf("%.1f%s", val, suffixes[i])
}
