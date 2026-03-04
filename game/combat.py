"""Turn-based combat system."""

from game.data import SKILLS


def resolve_skill(user, target, skill_name):
    """Execute a skill and return a description of what happened."""
    skill = SKILLS[skill_name]
    messages = []

    if skill["type"] == "physical":
        raw = int(user.effective_stat("atk") * skill["power"])
        dmg = target.take_damage(raw)
        messages.append(f"{user.name} uses {skill_name}! Deals {dmg} damage!")

    elif skill["type"] == "magic":
        if user.mp < skill["cost"]:
            messages.append(f"Not enough MP for {skill_name}!")
            return messages, False
        user.mp -= skill["cost"]
        raw = int(user.effective_stat("mag") * skill["power"])
        dmg = target.take_damage(raw)
        messages.append(f"{user.name} casts {skill_name}! Deals {dmg} magic damage! (MP -{skill['cost']})")

    elif skill["type"] == "magic_debuff":
        if user.mp < skill["cost"]:
            messages.append(f"Not enough MP for {skill_name}!")
            return messages, False
        user.mp -= skill["cost"]
        raw = int(user.effective_stat("mag") * skill["power"])
        dmg = target.take_damage(raw)
        target.buffs.append(("spd", 0.5, 3))
        messages.append(f"{user.name} casts {skill_name}! Deals {dmg} damage and slows the enemy!")

    elif skill["type"] == "buff_def":
        cost = skill.get("cost", 0)
        if cost and user.mp < cost:
            messages.append(f"Not enough MP for {skill_name}!")
            return messages, False
        if cost:
            user.mp -= cost
        user.buffs.append(("defense", skill["power"], 3))
        messages.append(f"{user.name} uses {skill_name}! Defense boosted!")

    elif skill["type"] == "buff_atk":
        user.buffs.append(("atk", skill["power"], 3))
        messages.append(f"{user.name} uses {skill_name}! Attack boosted!")

    elif skill["type"] == "debuff_atk":
        target.buffs.append(("atk", skill["power"], 3))
        messages.append(f"{user.name} uses {skill_name}! Enemy attack reduced!")

    elif skill["type"] == "debuff_spd":
        target.buffs.append(("spd", skill["power"], 3))
        messages.append(f"{user.name} uses {skill_name}! Enemy slowed!")

    elif skill["type"] == "dot":
        dot_dmg = int(user.effective_stat("atk") * skill["power"])
        target.dots.append((dot_dmg, 3))
        messages.append(f"{user.name} uses {skill_name}! Enemy will take {dot_dmg} damage per turn for 3 turns!")

    elif skill["type"] == "heal":
        if user.mp < skill["cost"]:
            messages.append(f"Not enough MP for {skill_name}!")
            return messages, False
        user.mp -= skill["cost"]
        amount = int(user.effective_stat("mag") * skill["power"] + user.max_hp * 0.15)
        healed = user.heal(amount)
        messages.append(f"{user.name} uses {skill_name}! Restored {healed} HP!")

    return messages, True


def enemy_turn(enemy, player, rng):
    """Simple enemy AI."""
    raw = enemy.effective_stat("atk") + rng.randint(-2, 2)
    dmg = player.take_damage(raw)
    return f"{enemy.name} attacks! Deals {dmg} damage to you!"


def run_combat(player, enemy, rng, ui):
    """Full combat loop. Returns True if player wins, False if defeated, None if fled."""
    ui.print_line(f"\n{'='*50}")
    ui.print_line(f"  BATTLE: {player.name} vs {enemy.name} (Lv.{enemy.level})")
    ui.print_line(f"{'='*50}")
    ui.print_line(f"  Enemy HP: {enemy.hp}/{enemy.max_hp}")
    ui.print_line("")

    turn = 0
    while player.alive and enemy.alive:
        turn += 1
        ui.print_line(f"--- Turn {turn} ---")
        ui.print_line(f"  Your HP: {player.hp}/{player.max_hp}  MP: {player.mp}/{player.max_mp}")
        ui.print_line(f"  {enemy.name} HP: {enemy.hp}/{enemy.max_hp}")
        ui.print_line("")

        # Player turn
        options = ["Attack"]
        for s in player.skills:
            skill = SKILLS[s]
            cost_str = f" (MP: {skill['cost']})" if skill.get("cost", 0) else ""
            options.append(f"Skill: {s}{cost_str}")
        options.append("Use Item")
        options.append("Flee")

        ui.print_line("Actions:")
        for i, opt in enumerate(options, 1):
            ui.print_line(f"  {i}. {opt}")

        choice = ui.get_choice(len(options))

        if choice == 1:
            # Basic attack
            raw = player.effective_stat("atk") + rng.randint(-1, 3)
            dmg = enemy.take_damage(raw)
            ui.print_line(f"\nYou strike the {enemy.name} for {dmg} damage!")

        elif 2 <= choice <= 1 + len(player.skills):
            skill_name = player.skills[choice - 2]
            msgs, success = resolve_skill(player, enemy, skill_name)
            for m in msgs:
                ui.print_line(m)
            if not success:
                continue

        elif choice == len(options) - 1:
            # Use item
            consumables = [(i, item) for i, item in enumerate(player.inventory) if isinstance(item, dict) and item.get("consumable")]
            if not consumables:
                ui.print_line("You have no usable items!")
                continue
            ui.print_line("Choose an item:")
            for idx, (inv_idx, item) in enumerate(consumables, 1):
                ui.print_line(f"  {idx}. {item['name']} - {item['desc']}")
            ui.print_line(f"  {len(consumables)+1}. Cancel")
            item_choice = ui.get_choice(len(consumables) + 1)
            if item_choice == len(consumables) + 1:
                continue
            _, item = consumables[item_choice - 1]
            player.inventory.remove(item)
            if item["type"] == "heal":
                healed = player.heal(item["power"])
                ui.print_line(f"Used {item['name']}! Restored {healed} HP.")
            elif item["type"] == "mana":
                restored = player.restore_mp(item["power"])
                ui.print_line(f"Used {item['name']}! Restored {restored} MP.")
            elif item["type"] == "buff_atk":
                player.buffs.append(("atk", 1 + item["power"] / player.atk, 5))
                ui.print_line(f"Used {item['name']}! Attack boosted!")
            elif item["type"] == "buff_def":
                player.buffs.append(("defense", 1 + item["power"] / max(1, player.defense), 5))
                ui.print_line(f"Used {item['name']}! Defense boosted!")
            elif item["type"] == "escape":
                ui.print_line("You vanish in a cloud of smoke!")
                return None

        elif choice == len(options):
            # Flee
            flee_chance = 0.4 + (player.effective_stat("spd") - enemy.effective_stat("spd")) * 0.05
            if rng.random() < flee_chance:
                ui.print_line("You successfully flee from battle!")
                return None
            else:
                ui.print_line("You failed to escape!")

        if not enemy.alive:
            break

        # DoT damage on enemy
        dot_dmg = enemy.tick_dots()
        if dot_dmg > 0:
            ui.print_line(f"  Poison deals {dot_dmg} to {enemy.name}!")
        if not enemy.alive:
            break

        # Enemy turn
        msg = enemy_turn(enemy, player, rng)
        ui.print_line(msg)

        # DoT damage on player
        dot_dmg = player.tick_dots()
        if dot_dmg > 0:
            ui.print_line(f"  Poison deals {dot_dmg} to you!")

        # Tick buffs
        player.tick_buffs()
        enemy.tick_buffs()
        ui.print_line("")

    if not player.alive:
        return False

    # Victory!
    ui.print_line(f"\n*** Victory! {enemy.name} is defeated! ***")
    player.record_kill(enemy.name)
    xp_msgs = player.xp_gain(enemy.xp_reward)
    player.gold += enemy.gold_reward
    ui.print_line(f"  Gained {enemy.xp_reward} XP and {enemy.gold_reward} gold!")
    for msg in xp_msgs:
        ui.print_line(msg)
    return True
