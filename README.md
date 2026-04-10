# Ghost Hunter

A retro first-person shooter built from scratch in Go using the [Ebitengine](https://ebitengine.org) game engine. Inspired by classic 90s shooters like Doom and Wolfenstein 3D, Ghost Hunter features a custom software raycaster, wave-based combat, multiple enemy types, and a 5-level progression system.

---

## Features

### Engine
- Custom software raycaster written in Go
- Textured walls, floors, and ceilings with distance shading
- Billboard sprite rendering with depth buffering
- Screen shake and gun kick feedback
- Minimap with player direction indicator
- Automatic screen resolution detection

### Gameplay
- 5 unique maps with distinct themes and layouts
- Wave-based enemy spawning with increasing difficulty
- Portal system to travel between maps
- Health and ammo pickup system
- Persistent high score saved to disk

### Weapons (Doom-style state machine)
- **Pistol** — single shot, no cooldown
- **Shotgun** — 5-ray spread, 3 ammo per shot
- **Machinegun** — hold to fire, animated E11 blaster sprite

### Enemies
- **Ghost** (Map 1) — phases through walls, slow
- **Wizard** (Map 2) — medium speed, direct chase
- **Demon** (Map 3) — fast, zigzag movement
- **Wraith** (Map 4) — circles the player while closing in
- **Reaper** (Map 5) — teleports when far away, high health

### AI
- Line of sight detection — enemies only attack when they can see you
- Wall sliding — enemies navigate around obstacles
- State machine — Chase, Attack, Pain, Death states
- Doom-style random direction turning when blocked

---

## Controls

| Key | Action |
|-----|--------|
| Arrow Up / Down | Move forward / backward |
| Arrow Left / Right | Turn left / right |
| Space | Shoot |
| 1 / 2 / 3 | Switch weapon |
| ESC | Pause / Resume |
| R | Restart (game over screen) |
| Q | Quit |

---

## Getting Started

### Requirements
- Go 1.21 or higher
- Ebitengine dependencies (see below)

### macOS
```bash
xcode-select --install
```

### Linux
```bash
sudo apt install libc6-dev libgles2-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev
```

### Run
```bash
git clone https://github.com/chuma-beep/doom-go
cd doom-go
go run .
```

### Build
```bash
go build -o ghosthunter .
./ghosthunter
```

---

## Project Structure

```
doom-go/
├── main.go                  # Entry point, asset loading
├── maps/                    # JSON level files
│   ├── map1.json            # The Haunted Halls (16x16)
│   ├── map2.json            # The Wizard's Den (16x16)
│   ├── map3.json            # The Labyrinth (32x32)
│   ├── map4.json            # The Arena (32x32)
│   └── map5.json            # The Boss Chamber (32x32)
├── assets/                  # Textures, sprites, sounds
│   ├── gun_pistol/          # Pistol animation frames
│   ├── gun_machinegun/      # Machinegun animation frames
│   └── blaster/             # E11 blaster source frames
└── internal/engine/
    ├── game.go              # Game loop, player movement, entity AI
    ├── renderer.go          # Raycaster, sprite rendering, HUD
    ├── weapons.go           # Doom-style weapon state machine
    ├── sprite.go            # Entity and pickup types
    ├── texture.go           # Asset loading
    ├── maploader.go         # JSON map loader
    ├── world.go             # Map access functions
    ├── sound.go             # Audio playback
    └── save.go              # High score persistence
```

---

## Asset Credits

- **Enemy sprites** — [FPS Monster Enemies](https://opengameart.org/content/fps-monster-enemies) by Ragnar Random (CC0)
- **Machinegun sprite** — [FP Animated Weapons E11](https://whiteknightstudios.itch.io/fp-animated-weapons-e11) by W_K_Studio (CC0)

---

## License

MIT License — free to use, modify, and distribute.
