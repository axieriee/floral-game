# Floravale — A Botanical Text RPG

A cozy, procedural text RPG built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Press flowers into your journal, blend their resonances, and befriend the locals.

## Play

```bash
go build -o floragame . && ./floragame
```

## Features

- **Floral Magic System** — Forage flowers with unique resonances (Clarity, Venom, Charm, etc.) and press them into your journal
- **Blending** — Combine two pressed flowers to discover compound effects. 10 blends to find
- **4 Traditions** — Hedge Keeper, Thornwalker, Petal Scholar, Rootborn — each changes how the game plays
- **NPC Interactions** — Maren the Beekeeper reacts differently to every flower you offer. Her trust is visible and shifts based on your choices
- **Journal Mechanic** — Your journal is your inventory, spell book, and relationship tracker
- **Soft Pastel Terminal UI** — Lavender borders, rose highlights, sage accents — powered by Lip Gloss
- **Procedural** — Weather, forage results, and flower potency shift each playthrough

## Controls

- **Arrow keys / j/k** — Navigate menus
- **Enter** — Select
- **Esc / q** — Go back
- **Ctrl+C** — Quit

## Flower Rarity

| Rarity   | Symbol | Examples                          |
|----------|--------|-----------------------------------|
| Common   | ⚘      | Daisy, Clover, Dandelion, Chamomile |
| Uncommon | ❀      | Lavender, Thistle, Foxglove, Honeysuckle |
| Rare     | ✾      | Moonflower, Ghost Orchid, Bloodroot |
| Mythic   | ❁      | Everbloom                          |
