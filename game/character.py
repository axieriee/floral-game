"""Player character and enemy entities."""

from game.data import CLASSES, SKILLS


class Character:
    """Base class for any combat-capable entity."""

    def __init__(self, name, hp, atk, defense, spd, mag, level=1):
        self.name = name
        self.max_hp = hp
        self.hp = hp
        self.max_mp = 50
        self.mp = 50
        self.atk = atk
        self.defense = defense
        self.spd = spd
        self.mag = mag
        self.level = level
        self.skills = []
        self.buffs = []  # list of (stat, multiplier, turns_remaining)
        self.dots = []   # list of (damage_per_turn, turns_remaining)
        self.alive = True

    def take_damage(self, amount):
        actual = max(1, amount - self.effective_stat("defense") // 3)
        self.hp = max(0, self.hp - actual)
        if self.hp == 0:
            self.alive = False
        return actual

    def heal(self, amount):
        old_hp = self.hp
        self.hp = min(self.max_hp, self.hp + amount)
        return self.hp - old_hp

    def restore_mp(self, amount):
        old_mp = self.mp
        self.mp = min(self.max_mp, self.mp + amount)
        return self.mp - old_mp

    def effective_stat(self, stat):
        base = getattr(self, stat)
        for buff_stat, mult, _ in self.buffs:
            if buff_stat == stat:
                base = int(base * mult)
        return base

    def tick_buffs(self):
        """Decrease buff durations; remove expired ones."""
        self.buffs = [(s, m, t - 1) for s, m, t in self.buffs if t > 1]

    def tick_dots(self):
        """Apply damage-over-time effects. Returns total dot damage dealt."""
        total = 0
        for dmg, _ in self.dots:
            self.hp = max(0, self.hp - dmg)
            total += dmg
        if self.hp == 0:
            self.alive = False
        self.dots = [(d, t - 1) for d, t in self.dots if t > 1]
        return total

    def stat_block(self):
        return (
            f"  HP: {self.hp}/{self.max_hp}  |  MP: {self.mp}/{self.max_mp}\n"
            f"  ATK: {self.effective_stat('atk')}  DEF: {self.effective_stat('defense')}  "
            f"SPD: {self.effective_stat('spd')}  MAG: {self.effective_stat('mag')}"
        )


class Player(Character):
    """The player character with inventory, XP, and equipment."""

    def __init__(self, name, char_class):
        cls = CLASSES[char_class]
        super().__init__(
            name=name,
            hp=cls["base_hp"],
            atk=cls["base_atk"],
            defense=cls["base_def"],
            spd=cls["base_spd"],
            mag=cls["base_mag"],
        )
        self.char_class = char_class
        self.skills = list(cls["skills"])
        self.xp = 0
        self.xp_to_level = 100
        self.gold = 25
        self.inventory = []
        self.weapon = None
        self.armor = None
        self.quests = []
        self.completed_quests = []
        self.regions_visited = set()
        self.kills = {}

    def xp_gain(self, amount):
        """Add XP and handle level-ups. Returns list of level-up messages."""
        self.xp += amount
        messages = []
        while self.xp >= self.xp_to_level:
            self.xp -= self.xp_to_level
            self.level += 1
            self.xp_to_level = int(self.xp_to_level * 1.4)
            # Stat growth
            self.max_hp += 10
            self.hp = self.max_hp
            self.max_mp += 5
            self.mp = self.max_mp
            self.atk += 2
            self.defense += 2
            self.spd += 1
            self.mag += 2
            messages.append(f"*** LEVEL UP! You are now level {self.level}! ***")
            messages.append(f"    HP+10  MP+5  ATK+2  DEF+2  SPD+1  MAG+2")
        return messages

    def equip_weapon(self, weapon):
        old = self.weapon
        if old:
            self.atk -= old["bonus"]
            self.inventory.append(old)
        self.weapon = weapon
        self.atk += weapon["bonus"]
        return old

    def equip_armor(self, armor):
        old = self.armor
        if old:
            self.defense -= old["bonus"]
            self.inventory.append(old)
        self.armor = armor
        self.defense += armor["bonus"]
        return old

    def full_status(self):
        lines = [
            f"=== {self.name} the {self.char_class} (Level {self.level}) ===",
            self.stat_block(),
            f"  XP: {self.xp}/{self.xp_to_level}  |  Gold: {self.gold}",
            f"  Weapon: {self.weapon['name'] if self.weapon else 'None'}",
            f"  Armor:  {self.armor['name'] if self.armor else 'None'}",
        ]
        return "\n".join(lines)

    def record_kill(self, enemy_name):
        self.kills[enemy_name] = self.kills.get(enemy_name, 0) + 1


class Enemy(Character):
    """A procedurally generated enemy."""

    def __init__(self, name, level, danger, rng):
        scale = 1 + (level - 1) * 0.3 + danger * 0.2
        hp = int((40 + rng.randint(0, 20)) * scale)
        atk = int((8 + rng.randint(0, 4)) * scale)
        defense = int((4 + rng.randint(0, 3)) * scale)
        spd = int((5 + rng.randint(0, 3)) * scale)
        mag = int((3 + rng.randint(0, 3)) * scale)
        super().__init__(name, hp, atk, defense, spd, mag, level)
        self.xp_reward = int(20 * scale)
        self.gold_reward = rng.randint(int(5 * scale), int(15 * scale))


class Boss(Enemy):
    """A more powerful enemy that guards dungeon ends."""

    def __init__(self, name, level, danger, rng):
        super().__init__(name, level, danger, rng)
        self.max_hp = int(self.max_hp * 2.5)
        self.hp = self.max_hp
        self.atk = int(self.atk * 1.5)
        self.defense = int(self.defense * 1.3)
        self.xp_reward = int(self.xp_reward * 3)
        self.gold_reward = int(self.gold_reward * 3)
        self.is_boss = True
