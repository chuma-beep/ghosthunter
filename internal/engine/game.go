package engine

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type Game struct {
	Pixels         []byte
	PlayerX        float64
	PlayerY        float64
	Angle          float64
	Score          int
	RespawnTimer   int
	Health         int
	DamageFlash    int
	Wave           int
	GunKick        int
	Ammo           int
	AmmoPickups    []AmmoPickup
	GameState      int
	HighScore      int
	CurrentMap     int
	Entities       []Entity
	LevelNameTimer int
    WaveTransition int 
	Paused         bool
}

func NewGame() *Game {
	return &Game{
		Pixels:     make([]byte, ScreenWidth*ScreenHeight*4),
		CurrentMap: 0,
		PlayerX:    8.0,
		PlayerY:    8.0,
		Angle:      0.0,
		Wave:       1,
		Ammo:       10,
		GameState:  0,
		Health:     100,
		AmmoPickups: []AmmoPickup{
			{X: 5.0, Y: 5.0, Active: true},
			{X: 11.0, Y: 11.0, Active: true},
			{X: 3.0, Y: 9.0, Active: true},
		},
		Entities: []Entity{
			{X: 2.0, Y: 2.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
			{X: 28.0, Y: 2.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
			{X: 2.0, Y: 28.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
		},
	}
}

func (g *Game) Update() error {

	fmt.Println("Update called, GameState:", g.GameState, "Entities:", len(g.Entities), "Wave:", g.Wave, "eespawnTimer:", g.RespawnTimer)
	// fmt.Println("Update called, GameState:", g.GameState, "Entities:", len(g.Entities), "Wave:", g.Wave)
	// start screen
	if g.GameState == 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.GameState = 1
		}

		if g.LevelNameTimer > 0 {
			g.LevelNameTimer--
		}
		return nil
	 }
    
	 //pause menu 
	 if inpututil.IsKeyJustPressed(ebiten.KeyEscape){
		 g.Paused = !g.Paused 
	 }
	 if g.Paused {
		 return nil
	 }


	//wave decay 
if g.WaveTransition > 0 {
    g.WaveTransition--
}

	// game over screen
	if g.GameState == 2 {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.Health = 100
			g.Score = 0
			g.Ammo = 10
			g.Wave = 1
			g.RespawnTimer = 0
			g.CurrentMap = 0
			g.PlayerX = 8.0
			g.PlayerY = 8.0
			g.Angle = 0.0
			g.GameState = 1
			g.Entities = []Entity{
				{X: 14.0, Y: 14.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
				{X: 14.0, Y: 2.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
				{X: 2.0, Y: 14.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
			}
		}
		return nil
	}

	// portal detection
	portalX, portalY := 13.0, 1.0
	pdx := g.PlayerX - portalX
	pdy := g.PlayerY - portalY
	portalDist := math.Sqrt(pdx*pdx + pdy*pdy)
	if portalDist < 0.8 {
		g.CurrentMap = (g.CurrentMap + 1) % 5
		g.PlayerX = 2.0
		g.PlayerY = 1.5
		g.RespawnTimer = 0
		g.LevelNameTimer = 60
		if g.CurrentMap == 0 || g.CurrentMap == 1 {
			g.Entities = []Entity{
				{X: 6.0, Y: 6.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
				{X: 10.0, Y: 4.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
				{X: 3.0, Y: 12.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
			}
		} else {
			g.Entities = []Entity{
				{X: 6.0, Y: 6.0, Type: EntityWizard, Health: 2, Speed: 0.003, Damage: 2},
				{X: 10.0, Y: 4.0, Type: EntityWizard, Health: 2, Speed: 0.003, Damage: 2},
				{X: 3.0, Y: 12.0, Type: EntityWizard, Health: 2, Speed: 0.003, Damage: 2},
			}
		}
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
		if int(newY) >= 0 && int(newY) < GetMapHeight(g.CurrentMap) &&
			int(newX) >= 0 && int(newX) < GetMapWidth(g.CurrentMap) &&
			GetMap(g.CurrentMap)[int(newY)][int(newX)] == 0 {
			g.PlayerX = newX
			g.PlayerY = newY
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		newX := g.PlayerX - math.Cos(g.Angle)*0.05
		newY := g.PlayerY - math.Sin(g.Angle)*0.05
		if int(newY) >= 0 && int(newY) < GetMapHeight(g.CurrentMap) &&
			int(newX) >= 0 && int(newX) < GetMapWidth(g.CurrentMap) &&
			GetMap(g.CurrentMap)[int(newY)][int(newX)] == 0 {
			g.PlayerX = newX
			g.PlayerY = newY
		}
	}

	// gun kick decay
	if g.GunKick > 0 {
		g.GunKick--
	}

	// move sprites toward player
	for i := range g.Entities {

		if g.Entities[i].Dead {
			if g.Entities[i].FadeTimer > 0 {
				g.Entities[i].FadeTimer--
			}
			continue
		}

		dx := g.PlayerX - g.Entities[i].X
		dy := g.PlayerY - g.Entities[i].Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > 0.5 {

			angle := math.Atan2(dy, dx)

			// slight wobble so enemies don't track perfectly
			// angle += (rand.Float64() - 0.5) * 0.2

			moveX := math.Cos(angle) * g.Entities[i].Speed
			moveY := math.Sin(angle) * g.Entities[i].Speed

			g.Entities[i].X += moveX
			g.Entities[i].Y += moveY
		}
	}

	// damage player when ghost touches them
	for _, sprite := range g.Entities {
		if sprite.Dead {
			continue
		}
		dx := sprite.X - g.PlayerX
		dy := sprite.Y - g.PlayerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 0.8 {
			g.Health -= sprite.Damage
			g.DamageFlash = 10
			newX := g.PlayerX - (dx/dist)*1.0
			newY := g.PlayerY - (dy/dist)*1.0
			if GetMap(g.CurrentMap)[int(newY)][int(newX)] == 0 {
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
			for i := len(g.Entities) - 1; i >= 0; i-- {
				if g.Entities[i].Dead {
					continue
				}
				dx := g.Entities[i].X - g.PlayerX
				dy := g.Entities[i].Y - g.PlayerY
				dist := math.Sqrt(dx*dx + dy*dy)
				spriteAngle := math.Atan2(dy, dx) - g.Angle
				for spriteAngle > math.Pi {
					spriteAngle -= 2 * math.Pi
				}
				for spriteAngle < -math.Pi {
					spriteAngle += 2 * math.Pi
				}
				if math.Abs(spriteAngle) < 0.2 && dist < 10 {
					g.Entities[i].Dead = true
					g.Entities[i].FadeTimer = 20
					g.Score++
					PlaySound("assets/ghost.wav")
				}
			}
		}
	}

	// remove fully faded dead sprites
	for i := len(g.Entities) - 1; i >= 0; i-- {
		if g.Entities[i].Dead && g.Entities[i].FadeTimer == 0 {
			g.Entities = append(g.Entities[:i], g.Entities[i+1:]...)
		}
	}

	// respawn when all dead
	if len(g.Entities) == 0 {
		g.RespawnTimer++
		if g.RespawnTimer > 180 {
			g.Wave++
            g.WaveTransition = 60
			g.Ammo += 3
			count := 3 + g.Wave
			var positions [][2]float64
			if g.CurrentMap <= 1 {
				positions = [][2]float64{
					{6.0, 6.0}, {10.0, 4.0}, {3.0, 12.0},
					{12.0, 12.0}, {8.0, 3.0}, {2.0, 8.0}, {13.0, 7.0},
				}
			} else {
				positions = [][2]float64{
					{15.0, 15.0}, {25.0, 5.0}, {5.0, 25.0},
					{25.0, 25.0}, {15.0, 5.0}, {5.0, 15.0}, {20.0, 20.0},
				}
			}
			g.Entities = []Entity{}
			for i := 0; i < count && i < len(positions); i++ {
				entityType := EntityGhost
				health := 1
				if g.CurrentMap > 1 {
					entityType = EntityWizard
					health = 2
				}
				g.Entities = append(g.Entities, Entity{
					X: positions[i][0], Y: positions[i][1],
					Type: entityType, Health: health,
					Speed:  0.003 + float64(g.Wave)*0.0005,
					Damage: 1,
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
