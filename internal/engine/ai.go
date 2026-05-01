package engine

import "math"

type AIController struct {
	Enabled        bool
	ShootRequested bool
	MoveForward    bool
	MoveBackward   bool
	ExploreTargetX float64
	ExploreTargetY float64
	ExploreTimer   int
	StuckCounter   int
	LastX          float64
	LastY          float64
}

func NewAIController() *AIController {
	return &AIController{Enabled: false}
}

func (ai *AIController) Update(g *Game) {
	if !ai.Enabled {
		return
	}

	ai.MoveForward = false
	ai.MoveBackward = false

	// Check if stuck
	dx := g.PlayerX - ai.LastX
	dy := g.PlayerY - ai.LastY
	distMoved := math.Sqrt(dx*dx + dy*dy)

	if distMoved < 0.01 {
		ai.StuckCounter++
		if ai.StuckCounter > 30 {
			// Stuck - turn randomly and pick new target
			g.Angle += 0.3
			ai.pickNewExplorationTarget(g)
			ai.ExploreTimer = 200
			ai.StuckCounter = 0
		}
	} else {
		ai.StuckCounter = 0
	}
	ai.LastX = g.PlayerX
	ai.LastY = g.PlayerY

	idx, dist := ai.findNearestEnemy(g)

	// Check for portal
	if ai.checkPortal(g) {
		return
	}

	if idx == -1 {
		// No enemies - explore
		ai.explore(g)
		return
	}

	enemy := &g.Entities[idx]
	targetAngle := math.Atan2(enemy.Y-g.PlayerY, enemy.X-g.PlayerX)
	diff := targetAngle - g.Angle

	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}

	// Turn toward enemy
	if diff > 0.05 {
		g.Angle += 0.05
	} else if diff < -0.05 {
		g.Angle -= 0.05
	}

	// Check if can move forward (no wall)
	canMoveForward := ai.canMoveInDirection(g, g.Angle, 0.08)

	// Movement based on distance
	if dist > 6 && canMoveForward {
		ai.MoveForward = true
	} else if dist < 2 {
		ai.MoveBackward = true
	}

	// Shoot if aligned and in range
	if math.Abs(diff) < 0.3 && dist < 8 {
		ai.ShootRequested = true
	}
}

func (ai *AIController) canMoveInDirection(g *Game, angle float64, dist float64) bool {
	newX := g.PlayerX + math.Cos(angle)*dist
	newY := g.PlayerY + math.Sin(angle)*dist
	mapH := GetMapHeight(g.CurrentMap)
	mapW := GetMapWidth(g.CurrentMap)

	if int(newY) < 0 || int(newY) >= mapH || int(newX) < 0 || int(newX) >= mapW {
		return false
	}
	if GetMap(g.CurrentMap)[int(newY)][int(newX)] == 1 {
		return false
	}
	return true
}

func (ai *AIController) checkPortal(g *Game) bool {
	var portalX, portalY float64
	if g.CurrentMap <= 1 {
		portalX, portalY = 13.0, 1.0
	} else {
		portalX, portalY = 28.0, 1.5
	}

	dx := g.PlayerX - portalX
	dy := g.PlayerY - portalY
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 2.0 {
		targetAngle := math.Atan2(portalY-g.PlayerY, portalX-g.PlayerX)
		diff := targetAngle - g.Angle
		for diff > math.Pi {
			diff -= 2 * math.Pi
		}
		for diff < -math.Pi {
			diff += 2 * math.Pi
		}
		if diff > 0.05 {
			g.Angle += 0.05
		} else if diff < -0.05 {
			g.Angle -= 0.05
		}
		if ai.canMoveInDirection(g, g.Angle, 0.08) {
			ai.MoveForward = true
		}
		return true
	}
	return false
}

func (ai *AIController) explore(g *Game) {
	// If no target or reached target, pick a new valid exploration point
	if ai.ExploreTimer <= 0 || ai.reachedTarget(g) {
		ai.pickNewExplorationTarget(g)
		ai.ExploreTimer = 300
	}
	ai.ExploreTimer--

	// Move toward exploration target
	targetAngle := math.Atan2(ai.ExploreTargetY-g.PlayerY, ai.ExploreTargetX-g.PlayerX)
	diff := targetAngle - g.Angle

	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}

	if diff > 0.1 {
		g.Angle += 0.05
	} else if diff < -0.1 {
		g.Angle -= 0.05
	} else if ai.canMoveInDirection(g, g.Angle, 0.08) {
		ai.MoveForward = true
	}
}

func (ai *AIController) reachedTarget(g *Game) bool {
	dx := g.PlayerX - ai.ExploreTargetX
	dy := g.PlayerY - ai.ExploreTargetY
	return math.Sqrt(dx*dx+dy*dy) < 1.0
}

func (ai *AIController) pickNewExplorationTarget(g *Game) {
	mapData := GetMap(g.CurrentMap)
	mapW := GetMapWidth(g.CurrentMap)
	mapH := GetMapHeight(g.CurrentMap)

	type spot struct{ x, y int }
	var validSpots []spot

	// Try to find spots near player (within 5 tiles)
	searchRadius := 5
	for dy := -searchRadius; dy <= searchRadius; dy++ {
		for dx := -searchRadius; dx <= searchRadius; dx++ {
			px := int(g.PlayerX) + dx
			py := int(g.PlayerY) + dy
			if py >= 0 && py < mapH && px >= 0 && px < mapW && mapData[py][px] == 0 {
				validSpots = append(validSpots, spot{px, py})
			}
		}
	}

	if len(validSpots) > 0 {
		// Pick a random spot not too close
		for i := 0; i < len(validSpots); i++ {
			idx := i + int(ai.ExploreTimer)%(len(validSpots)-i)
			validSpots[i], validSpots[idx] = validSpots[idx], validSpots[i]
		}

		for _, s := range validSpots {
			dist := math.Sqrt(math.Pow(float64(s.x)-g.PlayerX, 2) + math.Pow(float64(s.y)-g.PlayerY, 2))
			if dist > 2 {
				ai.ExploreTargetX = float64(s.x) + 0.5
				ai.ExploreTargetY = float64(s.y) + 0.5
				return
			}
		}

		// Just pick any valid spot
		ai.ExploreTargetX = float64(validSpots[0].x) + 0.5
		ai.ExploreTargetY = float64(validSpots[0].y) + 0.5
	} else {
		// Fallback - move in a pattern
		angle := float64(ai.ExploreTimer%20) * math.Pi / 10
		ai.ExploreTargetX = g.PlayerX + math.Cos(angle)*3
		ai.ExploreTargetY = g.PlayerY + math.Sin(angle)*3
	}
}

func (ai *AIController) findNearestEnemy(g *Game) (int, float64) {
	idx := -1
	minDist := math.MaxFloat64

	for i, e := range g.Entities {
		if e.Dead {
			continue
		}
		dx := e.X - g.PlayerX
		dy := e.Y - g.PlayerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < minDist && g.LineOfSight(g.PlayerX, g.PlayerY, e.X, e.Y) {
			minDist = dist
			idx = i
		}
	}

	return idx, minDist
}
