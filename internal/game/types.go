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
	FlowerType int       // index into FlowerTypes
	Planted    time.Time // when it was planted
	Quantity   int       // how many of this flower are growing
	AutoHarvest bool
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

// GameState holds the entire persistent game state.
type GameState struct {
	Petals       float64     `json:"petals"`
	TotalPetals  float64     `json:"total_petals"`  // lifetime petals earned
	Seeds        float64     `json:"seeds"`
	Nectar       float64     `json:"nectar"`         // prestige currency
	Plots        []PlotState `json:"plots"`
	Unlocked     []bool      `json:"unlocked"`       // which flower types are unlocked
	UpgradeLevels map[string]int `json:"upgrade_levels"`
	PrestigeCount int        `json:"prestige_count"`
	LastTick     time.Time   `json:"last_tick"`
	TotalHarvests int64      `json:"total_harvests"`
	PlayTime     time.Duration `json:"play_time"`
	CreatedAt    time.Time   `json:"created_at"`
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
}

func NewGameState() *GameState {
	now := time.Now()
	unlocked := make([]bool, len(FlowerTypes))
	unlocked[0] = true // Daisies are free

	plots := []PlotState{
		{FlowerType: 0, Planted: now, Quantity: 1},
		{FlowerType: 0, Planted: now, Quantity: 1},
	}

	return &GameState{
		Petals:        0,
		Seeds:         0,
		Nectar:        0,
		Plots:         plots,
		Unlocked:      unlocked,
		UpgradeLevels: make(map[string]int),
		LastTick:      now,
		CreatedAt:     now,
	}
}
