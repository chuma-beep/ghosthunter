package engine

import (
	"fmt"
    "math"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
    Pixels       []byte
    PlayerX      float64
    PlayerY      float64
    Angle        float64
    Sprites      []Sprite
    Score        int
    RespawnTimer int
    Health       int
    DamageFlash  int
    Wave         int
    GunKick      int
    Ammo         int
    AmmoPickups  []AmmoPickup
    GameState    int
    HighScore    int
    CurrentMap   int 
}

func NewGame() *Game {
    return &Game{
        Pixels:  make([]byte, ScreenWidth*ScreenHeight*4),
        CurrentMap: 0,
		PlayerX: 8.0,
        PlayerY: 8.0,
        Angle:   0.0,
        Wave:    1,
        Ammo:    10,
        GameState: 0,
        Health:  100,
        Sprites: []Sprite{
            {X: 6.0, Y: 6.0},
            {X: 10.0, Y: 4.0},
            {X: 3.0, Y: 12.0},
        },
        AmmoPickups: []AmmoPickup{
            {X: 5.0, Y: 5.0, Active: true},
            {X: 11.0, Y: 11.0, Active: true},
            {X: 3.0, Y: 9.0, Active: true},
        },
    }
}

func (g *Game) Update() error {
    // start screen
    if g.GameState == 0 {
        if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
            g.GameState = 1
        }
        return nil
    }



// portal detection
portalX, portalY := 13.0, 1.0
pdx := g.PlayerX - portalX
pdy := g.PlayerY - portalY
portalDist := math.Sqrt(pdx*pdx + pdy*pdy)
if portalDist < 0.8 {
    fmt.Println("PORTAL TRIGGERED", g.CurrentMap)
    if g.CurrentMap == 0 {
        g.CurrentMap = 1
    } else {
        g.CurrentMap = 0
    }
    g.PlayerX = 2.0
    g.PlayerY = 2.0
    g.Sprites = []Sprite{
        {X: 6.0, Y: 6.0},
        {X: 10.0, Y: 4.0},
        {X: 3.0, Y: 12.0},
    }
}



    // game over screen
    if g.GameState == 2 {
        if inpututil.IsKeyJustPressed(ebiten.KeyR) {
            g.Health = 100
            g.Score = 0
            g.Ammo = 10
            g.Wave = 1
            g.GameState = 1
            g.Sprites = []Sprite{
                {X: 6.0, Y: 6.0},
                {X: 10.0, Y: 4.0},
                {X: 3.0, Y: 12.0},
            }
        }
        return nil
    }

    // player movement
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        g.Angle -= 0.03
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        g.Angle += 0.03
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
        newX := g.PlayerX + math.Cos(g.Angle)*0.05
        newY := g.PlayerY + math.Sin(g.Angle)*0.05
        if GetMap(g.CurrentMap)[int(newY)][int(newX)] == 0 {
            g.PlayerX = newX
            g.PlayerY = newY
        }
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
        newX := g.PlayerX - math.Cos(g.Angle)*0.05
        newY := g.PlayerY - math.Sin(g.Angle)*0.05
        if GetMap(g.CurrentMap)[int(newY)][int(newX)] == 0 {
            g.PlayerX = newX
            g.PlayerY = newY
        }
    }

    // gun kick decay
    if g.GunKick > 0 {
        g.GunKick--
    }

    // move sprites toward player
    speed := 0.005 + float64(g.Wave)*0.002
    for i := range g.Sprites {
        if g.Sprites[i].Dead {
            if g.Sprites[i].FadeTimer > 0 {
                g.Sprites[i].FadeTimer--
            }
            continue
        }
        dx := g.PlayerX - g.Sprites[i].X
        dy := g.PlayerY - g.Sprites[i].Y
        dist := math.Sqrt(dx*dx + dy*dy)
        if dist > 0.5 {
            g.Sprites[i].X += (dx / dist) * speed
            g.Sprites[i].Y += (dy / dist) * speed
        }
    }

    // damage player when ghost touches them
    for _, sprite := range g.Sprites {
        if sprite.Dead {
            continue
        }
        dx := sprite.X - g.PlayerX
        dy := sprite.Y - g.PlayerY
        dist := math.Sqrt(dx*dx + dy*dy)
        if dist < 0.8 {
            g.Health -= 1
            g.DamageFlash = 10
            newX := g.PlayerX - (dx/dist) * 1.0
            newY := g.PlayerY - (dy/dist) * 1.0
            if GetMap(g.CurrentMap)[int(newY)][int(newX)] == 1 {
                g.PlayerX = newX
                g.PlayerY = newY
            }
        }
    }

    // ammo pickups
    for i := range g.AmmoPickups {
        if !g.AmmoPickups[i].Active {
            continue
        }
        dx := g.AmmoPickups[i].X - g.PlayerX
        dy := g.AmmoPickups[i].Y - g.PlayerY
        dist := math.Sqrt(dx*dx + dy*dy)
        if dist < 0.8 {
            g.Ammo += 5
            g.AmmoPickups[i].Active = false
        }
    }

    // shooting
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
        if g.Ammo > 0 {
            g.GunKick = 8
            PlaySound("assets/shoot.wav")
            g.Ammo--
            for i := len(g.Sprites) - 1; i >= 0; i-- {
                if g.Sprites[i].Dead {
                    continue
                }
                dx := g.Sprites[i].X - g.PlayerX
                dy := g.Sprites[i].Y - g.PlayerY
                dist := math.Sqrt(dx*dx + dy*dy)
                spriteAngle := math.Atan2(dy, dx) - g.Angle
                for spriteAngle > math.Pi { spriteAngle -= 2 * math.Pi }
                for spriteAngle < -math.Pi { spriteAngle += 2 * math.Pi }
                if math.Abs(spriteAngle) < 0.2 && dist < 10 {
                    g.Sprites[i].Dead = true
                    g.Sprites[i].FadeTimer = 20
                    g.Score++
                    PlaySound("assets/ghost.wav")
                }
            }
        }
    }

    // remove fully faded dead sprites
    for i := len(g.Sprites) - 1; i >= 0; i-- {
        if g.Sprites[i].Dead && g.Sprites[i].FadeTimer == 0 {
            g.Sprites = append(g.Sprites[:i], g.Sprites[i+1:]...)
        }
    }

    // respawn when all dead
    if len(g.Sprites) == 0 {
        g.RespawnTimer++
        if g.RespawnTimer > 180 {
            g.Wave++
            g.Ammo += 3
			count := 3 + g.Wave
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
            g.RespawnTimer = 0
        }
    }

    // respawn ammo pickups every 3 waves
    if g.RespawnTimer == 1 {
        for i := range g.AmmoPickups {
            g.AmmoPickups[i].Active = true
        }
    }

    // check health and save scores 
	if g.Health <= 0 {
    if g.Score > g.HighScore {
        g.HighScore = g.Score
        SaveHighScore(g.HighScore)
    }
    g.GameState = 2
    return nil
}

    return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ScreenWidth, ScreenHeight
}
