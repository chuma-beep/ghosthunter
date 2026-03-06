package main

import (
    "log"
    "doom-go/internal/engine"
    "github.com/hajimehoshi/ebiten/v2"
)

func safeInitAudio() {
    defer func() { recover() }()
    engine.InitAudio()
    engine.PlayMusic("assets/music.mp3")
}

func main() {
    engine.LoadTexture("assets/wall.png")
    engine.LoadTexture2("assets/wall2.png")
	engine.LoadSpriteTexture("assets/sprite_real.png")
	engine.LoadWizard("assets/wizard.png")
    engine.LoadGun("assets/gun.png")
    engine.LoadMap("maps/map1.json", 0)
	engine.LoadMap("maps/map1.json", 1)
	safeInitAudio()
    engine.LoadFloor("assets/floor.png", &engine.FloorTexture)
    engine.LoadFloor("assets/floor2.png", &engine.FloorTexture2)
    ebiten.SetWindowSize(engine.ScreenWidth*2, engine.ScreenHeight*2)
    ebiten.SetWindowTitle("doom-go")
    game := engine.NewGame()
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
