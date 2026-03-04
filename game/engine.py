"""Main game engine - ties all systems together."""

from game.rng import GameRNG
from game.ui import TerminalUI, clear_screen
from game.character import Player, Enemy, Boss
from game.combat import run_combat
from game.procgen import generate_world_map, generate_loot, generate_region
from game.events import trigger_random_event
from game.data import CLASSES, BIOMES


class GameEngine:
    """Core game loop and state management."""

    def __init__(self):
        self.ui = TerminalUI()
        self.rng = None
        self.player = None
        self.world = []
        self.current_region_idx = 0
        self.turn = 0
        self.game_over = False

    def run(self):
        """Main entry point."""
        clear_screen()
        self._title_screen()
        seed = self._setup_seed()
        self.rng = GameRNG(seed)
        self.ui.print_line(f"  World seed: {self.rng.seed}")
        self.ui.print_line("")
        self._create_character()
        self.world = generate_world_map(self.rng, self.player.level)
        self.ui.print_line(f"\n  A world of {len(self.world)} regions has been generated!")
        self.ui.print_line("  Type 'help' at any time for a list of commands.\n")

        self._main_loop()

    def _title_screen(self):
        self.ui.print_banner("FLORAL REALMS")
        self.ui.print_box([
            "A Procedural Text RPG",
            "",
            "Every world is unique. Every path is yours.",
            "Explore, fight, quest, and survive.",
            "",
            "Enter a seed for a reproducible world,",
            "or press Enter for a random adventure.",
        ])
        self.ui.print_line("")

    def _setup_seed(self):
        self.ui.print_line("Enter a world seed (or press Enter for random):")
        raw = self.ui.get_input("Seed: ")
        if raw and raw not in ("quit", "q"):
            try:
                return int(raw)
            except ValueError:
                # Hash string seeds
                return hash(raw) % (2**32)
        return None

    def _create_character(self):
        self.ui.print_separator()
        self.ui.print_line("Choose your class:\n")
        class_names = list(CLASSES.keys())
        for i, name in enumerate(class_names, 1):
            cls = CLASSES[name]
            self.ui.print_line(f"  {i}. {name} - {cls['desc']}")
            self.ui.print_line(f"     HP:{cls['base_hp']} ATK:{cls['base_atk']} DEF:{cls['base_def']} "
                              f"SPD:{cls['base_spd']} MAG:{cls['base_mag']}")
            self.ui.print_line(f"     Skills: {', '.join(cls['skills'])}")
            self.ui.print_line("")

        choice = self.ui.get_choice(len(class_names), "Class: ")
        if choice == -1:
            raise SystemExit
        chosen_class = class_names[choice - 1]

        self.ui.print_line(f"\nYou chose: {chosen_class}!")
        self.ui.print_line("What is your name, adventurer?")
        name = self.ui.get_input("Name: ")
        if not name or name in ("quit", "q"):
            name = "Hero"

        self.player = Player(name, chosen_class)
        self.ui.print_line(f"\nWelcome, {self.player.name} the {chosen_class}!")
        self.ui.print_line(self.player.full_status())

    def _main_loop(self):
        """Primary game loop: explore the world."""
        while not self.game_over:
            self.turn += 1
            region = self.world[self.current_region_idx]

            if not region["visited"]:
                region["visited"] = True
                self.player.regions_visited.add(region["name"])
                self.ui.print_line("")
                self.ui.print_banner(f"Entering: {region['name']}")
                self.ui.print_line(f"  {region['biome']['desc']}")
                self.ui.print_line("")

            self._region_menu(region)

            if not self.player.alive:
                self._game_over_screen()
                break

    def _region_menu(self, region):
        """Show options available in the current region."""
        self.ui.print_separator()
        self.ui.print_line(f"  Location: {region['name']}  |  Turn: {self.turn}")
        self.ui.print_line(f"  HP: {self.player.hp}/{self.player.max_hp}  "
                          f"MP: {self.player.mp}/{self.player.max_mp}  "
                          f"Gold: {self.player.gold}  Lv.{self.player.level}")
        self.ui.print_separator()

        options = []
        options.append(("Explore the wilds", "explore"))
        if region["npcs"]:
            options.append(("Talk to NPCs", "npc"))
        if region["dungeon"] and not region["dungeon"]["cleared"]:
            options.append((f"Enter dungeon: {region['dungeon']['name']}", "dungeon"))
        if region["shop"]:
            options.append(("Visit shop", "shop"))
        options.append(("Travel to another region", "travel"))
        options.append(("View status", "status"))
        options.append(("View inventory", "inventory"))
        options.append(("View quests", "quests"))
        options.append(("Save & Quit", "quit"))

        self.ui.print_line("\nWhat do you do?")
        for i, (label, _) in enumerate(options, 1):
            self.ui.print_line(f"  {i}. {label}")

        choice = self.ui.get_choice(len(options))
        if choice == -1:
            self.game_over = True
            return

        action = options[choice - 1][1]

        if action == "explore":
            self._explore(region)
        elif action == "npc":
            self._talk_to_npcs(region)
        elif action == "dungeon":
            self._enter_dungeon(region)
        elif action == "shop":
            self._visit_shop(region)
        elif action == "travel":
            self._travel()
        elif action == "status":
            self.ui.print_line("")
            self.ui.print_line(self.player.full_status())
        elif action == "inventory":
            self._show_inventory()
        elif action == "quests":
            self._show_quests()
        elif action == "quit":
            self.game_over = True

    def _explore(self, region):
        """Wander the wilds - random encounters and events."""
        self.ui.print_line(f"\n  You venture deeper into {region['name']}...")
        biome = region["biome"]

        # Random event
        event_result = trigger_random_event(self.player, self.rng, self.ui, self.player.level)

        # Enemy encounter
        if event_result == "ambush" or self.rng.random() < 0.5:
            enemy_name = self.rng.choice(biome["enemies"])
            enemy = Enemy(enemy_name, self.player.level, biome["danger"], self.rng)
            result = run_combat(self.player, enemy, self.rng, self.ui)
            if result is True:
                self._check_quest_progress("kill", enemy_name)
                # Loot drop
                if self.rng.random() < 0.35:
                    loot = generate_loot(self.rng, self.player.level, biome["danger"])
                    self.ui.print_line(f"  The {enemy_name} dropped: {loot['name']}!")
                    self._handle_loot(loot)
            elif result is False:
                return  # Player died
        else:
            # Gathering
            if self.rng.random() < 0.4:
                resource = self.rng.choice(biome["resources"])
                self.ui.print_line(f"  You found some {resource}!")
                self.player.inventory.append({"name": resource, "type": "resource", "desc": f"A resource: {resource}"})
                self._check_quest_progress("gather", resource)
            else:
                flavor = [
                    "You wander but find nothing of note.",
                    "The path is quiet. Too quiet.",
                    "You take in the scenery and rest briefly.",
                    "A gentle breeze carries the scent of distant flowers.",
                    "You hear distant sounds but cannot locate their source.",
                ]
                self.ui.print_line(f"  {self.rng.choice(flavor)}")
                # Small HP/MP regen from resting
                self.player.heal(5)
                self.player.restore_mp(3)

    def _talk_to_npcs(self, region):
        """Interact with NPCs in the region."""
        self.ui.print_line("\nWho do you want to talk to?")
        for i, npc in enumerate(region["npcs"], 1):
            self.ui.print_line(f"  {i}. {npc['name']}")
        self.ui.print_line(f"  {len(region['npcs']) + 1}. Go back")

        choice = self.ui.get_choice(len(region["npcs"]) + 1)
        if choice == -1 or choice == len(region["npcs"]) + 1:
            return

        npc = region["npcs"][choice - 1]
        self.ui.print_line(f"\n  {npc['name']}: \"{npc['greeting']}\"")

        if npc["quest"] and not npc["quest"]["completed"]:
            quest = npc["quest"]
            if quest not in self.player.quests and quest not in self.player.completed_quests:
                self.ui.print_line(f"\n  Quest available: {quest['name']}")
                self.ui.print_line(f"  {quest['desc']}")
                self.ui.print_line(f"  Reward: {quest['gold_reward']} gold, {quest['xp_reward']} XP")
                if self.ui.confirm("  Accept this quest?"):
                    self.player.quests.append(quest)
                    self.ui.print_line("  Quest accepted!")
            elif quest in self.player.quests and quest["progress"] >= quest["goal"]:
                self.ui.print_line(f"\n  {npc['name']}: \"You've done it! Here's your reward.\"")
                quest["completed"] = True
                self.player.quests.remove(quest)
                self.player.completed_quests.append(quest)
                self.player.gold += quest["gold_reward"]
                msgs = self.player.xp_gain(quest["xp_reward"])
                self.ui.print_line(f"  Received {quest['gold_reward']} gold and {quest['xp_reward']} XP!")
                for m in msgs:
                    self.ui.print_line(m)

    def _enter_dungeon(self, region):
        """Dungeon crawling through procedural rooms."""
        dungeon = region["dungeon"]
        self.ui.print_banner(f"Dungeon: {dungeon['name']}")
        self.ui.print_line(f"  {dungeon['theme_desc']}")
        self.ui.print_line(f"  Rooms: {len(dungeon['rooms'])}")
        self.ui.print_line("")

        for room in dungeon["rooms"]:
            if room["cleared"]:
                continue

            self.ui.print_separator()
            self.ui.print_line(f"  {room['desc']}")

            # Room enemy
            if room["enemy"]:
                if room["is_boss"]:
                    self.ui.print_line(f"\n  *** BOSS: {dungeon['boss_name']} blocks your path! ***")
                    enemy = Boss(dungeon["boss_name"], self.player.level, dungeon["danger"], self.rng)
                else:
                    self.ui.print_line(f"\n  A {room['enemy']} lurks here!")
                    enemy = Enemy(room["enemy"], self.player.level, dungeon["danger"], self.rng)

                result = run_combat(self.player, enemy, self.rng, self.ui)
                if result is True:
                    self._check_quest_progress("kill", room["enemy"])
                elif result is False:
                    return  # Dead
                elif result is None:
                    self.ui.print_line("  You retreat from the dungeon.")
                    return

            # Room loot
            if room["loot"]:
                self.ui.print_line(f"\n  You find: {room['loot']['name']}!")
                self._handle_loot(room["loot"])

            room["cleared"] = True

            if not room["is_boss"]:
                self.ui.print_line("\n  1. Continue deeper")
                self.ui.print_line("  2. Leave the dungeon")
                c = self.ui.get_choice(2)
                if c == 2 or c == -1:
                    self.ui.print_line("  You leave the dungeon.")
                    return

        dungeon["cleared"] = True
        self.ui.print_line("\n  *** Dungeon cleared! ***")
        bonus_gold = 30 + self.player.level * 10
        self.player.gold += bonus_gold
        msgs = self.player.xp_gain(50 + self.player.level * 20)
        self.ui.print_line(f"  Dungeon completion bonus: {bonus_gold} gold, {50 + self.player.level * 20} XP!")
        for m in msgs:
            self.ui.print_line(m)
        self._check_quest_progress("explore", dungeon["name"])

    def _visit_shop(self, region):
        """Buy items from a shop."""
        shop = region["shop"]
        self.ui.print_line(f"\n  Welcome to the shop! You have {self.player.gold} gold.\n")

        while True:
            for i, item in enumerate(shop["items"], 1):
                price = item["price"]
                desc = item.get("desc", "")
                self.ui.print_line(f"  {i}. {item['name']} - {price}g  {desc}")
            self.ui.print_line(f"  {len(shop['items']) + 1}. Sell items")
            self.ui.print_line(f"  {len(shop['items']) + 2}. Leave shop")

            choice = self.ui.get_choice(len(shop["items"]) + 2)
            if choice == -1 or choice == len(shop["items"]) + 2:
                break

            if choice == len(shop["items"]) + 1:
                self._sell_items()
                continue

            item = shop["items"][choice - 1]
            if self.player.gold >= item["price"]:
                self.player.gold -= item["price"]
                bought = dict(item)
                del bought["price"]
                if bought.get("slot") == "weapon":
                    if self.ui.confirm(f"  Equip {bought['name']}?"):
                        old = self.player.equip_weapon(bought)
                        if old:
                            self.ui.print_line(f"  Unequipped {old['name']}.")
                    else:
                        self.player.inventory.append(bought)
                elif bought.get("slot") == "armor":
                    if self.ui.confirm(f"  Equip {bought['name']}?"):
                        old = self.player.equip_armor(bought)
                        if old:
                            self.ui.print_line(f"  Unequipped {old['name']}.")
                    else:
                        self.player.inventory.append(bought)
                else:
                    self.player.inventory.append(bought)
                self.ui.print_line(f"  Bought {bought['name']}! Gold remaining: {self.player.gold}")
            else:
                self.ui.print_line("  Not enough gold!")
            self.ui.print_line("")

    def _sell_items(self):
        """Sell items from inventory."""
        sellable = [(i, item) for i, item in enumerate(self.player.inventory)
                    if isinstance(item, dict)]
        if not sellable:
            self.ui.print_line("  Nothing to sell!")
            return

        self.ui.print_line(f"\n  Your gold: {self.player.gold}")
        for idx, (inv_idx, item) in enumerate(sellable, 1):
            value = item.get("value", 5)
            self.ui.print_line(f"  {idx}. {item['name']} - sells for {value}g")
        self.ui.print_line(f"  {len(sellable) + 1}. Cancel")

        choice = self.ui.get_choice(len(sellable) + 1)
        if choice == -1 or choice == len(sellable) + 1:
            return
        _, item = sellable[choice - 1]
        value = item.get("value", 5)
        self.player.inventory.remove(item)
        self.player.gold += value
        self.ui.print_line(f"  Sold {item['name']} for {value} gold!")

    def _travel(self):
        """Travel to a different region."""
        self.ui.print_line("\nWhere would you like to go?")
        for i, region in enumerate(self.world, 1):
            status = " (visited)" if region["visited"] else " (unexplored)"
            danger = region["biome"]["danger"]
            danger_str = "*" * danger
            self.ui.print_line(f"  {i}. {region['name']} [Danger: {danger_str}]{status}")

        self.ui.print_line(f"  {len(self.world) + 1}. Expand the frontier (discover new region)")
        self.ui.print_line(f"  {len(self.world) + 2}. Stay here")

        choice = self.ui.get_choice(len(self.world) + 2)
        if choice == -1 or choice == len(self.world) + 2:
            return

        if choice == len(self.world) + 1:
            new_region = generate_region(self.rng, self.player.level)
            self.world.append(new_region)
            self.current_region_idx = len(self.world) - 1
            self.ui.print_line(f"\n  You discover a new region: {new_region['name']}!")
        else:
            self.current_region_idx = choice - 1

    def _handle_loot(self, loot):
        """Prompt the player about a loot drop."""
        if loot.get("slot") == "weapon":
            current = self.player.weapon
            self.ui.print_line(f"  {loot['desc']}")
            if current:
                self.ui.print_line(f"  Current weapon: {current['name']} (ATK +{current['bonus']})")
            if self.ui.confirm(f"  Equip {loot['name']}?"):
                old = self.player.equip_weapon(loot)
                if old:
                    self.ui.print_line(f"  Unequipped {old['name']}.")
            else:
                self.player.inventory.append(loot)
                self.ui.print_line("  Added to inventory.")
        elif loot.get("slot") == "armor":
            current = self.player.armor
            self.ui.print_line(f"  {loot['desc']}")
            if current:
                self.ui.print_line(f"  Current armor: {current['name']} (DEF +{current['bonus']})")
            if self.ui.confirm(f"  Equip {loot['name']}?"):
                old = self.player.equip_armor(loot)
                if old:
                    self.ui.print_line(f"  Unequipped {old['name']}.")
            else:
                self.player.inventory.append(loot)
                self.ui.print_line("  Added to inventory.")
        else:
            self.player.inventory.append(loot)
            self.ui.print_line("  Added to inventory.")

    def _check_quest_progress(self, quest_type, target):
        """Update quest progress when relevant actions happen."""
        for quest in self.player.quests:
            if quest["completed"]:
                continue
            if quest["type"] == quest_type and quest["target"] == target:
                quest["progress"] += 1
                if quest["progress"] >= quest["goal"]:
                    self.ui.print_line(f"\n  *** Quest objective complete: {quest['name']}! ***")
                    self.ui.print_line("  Return to the quest giver to claim your reward.")
                else:
                    self.ui.print_line(f"  Quest progress: {quest['name']} ({quest['progress']}/{quest['goal']})")

    def _show_inventory(self):
        """Display the player's inventory."""
        self.ui.print_line(f"\n=== Inventory ({len(self.player.inventory)} items) ===")
        if not self.player.inventory:
            self.ui.print_line("  Your pack is empty.")
            return

        for i, item in enumerate(self.player.inventory, 1):
            if isinstance(item, dict):
                desc = item.get("desc", "")
                self.ui.print_line(f"  {i}. {item['name']} - {desc}")
            else:
                self.ui.print_line(f"  {i}. {item}")

        self.ui.print_line(f"\n  1. Equip/Use an item")
        self.ui.print_line(f"  2. Go back")
        choice = self.ui.get_choice(2)
        if choice == 1:
            self._use_inventory_item()

    def _use_inventory_item(self):
        """Use or equip an item from inventory."""
        if not self.player.inventory:
            return
        self.ui.print_line("Which item? (number)")
        for i, item in enumerate(self.player.inventory, 1):
            name = item["name"] if isinstance(item, dict) else str(item)
            self.ui.print_line(f"  {i}. {name}")
        self.ui.print_line(f"  {len(self.player.inventory) + 1}. Cancel")

        choice = self.ui.get_choice(len(self.player.inventory) + 1)
        if choice == -1 or choice == len(self.player.inventory) + 1:
            return

        item = self.player.inventory[choice - 1]
        if not isinstance(item, dict):
            self.ui.print_line("  This item can't be used directly.")
            return

        if item.get("slot") == "weapon":
            old = self.player.equip_weapon(item)
            self.player.inventory.remove(item)
            self.ui.print_line(f"  Equipped {item['name']}!")
        elif item.get("slot") == "armor":
            old = self.player.equip_armor(item)
            self.player.inventory.remove(item)
            self.ui.print_line(f"  Equipped {item['name']}!")
        elif item.get("consumable"):
            self.player.inventory.remove(item)
            if item["type"] == "heal":
                healed = self.player.heal(item["power"])
                self.ui.print_line(f"  Used {item['name']}! Restored {healed} HP.")
            elif item["type"] == "mana":
                restored = self.player.restore_mp(item["power"])
                self.ui.print_line(f"  Used {item['name']}! Restored {restored} MP.")
            else:
                self.ui.print_line(f"  Used {item['name']}!")
        else:
            self.ui.print_line("  This item can't be used right now.")

    def _show_quests(self):
        """Display active and completed quests."""
        self.ui.print_line("\n=== Active Quests ===")
        if not self.player.quests:
            self.ui.print_line("  No active quests.")
        for q in self.player.quests:
            status = f"({q['progress']}/{q['goal']})"
            self.ui.print_line(f"  - {q['name']} {status}")
            self.ui.print_line(f"    {q['desc']}")

        if self.player.completed_quests:
            self.ui.print_line(f"\n  Completed: {len(self.player.completed_quests)} quests")

    def _game_over_screen(self):
        """Display game over information."""
        self.ui.print_line("")
        self.ui.print_banner("GAME OVER")
        self.ui.print_line(f"  {self.player.name} the {self.player.char_class} has fallen.")
        self.ui.print_line(f"  Level: {self.player.level}")
        self.ui.print_line(f"  Turns survived: {self.turn}")
        self.ui.print_line(f"  Regions explored: {len(self.player.regions_visited)}")
        self.ui.print_line(f"  Quests completed: {len(self.player.completed_quests)}")
        self.ui.print_line(f"  Enemies defeated: {sum(self.player.kills.values())}")
        if self.player.kills:
            self.ui.print_line("  Kill log:")
            for name, count in sorted(self.player.kills.items(), key=lambda x: -x[1]):
                self.ui.print_line(f"    {name}: {count}")
        self.ui.print_line(f"\n  World seed: {self.rng.seed}")
        self.ui.print_line("  Thanks for playing Floral Realms!")
