package engine

import "math"

type AIController struct {
	Enabled        bool
	ShootRequested bool
}

func NewAIController() *AIController {
	return &AIController{Enabled: false}
}

func (ai *AIController) Update(g *Game) {
	if !ai.Enabled {
		return
	}

	idx, dist := ai.findNearestEnemy(g)
	if idx == -1 {
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

	if diff > 0.05 {
		g.Angle += 0.05
	} else if diff < -0.05 {
		g.Angle -= 0.05
	}

	if math.Abs(diff) < 0.2 && dist < 8 {
		ai.ShootRequested = true
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
