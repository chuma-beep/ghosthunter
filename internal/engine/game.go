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
	GunFrame      int
  GunFrameTimer int
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
	ScreenShake    int 
	WeaponType     int 
	FireTimer      int
  HealthPickups []HealthPickup

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
		Pixels:     make([]byte, ScreenWidth*ScreenHeight*4),
		CurrentMap: 0,
		PlayerX:    8.0,
		PlayerY:    8.0,
		Angle:      0.0,
		Wave:       1,
		Ammo:       10,
		GameState:  0,
		Health:     100,
		HighScore:  highScore,
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
    HealthPickups: []HealthPickup{
      {X: 7.0, Y: 3.0, Active: true},
      {X: 3.0, Y: 7.0, Active: true},
      {X: 11.0, Y: 5.0, Active: true},
    },
	}
}

func (g *Game) Update() error {

	fmt.Println("Update called, GameState:", g.GameState, "Entities:", len(g.Entities), "Wave:", g.Wave, "eespawnTimer:", g.RespawnTimer)
// fmt.Println ("Update called, GameState:", g.GameState, "Entities:", len(g.Entities), "Wave:", g.Wave)
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Paused = !g.Paused
	}
	if g.Paused {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
    return ebiten.Termination
}
  	  
   //weapon switching 
	  if inpututil.IsKeyJustPressed(ebiten.Key1) {
    g.WeaponType = 0
}
if inpututil.IsKeyJustPressed(ebiten.Key2) {
    g.WeaponType = 1
}
if inpututil.IsKeyJustPressed(ebiten.Key3) {
    g.WeaponType = 2
}

if g.FireTimer > 0 {
    g.FireTimer--
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
    positions := [][2]float64{
        {6.0, 6.0}, {10.0, 4.0}, {3.0, 12.0},
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
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 0.8 {
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
    dist := math.Sqrt(dx*dx + dy*dy)
    if dist < 0.8 {
        g.Health += 25
        if g.Health > 100 {
            g.Health = 100
        }
        g.HealthPickups[i].Active = false
    }
}



	// shooting
canShoot := inpututil.IsKeyJustPressed(ebiten.KeySpace)
if g.WeaponType == 2 {
    canShoot = ebiten.IsKeyPressed(ebiten.KeySpace) && g.FireTimer == 0
}
if canShoot {
    ammoCost := 1
    if g.WeaponType == 1 {
        ammoCost = 3
    }
    if g.Ammo >= ammoCost {
        g.GunKick = 8
        g.ScreenShake = 8
       if g.WeaponType == 2 {
       g.GunFrame = 0
       g.GunFrameTimer = 1
      }
        PlaySound("assets/shoot.wav")
        g.Ammo -= ammoCost
        if g.WeaponType == 2 {
        g.FireTimer = 6 
       }
     if g.WeaponType == 1 {
      g.FireTimer = 30 
		}

    if g.GunFrameTimer > 0 {
    g.GunFrameTimer++
    if g.GunFrameTimer > 3 {
        g.GunFrameTimer = 1
        g.GunFrame++
        if g.GunFrame >= 8 {
            g.GunFrame = 0
            g.GunFrameTimer = 0
        }
    }
}

        switch g.WeaponType {
        case 0: // pistol - single ray
            g.shootRay(g.Angle, 1)
        case 1: // shotgun - 5 spread rays
            for s := -2; s <= 2; s++ {
                g.shootRay(g.Angle+float64(s)*0.05, 2)
            }
        case 2: // machinegun - single ray fast
            g.shootRay(g.Angle, 1)
        }
    }
}

// entity movement and animation
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
    if dist <= 0.5 {
        g.Health -= g.Entities[i].Damage
        g.DamageFlash = 10
        continue
    }
    angle := math.Atan2(dy, dx)
    switch g.Entities[i].Type {
    case EntityDemon:
        angle += math.Sin(float64(g.RespawnTimer+i)*0.3) * 0.5
       
               case EntityWraith:
    // circle player while slowly closing in
    perpAngle := angle + math.Pi/2
    // move 70% toward player, 30% perpendicular
    combinedX := math.Cos(angle)*0.7 + math.Cos(perpAngle)*0.3
    combinedY := math.Sin(angle)*0.7 + math.Sin(perpAngle)*0.3
    newEX := g.Entities[i].X + combinedX*g.Entities[i].Speed
    newEY := g.Entities[i].Y + combinedY*g.Entities[i].Speed
    if int(newEY) >= 0 && int(newEY) < GetMapHeight(g.CurrentMap) &&
        int(newEX) >= 0 && int(newEX) < GetMapWidth(g.CurrentMap) &&
        GetMap(g.CurrentMap)[int(newEY)][int(newEX)] == 0 {
        g.Entities[i].X = newEX
        g.Entities[i].Y = newEY
    }
    continue 
		       



    case EntityReaper:
        if dist > 8.0 {
            g.Entities[i].X = g.PlayerX - math.Cos(angle)*3.0
            g.Entities[i].Y = g.PlayerY - math.Sin(angle)*3.0
        }
    }
    newEX := g.Entities[i].X + math.Cos(angle)*g.Entities[i].Speed
    newEY := g.Entities[i].Y + math.Sin(angle)*g.Entities[i].Speed
    if int(newEY) >= 0 && int(newEY) < GetMapHeight(g.CurrentMap) &&
        int(newEX) >= 0 && int(newEX) < GetMapWidth(g.CurrentMap) &&
        GetMap(g.CurrentMap)[int(newEY)][int(newEX)] == 0 {
        g.Entities[i].X = newEX
        g.Entities[i].Y = newEY
    }
    // animate
    g.Entities[i].FrameTimer++
    if g.Entities[i].FrameTimer > 8 {
        g.Entities[i].FrameTimer = 0
        fc := enemyFrameCount(g.Entities[i].Type)
        if fc > 1 {
            g.Entities[i].Frame = (g.Entities[i].Frame + 1) % fc
        }
    }

// movement with wall avoidance
if g.Entities[i].Type == EntityGhost {
    // ghosts walk through walls
    g.Entities[i].X += math.Cos(angle) * g.Entities[i].Speed
    g.Entities[i].Y += math.Sin(angle) * g.Entities[i].Speed
} else {
    newEX := g.Entities[i].X + math.Cos(angle)*g.Entities[i].Speed
    newEY := g.Entities[i].Y + math.Sin(angle)*g.Entities[i].Speed
    movedX := false
    movedY := false
    // try moving on X axis
    if int(newEY) >= 0 && int(newEY) < GetMapHeight(g.CurrentMap) &&
        int(newEX) >= 0 && int(newEX) < GetMapWidth(g.CurrentMap) &&
        GetMap(g.CurrentMap)[int(newEY)][int(newEX)] == 0 {
        g.Entities[i].X = newEX
        g.Entities[i].Y = newEY
        movedX = true
        movedY = true
    }
    // if blocked try sliding along X only
    if !movedX {
        newEX2 := g.Entities[i].X + math.Cos(angle)*g.Entities[i].Speed
        if int(g.Entities[i].Y) >= 0 && int(g.Entities[i].Y) < GetMapHeight(g.CurrentMap) &&
            int(newEX2) >= 0 && int(newEX2) < GetMapWidth(g.CurrentMap) &&
            GetMap(g.CurrentMap)[int(g.Entities[i].Y)][int(newEX2)] == 0 {
            g.Entities[i].X = newEX2
            movedY = true
        }
    }
    // try sliding along Y only
    if !movedY {
        newEY2 := g.Entities[i].Y + math.Sin(angle)*g.Entities[i].Speed
        if int(newEY2) >= 0 && int(newEY2) < GetMapHeight(g.CurrentMap) &&
            int(g.Entities[i].X) >= 0 && int(g.Entities[i].X) < GetMapWidth(g.CurrentMap) &&
            GetMap(g.CurrentMap)[int(newEY2)][int(g.Entities[i].X)] == 0 {
            g.Entities[i].Y = newEY2
        }
    }
    // if still blocked rotate randomly to find a way around
    if !movedX && !movedY {
        if i%2 == 0 {
            angle += 0.3
        } else {
            angle -= 0.3
        }
        g.Entities[i].VX = math.Cos(angle)
        g.Entities[i].VY = math.Sin(angle)
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

	// respawn ammo pickups every 3 waves
if g.RespawnTimer == 1 {
    for i := range g.AmmoPickups {
        g.AmmoPickups[i].Active = true
    }
    for i := range g.HealthPickups {
        g.HealthPickups[i].Active = true
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



//enemyFrameCount
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


//shoot ray 
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
                g.Entities[i].FadeTimer = 20
                g.Score++
                PlaySound("assets/ghost.wav")
            }
        }
    }
}



func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
