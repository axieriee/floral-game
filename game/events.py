"""Random events, shrine interactions, riddles, and world encounters."""

from game.data import RANDOM_EVENTS, RIDDLES
from game.procgen import generate_loot


def trigger_random_event(player, rng, ui, player_level):
    """Roll for and execute a random event. Returns True if something happened."""
    if rng.random() > 0.4:
        return False

    event = rng.choice(RANDOM_EVENTS)
    ui.print_line(f"\n  * {event['text']}")

    if event["type"] == "loot":
        loot = generate_loot(rng, player_level, 1)
        ui.print_line(f"  You found: {loot['name']}!")
        if loot.get("slot") == "weapon":
            if ui.confirm(f"  Equip {loot['name']}? (ATK +{loot['bonus']})"):
                old = player.equip_weapon(loot)
                if old:
                    ui.print_line(f"  Unequipped {old['name']}.")
            else:
                player.inventory.append(loot)
        elif loot.get("slot") == "armor":
            if ui.confirm(f"  Equip {loot['name']}? (DEF +{loot['bonus']})"):
                old = player.equip_armor(loot)
                if old:
                    ui.print_line(f"  Unequipped {old['name']}.")
            else:
                player.inventory.append(loot)
        else:
            player.inventory.append(loot)

    elif event["type"] == "buff":
        ui.print_line("  The bard's melody fills you with vigor!")
        player.buffs.append(("atk", 1.2, 10))
        player.buffs.append(("defense", 1.2, 10))

    elif event["type"] == "trap":
        dmg = rng.randint(5, 15) + player_level * 2
        actual = max(1, dmg - player.effective_stat("defense") // 4)
        player.hp = max(1, player.hp - actual)
        ui.print_line(f"  You take {actual} damage from the trap!")

    elif event["type"] == "shrine":
        _handle_shrine(player, rng, ui)

    elif event["type"] == "heal":
        healed = player.heal(int(player.max_hp * 0.3))
        mp_restored = player.restore_mp(int(player.max_mp * 0.3))
        ui.print_line(f"  You drink deeply. Restored {healed} HP and {mp_restored} MP!")

    elif event["type"] == "riddle":
        _handle_riddle(player, rng, ui)

    elif event["type"] == "ambush":
        ui.print_line("  A powerful enemy appears! Prepare for battle!")
        return "ambush"

    elif event["type"] == "lore":
        ui.print_line("  You gain insight from the journal. +10 XP.")
        msgs = player.xp_gain(10)
        for m in msgs:
            ui.print_line(m)

    return True


def _handle_shrine(player, rng, ui):
    """Interactive shrine that can bless or curse."""
    ui.print_line("  1. Pray at the shrine")
    ui.print_line("  2. Leave it alone")
    choice = ui.get_choice(2)
    if choice == 1:
        if rng.random() < 0.7:
            bonus = rng.choice(["hp", "atk", "defense", "mag"])
            if bonus == "hp":
                player.max_hp += 10
                player.hp = min(player.hp + 10, player.max_hp)
                ui.print_line("  The shrine blesses you! Max HP +10!")
            elif bonus == "atk":
                player.atk += 2
                ui.print_line("  The shrine blesses you! ATK +2!")
            elif bonus == "defense":
                player.defense += 2
                ui.print_line("  The shrine blesses you! DEF +2!")
            else:
                player.mag += 2
                ui.print_line("  The shrine blesses you! MAG +2!")
        else:
            dmg = rng.randint(10, 25)
            player.hp = max(1, player.hp - dmg)
            ui.print_line(f"  The shrine crackles with dark energy! You take {dmg} damage!")
    else:
        ui.print_line("  You leave the shrine undisturbed.")


def _handle_riddle(player, rng, ui):
    """Riddle challenge for bonus rewards."""
    riddle, answer = rng.choice(RIDDLES)
    ui.print_line(f"\n  The hermit speaks: \"{riddle}\"")
    ui.print_line("  What is your answer?")
    response = ui.get_input("  Answer: ").lower().strip()
    if response == answer:
        gold = rng.randint(15, 40)
        player.gold += gold
        ui.print_line(f"  \"Correct!\" The hermit tosses you a pouch of {gold} gold!")
        msgs = player.xp_gain(25)
        for m in msgs:
            ui.print_line(m)
    else:
        ui.print_line(f"  \"Wrong! The answer was '{answer}'.\" The hermit vanishes.")
        ui.print_line("  You feel slightly drained. -5 HP.")
        player.hp = max(1, player.hp - 5)
