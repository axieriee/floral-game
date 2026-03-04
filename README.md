# Floral Realms - A Procedural Text RPG

A fully procedural text-based RPG built in Python with no external dependencies.

## How to Play

```bash
python main.py
```

## Features

- **Procedural World Generation** - Every playthrough generates a unique world from a seed. Use the same seed to replay the same world.
- **8 Unique Biomes** - Verdant Glade, Ashen Wastes, Frozen Hollow, Sunken Marsh, Crystal Caverns, Whispering Dunes, Bloomwood Forest, and Obsidian Peaks.
- **4 Character Classes** - Warrior, Mage, Rogue, and Ranger, each with unique stats and 3 skills.
- **Turn-Based Combat** - Skills, buffs, debuffs, damage-over-time, items, and flee mechanics.
- **Procedural Dungeons** - Multi-room dungeons with loot, enemies, and boss encounters.
- **NPC Quests** - Procedurally generated NPCs offer kill, gather, explore, escort, and fetch quests.
- **Shops** - Buy weapons, armor, and consumables. Sell your loot for gold.
- **Random Events** - Shrines, riddle challenges, traps, traveling bards, hidden caches, and lore discoveries.
- **Expandable Frontier** - Discover new regions beyond the starting world.
- **Character Progression** - Level up to increase stats, find better gear, complete quests.

## Project Structure

```
floral-game/
├── main.py              # Entry point
├── game/
│   ├── __init__.py
│   ├── rng.py           # Seeded random number generator
│   ├── data.py          # Static data tables (biomes, classes, items, etc.)
│   ├── character.py     # Player, Enemy, and Boss entities
│   ├── combat.py        # Turn-based combat system
│   ├── procgen.py       # Procedural generation (world, dungeons, NPCs, items)
│   ├── events.py        # Random events, shrines, riddles
│   └── engine.py        # Main game engine and loop
└── requirements.txt     # No dependencies needed
```

## Requirements

- Python 3.6+
- No external packages required
