package engine

import (
	"errors"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Pixels          []byte
	PlayerX         float64
	PlayerY         float64
	Angle           float64
	Score           int
	RespawnTimer    int
	Health          int
	DamageFlash     int
	Wave            int
	GunKick         int
	Ammo            int
	AmmoPickups     []AmmoPickup
	HealthPickups   []HealthPickup
	GameState       int
	HighScore       int
	CurrentMap      int
	Entities        []Entity
	LevelNameTimer  int
	WaveTransition  int
	Paused          bool
	ScreenShake     int
	WeaponType      int
	WeaponStateID   WeaponStateID
	WeaponStateTics int
    PauseMenuSelection int 
	ShowControls    bool 
}

func enemyForMap(mapIndex int) EntityType {
	switch mapIndex {
	case 0:
		return EntityGhost
	case 1:
		return EntityWizard
	case 2:
		return EntityDemon
	case 3:
		return EntityWraith
	case 4:
		return EntityReaper
	default:
		return EntityGhost
	}
}

func NewGame() *Game {
	highScore := LoadHighScore()
	return &Game{
		Pixels:          make([]byte, ScreenWidth*ScreenHeight*4),
		CurrentMap:      0,
		PlayerX:         8.0,
		PlayerY:         8.0,
		Angle:           0.0,
		Wave:            1,
		Ammo:            30,
		GameState:       0,
		Health:          100,
		HighScore:       highScore,
		WeaponStateID:   S_PISTOL_READY,
		WeaponStateTics: -1,
		AmmoPickups: []AmmoPickup{
			{X: 5.0, Y: 5.0, Active: true},
			{X: 11.0, Y: 11.0, Active: true},
			{X: 3.0, Y: 9.0, Active: true},
		},
		HealthPickups: []HealthPickup{
			{X: 7.0, Y: 3.0, Active: true},
			{X: 3.0, Y: 7.0, Active: true},
			{X: 11.0, Y: 5.0, Active: true},
		},
    Entities: []Entity{
      {X: 3.0, Y: 1.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
      {X: 7.0, Y: 1.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
      {X: 12.0, Y: 1.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
    },
	}
}

// LineOfSight casts a ray from (x1,y1) to (x2,y2) and returns true if no wall blocks it
func (g *Game) LineOfSight(x1, y1, x2, y2 float64) bool {
    dx := x2 - x1
    dy := y2 - y1
    dist := math.Sqrt(dx*dx + dy*dy)
    if dist > 15.0 {
        return false
    }
    steps := int(dist / 0.5) // coarser steps
    if steps < 2 {
        steps = 2
    }
    for i := 1; i < steps; i++ {
        t := float64(i) / float64(steps)
        rx := x1 + dx*t
        ry := y1 + dy*t
        if int(ry) < 0 || int(ry) >= GetMapHeight(g.CurrentMap) ||
            int(rx) < 0 || int(rx) >= GetMapWidth(g.CurrentMap) {
            return false
        }
        if GetMap(g.CurrentMap)[int(ry)][int(rx)] == 1 {
            return false
        }
    }
    return true
}



// moveEntity moves entity i toward angle, with wall sliding and type-specific behavior
func (g *Game) moveEntity(i int, angle float64) {
	switch g.Entities[i].Type {
	case EntityGhost:
		// ghosts walk through walls
		g.Entities[i].X += math.Cos(angle) * g.Entities[i].Speed
		g.Entities[i].Y += math.Sin(angle) * g.Entities[i].Speed
		return
	case EntityDemon:
		angle += math.Sin(float64(g.RespawnTimer+i)*0.5) * 0.8
	case EntityWraith:
		perpAngle := angle + math.Pi/2
		combinedX := math.Cos(angle)*0.6 + math.Cos(perpAngle)*0.4
		combinedY := math.Sin(angle)*0.6 + math.Sin(perpAngle)*0.4
		newEX := g.Entities[i].X + combinedX*g.Entities[i].Speed*4
		newEY := g.Entities[i].Y + combinedY*g.Entities[i].Speed*4
		if int(newEY) >= 0 && int(newEY) < GetMapHeight(g.CurrentMap) &&
			int(newEX) >= 0 && int(newEX) < GetMapWidth(g.CurrentMap) &&
			GetMap(g.CurrentMap)[int(newEY)][int(newEX)] == 0 {
			g.Entities[i].X = newEX
			g.Entities[i].Y = newEY
		}
		return
	case EntityReaper:
		dx := g.PlayerX - g.Entities[i].X
		dy := g.PlayerY - g.Entities[i].Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 8.0 {
			g.Entities[i].X = g.PlayerX - math.Cos(angle)*3.0
			g.Entities[i].Y = g.PlayerY - math.Sin(angle)*3.0
			return
		}
	}

	// default movement with wall sliding
	newEX := g.Entities[i].X + math.Cos(angle)*g.Entities[i].Speed
	newEY := g.Entities[i].Y + math.Sin(angle)*g.Entities[i].Speed
	movedX, movedY := false, false

	if int(newEY) >= 0 && int(newEY) < GetMapHeight(g.CurrentMap) &&
		int(newEX) >= 0 && int(newEX) < GetMapWidth(g.CurrentMap) &&
		GetMap(g.CurrentMap)[int(newEY)][int(newEX)] == 0 {
		g.Entities[i].X = newEX
		g.Entities[i].Y = newEY
		movedX, movedY = true, true
	}

	if !movedX {
		newEX2 := g.Entities[i].X + math.Cos(angle)*g.Entities[i].Speed
		if int(g.Entities[i].Y) >= 0 && int(g.Entities[i].Y) < GetMapHeight(g.CurrentMap) &&
			int(newEX2) >= 0 && int(newEX2) < GetMapWidth(g.CurrentMap) &&
			GetMap(g.CurrentMap)[int(g.Entities[i].Y)][int(newEX2)] == 0 {
			g.Entities[i].X = newEX2
			movedY = true
		}
	}

	if !movedY {
		newEY2 := g.Entities[i].Y + math.Sin(angle)*g.Entities[i].Speed
		if int(newEY2) >= 0 && int(newEY2) < GetMapHeight(g.CurrentMap) &&
			int(g.Entities[i].X) >= 0 && int(g.Entities[i].X) < GetMapWidth(g.CurrentMap) &&
			GetMap(g.CurrentMap)[int(newEY2)][int(g.Entities[i].X)] == 0 {
			g.Entities[i].Y = newEY2
		}
	}

	// rotate if fully blocked
	if !movedX && !movedY {
		if i%2 == 0 {
			g.Entities[i].FacingAngle += 0.3
		} else {
			g.Entities[i].FacingAngle -= 0.3
		}
	}
}

func (g *Game) Update() error {
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

	// quit
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

    if inpututil.IsKeyJustPressed(ebiten.KeyEscape){
		if g.ShowControls{
			g.ShowControls = false 
		}else {
			g.Paused = !g.Paused
		}
}      


    if g.Paused {
        // When a controls submenu is active, handle its specific input
        if g.ShowControls {
            // Press Escape to close the controls screen, returning to the pause menu
            if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
                g.ShowControls = false
            }
            // You might also add other keys to scroll through controls info if needed
            return nil
        }

        // --- Standard pause menu navigation ---
        // Up/Down keys change the selected menu item
        if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
            g.PauseMenuSelection--
            if g.PauseMenuSelection < 0 {
                g.PauseMenuSelection = 2 // last item index (assuming 3 options: 0,1,2)
            }
        }
        if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
            g.PauseMenuSelection++
            if g.PauseMenuSelection > 2 {
                g.PauseMenuSelection = 0
            }
        }

        // Enter key confirms the selection
        if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
            switch g.PauseMenuSelection {
            case 0: // Resume
                g.Paused = false
                g.PauseMenuSelection = 0 // reset for next pause
            case 1: // Show controls
                g.ShowControls = true
                // Keep paused, selection unchanged – when controls close, we return to pause menu
            case 2: // Quit
                // e.g., os.Exit(0) or signal a quit
                return errors.New("quit requested")
            }
        }


        return nil
    }


    
	// wave transition decay
	if g.WaveTransition > 0 {
		g.WaveTransition--
	}

	// game over screen
	if g.GameState == 2 {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.Health = 100
			g.Score = 0
			g.Ammo = 30
			g.Wave = 1
			g.RespawnTimer = 0
			g.CurrentMap = 0
			g.PlayerX = 8.0
			g.PlayerY = 8.0
			g.Angle = 0.0
			g.GameState = 1
			g.WeaponStateID = S_PISTOL_READY
			g.WeaponStateTics = -1
      g.Entities = []Entity{
        {X: 3.0, Y: 1.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
        {X: 7.0, Y: 1.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
        {X: 12.0, Y: 1.0, Type: EntityGhost, Health: 1, Speed: 0.003, Damage: 1},
      }
		}
		return nil
	}

	// weapon switching
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.WeaponType = 0
		g.SetWeaponState(S_PISTOL_READY)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.WeaponType = 1
		g.SetWeaponState(S_SHOTGUN_READY)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.WeaponType = 2
		g.SetWeaponState(S_MACHINEGUN_READY)
	}

	// portal detection
	var portalX, portalY float64
	if g.CurrentMap <= 1 {
		portalX, portalY = 13.0, 1.0
	} else {
		portalX, portalY = 28.0, 1.5
	}
	pdx := g.PlayerX - portalX
	pdy := g.PlayerY - portalY
	if math.Sqrt(pdx*pdx+pdy*pdy) < 0.8 {
		g.CurrentMap = (g.CurrentMap + 1) % 5
		g.PlayerX = 2.0
		g.PlayerY = 1.5
		g.RespawnTimer = 0
		g.LevelNameTimer = 60
    positions := [][2]float64{
       {3.0, 1.0}, {7.0, 1.0}, {12.0, 1.0},
    }
    if g.CurrentMap >= 2 {
      positions = [][2]float64{
        {15.0, 15.0}, {25.0, 5.0}, {5.0, 25.0},
      }
    }
		g.Entities = []Entity{}
		for _, pos := range positions {
			g.Entities = append(g.Entities, Entity{
				X:      pos[0],
				Y:      pos[1],
				Type:   enemyForMap(g.CurrentMap),
				Health: 1 + g.CurrentMap,
				Speed:  0.003 + float64(g.CurrentMap)*0.001,
				Damage: 1 + g.CurrentMap/2,
			})
		}
	}

	// player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.Angle -= 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.Angle += 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		newX := g.PlayerX + math.Cos(g.Angle)*0.08
		newY := g.PlayerY + math.Sin(g.Angle)*0.08
		if int(newY) >= 0 && int(newY) < GetMapHeight(g.CurrentMap) &&
			int(newX) >= 0 && int(newX) < GetMapWidth(g.CurrentMap) &&
			GetMap(g.CurrentMap)[int(newY)][int(newX)] == 0 {
			g.PlayerX = newX
			g.PlayerY = newY
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		newX := g.PlayerX - math.Cos(g.Angle)*0.08
		newY := g.PlayerY - math.Sin(g.Angle)*0.08
		if int(newY) >= 0 && int(newY) < GetMapHeight(g.CurrentMap) &&
			int(newX) >= 0 && int(newX) < GetMapWidth(g.CurrentMap) &&
			GetMap(g.CurrentMap)[int(newY)][int(newX)] == 0 {
			g.PlayerX = newX
			g.PlayerY = newY
		}
	}

	// ammo pickups
	for i := range g.AmmoPickups {
		if !g.AmmoPickups[i].Active {
			continue
		}
		dx := g.AmmoPickups[i].X - g.PlayerX
		dy := g.AmmoPickups[i].Y - g.PlayerY
		if math.Sqrt(dx*dx+dy*dy) < 0.8 {
			g.Ammo += 5
			g.AmmoPickups[i].Active = false
		}
	}

	// health pickups
	for i := range g.HealthPickups {
		if !g.HealthPickups[i].Active {
			continue
		}
		dx := g.HealthPickups[i].X - g.PlayerX
		dy := g.HealthPickups[i].Y - g.PlayerY
		if math.Sqrt(dx*dx+dy*dy) < 0.8 {
			g.Health += 25
			if g.Health > 100 {
				g.Health = 100
			}
			g.HealthPickups[i].Active = false
		}
	}

	// weapon state machine
	g.TickWeapon()

	// gun kick decay
	if g.GunKick > 0 {
		g.GunKick--
	}

	// screen shake decay
	if g.ScreenShake > 0 {
		g.ScreenShake--
	}

	// entity update — state machine
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
		angle := math.Atan2(dy, dx)
		// hasLOS := g.LineOfSight(g.Entities[i].X, g.Entities[i].Y, g.PlayerX, g.PlayerY)

    // only recheck LOS every 10 frames
  // g.Entities[i].LOSTimer++
  // if g.Entities[i].LOSTimer >= 20 {
  //   g.Entities[i].LOSTimer = 0
  //   g.Entities[i].HasLOS = g.LineOfSight(
  //       g.Entities[i].X, g.Entities[i].Y,
  //       g.PlayerX, g.PlayerY,
  //   )
  // }
  // hasLOS := g.Entities[i].HasLOS
  hasLOS := true



		switch g.Entities[i].State {
		case StateChase:
			if dist < 0.8 && hasLOS {
				g.Entities[i].State = StateAttack
				g.Entities[i].StateTimer = 30
			} else if hasLOS {
				g.Entities[i].FacingAngle = angle
				g.moveEntity(i, angle)
			} else {
				g.moveEntity(i, g.Entities[i].FacingAngle)
			}

		case StateAttack:
			g.Entities[i].StateTimer--
			if g.Entities[i].StateTimer == 15 {
				if dist < 1.5 && hasLOS {
					g.Health -= g.Entities[i].Damage
					g.DamageFlash = 10
				}
			}
			if g.Entities[i].StateTimer <= 0 {
				g.Entities[i].State = StateChase
			}

		case StatePain:
			g.Entities[i].StateTimer--
			if g.Entities[i].StateTimer <= 0 {
				g.Entities[i].State = StateChase
			}

		case StateDeath:
			if g.Entities[i].FadeTimer > 0 {
				g.Entities[i].FadeTimer--
			}
		}

		// animate frame
		g.Entities[i].FrameTimer++
		if g.Entities[i].FrameTimer > 8 {
			g.Entities[i].FrameTimer = 0
			fc := enemyFrameCount(g.Entities[i].Type)
			if fc > 1 {
				g.Entities[i].Frame = (g.Entities[i].Frame + 1) % fc
			}
		}
	}

	// remove fully faded dead entities
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
         {3.0, 1.0}, {7.0, 1.0}, {12.0, 1.0},
         {5.0, 5.0}, {10.0, 5.0}, {3.0, 8.0}, {12.0, 8.0},
      }
    } else {
      positions = [][2]float64{
         {15.0, 15.0}, {25.0, 5.0}, {5.0, 25.0},
         {25.0, 25.0}, {15.0, 5.0}, {5.0, 15.0}, {20.0, 20.0},
      }
    }
			g.Entities = []Entity{}
			for i := 0; i < count && i < len(positions); i++ {
				g.Entities = append(g.Entities, Entity{
					X:      positions[i][0],
					Y:      positions[i][1],
					Type:   enemyForMap(g.CurrentMap),
					Health: 1 + g.CurrentMap,
					Speed:  0.003 + float64(g.Wave)*0.0005,
					Damage: 1 + g.CurrentMap/2,
				})
			}
			g.RespawnTimer = 0
		}
	}

	// respawn pickups
	if g.RespawnTimer == 1 {
		for i := range g.AmmoPickups {
			g.AmmoPickups[i].Active = true
		}
		for i := range g.HealthPickups {
			g.HealthPickups[i].Active = true
		}
	}

	// damage flash decay
	if g.DamageFlash > 0 {
		g.DamageFlash--
	}

	// check health
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

func enemyFrameCount(t EntityType) int {
	switch t {
	case EntityDemon:
		return demonFrames
	case EntityWraith:
		return wraithFrames
	case EntityReaper:
		return reaperFrames
	default:
		return 1
	}
}

func (g *Game) shootRay(angle float64, damage int) {
	for i := len(g.Entities) - 1; i >= 0; i-- {
		if g.Entities[i].Dead {
			continue
		}
		dx := g.Entities[i].X - g.PlayerX
		dy := g.Entities[i].Y - g.PlayerY
		dist := math.Sqrt(dx*dx + dy*dy)
		spriteAngle := math.Atan2(dy, dx) - angle
		for spriteAngle > math.Pi {
			spriteAngle -= 2 * math.Pi
		}
		for spriteAngle < -math.Pi {
			spriteAngle += 2 * math.Pi
		}
		if math.Abs(spriteAngle) < 0.2 && dist < 10 {
			g.Entities[i].Health -= damage
			if g.Entities[i].Health <= 0 {
				g.Entities[i].Dead = true
				g.Entities[i].State = StateDeath
				g.Entities[i].FadeTimer = 20
				g.Score++
				PlaySound("assets/ghost.wav")
			} else {
				g.Entities[i].State = StatePain
				g.Entities[i].StateTimer = 10
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func isSpacePressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

func isSpaceJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}
