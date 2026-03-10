package game

import "time"

// Season represents a garden season.
type Season int

const (
	Spring Season = iota
	Summer
	Autumn
	Winter
	NumSeasons
)

// SeasonDuration is how long each season lasts in real time.
const SeasonDuration = 2 * time.Minute

// SeasonInfo holds display and gameplay data for a season.
type SeasonInfo struct {
	Name       string
	Emoji      string
	GrowthMult map[int]float64 // flower tier -> growth multiplier
	YieldMult  map[int]float64 // flower tier -> yield multiplier
	Special    string          // flavor text
}

var Seasons = [NumSeasons]SeasonInfo{
	{
		Name:  "Spring",
		Emoji: "🌱",
		GrowthMult: map[int]float64{
			0: 1.3, 1: 1.5, 2: 1.0, 3: 0.8,
		},
		YieldMult: map[int]float64{
			0: 1.2, 1: 1.3, 2: 1.0, 3: 1.0,
		},
		Special: "New growth blooms swiftly",
	},
	{
		Name:  "Summer",
		Emoji: "☀️",
		GrowthMult: map[int]float64{
			0: 1.0, 1: 1.0, 2: 1.5, 3: 1.3,
		},
		YieldMult: map[int]float64{
			0: 1.0, 1: 1.0, 2: 1.5, 3: 1.2,
		},
		Special: "The sun empowers exotic blooms",
	},
	{
		Name:  "Autumn",
		Emoji: "🍂",
		GrowthMult: map[int]float64{
			0: 0.8, 1: 1.0, 2: 1.2, 3: 1.5,
		},
		YieldMult: map[int]float64{
			0: 1.5, 1: 1.5, 2: 1.0, 3: 1.0,
		},
		Special: "Harvest season — common flowers yield more",
	},
	{
		Name:  "Winter",
		Emoji: "❄️",
		GrowthMult: map[int]float64{
			0: 0.5, 1: 0.5, 2: 0.7, 3: 1.0,
		},
		YieldMult: map[int]float64{
			0: 0.7, 1: 0.7, 2: 0.8, 3: 2.0,
		},
		Special: "Only the rarest flowers thrive in the cold",
	},
}

// CurrentSeason calculates the current season based on game time.
func (g *GameState) CurrentSeason() Season {
	elapsed := time.Since(g.CreatedAt)
	cyclePos := elapsed % (SeasonDuration * time.Duration(NumSeasons))
	return Season(cyclePos / SeasonDuration)
}

// SeasonProgress returns 0.0-1.0 how far through the current season we are.
func (g *GameState) SeasonProgress() float64 {
	elapsed := time.Since(g.CreatedAt)
	cyclePos := elapsed % (SeasonDuration * time.Duration(NumSeasons))
	inSeason := cyclePos % SeasonDuration
	return float64(inSeason) / float64(SeasonDuration)
}

// SeasonGrowthMult returns the seasonal growth multiplier for a flower tier.
func (g *GameState) SeasonGrowthMult(tier int) float64 {
	s := Seasons[g.CurrentSeason()]
	if mult, ok := s.GrowthMult[tier]; ok {
		return mult
	}
	return 1.0
}

// SeasonYieldMult returns the seasonal yield multiplier for a flower tier.
func (g *GameState) SeasonYieldMult(tier int) float64 {
	s := Seasons[g.CurrentSeason()]
	if mult, ok := s.YieldMult[tier]; ok {
		return mult
	}
	return 1.0
}

// HasGreenhouse returns whether the player has a greenhouse (immune to seasons).
func (g *GameState) HasGreenhouse() bool {
	return g.EssenceUpgrades["greenhouse"] > 0
}
