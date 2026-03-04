package main

import (
    "log"
    "doom-go/internal/engine"
    "github.com/hajimehoshi/ebiten/v2"
)

func main() {
    engine.LoadTexture("assets/wall.png")
    engine.LoadSpritTexture("assets/sprite_real.png")
    ebiten.SetWindowSize(engine.ScreenWidth*2, engine.ScreenHeight*2)
    ebiten.SetWindowTitle("doom-go")

    game := engine.NewGame()

    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}


