# Ghost Hunter

A retro-style 2.5D first-person shooter inspired by classic 90s games like DOOM. Built entirely in Go using the Ebitengine game engine.

![Gameplay](assets/screenshot.png)

## Features

- **Raycasting Engine**: Custom 2.5D raycasting renderer for walls, floors, and ceilings
- **Wave-based Survival**: Survive escalating waves of enemies across 5 unique maps
- **5 Enemy Types**: Ghost, Wizard, Demon, Wraith, and Reaper - each with unique AI behaviors
- **3 Weapons**: Pistol, Shotgun, and Machinegun with state-machine animations
- **Pickup System**: Health packs and ammo pickups scattered throughout levels
- **High Score Tracking**: Persisted to file for competitive play
- **Full Audio**: Sound effects and background music

## Controls

| Key | Action |
|-----|--------|
| Arrow Up/Down | Move forward/backward |
| Arrow Left/Right | Turn |
| Space | Fire weapon |
| 1, 2, 3 | Switch weapons |
| ESC | Pause/Resume |
| R | Restart (game over) |
| Q | Quit |

## Requirements

- Go 1.25+
- Sound card for audio (optional)

## Installation

```bash
git clone https://github.com/yourusername/doom-go.git
cd doom-go
go run .
```

Or run the pre-built binary:
```bash
./doom-go
```

## Project Structure

```
doom-go/
├── main.go              # Entry point
├── internal/engine/     # Game engine modules
│   ├── game.go          # Main game logic
│   ├── renderer.go      # Raycasting renderer & HUD
│   ├── weapons.go       # Weapon system
│   ├── world.go         # Map data access
│   ├── maploader.go     # Map file loader
│   └── ...
├── maps/                # JSON level definitions
└── assets/              # Sprites, textures, sounds
```

## License

MIT
