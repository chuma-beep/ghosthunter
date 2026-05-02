# Ghost Hunter

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://ebitengine.org/)
[![WASM](https://img.shields.io/badge/WASM-Supported-green)](https://github.com/golang/go/tree/master/lib/wasm)

A retro first-person shooter built from scratch in Go using the [Ebitengine](https://ebitengine.org) game engine. Inspired by classic 90s shooters like Doom and Wolfenstein 3D, Ghost Hunter features a custom software raycaster, wave-based combat, multiple enemy types, and a 5-level progression system.

> This project was built as a learning exercise to explore and understand Go, game development, WebAssembly, and machine learning integration.

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

**Enemy AI**
- Line of sight detection — enemies only attack when they can see you
- Wall sliding — enemies navigate around obstacles
- State machine — Chase, Attack, Pain, Death states
- Doom-style random direction turning when blocked

**Neural Network AI**
- End-to-end ML pipeline: collect training data → train model → run inference in Go
- 3-layer neural network with ReLU activations and sigmoid output
- Press **N** to toggle neural network AI during gameplay
- Model learns from player demonstrations via imitation learning

---

## Controls

| Key | Action |
|-----|--------|
| Arrow Up / Down | Move forward / backward |
| Arrow Left / Right | Turn left / right |
| Space | Shoot |
| 1 / 2 / 3 | Switch weapon |
| A | Toggle rule-based AI |
| N | Toggle neural network AI |
| D | Start/stop data collection |
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

### Windows
Install [MSYS2](https://www.msys2.org/), then run:
```bash
pacman -S mingw-w64-x86_64-gtk3 mingw-w64-x86_64-libvorbis
```

### Run
```bash
git clone https://github.com/chuma-beep/ghosthunter
cd ghosthunter
go run .
```

### Build
```bash
go build -o ghosthunter .
./ghosthunter
```

### Run in Browser (WASM)
```bash
# Build WASM
GOOS=js GOARCH=wasm go build -o ghosthunter.wasm .

# Copy wasm_exec.js
cp $(go env GOROOT)/lib/wasm/wasm_exec.js .

# Serve locally
go run github.com/hajimehoshi/wasmserve@latest .
```
Then open http://localhost:8080 in your browser.

---

## Project Structure

```
gosthunter/
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
├── ml/                      # Machine learning scripts
│   ├── train_numpy.py       # Train neural network (numpy only)
│   ├── train.py             # Train with PyTorch
│   ├── training_data.json   # Collected gameplay samples
│   └── model_weights.json   # Trained model (22 inputs → 5 outputs)
└── internal/engine/
    ├── game.go              # Game loop, player movement, entity AI
    ├── renderer.go          # Raycaster, sprite rendering, HUD
    ├── weapons.go           # Doom-style weapon state machine
    ├── sprite.go            # Entity and pickup types
    ├── texture.go           # Asset loading
    ├── maploader.go         # JSON map loader
    ├── world.go             # Map access functions
    ├── sound.go             # Audio playback
    ├── save.go              # High score persistence
    ├── ai.go                # Rule-based player AI
    └── neural_ai.go         # Neural network inference in Go
```

---

## Training Your Own AI

The neural network learns from your gameplay via imitation learning.

### Collect Training Data

1. Run the game: `go run .`
2. Play normally — move around, shoot enemies, collect pickups
3. Press **D** to start recording your actions
4. Play for a few minutes to gather diverse examples
5. Press **D** again to stop — data saves to `ml/training_data.json`

### Train the Model

```bash
cd ml
python train_numpy.py
```

This trains a 3-layer neural network (22 inputs → 64 hidden → 5 outputs) and exports weights to `model_weights.json`.

### Model Architecture

- **Input (22 features)**: player position, angle, health, ammo, weapon, enemy count, wave, map, portal distance, pickup distances, enemy distances/angles
- **Output (5 actions)**: move forward, move backward, turn left, turn right, shoot
- **Hidden layers**: 64 neurons each with ReLU activation
- **Output layer**: Sigmoid for multi-label classification

The trained model is embedded into the Go binary — inference runs entirely on CPU with no external dependencies.

---

## Asset Credits

- **Enemy sprites** — [FPS Monster Enemies](https://opengameart.org/content/fps-monster-enemies) by Ragnar Random (CC0)
- **Machinegun sprite** — [FP Animated Weapons E11](https://whiteknightstudios.itch.io/fp-animated-weapons-e11) by W_K_Studio (CC0)

---

## License

MIT License — free to use, modify, and distribute.
