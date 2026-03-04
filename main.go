package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"floragame/data"
	"floragame/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Game States ---

type gameScreen int

const (
	screenTitle gameScreen = iota
	screenCharCreate
	screenMeadow
	screenForage
	screenForageResult
	screenJournal
	screenBlend
	screenBlendResult
	screenNPCIntro
	screenNPCOffer
	screenNPCReaction
	screenSummary
)

// --- Pressed Flower (in journal) ---

type pressedFlower struct {
	flower  data.Flower
	damaged bool
}

// --- Main Model ---

type model struct {
	screen      gameScreen
	cursor      int
	width       int
	height      int
	rng         *rand.Rand
	day         int
	weather     string
	weathers    []string

	// Character
	tradition     *data.Tradition
	journalPages  []pressedFlower
	knownBlends   map[string]bool
	discoveredBlends []string

	// Meadow
	forageSpots   []data.ForageSpot
	currentSpot   int
	lastForaged   *pressedFlower
	forageMessage string

	// Blending
	blendSelect1  int
	blendSelect2  int
	blendPhase    int // 0=select first, 1=select second
	blendResult   string

	// NPC
	npc           data.NPC
	npcOffering   int // index into journalPages, -1 = none
	npcReaction   string
	npcTrustDelta int

	// Summary
	flowersFound  int
	blendsFound   int
}

func initialModel() model {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	weathers := []string{
		"Soft Rain", "Clear Skies", "Gentle Breeze", "Overcast",
		"Morning Mist", "Golden Hour", "Warm & Humid", "Cool & Crisp",
	}

	m := model{
		screen:       screenTitle,
		rng:          r,
		day:          1,
		weather:      weathers[r.Intn(len(weathers))],
		weathers:     weathers,
		knownBlends:  make(map[string]bool),
		npcOffering:  -1,
	}

	// Deep copy forage spots
	m.forageSpots = make([]data.ForageSpot, len(data.MeadowForageSpots))
	for i, s := range data.MeadowForageSpots {
		m.forageSpots[i] = data.ForageSpot{
			Name:     s.Name,
			Desc:     s.Desc,
			Flowers:  append([]string{}, s.Flowers...),
			Searched: false,
		}
	}

	// Deep copy NPC
	m.npc = data.NPC{
		Name:        data.Maren.Name,
		Title:       data.Maren.Title,
		Greeting:    data.Maren.Greeting,
		Trust:       data.Maren.Trust,
		Disposition: data.Maren.Disposition,
		Dialogue:    make(map[string]string),
	}
	for k, v := range data.Maren.Dialogue {
		m.npc.Dialogue[k] = v
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

// --- Update ---

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.screen {
	case screenTitle:
		return m.updateTitle(msg)
	case screenCharCreate:
		return m.updateCharCreate(msg)
	case screenMeadow:
		return m.updateMeadow(msg)
	case screenForage:
		return m.updateForage(msg)
	case screenForageResult:
		return m.updateForageResult(msg)
	case screenJournal:
		return m.updateJournal(msg)
	case screenBlend:
		return m.updateBlend(msg)
	case screenBlendResult:
		return m.updateBlendResult(msg)
	case screenNPCIntro:
		return m.updateNPCIntro(msg)
	case screenNPCOffer:
		return m.updateNPCOffer(msg)
	case screenNPCReaction:
		return m.updateNPCReaction(msg)
	case screenSummary:
		return m.updateSummary(msg)
	}
	return m, nil
}

// --- Title Screen ---

func (m model) updateTitle(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" || key.String() == " " {
			m.screen = screenCharCreate
			m.cursor = 2 // Pre-select Petal Scholar (index 2)
		}
	}
	return m, nil
}

func (m model) viewTitle() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	flower := `
        *
       /|\
      / | \
     /  |  \
    *.  |  .*
   /  '.|.'  \
  *    .|.    *
   \  '.|.'  /
    *.  |  .*
     \  |  /
      \ | /
       \|/
        |
       /|\
      /_|_\
`

	art := ui.TitleStyle.Width(w).Render(flower)
	title := ui.TitleStyle.Width(w).Render("✿  F L O R A V A L E  ✿")
	sub := ui.SubtitleStyle.Width(w).Render("A Botanical Text RPG")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))
	flavor := ui.FlavorStyle.Width(w).Render(
		"Every world is new. Every garden, yours.\n" +
			"Press flowers into your journal. Blend their resonances.\n" +
			"Befriend or beware. The meadow remembers.")
	prompt := ui.StatusBarStyle.Width(w).Render("\n[ Press Enter to begin your journey ]")

	return lipgloss.JoinVertical(lipgloss.Center,
		"", art, title, sub, sep, "", flavor, "", prompt, "")
}

// --- Character Creation ---

func (m model) updateCharCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(data.AllTraditions)-1 {
				m.cursor++
			}
		case "enter":
			m.tradition = &data.AllTraditions[m.cursor]
			// Grant starting flowers
			for _, name := range m.tradition.StartFlowers {
				if f, ok := data.AllFlowers[name]; ok {
					m.journalPages = append(m.journalPages, pressedFlower{flower: f})
				}
			}
			// Grant known blends
			for _, name := range m.tradition.KnownBlends {
				m.knownBlends[name] = true
			}
			m.screen = screenMeadow
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) viewCharCreate() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render("✿ Choose Your Tradition ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	var choices []string
	for i, t := range data.AllTraditions {
		nameStr := t.Symbol + " " + strings.ToUpper(t.Name)
		if i == m.cursor {
			nameStr = ui.SelectedStyle.Render("› " + nameStr)
		} else {
			nameStr = ui.ChoiceStyle.Render("  " + nameStr)
		}
		choices = append(choices, nameStr)

		if i == m.cursor {
			desc := ui.ProseStyle.Render(t.Desc)
			flavor := ui.FlavorStyle.Render(t.Flavor)
			perk := ui.PerkStyle.Render("    ✓ " + t.Perk)
			drawback := ui.DrawbackStyle.Render("    ✗ " + t.Drawback)

			startFlowers := "    Starts with: "
			for j, fn := range t.StartFlowers {
				if j > 0 {
					startFlowers += ", "
				}
				if f, ok := data.AllFlowers[fn]; ok {
					startFlowers += f.Rarity.Symbol() + " " + fn
				}
			}
			startF := ui.DimStyle.Render(startFlowers)

			blendStr := "    Known blends: "
			for j, bn := range t.KnownBlends {
				if j > 0 {
					blendStr += ", "
				}
				blendStr += bn
			}
			blendS := ui.DimStyle.Render(blendStr)

			choices = append(choices, desc, flavor, perk, drawback, startF, blendS, "")
		}
	}

	prompt := ui.StatusBarStyle.Width(w).Render("↑/↓ navigate · Enter to choose")

	sections := []string{"", header, sep, ""}
	sections = append(sections, choices...)
	sections = append(sections, "", prompt)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- Meadow (Main Hub) ---

func (m model) updateMeadow(msg tea.Msg) (tea.Model, tea.Cmd) {
	maxChoices := m.meadowChoiceCount()
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < maxChoices-1 {
				m.cursor++
			}
		case "enter":
			action := m.meadowAction(m.cursor)
			switch action {
			case "forage":
				m.screen = screenForage
				m.cursor = 0
			case "journal":
				m.screen = screenJournal
				m.cursor = 0
			case "blend":
				m.screen = screenBlend
				m.cursor = 0
				m.blendPhase = 0
				m.blendSelect1 = -1
				m.blendSelect2 = -1
			case "npc":
				m.screen = screenNPCIntro
				m.cursor = 0
			case "endday":
				m.day++
				m.weather = m.weathers[m.rng.Intn(len(m.weathers))]
				// Potency decay
				for i := range m.journalPages {
					if m.journalPages[i].flower.Potency > 0 && m.journalPages[i].flower.Rarity != data.Mythic {
						if m.rng.Float64() < 0.3 {
							m.journalPages[i].flower.Potency--
						}
					}
				}
			case "summary":
				m.screen = screenSummary
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m model) meadowChoiceCount() int {
	count := 0
	// Forage spots that aren't searched
	hasForage := false
	for _, s := range m.forageSpots {
		if !s.Searched {
			hasForage = true
			break
		}
	}
	if hasForage {
		count++ // forage
	}
	count++ // journal
	if len(m.journalPages) >= 2 {
		count++ // blend
	}
	count++ // visit maren
	count++ // end day
	count++ // end journey
	return count
}

func (m model) meadowAction(idx int) string {
	actions := []string{}
	hasForage := false
	for _, s := range m.forageSpots {
		if !s.Searched {
			hasForage = true
			break
		}
	}
	if hasForage {
		actions = append(actions, "forage")
	}
	actions = append(actions, "journal")
	if len(m.journalPages) >= 2 {
		actions = append(actions, "blend")
	}
	actions = append(actions, "npc", "endday", "summary")
	if idx < len(actions) {
		return actions[idx]
	}
	return ""
}

func (m model) viewMeadow() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render(
		fmt.Sprintf("✿ Clover Meadows ─── Day %d · %s ✿", m.day, m.weather))
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	var prose string
	switch {
	case m.day == 1:
		prose = "The path opens into a wide meadow. Wildflowers nod in a\n" +
			"  gentle breeze. To the south, a cottage with a thin plume\n" +
			"  of chimney smoke. The air smells of honey and warm grass.\n\n" +
			"  Your journal is open in your hands — a few pressed\n" +
			"  specimens from your training, but mostly blank pages\n" +
			"  waiting to be filled."
	case strings.Contains(m.weather, "Rain"):
		prose = "Soft rain dimples the stream and beads on flower petals.\n" +
			"  The meadow smells richer today — earth and green things\n" +
			"  drinking deep. Maren's chimney smoke is thicker,\n" +
			"  fighting the damp."
	default:
		prose = "The meadow stretches around you, familiar now but still\n" +
			"  full of surprises. Birds call from the hedgerow.\n" +
			"  You feel the weight of your journal, its pages growing\n" +
			"  heavier with pressed flowers."
	}
	proseR := ui.ProseStyle.Render(prose)

	// Status line
	flowerCount := len(m.journalPages)
	statusStr := fmt.Sprintf("%s %s · ⚘ %d flowers · %s",
		m.tradition.Symbol, m.tradition.Name, flowerCount, m.weather)
	status := ui.StatusBarStyle.Width(w).Render(statusStr)

	sep2 := ui.SepStyle.Render(ui.FloralSeparator(w))

	// Choices
	choiceLabel := ui.FlavorStyle.Render("  ❀ What do you do?")
	var choices []string
	idx := 0

	hasForage := false
	for _, s := range m.forageSpots {
		if !s.Searched {
			hasForage = true
			break
		}
	}
	if hasForage {
		choices = append(choices, m.renderChoice(idx, "Forage in the meadow"))
		idx++
	}
	choices = append(choices, m.renderChoice(idx, "Open your journal"))
	idx++
	if len(m.journalPages) >= 2 {
		choices = append(choices, m.renderChoice(idx, "Try blending flowers"))
		idx++
	}
	choices = append(choices, m.renderChoice(idx, "Visit Maren's cottage"))
	idx++
	choices = append(choices, m.renderChoice(idx, "Rest until tomorrow"))
	idx++
	choices = append(choices, m.renderChoice(idx, "End your journey"))

	prompt := ui.DimStyle.Render("    ↑/↓ navigate · Enter to choose")

	sections := []string{"", header, sep, "", proseR, "", status, sep2, "", choiceLabel, ""}
	sections = append(sections, choices...)
	sections = append(sections, "", prompt)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- Forage Screen ---

func (m model) updateForage(msg tea.Msg) (tea.Model, tea.Cmd) {
	available := m.availableForageSpots()
	maxChoices := len(available) + 1 // +1 for "go back"

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < maxChoices-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(available) {
				// Go back
				m.screen = screenMeadow
				m.cursor = 0
			} else {
				spotIdx := available[m.cursor]
				m.currentSpot = spotIdx
				spot := &m.forageSpots[spotIdx]
				spot.Searched = true

				// Pick a random flower from the spot
				flowerName := spot.Flowers[m.rng.Intn(len(spot.Flowers))]
				f := data.AllFlowers[flowerName]

				// Petal Scholar drawback: chance of damaging first forages
				damaged := false
				if m.tradition.Name == "Petal Scholar" && m.flowersFound < 2 {
					if m.rng.Float64() < 0.4 {
						damaged = true
						f.Potency = max(1, f.Potency-2)
					}
				}

				pf := pressedFlower{flower: f, damaged: damaged}
				m.lastForaged = &pf
				m.journalPages = append(m.journalPages, pf)
				m.flowersFound++

				if damaged {
					m.forageMessage = "Your hands are clumsy — ink-stained, not earth-stained.\nThe specimen is damaged, but salvageable."
				} else {
					messages := []string{
						"You kneel carefully and press it between the pages.",
						"It comes away cleanly. A perfect specimen.",
						"You hold it up to the light before pressing it gently.",
						"The petals are soft between your fingers as you lay it flat.",
					}
					m.forageMessage = messages[m.rng.Intn(len(messages))]
				}

				m.screen = screenForageResult
				m.cursor = 0
			}
		case "esc", "q":
			m.screen = screenMeadow
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) availableForageSpots() []int {
	var available []int
	for i, s := range m.forageSpots {
		if !s.Searched {
			available = append(available, i)
		}
	}
	return available
}

func (m model) viewForage() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render("✿ Foraging ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	prose := ui.ProseStyle.Render(
		"You scan the meadow for interesting growth.\n" +
			"  Several spots catch your eye...")

	available := m.availableForageSpots()
	var choices []string
	for i, spotIdx := range available {
		spot := m.forageSpots[spotIdx]
		if i == m.cursor {
			choices = append(choices, ui.SelectedStyle.Render("› "+spot.Name))
			choices = append(choices, ui.FlavorStyle.Render("  "+spot.Desc))
			choices = append(choices, "")
		} else {
			choices = append(choices, ui.ChoiceStyle.Render("  "+spot.Name))
		}
	}

	backIdx := len(available)
	if m.cursor == backIdx {
		choices = append(choices, ui.SelectedStyle.Render("› Go back"))
	} else {
		choices = append(choices, ui.ChoiceStyle.Render("  Go back"))
	}

	sections := []string{"", header, sep, "", prose, ""}
	sections = append(sections, choices...)
	sections = append(sections, "", ui.DimStyle.Render("    ↑/↓ navigate · Enter to choose"))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- Forage Result ---

func (m model) updateForageResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" || key.String() == " " {
			m.screen = screenMeadow
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) viewForageResult() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render("✿ You found something! ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	pf := m.lastForaged
	f := pf.flower

	nameStyle := flowerNameStyle(f.Rarity)
	name := nameStyle.Render(f.Rarity.Symbol() + " " + f.Name)
	res := ui.ResonanceStyle.Render("Resonance: " + f.Resonance)
	desc := ui.ProseStyle.Render(f.Desc)
	rarity := ui.DimStyle.Render("Rarity: " + f.Rarity.String())
	potency := "  Potency: " + ui.RenderPotency(f.Potency, f.MaxPotency)

	var msgR string
	if pf.damaged {
		msgR = ui.WarnStyle.Render("  " + m.forageMessage)
	} else {
		msgR = ui.FlavorStyle.Render("  " + m.forageMessage)
	}

	flowerBox := ui.InnerBoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			"  "+name,
			"  "+res,
			"  "+rarity,
			potency,
		))

	prompt := ui.StatusBarStyle.Width(w).Render("[ Press Enter to continue ]")

	return lipgloss.JoinVertical(lipgloss.Left,
		"", header, sep, "", desc, "", flowerBox, "", msgR, "", prompt)
}

// --- Journal ---

func (m model) updateJournal(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc", "q", "enter":
			m.screen = screenMeadow
			m.cursor = 0
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right", "l":
			maxPage := (len(m.journalPages) - 1) / 2
			if m.cursor < maxPage {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m model) viewJournal() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render("✿ Your Journal ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	if len(m.journalPages) == 0 {
		empty := ui.FlavorStyle.Render("  Your journal pages are blank.\n  Go forage to find flowers to press!")
		prompt := ui.StatusBarStyle.Width(w).Render("[ Press Enter to go back ]")
		return lipgloss.JoinVertical(lipgloss.Left, "", header, sep, "", empty, "", prompt)
	}

	// Show two pages at a time
	pageStart := m.cursor * 2
	pageW := 26

	renderPage := func(idx int) string {
		if idx >= len(m.journalPages) {
			return ui.InnerBoxStyle.Width(pageW).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					"",
					ui.DimStyle.Render(fmt.Sprintf("  Page %d", idx+1)),
					"",
					ui.DimStyle.Render("  (empty page)"),
					"",
					ui.DimStyle.Render("  press a flower"),
					ui.DimStyle.Render("  to this page"),
					"",
				))
		}
		pf := m.journalPages[idx]
		f := pf.flower
		nameS := flowerNameStyle(f.Rarity)

		lines := []string{
			"",
			"  " + nameS.Render(f.Rarity.Symbol()+" "+f.Name),
			"  " + ui.ResonanceStyle.Render(f.Resonance),
			"  " + ui.DimStyle.Render(fmt.Sprintf("Pressed Day %d", max(1, m.day))),
			"  " + ui.RenderPotency(f.Potency, f.MaxPotency) + " potency",
			"",
		}
		if pf.damaged {
			lines = append(lines, "  "+ui.WarnStyle.Render("(damaged)"))
		}
		if f.Potency == 0 {
			lines = append(lines, "  "+ui.WarnStyle.Render("(faded)"))
		}
		lines = append(lines, "")
		return ui.InnerBoxStyle.Width(pageW).Render(
			lipgloss.JoinVertical(lipgloss.Left, lines...))
	}

	page1 := renderPage(pageStart)
	page2 := renderPage(pageStart + 1)

	pages := lipgloss.JoinHorizontal(lipgloss.Top, "  ", page1, "  ", page2)

	// Known blends
	var blendLines []string
	if len(m.knownBlends) > 0 {
		blendLines = append(blendLines, ui.FlavorStyle.Render("  Known blends:"))
		for _, b := range data.AllBlends {
			if m.knownBlends[b.Name] {
				blendLines = append(blendLines,
					"  "+ui.BlendStyle.Render("  ✾ "+b.Name)+
						ui.DimStyle.Render(" ("+b.Flower1+" + "+b.Flower2+")"))
			}
		}
	}

	totalPages := len(m.journalPages)
	pageStr := fmt.Sprintf("Page %d-%d of %d", pageStart+1, min(pageStart+2, totalPages), totalPages)
	nav := ui.DimStyle.Render("    ‹ prev page    " + pageStr + "    next page ›")
	prompt := ui.StatusBarStyle.Width(w).Render("←/→ flip pages · Enter/Esc to close")

	sections := []string{"", header, sep, "", pages, ""}
	sections = append(sections, blendLines...)
	sections = append(sections, "", nav, "", prompt)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- Blend Screen ---

func (m model) updateBlend(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc", "q":
			m.screen = screenMeadow
			m.cursor = 0
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			maxC := len(m.journalPages) // flowers + cancel
			if m.cursor < maxC {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(m.journalPages) {
				// Cancel
				m.screen = screenMeadow
				m.cursor = 0
				return m, nil
			}
			if m.blendPhase == 0 {
				m.blendSelect1 = m.cursor
				m.blendPhase = 1
				m.cursor = 0
			} else {
				if m.cursor == m.blendSelect1 {
					// Can't blend with self
					return m, nil
				}
				m.blendSelect2 = m.cursor
				// Check for blend
				f1 := m.journalPages[m.blendSelect1].flower.Name
				f2 := m.journalPages[m.blendSelect2].flower.Name

				found := false
				for _, b := range data.AllBlends {
					if (b.Flower1 == f1 && b.Flower2 == f2) || (b.Flower1 == f2 && b.Flower2 == f1) {
						found = true
						if m.knownBlends[b.Name] {
							m.blendResult = ui.BlendStyle.Render("✾ "+b.Name) + "\n\n" +
								ui.FlavorStyle.Render("  "+b.Desc) + "\n\n" +
								ui.ProseStyle.Render("  Effect: "+b.Effect)
						} else {
							m.knownBlends[b.Name] = true
							m.discoveredBlends = append(m.discoveredBlends, b.Name)
							m.blendsFound++
							m.blendResult = ui.SuccessStyle.Render("✾ New blend discovered!") + "\n\n" +
								ui.BlendStyle.Render("  "+b.Name) + "\n\n" +
								ui.FlavorStyle.Render("  "+b.Desc) + "\n\n" +
								ui.ProseStyle.Render("  Effect: "+b.Effect)
						}
						// Reduce potency of used flowers
						m.journalPages[m.blendSelect1].flower.Potency = max(0,
							m.journalPages[m.blendSelect1].flower.Potency-1)
						m.journalPages[m.blendSelect2].flower.Potency = max(0,
							m.journalPages[m.blendSelect2].flower.Potency-1)
						break
					}
				}
				if !found {
					m.blendResult = ui.WarnStyle.Render("  The flowers resist each other.") + "\n" +
						ui.FlavorStyle.Render("  Nothing happens. Perhaps a different combination...")
				}
				m.screen = screenBlendResult
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m model) viewBlend() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	var phaseStr string
	if m.blendPhase == 0 {
		phaseStr = "Choose the first flower:"
	} else {
		f1Name := m.journalPages[m.blendSelect1].flower.Name
		phaseStr = fmt.Sprintf("Blending with %s — choose the second:", f1Name)
	}

	header := ui.HeaderStyle.Width(w).Render("✿ Blending ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))
	phase := ui.ProseStyle.Render("  " + phaseStr)

	var choices []string
	for i, pf := range m.journalPages {
		f := pf.flower
		nameS := flowerNameStyle(f.Rarity)
		label := nameS.Render(f.Rarity.Symbol()+" "+f.Name) +
			" " + ui.DimStyle.Render("("+f.Resonance+")")
		if i == m.blendSelect1 {
			label += ui.SuccessStyle.Render(" ← selected")
		}
		if i == m.cursor {
			choices = append(choices, ui.SelectedStyle.Render("› ")+label)
		} else {
			choices = append(choices, ui.ChoiceStyle.Render("  ")+label)
		}
	}
	if m.cursor == len(m.journalPages) {
		choices = append(choices, ui.SelectedStyle.Render("› Cancel"))
	} else {
		choices = append(choices, ui.ChoiceStyle.Render("  Cancel"))
	}

	prompt := ui.DimStyle.Render("    ↑/↓ navigate · Enter to choose · Esc to cancel")

	sections := []string{"", header, sep, "", phase, ""}
	sections = append(sections, choices...)
	sections = append(sections, "", prompt)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- Blend Result ---

func (m model) updateBlendResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" || key.String() == " " {
			m.screen = screenMeadow
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) viewBlendResult() string {
	w := m.width
	if w < 20 {
		w = 60
	}
	header := ui.HeaderStyle.Width(w).Render("✿ Blend Result ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))
	prompt := ui.StatusBarStyle.Width(w).Render("[ Press Enter to continue ]")

	return lipgloss.JoinVertical(lipgloss.Left,
		"", header, sep, "", m.blendResult, "", prompt)
}

// --- NPC Intro ---

func (m model) updateNPCIntro(msg tea.Msg) (tea.Model, tea.Cmd) {
	maxC := 2 // offer a flower, say goodbye
	if len(m.journalPages) == 0 {
		maxC = 1
	}
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < maxC-1 {
				m.cursor++
			}
		case "enter":
			if len(m.journalPages) > 0 && m.cursor == 0 {
				m.screen = screenNPCOffer
				m.cursor = 0
			} else {
				m.screen = screenMeadow
				m.cursor = 0
			}
		case "esc", "q":
			m.screen = screenMeadow
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) viewNPCIntro() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render("✿ Maren's Cottage ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	prose := ui.ProseStyle.Render(
		"The cottage door is weathered oak, decorated with dried\n" +
			"  lavender bundles. Bee hives line the garden wall.\n" +
			"  Maren opens the door before you knock.")

	greeting := ui.DialogueStyle.Render(
		fmt.Sprintf("\n  %s: \"%s\"",
			ui.NPCNameStyle.Render(m.npc.Name), m.npc.Greeting))

	trust := fmt.Sprintf("  %s %s: %s  %s",
		ui.NPCNameStyle.Render(m.npc.Name),
		ui.DimStyle.Render(m.npc.Title),
		ui.RenderTrust(m.npc.Trust),
		ui.DimStyle.Render(m.npc.Disposition))

	sep2 := ui.SepStyle.Render(ui.FloralSeparator(w))

	var choices []string
	idx := 0
	if len(m.journalPages) > 0 {
		choices = append(choices, m.renderChoice(idx, "⚘ Offer her a flower from your journal"))
		idx++
	}
	choices = append(choices, m.renderChoice(idx, "Say goodbye"))

	sections := []string{"", header, sep, "", prose, "", greeting, "", trust, sep2, ""}
	sections = append(sections, choices...)
	sections = append(sections, "", ui.DimStyle.Render("    ↑/↓ navigate · Enter to choose"))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- NPC Offer ---

func (m model) updateNPCOffer(msg tea.Msg) (tea.Model, tea.Cmd) {
	maxC := len(m.journalPages) + 1 // +1 for cancel
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < maxC-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(m.journalPages) {
				m.screen = screenNPCIntro
				m.cursor = 0
			} else {
				m.npcOffering = m.cursor
				f := m.journalPages[m.cursor].flower
				// Get reaction
				if reaction, ok := m.npc.Dialogue[f.Name]; ok {
					m.npcReaction = reaction
				} else {
					m.npcReaction = "She looks at the " + f.Name + " thoughtfully. \"Thank you. It's... lovely.\""
				}
				// Trust changes based on flower
				m.npcTrustDelta = 0
				switch f.Name {
				case "Foxglove":
					m.npcTrustDelta = -2
					m.npc.Disposition = "wary · distant"
				case "Ghost Orchid", "Moonflower":
					m.npcTrustDelta = 2
					m.npc.Disposition = "awed · open"
				case "Honeysuckle", "Clover":
					m.npcTrustDelta = 2
					m.npc.Disposition = "warm · friendly"
				case "Daisy", "Lavender", "Chamomile":
					m.npcTrustDelta = 1
					m.npc.Disposition = "grateful · warm"
				case "Bloodroot":
					m.npcTrustDelta = -1
					m.npc.Disposition = "cautious · assessing"
				default:
					m.npcTrustDelta = 1
					m.npc.Disposition = "pleased · curious"
				}
				m.npc.Trust = max(0, min(5, m.npc.Trust+m.npcTrustDelta))

				// Remove the flower (it was given away)
				m.journalPages = append(m.journalPages[:m.cursor], m.journalPages[m.cursor+1:]...)

				m.screen = screenNPCReaction
				m.cursor = 0
			}
		case "esc", "q":
			m.screen = screenNPCIntro
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) viewNPCOffer() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render("✿ Offer a Flower ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))
	prose := ui.FlavorStyle.Render("  Choose a flower from your journal to offer Maren.\n  She'll react to what you give — and remember it.")

	var choices []string
	for i, pf := range m.journalPages {
		f := pf.flower
		nameS := flowerNameStyle(f.Rarity)
		label := nameS.Render(f.Rarity.Symbol()+" "+f.Name) +
			" " + ui.DimStyle.Render("("+f.Resonance+")")
		if f.Potency == 0 {
			label += " " + ui.WarnStyle.Render("[faded]")
		}
		if i == m.cursor {
			choices = append(choices, ui.SelectedStyle.Render("› ")+label)
			choices = append(choices, ui.FlavorStyle.Render("    "+f.Desc))
			choices = append(choices, "")
		} else {
			choices = append(choices, ui.ChoiceStyle.Render("  ")+label)
		}
	}
	if m.cursor == len(m.journalPages) {
		choices = append(choices, ui.SelectedStyle.Render("› Keep your flowers"))
	} else {
		choices = append(choices, ui.ChoiceStyle.Render("  Keep your flowers"))
	}

	sections := []string{"", header, sep, "", prose, ""}
	sections = append(sections, choices...)
	sections = append(sections, "", ui.DimStyle.Render("    ↑/↓ navigate · Enter to choose"))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- NPC Reaction ---

func (m model) updateNPCReaction(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" || key.String() == " " {
			m.screen = screenMeadow
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) viewNPCReaction() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.HeaderStyle.Width(w).Render("✿ Maren's Cottage ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	reaction := ui.DialogueStyle.Render("  " + m.npcReaction)

	var trustChange string
	if m.npcTrustDelta > 0 {
		trustChange = ui.SuccessStyle.Render(fmt.Sprintf("  Trust increased by %d", m.npcTrustDelta))
	} else if m.npcTrustDelta < 0 {
		trustChange = ui.WarnStyle.Render(fmt.Sprintf("  Trust decreased by %d", -m.npcTrustDelta))
	}

	trust := fmt.Sprintf("  %s %s: %s  %s",
		ui.NPCNameStyle.Render(m.npc.Name),
		ui.DimStyle.Render(m.npc.Title),
		ui.RenderTrust(m.npc.Trust),
		ui.DimStyle.Render(m.npc.Disposition))

	prompt := ui.StatusBarStyle.Width(w).Render("[ Press Enter to continue ]")

	return lipgloss.JoinVertical(lipgloss.Left,
		"", header, sep, "", reaction, "", trustChange, "", trust, "", prompt)
}

// --- Summary ---

func (m model) updateSummary(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "enter" || key.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) viewSummary() string {
	w := m.width
	if w < 20 {
		w = 60
	}

	header := ui.TitleStyle.Width(w).Render("✿ Journey's End ✿")
	sep := ui.SepStyle.Render(ui.FloralSeparator(w))

	lines := []string{
		ui.ProseStyle.Render(fmt.Sprintf("  You close your journal after %d days in the meadow.", m.day)),
		"",
		ui.FlavorStyle.Render("  The pages are heavier now — pressed flowers,"),
		ui.FlavorStyle.Render("  sketched maps, and the memory of honey and smoke."),
		"",
		ui.SepStyle.Render("  " + ui.FloralSeparator(40)),
		"",
		ui.ProseStyle.Render(fmt.Sprintf("  Tradition: %s %s", m.tradition.Symbol, m.tradition.Name)),
		ui.ProseStyle.Render(fmt.Sprintf("  Days spent: %d", m.day)),
		ui.ProseStyle.Render(fmt.Sprintf("  Flowers found: %d", m.flowersFound)),
		ui.ProseStyle.Render(fmt.Sprintf("  Flowers in journal: %d", len(m.journalPages))),
		ui.ProseStyle.Render(fmt.Sprintf("  Blends discovered: %d", m.blendsFound)),
		ui.ProseStyle.Render(fmt.Sprintf("  Maren's trust: %s", ui.RenderTrust(m.npc.Trust))),
		"",
	}

	if len(m.discoveredBlends) > 0 {
		lines = append(lines, ui.FlavorStyle.Render("  Blends discovered:"))
		for _, b := range m.discoveredBlends {
			lines = append(lines, ui.BlendStyle.Render("    ✾ "+b))
		}
		lines = append(lines, "")
	}

	lines = append(lines,
		ui.FlavorStyle.Render("  In the place where you rested, a small garden grows."),
		ui.FlavorStyle.Render("  The next traveler will find it waiting."),
		"",
		ui.StatusBarStyle.Width(w).Render("[ Press Enter to close ]"),
	)

	sections := []string{"", header, sep}
	sections = append(sections, lines...)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// --- View Router ---

func (m model) View() string {
	switch m.screen {
	case screenTitle:
		return m.viewTitle()
	case screenCharCreate:
		return m.viewCharCreate()
	case screenMeadow:
		return m.viewMeadow()
	case screenForage:
		return m.viewForage()
	case screenForageResult:
		return m.viewForageResult()
	case screenJournal:
		return m.viewJournal()
	case screenBlend:
		return m.viewBlend()
	case screenBlendResult:
		return m.viewBlendResult()
	case screenNPCIntro:
		return m.viewNPCIntro()
	case screenNPCOffer:
		return m.viewNPCOffer()
	case screenNPCReaction:
		return m.viewNPCReaction()
	case screenSummary:
		return m.viewSummary()
	}
	return ""
}

// --- Helpers ---

func (m model) renderChoice(idx int, label string) string {
	if idx == m.cursor {
		return ui.SelectedStyle.Render("› " + label)
	}
	return ui.ChoiceStyle.Render("  " + label)
}

func flowerNameStyle(r data.Rarity) lipgloss.Style {
	switch r {
	case data.Common:
		return ui.CommonFlowerStyle
	case data.Uncommon:
		return ui.UncommonFlowerStyle
	case data.Rare:
		return ui.RareFlowerStyle
	case data.Mythic:
		return ui.MythicFlowerStyle
	default:
		return ui.ProseStyle
	}
}

// --- Main ---

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
