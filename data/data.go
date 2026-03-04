package data

// Rarity levels for flowers
type Rarity int

const (
	Common   Rarity = iota
	Uncommon
	Rare
	Mythic
)

func (r Rarity) String() string {
	switch r {
	case Common:
		return "Common"
	case Uncommon:
		return "Uncommon"
	case Rare:
		return "Rare"
	case Mythic:
		return "Mythic"
	default:
		return "Unknown"
	}
}

func (r Rarity) Symbol() string {
	switch r {
	case Common:
		return "⚘"
	case Uncommon:
		return "❀"
	case Rare:
		return "✾"
	case Mythic:
		return "❁"
	default:
		return "?"
	}
}

// Flower represents a forageable plant
type Flower struct {
	Name      string
	Resonance string // magical property
	Desc      string
	Rarity    Rarity
	Potency   int // 1-5, ticks down over time
	MaxPotency int
	Biomes    []string // where it can be found
}

// Tradition is a character background/class
type Tradition struct {
	Name        string
	Symbol      string
	Desc        string
	Flavor      string
	Perk        string
	Drawback    string
	StartFlowers []string // names of starting flowers
	KnownBlends  []string // names of starting known blends
}

// Blend is a combination of two flowers
type Blend struct {
	Name       string
	Flower1    string
	Flower2    string
	Effect     string
	Desc       string
	Discovered bool
}

// NPC represents a non-player character
type NPC struct {
	Name        string
	Title       string
	Greeting    string
	Trust       int // 0-5
	Disposition string
	Dialogue    map[string]string // flower name -> reaction
}

// ForageSpot is a place to find flowers
type ForageSpot struct {
	Name     string
	Desc     string
	Flowers  []string // possible flower names
	Searched bool
}

// AllFlowers is the master flower catalog
var AllFlowers = map[string]Flower{
	// Common
	"Daisy": {
		Name: "Daisy", Resonance: "Clarity", Rarity: Common,
		Desc: "A simple white flower with a golden heart. Reveals what is hidden.",
		MaxPotency: 5, Potency: 5,
		Biomes: []string{"meadow", "village", "riverbank"},
	},
	"Clover": {
		Name: "Clover", Resonance: "Luck", Rarity: Common,
		Desc: "Three green leaves pressed tight. Fortune favors the finder.",
		MaxPotency: 5, Potency: 5,
		Biomes: []string{"meadow", "hillside", "village"},
	},
	"Dandelion": {
		Name: "Dandelion", Resonance: "Whisper", Rarity: Common,
		Desc: "Gone to seed, each filament carries a word on the wind.",
		MaxPotency: 4, Potency: 4,
		Biomes: []string{"meadow", "roadside", "ruins"},
	},
	"Chamomile": {
		Name: "Chamomile", Resonance: "Calm", Rarity: Common,
		Desc: "Tiny white petals around a dome of gold. Soothes all things.",
		MaxPotency: 5, Potency: 5,
		Biomes: []string{"meadow", "garden", "riverbank"},
	},
	// Uncommon
	"Lavender": {
		Name: "Lavender", Resonance: "Drowse", Rarity: Uncommon,
		Desc: "Purple spires that hum with quiet power. Eases the restless.",
		MaxPotency: 4, Potency: 4,
		Biomes: []string{"hillside", "garden", "village"},
	},
	"Thistle": {
		Name: "Thistle", Resonance: "Ward", Rarity: Uncommon,
		Desc: "Sharp and proud. Nothing passes a thistle's guard.",
		MaxPotency: 4, Potency: 4,
		Biomes: []string{"hillside", "ruins", "roadside"},
	},
	"Foxglove": {
		Name: "Foxglove", Resonance: "Venom", Rarity: Uncommon,
		Desc: "Beautiful and deadly. The only flower that cuts.",
		MaxPotency: 3, Potency: 3,
		Biomes: []string{"deep_wood", "ruins", "shadow"},
	},
	"Honeysuckle": {
		Name: "Honeysuckle", Resonance: "Charm", Rarity: Uncommon,
		Desc: "Sweet nectar that opens hearts and loosens tongues.",
		MaxPotency: 4, Potency: 4,
		Biomes: []string{"garden", "village", "meadow"},
	},
	"Forget-me-not": {
		Name: "Forget-me-not", Resonance: "Memory", Rarity: Uncommon,
		Desc: "Tiny blue stars. They hold what the mind lets go.",
		MaxPotency: 3, Potency: 3,
		Biomes: []string{"riverbank", "ruins", "garden"},
	},
	// Rare
	"Moonflower": {
		Name: "Moonflower", Resonance: "Sight", Rarity: Rare,
		Desc: "Blooms only under moonlight. Reveals what daylight cannot.",
		MaxPotency: 2, Potency: 2,
		Biomes: []string{"deep_wood", "shadow", "ruins"},
	},
	"Edelweiss": {
		Name: "Edelweiss", Resonance: "Endure", Rarity: Rare,
		Desc: "Born of frost and altitude. Grants strength to persist.",
		MaxPotency: 2, Potency: 2,
		Biomes: []string{"mountain", "hillside"},
	},
	"Ghost Orchid": {
		Name: "Ghost Orchid", Resonance: "Veil", Rarity: Rare,
		Desc: "Almost translucent. Makes the bearer hard to notice.",
		MaxPotency: 2, Potency: 2,
		Biomes: []string{"deep_wood", "shadow", "ruins"},
	},
	"Bloodroot": {
		Name: "Bloodroot", Resonance: "Bond", Rarity: Rare,
		Desc: "Red sap seeps from the stem. Forges deep connections.",
		MaxPotency: 2, Potency: 2,
		Biomes: []string{"deep_wood", "riverbank"},
	},
	// Mythic
	"Everbloom": {
		Name: "Everbloom", Resonance: "Eternal", Rarity: Mythic,
		Desc: "It does not fade. Its resonance depends on where you find it.",
		MaxPotency: 99, Potency: 99,
		Biomes: []string{"mythic"},
	},
}

// AllBlends is the master blend catalog
var AllBlends = []Blend{
	{Name: "Deep Rest", Flower1: "Lavender", Flower2: "Chamomile",
		Effect: "Full heal, but you lose half a day",
		Desc:   "The two scents mingle into something deeper than sleep."},
	{Name: "True Sight", Flower1: "Daisy", Flower2: "Moonflower",
		Effect: "Reveals the core secret of a location",
		Desc:   "Clarity meets moonlight. The world unfolds."},
	{Name: "Bittersweet", Flower1: "Foxglove", Flower2: "Honeysuckle",
		Effect: "Force an NPC truth, damages relationship",
		Desc:   "Sweet poison. The truth comes out, but it stings."},
	{Name: "Rootbound", Flower1: "Forget-me-not", Flower2: "Bloodroot",
		Effect: "An NPC shares their deepest memory",
		Desc:   "Memory and bond intertwine. Something ancient surfaces."},
	{Name: "Ironbark", Flower1: "Thistle", Flower2: "Edelweiss",
		Effect: "Immovable for one scene, flowers crumble after",
		Desc:   "Ward meets endurance. You become the mountain."},
	{Name: "Rumor Wind", Flower1: "Ghost Orchid", Flower2: "Dandelion",
		Effect: "Hear what NPCs say about you",
		Desc:   "Invisible whispers carried on dandelion seeds."},
	{Name: "Dreamwalk", Flower1: "Ghost Orchid", Flower2: "Forget-me-not",
		Effect: "Dream fragments reveal lore at night",
		Desc:   "The veil between memory and dream grows thin."},
	{Name: "Lucky Charm", Flower1: "Clover", Flower2: "Honeysuckle",
		Effect: "Better foraging and NPC reactions for 3 days",
		Desc:   "Luck and charm dance together. Doors open."},
	{Name: "Poison Ward", Flower1: "Foxglove", Flower2: "Thistle",
		Effect: "Creates a toxic barrier nothing can cross",
		Desc:   "Venom meets thorn. A deadly perimeter grows."},
	{Name: "Gentle Sight", Flower1: "Chamomile", Flower2: "Daisy",
		Effect: "Reveals hidden paths without alarming wildlife",
		Desc:   "Calm clarity. The world shows you its gentle secrets."},
}

// AllTraditions is the character class list
var AllTraditions = []Tradition{
	{
		Name:   "Hedge Keeper",
		Symbol: "⚘",
		Desc:   "Raised at the border between the wild and the tended.",
		Flavor: "You know the names of common flowers like old friends. The kettle is always on. The garden gate is always open.",
		Perk:   "Common flowers last longer in your journal. Villagers trust you.",
		Drawback: "The deep woods feel wrong. Rare flowers wilt faster in your hands.",
		StartFlowers: []string{"Daisy", "Chamomile"},
		KnownBlends:  []string{"Deep Rest"},
	},
	{
		Name:   "Thornwalker",
		Symbol: "❀",
		Desc:   "You came from the overgrown places.",
		Flavor: "Brambles part for you. The dark between the trees is just another shade of green. You've never owned a door.",
		Perk:   "Can forage in dangerous biomes. Recognizes rare species on sight.",
		Drawback: "Settlements find you unsettling. Shopkeepers charge you more.",
		StartFlowers: []string{"Thistle", "Foxglove"},
		KnownBlends:  []string{"Poison Ward"},
	},
	{
		Name:   "Petal Scholar",
		Symbol: "✾",
		Desc:   "Trained in a library of pressed specimens.",
		Flavor: "Your fingers are ink-stained, not earth-stained. You've read about every flower. You've touched very few.",
		Perk:   "Begin knowing 3 blends. Can carry more journal pages.",
		Drawback: "First foraging attempts may damage what you find.",
		StartFlowers: []string{"Lavender", "Forget-me-not"},
		KnownBlends:  []string{"Deep Rest", "True Sight", "Dreamwalk"},
	},
	{
		Name:   "Rootborn",
		Symbol: "❁",
		Desc:   "You are not entirely human. A symbiosis event in childhood left you part-plant.",
		Flavor: "Your veins run green in certain lights. Flowers lean toward you when you pass. You dream in seasons.",
		Perk:   "You don't press flowers — you GROW them. Unique interactions.",
		Drawback: "Bound to seasons. In winter, you weaken.",
		StartFlowers: []string{"Bloodroot", "Clover"},
		KnownBlends:  []string{"Rootbound"},
	},
}

// MeadowForageSpots for the prototype meadow biome
var MeadowForageSpots = []ForageSpot{
	{
		Name:    "Sun-warmed stones",
		Desc:    "A cluster of flat, pale stones baking in the sunlight. Small yellow and white flowers peek from the cracks between them.",
		Flowers: []string{"Daisy", "Chamomile", "Dandelion"},
	},
	{
		Name:    "The stream bank",
		Desc:    "The brook murmurs over smooth pebbles. The bank is soft and damp, crowded with delicate growth. Something blue catches your eye.",
		Flowers: []string{"Forget-me-not", "Chamomile", "Clover"},
	},
	{
		Name:    "An old stump",
		Desc:    "A gnarled stump draped in moss, surrounded by a ring of mushrooms. Climbing vines bear small, sweet-smelling flowers.",
		Flowers: []string{"Honeysuckle", "Clover", "Foxglove"},
	},
	{
		Name:    "The tall grass",
		Desc:    "Waist-high grass sways in the breeze, heavy with seed heads. Purple spires rise above the green, humming with bees.",
		Flowers: []string{"Lavender", "Thistle", "Dandelion"},
	},
}

// Maren is the prototype NPC
var Maren = NPC{
	Name:        "Maren",
	Title:       "the Beekeeper",
	Greeting:    "Oh! A visitor. The bees told me someone was coming, but I thought they were just being dramatic.",
	Trust:       2,
	Disposition: "cautious · curious",
	Dialogue: map[string]string{
		"Daisy":        "\"A daisy.\" She takes it gently, turning it in her calloused fingers. \"My daughter used to make chains of these.\" Her eyes go somewhere far away for a moment. Trust deepens.",
		"Lavender":     "\"Lavender...\" She breathes it in and her shoulders drop. \"I haven't slept well in weeks. The hives have been restless. Thank you.\" She looks at you differently now.",
		"Ghost Orchid":  "Her eyes widen. \"Where did you — I've only ever seen drawings of these. In the old texts.\" She reaches out, then pulls her hand back. \"You'd give this to me? Do you know what this is worth?\"",
		"Foxglove":     "She steps back. Her hand finds the doorframe. \"That's foxglove.\" Her voice has changed. \"My mother — someone brought her foxglove tea once. She never woke up.\" The warmth is gone from the room.",
		"Honeysuckle":  "She laughs — a surprised, bright sound. \"Honeysuckle! Come in, come in. I'll put the kettle on. My bees go mad for this stuff.\" She's already halfway to the kitchen.",
		"Chamomile":    "\"Chamomile. Sensible choice.\" She nods approvingly. \"I can always use more of this. The bees have been anxious lately — I brew it into a mist for the hives.\"",
		"Thistle":      "She raises an eyebrow. \"A thistle? You're either very practical or very prickly.\" A half-smile. \"I respect that. My garden could use some guardians.\"",
		"Clover":       "\"Clover! Oh, my bees will love you.\" She brightens immediately. \"There's a patch behind the cottage they've nearly stripped bare. This will help.\"",
		"Dandelion":    "\"A dandelion.\" She blows on it gently and the seeds scatter inside the cottage. \"Oops.\" She doesn't look sorry. \"Make a wish, I suppose.\"",
		"Forget-me-not": "She takes the tiny blue flower and holds it up to the light. \"These grew by the river where I grew up.\" Quiet for a moment. \"I haven't been back in a long time.\"",
		"Bloodroot":    "Her face goes still. \"That's bloodroot. You're either a healer or a fool, and healers usually know better than to offer it to a stranger.\" She studies you carefully.",
		"Moonflower":   "\"A moonflower.\" She whispers it like a secret. \"I've read about these. They say if you press one under a full moon, it never fully fades.\" She looks at you with genuine awe.",
	},
}
