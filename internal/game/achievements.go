package game

import "fmt"

// AchievementID uniquely identifies an achievement.
type AchievementID string

// Achievement defines a trackable milestone.
type Achievement struct {
	ID          AchievementID
	Name        string
	Description string
	Emoji       string
	RewardType  string  // "petals", "seeds", "nectar"
	Reward      float64 // amount of reward
	Hidden      bool    // hidden until unlocked
}

// AchievementCheck is a function that returns true if the achievement condition is met.
type AchievementCheck func(g *GameState) bool

var Achievements = []Achievement{
	// Harvest milestones
	{ID: "first_harvest", Name: "First Bloom", Description: "Harvest your first flower", Emoji: "🌸", RewardType: "petals", Reward: 5},
	{ID: "harvest_50", Name: "Green Thumb", Description: "Harvest 50 flowers", Emoji: "👍", RewardType: "petals", Reward: 50},
	{ID: "harvest_500", Name: "Master Gardener", Description: "Harvest 500 flowers", Emoji: "🏅", RewardType: "petals", Reward: 500},
	{ID: "harvest_5000", Name: "Harvest Legend", Description: "Harvest 5,000 flowers", Emoji: "🏆", RewardType: "seeds", Reward: 200},

	// Petal milestones
	{ID: "petals_100", Name: "Petal Collector", Description: "Earn 100 lifetime petals", Emoji: "💐", RewardType: "seeds", Reward: 10},
	{ID: "petals_1000", Name: "Petal Hoarder", Description: "Earn 1,000 lifetime petals", Emoji: "💰", RewardType: "seeds", Reward: 50},
	{ID: "petals_10000", Name: "Petal Tycoon", Description: "Earn 10,000 lifetime petals", Emoji: "👑", RewardType: "seeds", Reward: 200},
	{ID: "petals_100000", Name: "Petal Emperor", Description: "Earn 100,000 lifetime petals", Emoji: "💎", RewardType: "nectar", Reward: 5},

	// Flower unlocks
	{ID: "unlock_3", Name: "Budding Botanist", Description: "Unlock 3 flower types", Emoji: "📖", RewardType: "petals", Reward: 30},
	{ID: "unlock_all_base", Name: "Full Spectrum", Description: "Unlock all base flowers", Emoji: "🌈", RewardType: "seeds", Reward: 100},
	{ID: "hybrid_1", Name: "First Cross", Description: "Discover your first hybrid", Emoji: "🧬", RewardType: "petals", Reward: 100},
	{ID: "hybrid_all", Name: "Hybridization Master", Description: "Discover all hybrid flowers", Emoji: "🧪", RewardType: "nectar", Reward: 10},

	// Upgrade milestones
	{ID: "first_upgrade", Name: "Handy Work", Description: "Purchase your first upgrade", Emoji: "🔧", RewardType: "petals", Reward: 10},
	{ID: "auto_harvest", Name: "Automation Age", Description: "Unlock auto-harvest", Emoji: "🤖", RewardType: "seeds", Reward: 25},
	{ID: "max_plots", Name: "Land Baron", Description: "Unlock all garden plots", Emoji: "🏡", RewardType: "seeds", Reward: 500},

	// Prestige
	{ID: "first_prestige", Name: "Rebirth", Description: "Prestige for the first time", Emoji: "✨", RewardType: "seeds", Reward: 50},
	{ID: "first_transcend", Name: "Transcendence", Description: "Transcend for the first time", Emoji: "💫", RewardType: "nectar", Reward: 3},

	// Combo
	{ID: "combo_5", Name: "Quick Fingers", Description: "Reach a 5x harvest combo", Emoji: "⚡", RewardType: "petals", Reward: 25},
	{ID: "combo_10", Name: "Combo Master", Description: "Reach a 10x harvest combo", Emoji: "🔥", RewardType: "petals", Reward: 100},
	{ID: "combo_25", Name: "Unstoppable", Description: "Reach a 25x harvest combo", Emoji: "💥", RewardType: "seeds", Reward: 100},

	// Events
	{ID: "survive_drought", Name: "Drought Survivor", Description: "Harvest during a drought", Emoji: "🏜️", RewardType: "petals", Reward: 50},
	{ID: "golden_harvest", Name: "Golden Touch", Description: "Harvest during Golden Hour", Emoji: "🌟", RewardType: "seeds", Reward: 30},

	// Secret achievements
	{ID: "night_owl", Name: "Night Owl", Description: "Play during winter season", Emoji: "🦉", RewardType: "petals", Reward: 20, Hidden: true},
	{ID: "full_garden", Name: "Paradise", Description: "Have all plots growing at once", Emoji: "🌺", RewardType: "seeds", Reward: 50, Hidden: true},
}

// achievementChecks maps achievement IDs to their check functions.
var achievementChecks = map[AchievementID]AchievementCheck{
	"first_harvest":  func(g *GameState) bool { return g.TotalHarvests >= 1 },
	"harvest_50":     func(g *GameState) bool { return g.TotalHarvests >= 50 },
	"harvest_500":    func(g *GameState) bool { return g.TotalHarvests >= 500 },
	"harvest_5000":   func(g *GameState) bool { return g.TotalHarvests >= 5000 },
	"petals_100":     func(g *GameState) bool { return g.TotalPetals >= 100 },
	"petals_1000":    func(g *GameState) bool { return g.TotalPetals >= 1000 },
	"petals_10000":   func(g *GameState) bool { return g.TotalPetals >= 10000 },
	"petals_100000":  func(g *GameState) bool { return g.TotalPetals >= 100000 },
	"unlock_3":       func(g *GameState) bool { return g.CountUnlocked() >= 3 },
	"unlock_all_base": func(g *GameState) bool { return g.AllBaseUnlocked() },
	"hybrid_1":       func(g *GameState) bool { return len(g.DiscoveredHybrids) >= 1 },
	"hybrid_all":     func(g *GameState) bool { return len(g.DiscoveredHybrids) >= len(HybridFlowers) },
	"first_upgrade":  func(g *GameState) bool { return g.TotalUpgradesBought() >= 1 },
	"auto_harvest":   func(g *GameState) bool { return g.HasAutoHarvest() },
	"max_plots":      func(g *GameState) bool { return g.UpgradeLevels["new_plot"] >= 8 },
	"first_prestige": func(g *GameState) bool { return g.PrestigeCount >= 1 },
	"first_transcend": func(g *GameState) bool { return g.Prestige2Count >= 1 },
	"combo_5":        func(g *GameState) bool { return g.BestCombo >= 5 },
	"combo_10":       func(g *GameState) bool { return g.BestCombo >= 10 },
	"combo_25":       func(g *GameState) bool { return g.BestCombo >= 25 },
	"survive_drought": func(g *GameState) bool { return g.HarvestedDuringEvent("drought") },
	"golden_harvest":  func(g *GameState) bool { return g.HarvestedDuringEvent("golden_hour") },
	"night_owl":      func(g *GameState) bool { return g.CurrentSeason() == Winter },
	"full_garden": func(g *GameState) bool {
		if len(g.Plots) < 3 {
			return false
		}
		for i := range g.Plots {
			if g.PlotProgress(i) >= 1.0 {
				return false // a ready plot means it's not actively growing
			}
		}
		return true
	},
}

// CheckAchievements checks all achievements and returns newly completed ones.
func (g *GameState) CheckAchievements() []Achievement {
	var newlyCompleted []Achievement
	for _, a := range Achievements {
		if g.CompletedAchievements[a.ID] {
			continue
		}
		check, ok := achievementChecks[a.ID]
		if !ok {
			continue
		}
		if check(g) {
			g.CompletedAchievements[a.ID] = true
			// Grant reward
			switch a.RewardType {
			case "petals":
				g.Petals += a.Reward
				g.TotalPetals += a.Reward
			case "seeds":
				g.Seeds += a.Reward
			case "nectar":
				g.Nectar += a.Reward
			}
			newlyCompleted = append(newlyCompleted, a)
			g.AddLog(fmt.Sprintf("%s Achievement: %s!", a.Emoji, a.Name), "#FFD700")
		}
	}
	return newlyCompleted
}

// CountUnlocked returns the number of unlocked flowers.
func (g *GameState) CountUnlocked() int {
	count := 0
	for _, u := range g.Unlocked {
		if u {
			count++
		}
	}
	return count
}

// AllBaseUnlocked returns true if all base flowers are unlocked.
func (g *GameState) AllBaseUnlocked() bool {
	for i := range FlowerTypes {
		if i < len(g.Unlocked) && !g.Unlocked[i] {
			return false
		}
	}
	return true
}

// TotalUpgradesBought returns total upgrade levels purchased.
func (g *GameState) TotalUpgradesBought() int {
	total := 0
	for _, v := range g.UpgradeLevels {
		total += v
	}
	return total
}

// HarvestedDuringEvent returns whether the player has harvested during a specific event type.
func (g *GameState) HarvestedDuringEvent(eventEffect string) bool {
	if g.EventHarvests == nil {
		return false
	}
	return g.EventHarvests[eventEffect]
}

// AchievementProgress returns (completed, total) counts.
func AchievementProgress(g *GameState) (int, int) {
	completed := 0
	for _, a := range Achievements {
		if g.CompletedAchievements[a.ID] {
			completed++
		}
	}
	return completed, len(Achievements)
}
