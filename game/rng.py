"""Seeded random number generation for reproducible procedural content."""

import random


class GameRNG:
    """Wrapper around random.Random for deterministic procedural generation."""

    def __init__(self, seed=None):
        self.seed = seed if seed is not None else random.randint(0, 2**32 - 1)
        self._rng = random.Random(self.seed)

    def randint(self, a, b):
        return self._rng.randint(a, b)

    def choice(self, seq):
        return self._rng.choice(seq)

    def choices(self, population, k=1):
        return self._rng.choices(population, k=k)

    def shuffle(self, seq):
        self._rng.shuffle(seq)
        return seq

    def random(self):
        return self._rng.random()

    def weighted_choice(self, options_weights):
        """Choose from list of (option, weight) tuples."""
        total = sum(w for _, w in options_weights)
        r = self._rng.random() * total
        cumulative = 0
        for option, weight in options_weights:
            cumulative += weight
            if r <= cumulative:
                return option
        return options_weights[-1][0]

    def sample(self, population, k):
        return self._rng.sample(population, k)
