package engine

import (
	"encoding/json"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState struct {
	PlayerX          float64   `json:"player_x"`
	PlayerY          float64   `json:"player_y"`
	PlayerAngle      float64   `json:"player_angle"`
	Health           int       `json:"health"`
	Ammo             int       `json:"ammo"`
	Weapon           int       `json:"weapon"`
	EnemyCount       int       `json:"enemy_count"`
	EnemyDistances   []float64 `json:"enemy_distances"`
	EnemyAngles      []float64 `json:"enemy_angles"`
	HasAmmoPickup    bool      `json:"has_ammo_pickup"`
	AmmoPickupDist   float64   `json:"ammo_pickup_dist"`
	HasHealthPickup  bool      `json:"has_health_pickup"`
	HealthPickupDist float64   `json:"health_pickup_dist"`
	PortalDist       float64   `json:"portal_dist"`
	Wave             int       `json:"wave"`
	CurrentMap       int       `json:"current_map"`
}

type PlayerAction struct {
	MoveForward  bool `json:"move_forward"`
	MoveBackward bool `json:"move_backward"`
	TurnLeft     bool `json:"turn_left"`
	TurnRight    bool `json:"turn_right"`
	Shoot        bool `json:"shoot"`
	SwitchWeapon int  `json:"switch_weapon"`
}

type TrainingSample struct {
	State  GameState    `json:"state"`
	Action PlayerAction `json:"action"`
}

type DataCollector struct {
	Enabled    bool
	Samples    []TrainingSample
	FrameCount int
	InputState InputState
	PrevAction PlayerAction
}

type InputState struct {
	ArrowUp    bool
	ArrowDown  bool
	ArrowLeft  bool
	ArrowRight bool
	Space      bool
	Key1       bool
	Key2       bool
	Key3       bool
}

func NewDataCollector() *DataCollector {
	return &DataCollector{
		Enabled: false,
		Samples: make([]TrainingSample, 0),
	}
}

func (dc *DataCollector) Start() {
	dc.Enabled = true
	dc.Samples = dc.Samples[:0]
	dc.FrameCount = 0
	println("Data collection started")
}

func (dc *DataCollector) Stop(filename string) {
	dc.Enabled = false

	if len(dc.Samples) == 0 {
		println("No samples collected")
		return
	}

	data, err := json.MarshalIndent(dc.Samples, "", "  ")
	if err != nil {
		println("Failed to marshal:", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		println("Failed to write:", err)
		return
	}

	println("Saved", len(dc.Samples), "samples to", filename)
}

func (dc *DataCollector) Collect(g *Game) {
	if !dc.Enabled {
		return
	}

	dc.FrameCount++

	// Only collect every 3rd frame to reduce redundancy
	if dc.FrameCount%3 != 0 {
		return
	}

	// Get current input state
	currentInput := dc.getInputState(g)

	// Skip if no input (player idle)
	if !currentInput.ArrowUp && !currentInput.ArrowDown &&
		!currentInput.ArrowLeft && !currentInput.ArrowRight &&
		!currentInput.Space && !currentInput.Key1 && !currentInput.Key2 && !currentInput.Key3 {
		// Still record idle states occasionally
		if dc.FrameCount%30 != 0 {
			return
		}
	}

	// Build game state
	state := dc.extractState(g)

	// Build action
	action := dc.inputToAction(currentInput, g)

	dc.Samples = append(dc.Samples, TrainingSample{
		State:  state,
		Action: action,
	})

	// Auto-save every 1000 samples
	if len(dc.Samples) >= 1000 {
		dc.SaveProgress("training_data_partial.json")
	}

	dc.PrevAction = action
}

func (dc *DataCollector) SaveProgress(filename string) {
	data, _ := json.MarshalIndent(dc.Samples, "", "  ")
	os.WriteFile(filename, data, 0644)
	println("Progress saved:", len(dc.Samples), "samples")
}

func (dc *DataCollector) getInputState(g *Game) InputState {
	return InputState{
		ArrowUp:    ebiten.IsKeyPressed(ebiten.KeyArrowUp),
		ArrowDown:  ebiten.IsKeyPressed(ebiten.KeyArrowDown),
		ArrowLeft:  ebiten.IsKeyPressed(ebiten.KeyArrowLeft),
		ArrowRight: ebiten.IsKeyPressed(ebiten.KeyArrowRight),
		Space:      ebiten.IsKeyPressed(ebiten.KeySpace),
		Key1:       inpututil.IsKeyJustPressed(ebiten.Key1),
		Key2:       inpututil.IsKeyJustPressed(ebiten.Key2),
		Key3:       inpututil.IsKeyJustPressed(ebiten.Key3),
	}
}

func (dc *DataCollector) extractState(g *Game) GameState {
	state := GameState{
		PlayerX:     g.PlayerX,
		PlayerY:     g.PlayerY,
		PlayerAngle: g.Angle,
		Health:      g.Health,
		Ammo:        g.Ammo,
		Weapon:      g.WeaponType,
		Wave:        g.Wave,
		CurrentMap:  g.CurrentMap,
	}

	// Count alive enemies
	aliveCount := 0
	for _, e := range g.Entities {
		if !e.Dead {
			aliveCount++
		}
	}
	state.EnemyCount = aliveCount

	// Get enemy distances and angles
	state.EnemyDistances = make([]float64, 0, 5)
	state.EnemyAngles = make([]float64, 0, 5)
	for _, e := range g.Entities {
		if !e.Dead {
			dx := e.X - g.PlayerX
			dy := e.Y - g.PlayerY
			dist := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx) - g.Angle

			// Normalize angle
			for angle > math.Pi {
				angle -= 2 * math.Pi
			}
			for angle < -math.Pi {
				angle += 2 * math.Pi
			}

			state.EnemyDistances = append(state.EnemyDistances, dist)
			state.EnemyAngles = append(state.EnemyAngles, angle)
		}
	}

	// Check ammo pickups
	minAmmoDist := math.MaxFloat64
	for _, p := range g.AmmoPickups {
		if p.Active {
			dx := p.X - g.PlayerX
			dy := p.Y - g.PlayerY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < minAmmoDist {
				minAmmoDist = dist
			}
		}
	}
	state.HasAmmoPickup = minAmmoDist < math.MaxFloat64
	state.AmmoPickupDist = minAmmoDist

	// Check health pickups
	minHealthDist := math.MaxFloat64
	for _, p := range g.HealthPickups {
		if p.Active {
			dx := p.X - g.PlayerX
			dy := p.Y - g.PlayerY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < minHealthDist {
				minHealthDist = dist
			}
		}
	}
	state.HasHealthPickup = minHealthDist < math.MaxFloat64
	state.HealthPickupDist = minHealthDist

	// Portal distance
	var portalX, portalY float64
	if g.CurrentMap <= 1 {
		portalX, portalY = 13.0, 1.0
	} else {
		portalX, portalY = 28.0, 1.5
	}
	dx := portalX - g.PlayerX
	dy := portalY - g.PlayerY
	state.PortalDist = math.Sqrt(dx*dx + dy*dy)

	return state
}

func (dc *DataCollector) inputToAction(input InputState, g *Game) PlayerAction {
	action := PlayerAction{
		MoveForward:  input.ArrowUp,
		MoveBackward: input.ArrowDown,
		TurnLeft:     input.ArrowLeft,
		TurnRight:    input.ArrowRight,
		Shoot:        input.Space,
	}

	// Determine weapon switch
	if input.Key1 {
		action.SwitchWeapon = 1
	} else if input.Key2 {
		action.SwitchWeapon = 2
	} else if input.Key3 {
		action.SwitchWeapon = 3
	}

	return action
}
