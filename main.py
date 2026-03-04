#!/usr/bin/env python3
"""Floral Realms - A Procedural Text RPG.

Run this file to start the game:
    python main.py

Features:
    - Procedural world generation with seeded RNG
    - 8 unique biomes with distinct enemies and resources
    - 4 character classes with unique skill trees
    - Turn-based combat with skills, buffs, DoTs, and items
    - Procedural dungeons with boss encounters
    - NPC interactions and quest system
    - Shops for buying and selling gear
    - Random events: shrines, riddles, traps, lore, and more
    - Expandable frontier - discover new regions as you play
    - Game-over stats and kill tracking
"""

from game.engine import GameEngine


def main():
    try:
        engine = GameEngine()
        engine.run()
    except (KeyboardInterrupt, SystemExit):
        print("\n\n  Farewell, adventurer! Until next time.\n")


if __name__ == "__main__":
    main()
