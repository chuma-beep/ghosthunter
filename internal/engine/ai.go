package engine

import "math"

type AIController struct {
	Enabled        bool
	ShootRequested bool
	MoveForward    bool
	MoveBackward   bool
	TurnLeft       bool
	TurnRight      bool
	ExploreTargetX float64
	ExploreTargetY float64
	ExploreTimer   int
	StuckCounter   int
	LastX          float64
	LastY          float64
	LastAngle      float64
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
	ai.TurnLeft = false
	ai.TurnRight = false

	// Check if stuck (no movement for 20+ frames)
	dx := g.PlayerX - ai.LastX
	dy := g.PlayerY - ai.LastY
	distMoved := math.Sqrt(dx*dx + dy*dy)

	if distMoved < 0.02 {
		ai.StuckCounter++
		if ai.StuckCounter > 20 {
			// Stuck - try turning and moving differently
			if ai.StuckCounter%2 == 0 {
				g.Angle += 0.4 // Turn around
			} else {
				g.Angle -= 0.4
			}
			ai.MoveForward = true
			ai.StuckCounter = 0
		}
	} else {
		ai.StuckCounter = 0
	}
	ai.LastX = g.PlayerX
	ai.LastY = g.PlayerY
	ai.LastAngle = g.Angle

	// Check for ammo pickup priority
	if g.Ammo < 5 {
		if ai.seekAmmoPickup(g) {
			return
		}
	}

	// Check for portal
	if ai.checkPortal(g) {
		return
	}

	idx, dist := ai.findNearestEnemy(g)

	if idx == -1 {
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
	if diff > 0.1 {
		g.Angle += 0.05
	} else if diff < -0.1 {
		g.Angle -= 0.05
	}

	// Check if can move in facing direction
	canMove := ai.canMoveInDirection(g, g.Angle, 0.1)

	// Movement based on distance
	if dist > 5 && canMove {
		ai.MoveForward = true
	} else if dist < 2 {
		ai.MoveBackward = true
	}

	// If can't move forward but needs to, try turning
	if (dist > 5) && !canMove {
		if diff > 0 {
			g.Angle += 0.1
		} else {
			g.Angle -= 0.1
		}
	}

	// Shoot if aligned and in range
	if math.Abs(diff) < 0.4 && dist < 8 {
		ai.ShootRequested = true
	}
}

func (ai *AIController) canMoveInDirection(g *Game, angle float64, dist float64) bool {
	for checkDist := 0.2; checkDist <= dist; checkDist += 0.2 {
		newX := g.PlayerX + math.Cos(angle)*checkDist
		newY := g.PlayerY + math.Sin(angle)*checkDist
		mapH := GetMapHeight(g.CurrentMap)
		mapW := GetMapWidth(g.CurrentMap)

		if int(newY) < 0 || int(newY) >= mapH || int(newX) < 0 || int(newX) >= mapW {
			return false
		}
		if GetMap(g.CurrentMap)[int(newY)][int(newX)] == 1 {
			return false
		}
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

	if dist < 2.5 {
		targetAngle := math.Atan2(portalY-g.PlayerY, portalX-g.PlayerX)
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
		}
		if ai.canMoveInDirection(g, g.Angle, 0.1) {
			ai.MoveForward = true
		}
		return true
	}
	return false
}

func (ai *AIController) explore(g *Game) {
	// Pick new target periodically or when reached
	if ai.ExploreTimer <= 0 || ai.reachedTarget(g) {
		ai.pickNewExplorationTarget(g)
		ai.ExploreTimer = 250
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
	} else if ai.canMoveInDirection(g, g.Angle, 0.1) {
		ai.MoveForward = true
	} else {
		// Can't move forward - try turning
		g.Angle += 0.2
		ai.pickNewExplorationTarget(g)
		ai.ExploreTimer = 100
	}
}

func (ai *AIController) reachedTarget(g *Game) bool {
	dx := g.PlayerX - ai.ExploreTargetX
	dy := g.PlayerY - ai.ExploreTargetY
	return math.Sqrt(dx*dx+dy*dy) < 1.2
}

func (ai *AIController) pickNewExplorationTarget(g *Game) {
	mapData := GetMap(g.CurrentMap)
	mapW := GetMapWidth(g.CurrentMap)
	mapH := GetMapHeight(g.CurrentMap)

	// Try to find walkable tiles near player
	type spot struct{ x, y int }
	var validSpots []spot

	// Search in expanding radius from player
	for radius := 2; radius <= 8; radius++ {
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				px := int(g.PlayerX) + dx
				py := int(g.PlayerY) + dy
				if py >= 1 && py < mapH-1 && px >= 1 && px < mapW-1 && mapData[py][px] == 0 {
					validSpots = append(validSpots, spot{px, py})
				}
			}
		}
		if len(validSpots) >= 5 {
			break
		}
	}

	if len(validSpots) > 0 {
		// Pick a random spot
		idx := int(ai.ExploreTimer) % len(validSpots)
		ai.ExploreTargetX = float64(validSpots[idx].x) + 0.5
		ai.ExploreTargetY = float64(validSpots[idx].y) + 0.5
	} else {
		// Fallback - try to move in a different direction
		ai.ExploreTargetX = g.PlayerX + math.Cos(g.Angle+math.Pi/2)*3
		ai.ExploreTargetY = g.PlayerY + math.Sin(g.Angle+math.Pi/2)*3
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

func (ai *AIController) seekAmmoPickup(g *Game) bool {
	minDist := math.MaxFloat64
	var targetX, targetY float64
	found := false

	for _, pickup := range g.AmmoPickups {
		if !pickup.Active {
			continue
		}
		dx := pickup.X - g.PlayerX
		dy := pickup.Y - g.PlayerY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < minDist {
			minDist = dist
			targetX = pickup.X
			targetY = pickup.Y
			found = true
		}
	}

	if !found {
		return false
	}

	targetAngle := math.Atan2(targetY-g.PlayerY, targetX-g.PlayerX)
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
	} else if ai.canMoveInDirection(g, g.Angle, 0.1) {
		ai.MoveForward = true
	}

	return true
}
