"""Procedural generation for the game world."""

from game.data import (
    BIOMES, NPC_FIRST_NAMES, NPC_TITLES, NPC_GREETINGS,
    QUEST_TEMPLATES, LOST_ITEMS, WEAPON_PREFIXES, WEAPON_TYPES,
    ARMOR_PREFIXES, ARMOR_TYPES, CONSUMABLE_ITEMS, DUNGEON_THEMES,
)


def generate_region(rng, player_level):
    """Generate a region with a biome, points of interest, and connections."""
    biome = rng.choice(BIOMES)
    region = {
        "name": f"The {biome['name']}",
        "biome": biome,
        "visited": False,
        "npcs": [],
        "dungeon": None,
        "shop": None,
    }
    # NPCs
    num_npcs = rng.randint(1, 3)
    for _ in range(num_npcs):
        region["npcs"].append(generate_npc(rng, biome, player_level))

    # Chance of a dungeon
    if rng.random() < 0.6:
        region["dungeon"] = generate_dungeon(rng, biome, player_level)

    # Chance of a shop
    if rng.random() < 0.5:
        region["shop"] = generate_shop(rng, player_level)

    return region


def generate_npc(rng, biome, player_level):
    """Generate a named NPC with dialogue and possibly a quest."""
    name = f"{rng.choice(NPC_FIRST_NAMES)} {rng.choice(NPC_TITLES)}"
    greeting = rng.choice(NPC_GREETINGS)
    npc = {
        "name": name,
        "greeting": greeting,
        "quest": None,
    }
    if rng.random() < 0.6:
        npc["quest"] = generate_quest(rng, biome, player_level)
    return npc


def generate_quest(rng, biome, player_level):
    """Generate a quest from templates."""
    template = rng.choice(QUEST_TEMPLATES)
    enemy = rng.choice(biome["enemies"])
    resource = rng.choice(biome["resources"])
    item = rng.choice(LOST_ITEMS)

    name = template["name"].format(
        enemy=enemy, resource=resource, location=biome["name"], item=item,
    )
    desc = template["desc"].format(
        enemy=enemy, resource=resource, location=biome["name"], item=item,
    )

    base_gold = 20 + player_level * 10
    base_xp = 30 + player_level * 15

    return {
        "name": name,
        "desc": desc,
        "type": template["type"],
        "target": enemy if template["type"] == "kill" else resource if template["type"] == "gather" else item,
        "progress": 0,
        "goal": rng.randint(1, 3) if template["type"] in ("kill", "gather") else 1,
        "gold_reward": int(base_gold * template["reward_gold_mult"]),
        "xp_reward": int(base_xp * template["reward_xp_mult"]),
        "completed": False,
    }


def generate_weapon(rng, player_level):
    """Generate a random weapon scaled to player level."""
    tier = min(len(WEAPON_PREFIXES) - 1, (player_level - 1) // 2 + rng.randint(0, 1))
    prefix = WEAPON_PREFIXES[tier]
    wtype = rng.choice(WEAPON_TYPES)
    bonus = tier * 3 + rng.randint(1, 4)
    value = (tier + 1) * 15 + rng.randint(0, 10)
    return {
        "name": f"{prefix} {wtype}",
        "slot": "weapon",
        "bonus": bonus,
        "value": value,
        "desc": f"A {prefix.lower()} {wtype.lower()}. ATK +{bonus}.",
    }


def generate_armor(rng, player_level):
    """Generate a random armor piece scaled to player level."""
    tier = min(len(ARMOR_PREFIXES) - 1, (player_level - 1) // 2 + rng.randint(0, 1))
    prefix = ARMOR_PREFIXES[tier]
    atype = rng.choice(ARMOR_TYPES)
    bonus = tier * 2 + rng.randint(1, 3)
    value = (tier + 1) * 12 + rng.randint(0, 8)
    return {
        "name": f"{prefix} {atype}",
        "slot": "armor",
        "bonus": bonus,
        "value": value,
        "desc": f"A {prefix.lower()} {atype.lower()}. DEF +{bonus}.",
    }


def generate_loot(rng, player_level, danger):
    """Generate a random piece of loot (weapon, armor, or consumable)."""
    roll = rng.random()
    if roll < 0.3:
        return generate_weapon(rng, player_level + danger)
    elif roll < 0.6:
        return generate_armor(rng, player_level + danger)
    else:
        item = dict(rng.choice(CONSUMABLE_ITEMS))
        item["consumable"] = True
        return item


def generate_dungeon(rng, biome, player_level):
    """Generate a multi-room dungeon."""
    theme = rng.choice(DUNGEON_THEMES)
    num_rooms = rng.randint(3, 6)
    rooms = []
    for i in range(num_rooms):
        is_boss_room = (i == num_rooms - 1)
        room = {
            "index": i,
            "desc": _room_description(rng, theme, i, is_boss_room),
            "enemy": rng.choice(biome["enemies"]) if rng.random() < 0.7 or is_boss_room else None,
            "is_boss": is_boss_room,
            "loot": generate_loot(rng, player_level, biome["danger"]) if rng.random() < 0.5 else None,
            "cleared": False,
        }
        rooms.append(room)

    boss_names = [
        f"Ancient {rng.choice(biome['enemies'])}",
        f"Elder {rng.choice(biome['enemies'])}",
        f"{biome['name']} Guardian",
    ]

    return {
        "name": theme["name"],
        "theme_desc": theme["desc"],
        "rooms": rooms,
        "boss_name": rng.choice(boss_names),
        "cleared": False,
        "danger": biome["danger"],
    }


def _room_description(rng, theme, room_index, is_boss):
    """Generate a short room description."""
    ambient = [
        "Shadows dance on the walls.",
        "A cold draft chills your bones.",
        "Strange markings cover the floor.",
        "The air hums with faint energy.",
        "Water drips steadily from above.",
        "Cobwebs hang thick in every corner.",
        "Faded murals depict forgotten battles.",
        "A faint glow emanates from the far wall.",
    ]
    if is_boss:
        return f"Room {room_index + 1} - The final chamber. {rng.choice(ambient)} Something powerful waits here."
    return f"Room {room_index + 1} - {rng.choice(ambient)}"


def generate_shop(rng, player_level):
    """Generate a shop with items for sale."""
    items = []
    # Always have some potions
    for _ in range(rng.randint(2, 4)):
        item = dict(rng.choice(CONSUMABLE_ITEMS))
        item["consumable"] = True
        item["price"] = rng.randint(10, 30) + player_level * 3
        items.append(item)
    # Weapons and armor
    if rng.random() < 0.7:
        w = generate_weapon(rng, player_level)
        w["price"] = w["value"] + rng.randint(5, 15)
        items.append(w)
    if rng.random() < 0.7:
        a = generate_armor(rng, player_level)
        a["price"] = a["value"] + rng.randint(5, 15)
        items.append(a)
    return {"items": items}


def generate_world_map(rng, player_level, num_regions=5):
    """Generate a set of connected regions."""
    regions = []
    for _ in range(num_regions):
        regions.append(generate_region(rng, player_level))
    return regions
