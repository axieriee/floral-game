package game

import (
	"math/rand"
	"time"
)

// GardenEvent represents an active event affecting the garden.
type GardenEvent struct {
	Type      EventType `json:"type"`
	StartTime time.Time `json:"start_time"`
	Duration  time.Duration `json:"duration"`
}

// EventType identifies a garden event.
type EventType string

const (
	EventNone       EventType = ""
	EventPetalRain  EventType = "petal_rain"
	EventBeeSurge   EventType = "bee_surge"
	EventDrought    EventType = "drought"
	EventGoldenHour EventType = "golden_hour"
	EventWindfall   EventType = "windfall"
	EventFrostSnap  EventType = "frost_snap"
)

// EventInfo describes an event type.
type EventInfo struct {
	Type        EventType
	Name        string
	Emoji       string
	Description string
	Duration    time.Duration
	Weight      int // relative probability
	Positive    bool
}

var EventTypes = []EventInfo{
	{
		Type: EventPetalRain, Name: "Petal Rain", Emoji: "🌧️",
		Description: "A magical rain doubles petal yield!",
		Duration: 30 * time.Second, Weight: 20, Positive: true,
	},
	{
		Type: EventBeeSurge, Name: "Bee Surge", Emoji: "🐝",
		Description: "A swarm of bees speeds growth by 50%!",
		Duration: 25 * time.Second, Weight: 20, Positive: true,
	},
	{
		Type: EventDrought, Name: "Drought", Emoji: "🏜️",
		Description: "A dry spell slows growth by 40%",
		Duration: 20 * time.Second, Weight: 15, Positive: false,
	},
	{
		Type: EventGoldenHour, Name: "Golden Hour", Emoji: "🌟",
		Description: "Everything shines — 3x petal yield!",
		Duration: 15 * time.Second, Weight: 10, Positive: true,
	},
	{
		Type: EventWindfall, Name: "Windfall", Emoji: "🍃",
		Description: "Seeds blow in from afar!",
		Duration: 0, Weight: 15, Positive: true, // instant effect
	},
	{
		Type: EventFrostSnap, Name: "Frost Snap", Emoji: "🥶",
		Description: "A sudden frost — yield halved but seeds doubled!",
		Duration: 20 * time.Second, Weight: 10, Positive: false,
	},
}

// GetEventInfo returns info for an event type.
func GetEventInfo(t EventType) EventInfo {
	for _, e := range EventTypes {
		if e.Type == t {
			return e
		}
	}
	return EventInfo{}
}

// IsEventActive returns whether a garden event is currently active.
func (g *GameState) IsEventActive() bool {
	if g.ActiveEvent == nil {
		return false
	}
	if g.ActiveEvent.Type == EventNone {
		return false
	}
	return time.Since(g.ActiveEvent.StartTime) < g.ActiveEvent.Duration
}

// EventTimeRemaining returns how much time is left on the active event.
func (g *GameState) EventTimeRemaining() time.Duration {
	if !g.IsEventActive() {
		return 0
	}
	elapsed := time.Since(g.ActiveEvent.StartTime)
	remaining := g.ActiveEvent.Duration - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// EventGrowthMult returns the growth multiplier from the active event.
func (g *GameState) EventGrowthMult() float64 {
	if !g.IsEventActive() {
		return 1.0
	}
	switch g.ActiveEvent.Type {
	case EventBeeSurge:
		return 1.5
	case EventDrought:
		return 0.6
	case EventFrostSnap:
		return 0.8
	}
	return 1.0
}

// EventYieldMult returns the yield multiplier from the active event.
func (g *GameState) EventYieldMult() float64 {
	if !g.IsEventActive() {
		return 1.0
	}
	switch g.ActiveEvent.Type {
	case EventPetalRain:
		return 2.0
	case EventGoldenHour:
		return 3.0
	case EventFrostSnap:
		return 0.5
	}
	return 1.0
}

// EventSeedMult returns the seed multiplier from the active event.
func (g *GameState) EventSeedMult() float64 {
	if !g.IsEventActive() {
		return 1.0
	}
	switch g.ActiveEvent.Type {
	case EventFrostSnap:
		return 2.0
	}
	return 1.0
}

// TryTriggerEvent rolls for a random event. Called each tick.
// Returns the event info if one triggered, nil otherwise.
func (g *GameState) TryTriggerEvent() *EventInfo {
	// Don't trigger if one is already active
	if g.IsEventActive() {
		return nil
	}

	// Cooldown: at least 30 seconds between events
	if g.ActiveEvent != nil && time.Since(g.ActiveEvent.StartTime) < 30*time.Second {
		return nil
	}

	// ~2% chance per tick (ticks are 200ms, so roughly once per 10 seconds on average)
	if rand.Float64() > 0.004 {
		return nil
	}

	// Weighted random selection
	totalWeight := 0
	for _, e := range EventTypes {
		totalWeight += e.Weight
	}
	roll := rand.Intn(totalWeight)
	cumulative := 0
	for _, e := range EventTypes {
		cumulative += e.Weight
		if roll < cumulative {
			// Handle instant events
			if e.Type == EventWindfall {
				seedBonus := 10.0 + float64(len(g.Plots))*5
				g.Seeds += seedBonus
				g.ActiveEvent = &GardenEvent{
					Type:      e.Type,
					StartTime: time.Now(),
					Duration:  3 * time.Second, // brief display
				}
				return &e
			}

			g.ActiveEvent = &GardenEvent{
				Type:      e.Type,
				StartTime: time.Now(),
				Duration:  e.Duration,
			}
			return &e
		}
	}

	return nil
}
