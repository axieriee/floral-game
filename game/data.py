"""Static data tables for procedural generation."""

# --- Biomes & Region Flavor ---

BIOMES = [
    {
        "name": "Verdant Glade",
        "desc": "Sunlight filters through a thick canopy of ancient oaks and flowering vines.",
        "enemies": ["Thorn Sprite", "Wild Boar", "Moss Golem"],
        "resources": ["Healing Herb", "Ironwood Branch", "Dewdrop Nectar"],
        "danger": 1,
    },
    {
        "name": "Ashen Wastes",
        "desc": "Cracked earth stretches endlessly under a haze of volcanic ash.",
        "enemies": ["Cinder Wraith", "Lava Beetle", "Ash Stalker"],
        "resources": ["Ember Shard", "Volcanic Glass", "Phoenix Feather"],
        "danger": 3,
    },
    {
        "name": "Frozen Hollow",
        "desc": "Ice crystals hang in the still air; your breath forms clouds of mist.",
        "enemies": ["Frost Wolf", "Ice Revenant", "Snow Harpy"],
        "resources": ["Frost Lotus", "Glacial Ore", "Yeti Fur"],
        "danger": 2,
    },
    {
        "name": "Sunken Marsh",
        "desc": "Murky water laps at gnarled roots; will-o'-wisps drift in the fog.",
        "enemies": ["Bog Leech", "Swamp Troll", "Mire Phantom"],
        "resources": ["Swamp Reed", "Luminous Fungus", "Toad Venom"],
        "danger": 2,
    },
    {
        "name": "Crystal Caverns",
        "desc": "Enormous crystals hum with arcane energy, casting prismatic light on the walls.",
        "enemies": ["Crystal Golem", "Shadow Bat", "Cave Lurker"],
        "resources": ["Mana Crystal", "Glowstone", "Stalagmite Iron"],
        "danger": 3,
    },
    {
        "name": "Whispering Dunes",
        "desc": "Sand shifts beneath your feet as an ancient wind carries forgotten words.",
        "enemies": ["Sand Wurm", "Dust Devil", "Scorpion King"],
        "resources": ["Desert Rose", "Sun Gold", "Cactus Water"],
        "danger": 2,
    },
    {
        "name": "Bloomwood Forest",
        "desc": "Every tree is in full bloom, petals drifting like pastel snow.",
        "enemies": ["Petal Dancer", "Briar Beast", "Fae Trickster"],
        "resources": ["Moonpetal", "Enchanted Sap", "Fairy Dust"],
        "danger": 1,
    },
    {
        "name": "Obsidian Peaks",
        "desc": "Jagged black mountains pierce the sky; lightning crackles overhead.",
        "enemies": ["Storm Raptor", "Rock Titan", "Thunder Elemental"],
        "resources": ["Thunderstone", "Obsidian Shard", "Eagle Feather"],
        "danger": 4,
    },
]

# --- Character Classes ---

CLASSES = {
    "Warrior": {
        "desc": "A battle-hardened fighter with unmatched strength.",
        "base_hp": 120,
        "base_atk": 14,
        "base_def": 10,
        "base_spd": 6,
        "base_mag": 3,
        "skills": ["Power Strike", "Shield Wall", "War Cry"],
    },
    "Mage": {
        "desc": "A wielder of arcane forces, fragile but devastating.",
        "base_hp": 70,
        "base_atk": 6,
        "base_def": 5,
        "base_spd": 7,
        "base_mag": 16,
        "skills": ["Fireball", "Frost Nova", "Arcane Shield"],
    },
    "Rogue": {
        "desc": "A swift shadow who strikes where it hurts most.",
        "base_hp": 85,
        "base_atk": 12,
        "base_def": 6,
        "base_spd": 14,
        "base_mag": 5,
        "skills": ["Backstab", "Smoke Bomb", "Poison Blade"],
    },
    "Ranger": {
        "desc": "A keen-eyed wanderer at home in any wilderness.",
        "base_hp": 95,
        "base_atk": 11,
        "base_def": 7,
        "base_spd": 10,
        "base_mag": 8,
        "skills": ["Aimed Shot", "Nature's Mend", "Trap"],
    },
}

# --- Skill Definitions ---

SKILLS = {
    "Power Strike": {"type": "physical", "power": 2.0, "cost": 0, "desc": "A crushing blow dealing double damage."},
    "Shield Wall": {"type": "buff_def", "power": 1.5, "cost": 0, "desc": "Raise your guard, boosting defense for 3 turns."},
    "War Cry": {"type": "buff_atk", "power": 1.3, "cost": 0, "desc": "Rally yourself, boosting attack for 3 turns."},
    "Fireball": {"type": "magic", "power": 2.5, "cost": 15, "desc": "Hurl a ball of flame at your foe."},
    "Frost Nova": {"type": "magic_debuff", "power": 1.5, "cost": 12, "desc": "Blast of frost that damages and slows."},
    "Arcane Shield": {"type": "buff_def", "power": 2.0, "cost": 10, "desc": "Conjure a magical barrier."},
    "Backstab": {"type": "physical", "power": 2.5, "cost": 0, "desc": "Strike from the shadows for massive damage."},
    "Smoke Bomb": {"type": "debuff_atk", "power": 0.6, "cost": 0, "desc": "Blind the enemy, reducing their attack."},
    "Poison Blade": {"type": "dot", "power": 0.8, "cost": 0, "desc": "Coat your blade in poison; enemy takes damage over 3 turns."},
    "Aimed Shot": {"type": "physical", "power": 2.2, "cost": 0, "desc": "A precise shot to a weak point."},
    "Nature's Mend": {"type": "heal", "power": 0.4, "cost": 10, "desc": "Channel nature to restore health."},
    "Trap": {"type": "debuff_spd", "power": 0.5, "cost": 0, "desc": "Set a trap that slows the enemy."},
}

# --- NPC Names & Dialogue ---

NPC_FIRST_NAMES = [
    "Aldric", "Brenna", "Cedric", "Dahlia", "Elara", "Fern", "Gareth",
    "Hazel", "Iris", "Jasper", "Kira", "Linden", "Mira", "Nolan",
    "Orin", "Petra", "Quinn", "Rowan", "Sage", "Thalia",
]

NPC_TITLES = [
    "the Wanderer", "the Wise", "the Bold", "the Herbalist",
    "the Merchant", "the Smith", "the Seer", "the Lost",
    "the Brave", "the Silent", "the Elder", "the Young",
]

NPC_GREETINGS = [
    "Well met, traveler. These roads grow more dangerous by the day.",
    "Ah, a fresh face! You look like someone who can handle themselves.",
    "Careful where you step... this land has teeth.",
    "Blessings upon you. I sense a great journey ahead of you.",
    "You there! I've been waiting for someone brave enough to help.",
    "The flowers whisper of your coming. Yes, I listen to flowers.",
    "Another adventurer? The last one didn't fare so well...",
    "Welcome! My wares are the finest in all the realms.",
]

# --- Quest Templates ---

QUEST_TEMPLATES = [
    {
        "name": "Slay the {enemy}",
        "desc": "A fearsome {enemy} terrorizes the region. Defeat it to restore peace.",
        "type": "kill",
        "reward_gold_mult": 1.5,
        "reward_xp_mult": 1.5,
    },
    {
        "name": "Gather {resource}",
        "desc": "A local herbalist needs {resource} found only in dangerous territory.",
        "type": "gather",
        "reward_gold_mult": 1.0,
        "reward_xp_mult": 1.2,
    },
    {
        "name": "Explore the {location}",
        "desc": "Rumors speak of treasure hidden deep within the {location}.",
        "type": "explore",
        "reward_gold_mult": 2.0,
        "reward_xp_mult": 1.8,
    },
    {
        "name": "Escort through {location}",
        "desc": "A merchant needs safe passage through the treacherous {location}.",
        "type": "escort",
        "reward_gold_mult": 1.3,
        "reward_xp_mult": 1.3,
    },
    {
        "name": "Retrieve the Lost {item}",
        "desc": "A precious {item} was stolen and hidden. Track it down and return it.",
        "type": "fetch",
        "reward_gold_mult": 1.2,
        "reward_xp_mult": 1.4,
    },
]

LOST_ITEMS = [
    "Amulet", "Crown", "Tome", "Blade", "Chalice",
    "Locket", "Staff", "Ring", "Map", "Relic",
]

# --- Item Generation ---

WEAPON_PREFIXES = ["Rusty", "Iron", "Steel", "Enchanted", "Legendary", "Cursed", "Ancient", "Gleaming"]
WEAPON_TYPES = ["Sword", "Axe", "Dagger", "Staff", "Bow", "Mace", "Spear", "Wand"]

ARMOR_PREFIXES = ["Tattered", "Leather", "Chain", "Plate", "Mystic", "Dragon", "Shadow", "Crystal"]
ARMOR_TYPES = ["Helmet", "Chestplate", "Leggings", "Boots", "Shield", "Cloak", "Gauntlets", "Bracers"]

CONSUMABLE_ITEMS = [
    {"name": "Health Potion", "type": "heal", "power": 30, "desc": "Restores 30 HP."},
    {"name": "Greater Health Potion", "type": "heal", "power": 60, "desc": "Restores 60 HP."},
    {"name": "Mana Elixir", "type": "mana", "power": 25, "desc": "Restores 25 MP."},
    {"name": "Strength Tonic", "type": "buff_atk", "power": 5, "desc": "Temporarily boosts attack by 5."},
    {"name": "Iron Skin Potion", "type": "buff_def", "power": 5, "desc": "Temporarily boosts defense by 5."},
    {"name": "Antidote", "type": "cure", "power": 0, "desc": "Cures poison and other ailments."},
    {"name": "Smoke Pellet", "type": "escape", "power": 0, "desc": "Guarantees escape from battle."},
]

# --- Dungeon Themes ---

DUNGEON_THEMES = [
    {"name": "Ruined Temple", "desc": "Crumbling pillars line a hall reclaimed by nature."},
    {"name": "Bandit Hideout", "desc": "A network of caves littered with stolen goods."},
    {"name": "Cursed Crypt", "desc": "The dead do not rest easily in this place."},
    {"name": "Dragon's Lair", "desc": "The air is thick with heat; bones crunch underfoot."},
    {"name": "Enchanted Grove", "desc": "A place where magic has gone wild and dangerous."},
    {"name": "Abandoned Mine", "desc": "Creaking timbers and echoing drips; something stirs in the dark."},
    {"name": "Flooded Ruins", "desc": "Knee-deep water fills ancient corridors carved with runes."},
    {"name": "Sky Fortress", "desc": "A crumbling keep floating among the clouds, defying gravity."},
]

# --- Event Pool ---

RANDOM_EVENTS = [
    {
        "text": "You find a hidden cache behind a loose stone!",
        "type": "loot",
    },
    {
        "text": "A traveling bard offers to play you a song of power.",
        "type": "buff",
    },
    {
        "text": "You stumble into a trap! Poison darts fly from the walls!",
        "type": "trap",
    },
    {
        "text": "A mysterious shrine glows before you. Do you pray?",
        "type": "shrine",
    },
    {
        "text": "You discover a bubbling spring of crystal-clear water.",
        "type": "heal",
    },
    {
        "text": "An old hermit blocks your path and demands a riddle contest.",
        "type": "riddle",
    },
    {
        "text": "The ground trembles... something enormous approaches!",
        "type": "ambush",
    },
    {
        "text": "You find a withered journal. The last entry reads: 'Do not trust the flowers.'",
        "type": "lore",
    },
]

RIDDLES = [
    ("I have cities, but no houses. I have mountains, but no trees. I have water, but no fish. What am I?", "map"),
    ("The more you take, the more you leave behind. What am I?", "footsteps"),
    ("I speak without a mouth and hear without ears. I have no body, but I come alive with the wind. What am I?", "echo"),
    ("What has roots that nobody sees, is taller than trees, up up it goes, and yet never grows?", "mountain"),
    ("I can be cracked, made, told, and played. What am I?", "joke"),
]
