package game

import (
	"math/rand"
	"time"
)

// HybridRecipe defines how to breed a hybrid flower.
type HybridRecipe struct {
	Parent1    int    // index into FlowerTypes
	Parent2    int    // index into FlowerTypes
	ResultName string // name of hybrid
	ResultIdx  int    // index into HybridFlowers
	Chance     float64
}

// HybridFlower is a discovered hybrid with unique properties.
type HybridFlower struct {
	Name       string
	Emoji      string
	GrowTime   time.Duration
	PetalYield float64
	SeedCost   float64
	Tier       int
	Lore       string
}

var HybridFlowers = []HybridFlower{
	{
		Name: "Moonpetal", Emoji: "🌙", GrowTime: 12 * time.Second,
		PetalYield: 8, SeedCost: 10, Tier: 1,
		Lore: "Born from daisy innocence and tulip grace",
	},
	{
		Name: "Blazebloom", Emoji: "🔥", GrowTime: 30 * time.Second,
		PetalYield: 25, SeedCost: 40, Tier: 2,
		Lore: "Tulip passion fused with rose thorns",
	},
	{
		Name: "Starweaver", Emoji: "⭐", GrowTime: 60 * time.Second,
		PetalYield: 60, SeedCost: 100, Tier: 2,
		Lore: "Rose elegance meets sunflower radiance",
	},
	{
		Name: "Voidlily", Emoji: "🌑", GrowTime: 120 * time.Second,
		PetalYield: 200, SeedCost: 300, Tier: 3,
		Lore: "Orchid mystery entwined with lotus serenity",
	},
	{
		Name: "Eternia", Emoji: "💫", GrowTime: 240 * time.Second,
		PetalYield: 600, SeedCost: 1500, Tier: 3,
		Lore: "The crystalline light bends through sacred waters",
	},
	{
		Name: "Prismatic Rose", Emoji: "🌈", GrowTime: 15 * time.Second,
		PetalYield: 12, SeedCost: 15, Tier: 1,
		Lore: "A daisy's simplicity painted with rose hues",
	},
	{
		Name: "Thunderpetal", Emoji: "⚡", GrowTime: 50 * time.Second,
		PetalYield: 45, SeedCost: 80, Tier: 2,
		Lore: "Sunflower power channeled through orchid precision",
	},
}

var HybridRecipes = []HybridRecipe{
	{Parent1: 0, Parent2: 1, ResultIdx: 0, Chance: 0.20}, // Daisy + Tulip = Moonpetal
	{Parent1: 1, Parent2: 2, ResultIdx: 1, Chance: 0.15}, // Tulip + Rose = Blazebloom
	{Parent1: 2, Parent2: 3, ResultIdx: 2, Chance: 0.12}, // Rose + Sunflower = Starweaver
	{Parent1: 4, Parent2: 5, ResultIdx: 3, Chance: 0.08}, // Orchid + Lotus = Voidlily
	{Parent1: 5, Parent2: 6, ResultIdx: 4, Chance: 0.05}, // Lotus + Crystal Bloom = Eternia
	{Parent1: 0, Parent2: 2, ResultIdx: 5, Chance: 0.18}, // Daisy + Rose = Prismatic Rose
	{Parent1: 3, Parent2: 4, ResultIdx: 6, Chance: 0.10}, // Sunflower + Orchid = Thunderpetal
}

// TotalFlowerCount returns the total number of plantable flowers (base + hybrids).
func TotalFlowerCount() int {
	return len(FlowerTypes) + len(HybridFlowers)
}

// GetFlowerInfo returns name, emoji, grow time, yield, seed cost, and tier for any flower index.
// Indices 0..len(FlowerTypes)-1 are base flowers; beyond that are hybrids.
func GetFlowerInfo(idx int) (name, emoji string, growTime time.Duration, yield, seedCost float64, tier int) {
	if idx < len(FlowerTypes) {
		ft := FlowerTypes[idx]
		return ft.Name, ft.Emoji, ft.GrowTime, ft.PetalYield, ft.SeedCost, ft.Tier
	}
	hi := idx - len(FlowerTypes)
	if hi < len(HybridFlowers) {
		h := HybridFlowers[hi]
		return h.Name, h.Emoji, h.GrowTime, h.PetalYield, h.SeedCost, h.Tier
	}
	return "Unknown", "?", time.Second, 0, 0, 0
}

// CheckHybridBreeding checks if adjacent plots can produce a hybrid.
// Called on harvest. Returns hybrid index or -1.
func (g *GameState) CheckHybridBreeding(plotIdx int) int {
	if plotIdx >= len(g.Plots) {
		return -1
	}
	plot := g.Plots[plotIdx]

	// Check neighbors (plot before and after)
	neighbors := []int{}
	if plotIdx > 0 {
		neighbors = append(neighbors, plotIdx-1)
	}
	if plotIdx < len(g.Plots)-1 {
		neighbors = append(neighbors, plotIdx+1)
	}

	for _, ni := range neighbors {
		neighbor := g.Plots[ni]
		for _, recipe := range HybridRecipes {
			p1, p2 := recipe.Parent1, recipe.Parent2
			if (plot.FlowerType == p1 && neighbor.FlowerType == p2) ||
				(plot.FlowerType == p2 && neighbor.FlowerType == p1) {

				// Already discovered?
				hybridGlobalIdx := len(FlowerTypes) + recipe.ResultIdx
				if g.Unlocked[hybridGlobalIdx] {
					continue // already known, skip
				}

				// Pollination bonus increases chance
				chance := recipe.Chance
				if g.UpgradeLevels["pollination"] > 0 {
					chance *= 1.0 + float64(g.UpgradeLevels["pollination"])*0.25
				}

				if rand.Float64() < chance {
					return recipe.ResultIdx
				}
			}
		}
	}
	return -1
}

// DiscoverHybrid unlocks a hybrid flower.
func (g *GameState) DiscoverHybrid(hybridIdx int) {
	globalIdx := len(FlowerTypes) + hybridIdx
	if globalIdx < len(g.Unlocked) {
		g.Unlocked[globalIdx] = true
	}
	g.DiscoveredHybrids = append(g.DiscoveredHybrids, hybridIdx)
}
