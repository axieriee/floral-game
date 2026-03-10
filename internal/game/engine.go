package game

import (
	"math"
	"math/rand"
	"time"
)

// GrowthMultiplier returns the speed multiplier from upgrades.
func (g *GameState) GrowthMultiplier() float64 {
	level := g.UpgradeLevels["grow_speed"]
	return 1.0 + float64(level)*0.10
}

// PetalMultiplier returns the petal yield multiplier.
func (g *GameState) PetalMultiplier() float64 {
	mult := 1.0
	mult += float64(g.UpgradeLevels["petal_mult"]) * 0.5
	// Prestige bonus: +10% per prestige
	mult += float64(g.PrestigeCount) * 0.10
	return mult
}

// FlatBonus returns flat bonus petals per harvest.
func (g *GameState) FlatBonus() float64 {
	return float64(g.UpgradeLevels["flat_bonus"])
}

// DoubleChance returns the probability of a double harvest.
func (g *GameState) DoubleChance() float64 {
	level := g.UpgradeLevels["double_chance"]
	return math.Min(float64(level)*0.15, 0.90)
}

// HasAutoHarvest returns whether auto-harvest is unlocked.
func (g *GameState) HasAutoHarvest() bool {
	return g.UpgradeLevels["auto_harvest"] >= 1
}

// SeedGenRate returns seeds earned per harvest.
func (g *GameState) SeedGenRate() float64 {
	return float64(g.UpgradeLevels["seed_gen"]) * 0.5
}

// EffectiveGrowTime returns adjusted grow time for a flower type.
func (g *GameState) EffectiveGrowTime(ft FlowerType) time.Duration {
	ms := float64(ft.GrowTime.Milliseconds()) / g.GrowthMultiplier()
	return time.Duration(ms) * time.Millisecond
}

// PlotProgress returns 0.0-1.0 progress for a plot.
func (g *GameState) PlotProgress(plotIdx int) float64 {
	if plotIdx >= len(g.Plots) {
		return 0
	}
	plot := g.Plots[plotIdx]
	ft := FlowerTypes[plot.FlowerType]
	growTime := g.EffectiveGrowTime(ft)
	elapsed := time.Since(plot.Planted)
	progress := float64(elapsed) / float64(growTime)
	if progress > 1.0 {
		progress = 1.0
	}
	return progress
}

// IsReady returns whether a plot is ready to harvest.
func (g *GameState) IsReady(plotIdx int) bool {
	return g.PlotProgress(plotIdx) >= 1.0
}

// Harvest collects petals from a ready plot and replants.
func (g *GameState) Harvest(plotIdx int) (petals float64, seeds float64, doubled bool) {
	if plotIdx >= len(g.Plots) || !g.IsReady(plotIdx) {
		return 0, 0, false
	}
	plot := &g.Plots[plotIdx]
	ft := FlowerTypes[plot.FlowerType]

	petals = (ft.PetalYield + g.FlatBonus()) * g.PetalMultiplier()
	seeds = g.SeedGenRate()

	if rand.Float64() < g.DoubleChance() {
		petals *= 2
		seeds *= 2
		doubled = true
	}

	g.Petals += petals
	g.TotalPetals += petals
	g.Seeds += seeds
	g.TotalHarvests++

	// Replant
	plot.Planted = time.Now()

	return petals, seeds, doubled
}

// UpgradeCost returns the cost of the next level of an upgrade.
func UpgradeCost(u Upgrade, currentLevel int) float64 {
	return u.BaseCost * math.Pow(u.CostScale, float64(currentLevel))
}

// BuyUpgrade attempts to purchase an upgrade.
func (g *GameState) BuyUpgrade(upgradeIdx int) bool {
	if upgradeIdx >= len(Upgrades) {
		return false
	}
	u := Upgrades[upgradeIdx]
	currentLevel := g.UpgradeLevels[u.Effect]
	if u.MaxLevel > 0 && currentLevel >= u.MaxLevel {
		return false
	}
	cost := UpgradeCost(u, currentLevel)
	if g.Petals < cost {
		return false
	}
	g.Petals -= cost
	g.UpgradeLevels[u.Effect] = currentLevel + 1

	// If it's the "new_plot" upgrade, add a plot
	if u.Effect == "new_plot" {
		g.Plots = append(g.Plots, PlotState{
			FlowerType: 0,
			Planted:    time.Now(),
			Quantity:   1,
		})
	}

	return true
}

// UnlockFlower attempts to unlock a flower type.
func (g *GameState) UnlockFlower(flowerIdx int) bool {
	if flowerIdx >= len(FlowerTypes) || g.Unlocked[flowerIdx] {
		return false
	}
	ft := FlowerTypes[flowerIdx]
	if g.Petals < ft.UnlockCost {
		return false
	}
	g.Petals -= ft.UnlockCost
	g.Unlocked[flowerIdx] = true
	return true
}

// PlantFlower changes what's growing in a plot.
func (g *GameState) PlantFlower(plotIdx, flowerIdx int) bool {
	if plotIdx >= len(g.Plots) || flowerIdx >= len(FlowerTypes) {
		return false
	}
	if !g.Unlocked[flowerIdx] {
		return false
	}
	ft := FlowerTypes[flowerIdx]
	if g.Seeds < ft.SeedCost && ft.SeedCost > 0 {
		return false
	}
	if ft.SeedCost > 0 {
		g.Seeds -= ft.SeedCost
	}
	g.Plots[plotIdx] = PlotState{
		FlowerType: flowerIdx,
		Planted:    time.Now(),
		Quantity:   1,
	}
	return true
}

// Tick processes auto-harvest and offline progress.
func (g *GameState) Tick() {
	if g.HasAutoHarvest() {
		for i := range g.Plots {
			if g.IsReady(i) {
				g.Harvest(i)
			}
		}
	}
}

// CalculateOfflineProgress processes time that passed while the game was closed.
func (g *GameState) CalculateOfflineProgress() (petals float64, harvests int) {
	elapsed := time.Since(g.LastTick)
	if elapsed < time.Second {
		return 0, 0
	}

	if !g.HasAutoHarvest() {
		g.LastTick = time.Now()
		return 0, 0
	}

	// Simulate offline harvests
	for i := range g.Plots {
		ft := FlowerTypes[g.Plots[i].FlowerType]
		growTime := g.EffectiveGrowTime(ft)
		if growTime <= 0 {
			continue
		}
		numHarvests := int(elapsed / growTime)
		for j := 0; j < numHarvests; j++ {
			p, _, _ := g.Harvest(i)
			petals += p
			harvests++
		}
	}

	g.LastTick = time.Now()
	return petals, harvests
}

// NectarFromPrestige calculates nectar earned from a prestige reset.
func (g *GameState) NectarFromPrestige() float64 {
	if g.TotalPetals < 10000 {
		return 0
	}
	return math.Floor(math.Sqrt(g.TotalPetals / 1000))
}

// Prestige resets progress and grants nectar.
func (g *GameState) Prestige() float64 {
	nectar := g.NectarFromPrestige()
	if nectar <= 0 {
		return 0
	}

	g.Nectar += nectar
	g.PrestigeCount++

	// Reset
	g.Petals = 0
	g.TotalPetals = 0
	g.Seeds = 0
	g.UpgradeLevels = make(map[string]int)
	g.Unlocked = make([]bool, len(FlowerTypes))
	g.Unlocked[0] = true

	now := time.Now()
	g.Plots = []PlotState{
		{FlowerType: 0, Planted: now, Quantity: 1},
		{FlowerType: 0, Planted: now, Quantity: 1},
	}

	return nectar
}

// PetalsPerSecond estimates current petal generation rate.
func (g *GameState) PetalsPerSecond() float64 {
	if !g.HasAutoHarvest() {
		return 0
	}
	total := 0.0
	for _, plot := range g.Plots {
		ft := FlowerTypes[plot.FlowerType]
		growTime := g.EffectiveGrowTime(ft)
		if growTime <= 0 {
			continue
		}
		harvestsPerSec := 1.0 / growTime.Seconds()
		yield := (ft.PetalYield + g.FlatBonus()) * g.PetalMultiplier()
		total += harvestsPerSec * yield
	}
	return total
}
