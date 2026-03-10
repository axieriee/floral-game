package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// GrowthMultiplier returns the speed multiplier from upgrades + essence.
func (g *GameState) GrowthMultiplier() float64 {
	mult := 1.0 + float64(g.UpgradeLevels["grow_speed"])*0.10
	mult += float64(g.EssenceUpgrades["eternal_growth"]) * 0.25
	return mult
}

// PetalMultiplier returns the petal yield multiplier.
func (g *GameState) PetalMultiplier() float64 {
	mult := 1.0
	mult += float64(g.UpgradeLevels["petal_mult"]) * 0.5
	mult += float64(g.PrestigeCount) * 0.10
	mult += float64(g.EssenceUpgrades["soul_yield"]) * 0.50
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

// EffectiveGrowTime returns adjusted grow time for a flower at a given plot.
func (g *GameState) EffectiveGrowTime(flowerIdx int) time.Duration {
	_, _, baseGrow, _, _, tier := GetFlowerInfo(flowerIdx)

	ms := float64(baseGrow.Milliseconds()) / g.GrowthMultiplier()

	// Season modifier (unless greenhouse)
	seasonMult := g.SeasonGrowthMult(tier)
	if seasonMult > 0 {
		ms /= seasonMult
	}

	return time.Duration(ms) * time.Millisecond
}

// EffectiveGrowTimeForPlot returns grow time accounting for greenhouse status.
func (g *GameState) EffectiveGrowTimeForPlot(plotIdx int) time.Duration {
	if plotIdx >= len(g.Plots) {
		return time.Second
	}
	plot := g.Plots[plotIdx]
	_, _, baseGrow, _, _, tier := GetFlowerInfo(plot.FlowerType)

	ms := float64(baseGrow.Milliseconds()) / g.GrowthMultiplier()

	// Greenhouse plots ignore season
	if !plot.IsGreenhouse {
		seasonMult := g.SeasonGrowthMult(tier)
		if seasonMult > 0 {
			ms /= seasonMult
		}
	}

	return time.Duration(ms) * time.Millisecond
}

// PlotProgress returns 0.0-1.0 progress for a plot.
func (g *GameState) PlotProgress(plotIdx int) float64 {
	if plotIdx >= len(g.Plots) {
		return 0
	}
	plot := g.Plots[plotIdx]
	growTime := g.EffectiveGrowTimeForPlot(plotIdx)
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

// EffectiveYield returns the petal yield for a flower at a given plot.
func (g *GameState) EffectiveYield(plotIdx int) float64 {
	if plotIdx >= len(g.Plots) {
		return 0
	}
	plot := g.Plots[plotIdx]
	_, _, _, baseYield, _, tier := GetFlowerInfo(plot.FlowerType)

	yield := (baseYield + g.FlatBonus()) * g.PetalMultiplier()

	// Season yield modifier (unless greenhouse)
	if !plot.IsGreenhouse {
		yield *= g.SeasonYieldMult(tier)
	}

	return yield
}

// Harvest collects petals from a ready plot and replants.
// Returns petals, seeds, doubled, hybridDiscovered (-1 if none).
func (g *GameState) Harvest(plotIdx int) (petals float64, seeds float64, doubled bool, hybridIdx int) {
	hybridIdx = -1
	if plotIdx >= len(g.Plots) || !g.IsReady(plotIdx) {
		return 0, 0, false, -1
	}

	petals = g.EffectiveYield(plotIdx)
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

	// Check for hybrid breeding
	hybridIdx = g.CheckHybridBreeding(plotIdx)
	if hybridIdx >= 0 {
		g.DiscoverHybrid(hybridIdx)
		h := HybridFlowers[hybridIdx]
		g.AddLog(fmt.Sprintf("DISCOVERY! %s %s bred from your garden!", h.Emoji, h.Name), "#FFD700")
	}

	// Replant
	g.Plots[plotIdx].Planted = time.Now()

	return petals, seeds, doubled, hybridIdx
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

	if u.Effect == "new_plot" {
		g.Plots = append(g.Plots, PlotState{
			FlowerType: 0,
			Planted:    time.Now(),
			Quantity:   1,
		})
	}

	return true
}

// BuyEssenceUpgrade attempts to purchase an essence upgrade.
func (g *GameState) BuyEssenceUpgrade(idx int) bool {
	if idx >= len(EssenceUpgrades) {
		return false
	}
	eu := EssenceUpgrades[idx]
	currentLevel := g.EssenceUpgrades[eu.Effect]
	if eu.MaxLevel > 0 && currentLevel >= eu.MaxLevel {
		return false
	}
	if g.Essence < eu.Cost {
		return false
	}
	g.Essence -= eu.Cost
	g.EssenceUpgrades[eu.Effect] = currentLevel + 1

	// Greenhouse: mark the next available non-greenhouse plot
	if eu.Effect == "greenhouse" {
		for i := range g.Plots {
			if !g.Plots[i].IsGreenhouse {
				g.Plots[i].IsGreenhouse = true
				break
			}
		}
	}

	return true
}

// UnlockFlower attempts to unlock a flower type (base or hybrid).
func (g *GameState) UnlockFlower(flowerIdx int) bool {
	if flowerIdx >= len(g.Unlocked) || g.Unlocked[flowerIdx] {
		return false
	}
	// Only base flowers can be manually unlocked
	if flowerIdx >= len(FlowerTypes) {
		return false // hybrids are discovered through breeding
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
	if plotIdx >= len(g.Plots) || flowerIdx >= len(g.Unlocked) {
		return false
	}
	if !g.Unlocked[flowerIdx] {
		return false
	}
	_, _, _, _, seedCost, _ := GetFlowerInfo(flowerIdx)
	if g.Seeds < seedCost && seedCost > 0 {
		return false
	}
	if seedCost > 0 {
		g.Seeds -= seedCost
	}
	greenhouse := g.Plots[plotIdx].IsGreenhouse
	g.Plots[plotIdx] = PlotState{
		FlowerType:   flowerIdx,
		Planted:      time.Now(),
		Quantity:     1,
		IsGreenhouse: greenhouse,
	}
	return true
}

// Tick processes auto-harvest.
func (g *GameState) Tick() {
	if g.HasAutoHarvest() {
		for i := range g.Plots {
			if g.IsReady(i) {
				petals, _, doubled, hybridIdx := g.Harvest(i)
				if petals > 0 {
					_, emoji, _, _, _, _ := GetFlowerInfo(g.Plots[i].FlowerType)
					msg := fmt.Sprintf("%s +%.0f petals", emoji, petals)
					if doubled {
						msg += " (DOUBLE!)"
					}
					_ = hybridIdx // hybrid log is added in Harvest
					g.AddLog(msg, "#3CB371")
				}
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

	for i := range g.Plots {
		growTime := g.EffectiveGrowTimeForPlot(i)
		if growTime <= 0 {
			continue
		}
		numHarvests := int(elapsed / growTime)
		if numHarvests > 1000 {
			numHarvests = 1000 // cap offline harvests
		}
		for j := 0; j < numHarvests; j++ {
			p, _, _, _ := g.Harvest(i)
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

// Prestige resets progress and grants nectar (layer 1).
func (g *GameState) Prestige() float64 {
	nectar := g.NectarFromPrestige()
	if nectar <= 0 {
		return 0
	}

	g.Nectar += nectar
	g.PrestigeCount++
	g.AddLog(fmt.Sprintf("✨ Prestige %d! Earned %.0f nectar", g.PrestigeCount, nectar), "#FFD700")

	// Memory of Flowers: keep some unlocks
	memoryLevel := g.EssenceUpgrades["memory"]
	keptUnlocks := make([]bool, len(g.Unlocked))
	keptUnlocks[0] = true // Daisies always unlocked
	if memoryLevel > 0 {
		kept := 0
		for i := 1; i < len(g.Unlocked) && kept < memoryLevel; i++ {
			if g.Unlocked[i] {
				keptUnlocks[i] = true
				kept++
			}
		}
	}

	// Reset
	g.Petals = 0
	g.TotalPetals = 0
	startSeeds := float64(g.EssenceUpgrades["start_seeds"]) * 50
	g.Seeds = startSeeds
	g.UpgradeLevels = make(map[string]int)
	g.Unlocked = keptUnlocks

	now := time.Now()
	g.Plots = []PlotState{
		{FlowerType: 0, Planted: now, Quantity: 1},
		{FlowerType: 0, Planted: now, Quantity: 1},
	}

	// Restore greenhouse plots from essence upgrades
	ghCount := g.EssenceUpgrades["greenhouse"]
	for i := 0; i < ghCount && i < len(g.Plots); i++ {
		g.Plots[i].IsGreenhouse = true
	}

	return nectar
}

// EssenceFromPrestige2 calculates essence earned from a layer 2 prestige.
func (g *GameState) EssenceFromPrestige2() float64 {
	if g.Nectar < 50 {
		return 0
	}
	return math.Floor(math.Log2(g.Nectar/10) + float64(g.PrestigeCount)/5)
}

// Prestige2 performs a layer 2 prestige, resetting nectar and prestige count.
func (g *GameState) Prestige2() float64 {
	essence := g.EssenceFromPrestige2()
	if essence <= 0 {
		return 0
	}

	g.Essence += essence
	g.Prestige2Count++
	g.AddLog(fmt.Sprintf("💫 Transcendence %d! Earned %.0f essence", g.Prestige2Count, essence), "#FF69B4")

	// Reset layer 1 + game
	g.Nectar = 0
	g.PrestigeCount = 0
	g.Petals = 0
	g.TotalPetals = 0
	g.Seeds = 0
	g.UpgradeLevels = make(map[string]int)
	g.DiscoveredHybrids = nil

	totalFlowers := len(FlowerTypes) + len(HybridFlowers)
	g.Unlocked = make([]bool, totalFlowers)
	g.Unlocked[0] = true

	now := time.Now()
	g.Plots = []PlotState{
		{FlowerType: 0, Planted: now, Quantity: 1},
		{FlowerType: 0, Planted: now, Quantity: 1},
	}

	return essence
}

// PetalsPerSecond estimates current petal generation rate.
func (g *GameState) PetalsPerSecond() float64 {
	if !g.HasAutoHarvest() {
		return 0
	}
	total := 0.0
	for i := range g.Plots {
		growTime := g.EffectiveGrowTimeForPlot(i)
		if growTime <= 0 {
			continue
		}
		harvestsPerSec := 1.0 / growTime.Seconds()
		yield := g.EffectiveYield(i)
		total += harvestsPerSec * yield
	}
	return total
}
