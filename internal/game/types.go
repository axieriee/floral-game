package game

import "time"

// FlowerType represents a type of flower that can be grown.
type FlowerType struct {
	Name       string
	Emoji      string
	GrowTime   time.Duration // base time to grow one flower
	PetalYield float64       // petals per harvest
	SeedCost   float64       // cost to plant
	UnlockCost float64       // petals needed to unlock this flower
	Tier       int
}

// PlotState represents a single garden plot.
type PlotState struct {
	FlowerType  int       `json:"flower_type"`  // index into all flowers (base + hybrid)
	Planted     time.Time `json:"planted"`
	Quantity    int       `json:"quantity"`
	AutoHarvest bool      `json:"auto_harvest"`
	IsGreenhouse bool     `json:"is_greenhouse"` // immune to season effects
}

// Upgrade represents a purchasable upgrade.
type Upgrade struct {
	Name        string
	Description string
	BaseCost    float64
	CostScale   float64 // multiplicative cost increase per level
	MaxLevel    int     // 0 = unlimited
	Effect      string  // key for what it does
}

// LogEntry represents an event in the game log.
type LogEntry struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Color   string    `json:"color"` // hex color for display
}

// GameState holds the entire persistent game state.
type GameState struct {
	// Resources
	Petals      float64 `json:"petals"`
	TotalPetals float64 `json:"total_petals"`
	Seeds       float64 `json:"seeds"`
	Nectar      float64 `json:"nectar"`  // prestige 1 currency
	Essence     float64 `json:"essence"` // prestige 2 currency

	// Garden
	Plots   []PlotState `json:"plots"`
	Unlocked []bool     `json:"unlocked"` // base flowers + hybrids

	// Upgrades
	UpgradeLevels  map[string]int `json:"upgrade_levels"`
	EssenceUpgrades map[string]int `json:"essence_upgrades"` // persistent across prestige 1

	// Progression
	PrestigeCount    int   `json:"prestige_count"`
	Prestige2Count   int   `json:"prestige2_count"`
	TotalHarvests    int64 `json:"total_harvests"`
	DiscoveredHybrids []int `json:"discovered_hybrids"`

	// Timing
	LastTick  time.Time     `json:"last_tick"`
	PlayTime  time.Duration `json:"play_time"`
	CreatedAt time.Time     `json:"created_at"`
	LastSave  time.Time     `json:"last_save"`

	// Event log (not persisted — rebuilt on load)
	Log []LogEntry `json:"-"`
}

var FlowerTypes = []FlowerType{
	{Name: "Daisy", Emoji: "🌼", GrowTime: 3 * time.Second, PetalYield: 1, SeedCost: 0, UnlockCost: 0, Tier: 0},
	{Name: "Tulip", Emoji: "🌷", GrowTime: 8 * time.Second, PetalYield: 3, SeedCost: 5, UnlockCost: 25, Tier: 1},
	{Name: "Rose", Emoji: "🌹", GrowTime: 20 * time.Second, PetalYield: 10, SeedCost: 20, UnlockCost: 100, Tier: 1},
	{Name: "Sunflower", Emoji: "🌻", GrowTime: 45 * time.Second, PetalYield: 30, SeedCost: 50, UnlockCost: 500, Tier: 2},
	{Name: "Orchid", Emoji: "🪻", GrowTime: 90 * time.Second, PetalYield: 80, SeedCost: 150, UnlockCost: 2000, Tier: 2},
	{Name: "Lotus", Emoji: "🪷", GrowTime: 180 * time.Second, PetalYield: 250, SeedCost: 500, UnlockCost: 10000, Tier: 3},
	{Name: "Crystal Bloom", Emoji: "💎", GrowTime: 300 * time.Second, PetalYield: 800, SeedCost: 2000, UnlockCost: 50000, Tier: 3},
}

var Upgrades = []Upgrade{
	{Name: "Fertile Soil", Description: "Flowers grow 10% faster", BaseCost: 15, CostScale: 1.8, MaxLevel: 20, Effect: "grow_speed"},
	{Name: "Golden Trowel", Description: "+1 petal per harvest", BaseCost: 10, CostScale: 1.5, MaxLevel: 50, Effect: "flat_bonus"},
	{Name: "Bee Colony", Description: "Auto-harvest ready flowers", BaseCost: 100, CostScale: 2.5, MaxLevel: 1, Effect: "auto_harvest"},
	{Name: "New Plot", Description: "Add a garden plot", BaseCost: 30, CostScale: 2.0, MaxLevel: 8, Effect: "new_plot"},
	{Name: "Seed Pouch", Description: "Earn seeds from harvests", BaseCost: 50, CostScale: 2.0, MaxLevel: 10, Effect: "seed_gen"},
	{Name: "Petal Multiplier", Description: "x1.5 petal yield", BaseCost: 200, CostScale: 3.0, MaxLevel: 5, Effect: "petal_mult"},
	{Name: "Compost Bin", Description: "15% chance of double harvest", BaseCost: 75, CostScale: 2.2, MaxLevel: 10, Effect: "double_chance"},
	{Name: "Pollination", Description: "+25% hybrid discovery chance", BaseCost: 150, CostScale: 2.0, MaxLevel: 5, Effect: "pollination"},
}

// EssenceUpgradeInfo defines upgrades bought with Essence (prestige 2 currency).
type EssenceUpgradeInfo struct {
	Name        string
	Description string
	Cost        float64
	MaxLevel    int
	Effect      string
}

var EssenceUpgrades = []EssenceUpgradeInfo{
	{Name: "Greenhouse", Description: "One plot ignores season penalties", Cost: 5, MaxLevel: 5, Effect: "greenhouse"},
	{Name: "Ancient Seeds", Description: "Start each run with 50 seeds", Cost: 3, MaxLevel: 10, Effect: "start_seeds"},
	{Name: "Eternal Bloom", Description: "+25% base growth speed per level", Cost: 8, MaxLevel: 5, Effect: "eternal_growth"},
	{Name: "Soul of the Garden", Description: "+50% petal yield per level", Cost: 10, MaxLevel: 5, Effect: "soul_yield"},
	{Name: "Memory of Flowers", Description: "Keep 1 flower unlock per level on prestige", Cost: 15, MaxLevel: 3, Effect: "memory"},
}

func NewGameState() *GameState {
	now := time.Now()
	totalFlowers := len(FlowerTypes) + len(HybridFlowers)
	unlocked := make([]bool, totalFlowers)
	unlocked[0] = true // Daisies are free

	plots := []PlotState{
		{FlowerType: 0, Planted: now, Quantity: 1},
		{FlowerType: 0, Planted: now, Quantity: 1},
	}

	return &GameState{
		Petals:          0,
		Seeds:           0,
		Nectar:          0,
		Essence:         0,
		Plots:           plots,
		Unlocked:        unlocked,
		UpgradeLevels:   make(map[string]int),
		EssenceUpgrades: make(map[string]int),
		LastTick:        now,
		CreatedAt:       now,
		Log:             []LogEntry{},
	}
}

// AddLog adds an event to the game log, keeping max 50 entries.
func (g *GameState) AddLog(msg, color string) {
	g.Log = append(g.Log, LogEntry{
		Time:    time.Now(),
		Message: msg,
		Color:   color,
	})
	if len(g.Log) > 50 {
		g.Log = g.Log[len(g.Log)-50:]
	}
}
