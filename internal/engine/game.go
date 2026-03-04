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
    Sprites []Sprite
}

func NewGame() *Game {
    return &Game{
        Pixels:  make([]byte, ScreenWidth*ScreenHeight*4),
        PlayerX: 8.0,
        PlayerY: 8.0,
        Angle:   0.0,
        Sprites: []Sprite{
         {X: 6.0, Y: 6.0, VX: 0.0, VY: 0.0},
         {X: 10.0, Y: 4.0, VX: 0.0, VY: 0.0},
         {X: 3.0, Y: 12.0, VX: 0.0, VY: 0.0},
       },
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


// move sprites toward player
for i := range g.Sprites {
    dx := g.PlayerX - g.Sprites[i].X
    dy := g.PlayerY - g.Sprites[i].Y
    dist := math.Sqrt(dx*dx + dy*dy)
    if dist > 0.5 {
        g.Sprites[i].X += (dx / dist) * 0.005
        g.Sprites[i].Y += (dy / dist) * 0.005
    }
}


    return nil
}


// shooting
if ebiten.IsKeyPressed(ebiten.KeySpace) {
    for i := len(g.Sprites) - 1; i >= 0; i-- {
        dx := g.Sprites[i].X - g.PlayerX
        dy := g.Sprites[i].Y - g.PlayerY
        dist := math.Sqrt(dx*dx + dy*dy)

        spriteAngle := math.Atan2(dy, dx) - g.Angle
        for spriteAngle > math.Pi { spriteAngle -= 2 * math.Pi }
        for spriteAngle < -math.Pi { spriteAngle += 2 * math.Pi }

        if math.Abs(spriteAngle) < 0.2 && dist < 10 {
            g.Sprites = append(g.Sprites[:i], g.Sprites[i+1:]...)
        }
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ScreenWidth, ScreenHeight
}
