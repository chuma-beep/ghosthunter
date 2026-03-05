package engine

import (
    "math"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
    Pixels  []byte
    PlayerX float64
    PlayerY float64
    Angle   float64
    Sprites []Sprite
    Score     int
    RespawnTimer int
    Health      int
		DamageFlash int
		Wave     int 
}

func NewGame() *Game {
    return &Game{
        Pixels:  make([]byte, ScreenWidth*ScreenHeight*4),
        PlayerX: 8.0,
        PlayerY: 8.0,
        Angle:   0.0,
				Wave:    1,
		Health: 100,
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

for _, sprite := range g.Sprites {
	dx := sprite.X - g.PlayerX 
	dy := sprite.Y - g.PlayerY
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 0.8{
	    g.Health -= 1
      g.DamageFlash = 10
	}
}


if g.Health <= 0{
	if ebiten.IsKeyPressed(ebiten.KeyR){
		g.Health = 100
		g.Score = 0
		g.Sprites =  []Sprite{
			{X: 6.0, Y: 6.0},
			{X: 10.0, Y: 4.0},
			{X: 3.0, Y: 12.0},
		}
	}
   return nil
}


// shooting
if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
    PlaySound("assets/shoot.wav")
    for i := len(g.Sprites) - 1; i >= 0; i-- {
        dx := g.Sprites[i].X - g.PlayerX
        dy := g.Sprites[i].Y - g.PlayerY
        dist := math.Sqrt(dx*dx + dy*dy)
        spriteAngle := math.Atan2(dy, dx) - g.Angle
        for spriteAngle > math.Pi { spriteAngle -= 2 * math.Pi }
        for spriteAngle < -math.Pi { spriteAngle += 2 * math.Pi }
        if math.Abs(spriteAngle) < 0.2 && dist < 10 {
            g.Sprites = append(g.Sprites[:i], g.Sprites[i+1:]...)
            g.Score++
						PlaySound("assets/ghost.wav")
        }
    }
}



  // respawn ghosts when all are dead
if len(g.Sprites) == 0 {
    g.RespawnTimer++
    if g.RespawnTimer > 180 {
        g.Wave++
        count := 3 + g.Wave
        speed := 0.005 + float64(g.Wave)*0.002
        positions := [][2]float64{
            {6.0, 6.0},
            {10.0, 4.0},
            {3.0, 12.0},
            {12.0, 12.0},
            {8.0, 3.0},
            {2.0, 8.0},
            {13.0, 7.0},
        }
        g.Sprites = []Sprite{}
        for i := 0; i < count && i < len(positions); i++ {
            g.Sprites = append(g.Sprites, Sprite{
                X: positions[i][0],
                Y: positions[i][1],
            })
        }
        _ = speed
        g.RespawnTimer = 0
    }
}
 if g.Health <= 0 {
	 g.Health = 0 
	 g.Sprites = []Sprite{}
	  }
 


    return nil
}




func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ScreenWidth, ScreenHeight
}
