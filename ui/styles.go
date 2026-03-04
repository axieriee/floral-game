package ui

import "github.com/charmbracelet/lipgloss"

// Floral color palette — soft pastels
var (
	Rose     = lipgloss.Color("#e8a0bf")
	Lavender = lipgloss.Color("#b4a7d6")
	Sage     = lipgloss.Color("#a3be8c")
	Peach    = lipgloss.Color("#f2b88a")
	Sky      = lipgloss.Color("#89b4c4")
	Cream    = lipgloss.Color("#e8dcc8")
	Dust     = lipgloss.Color("#8a7f72")
	Plum     = lipgloss.Color("#6b5b7b")
	Muted    = lipgloss.Color("#9a8f82")
	Faint    = lipgloss.Color("#5c5555")
	Dim      = lipgloss.Color("#3d3636")
	BgDark   = lipgloss.Color("#1e1e2e")
	White    = lipgloss.Color("#d4c4b0")
)

// Styles
var (
	// Title screen
	TitleStyle = lipgloss.NewStyle().
			Foreground(Rose).
			Bold(true).
			Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Lavender).
			Italic(true).
			Align(lipgloss.Center)

	// Journal box border
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Lavender).
			Padding(1, 2)

	// Inner content box
	InnerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Dust).
			Padding(0, 1)

	// Prose / narrative text
	ProseStyle = lipgloss.NewStyle().
			Foreground(Cream).
			PaddingLeft(2)

	// Soft italic flavor text
	FlavorStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true).
			PaddingLeft(2)

	// Menu choice (unselected)
	ChoiceStyle = lipgloss.NewStyle().
			Foreground(Dust).
			PaddingLeft(4)

	// Menu choice (selected / cursor)
	SelectedStyle = lipgloss.NewStyle().
			Foreground(Rose).
			Bold(true).
			PaddingLeft(2)

	// Header bar
	HeaderStyle = lipgloss.NewStyle().
			Foreground(Lavender).
			Bold(true).
			Align(lipgloss.Center)

	// Flower name styles by rarity
	CommonFlowerStyle = lipgloss.NewStyle().
				Foreground(Sage)

	UncommonFlowerStyle = lipgloss.NewStyle().
				Foreground(Lavender)

	RareFlowerStyle = lipgloss.NewStyle().
			Foreground(Rose).
			Bold(true)

	MythicFlowerStyle = lipgloss.NewStyle().
				Foreground(Peach).
				Bold(true)

	// Resonance label
	ResonanceStyle = lipgloss.NewStyle().
			Foreground(Sky).
			Italic(true)

	// NPC dialogue
	DialogueStyle = lipgloss.NewStyle().
			Foreground(Peach).
			Italic(true).
			PaddingLeft(2)

	// NPC name
	NPCNameStyle = lipgloss.NewStyle().
			Foreground(Peach).
			Bold(true)

	// Trust pips
	TrustFullStyle = lipgloss.NewStyle().
			Foreground(Rose)

	TrustEmptyStyle = lipgloss.NewStyle().
			Foreground(Faint)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Align(lipgloss.Center)

	// Potency pips
	PotencyFullStyle = lipgloss.NewStyle().
				Foreground(Sage)

	PotencyEmptyStyle = lipgloss.NewStyle().
				Foreground(Faint)

	// Separator
	SepStyle = lipgloss.NewStyle().
			Foreground(Dust)

	// Blend name
	BlendStyle = lipgloss.NewStyle().
			Foreground(Rose).
			Bold(true)

	// Warning / damage text
	WarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c97070")).
			Italic(true)

	// Success text
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Sage).
			Bold(true)

	// Dim / disabled
	DimStyle = lipgloss.NewStyle().
			Foreground(Faint)

	// Perk text
	PerkStyle = lipgloss.NewStyle().
			Foreground(Sage)

	// Drawback text
	DrawbackStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c97070"))
)

// Helper to render potency pips
func RenderPotency(current, max int) string {
	s := ""
	for i := 0; i < max; i++ {
		if i < current {
			s += PotencyFullStyle.Render("❀")
		} else {
			s += PotencyEmptyStyle.Render("◌")
		}
	}
	return s
}

// Helper to render trust hearts
func RenderTrust(level int) string {
	s := ""
	for i := 0; i < 5; i++ {
		if i < level {
			s += TrustFullStyle.Render("♥")
		} else {
			s += TrustEmptyStyle.Render("♡")
		}
	}
	return s
}

// Separator line with flowers
func FloralSeparator(width int) string {
	if width < 10 {
		width = 50
	}
	left := "─✿─"
	right := "─✿─"
	mid := ""
	remaining := width - 8
	for i := 0; i < remaining; i++ {
		mid += "─"
	}
	return SepStyle.Render(left + mid + right)
}
