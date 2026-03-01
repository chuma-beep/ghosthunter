package engine

import (
    "math"
    "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    Pixels  []byte
    PlayerX float64
    PlayerY float64
    Angle   float64
    SpriteX float64
    SpriteY float64
}

func NewGame() *Game {
    return &Game{
        Pixels:  make([]byte, ScreenWidth*ScreenHeight*4),
        PlayerX: 8.0,
        PlayerY: 8.0,
        Angle:   0.0,
        SpriteX: 6.0,
        SpriteY: 6.0,
    }
}

func (g *Game) Update() error {
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        g.Angle -= 0.03
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        g.Angle += 0.03
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
        newX := g.PlayerX + math.Cos(g.Angle)*0.05
        newY := g.PlayerY + math.Sin(g.Angle)*0.05
        if WorldMap[int(newY)][int(newX)] == 0 {
            g.PlayerX = newX
            g.PlayerY = newY
        }
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
        newX := g.PlayerX - math.Cos(g.Angle)*0.05
        newY := g.PlayerY - math.Sin(g.Angle)*0.05
        if WorldMap[int(newY)][int(newX)] == 0 {
            g.PlayerX = newX
            g.PlayerY = newY
        }
    }
    return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ScreenWidth, ScreenHeight
}
