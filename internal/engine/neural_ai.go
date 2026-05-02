package engine

import (
	"encoding/json"
	"math"
	"os"

	_ "embed"
)

//go:embed model_weights.json
var embeddedModel []byte

type FlatModel struct {
	InputSize  int       `json:"input_size"`
	OutputSize int       `json:"output_size"`
	HiddenSize int       `json:"hidden_size"`
	W1         []float64 `json:"W1"`
	B1         []float64 `json:"b1"`
	W2         []float64 `json:"W2"`
	B2         []float64 `json:"b2"`
	W3         []float64 `json:"W3"`
	B3         []float64 `json:"b3"`
	Mean       []float64 `json:"mean"`
	Std        []float64 `json:"std"`
}

type NeuralAI struct {
	Enabled bool
	model   FlatModel

	ShootRequested bool
	MoveForward    bool
	MoveBackward   bool
}

func NewNeuralAI() *NeuralAI {
	return &NeuralAI{Enabled: false}
}

func (n *NeuralAI) LoadModel(filename string) error {
	var data []byte
	var err error

	// Try embedded first (for WASM)
	if len(embeddedModel) > 0 {
		data = embeddedModel
		println("Using embedded model")
	} else {
		// Fall back to file loading (for desktop)
		data, err = os.ReadFile(filename)
		if err != nil {
			return err
		}
		println("Using file model")
	}

	err = json.Unmarshal(data, &n.model)
	if err != nil {
		return err
	}

	println("Neural AI loaded:", n.model.InputSize, "->", n.model.OutputSize, "hidden:", n.model.HiddenSize)
	return nil
}

func (n *NeuralAI) forward(input []float64) []float64 {
	inputSize := n.model.InputSize
	hiddenSize := n.model.HiddenSize
	outputSize := n.model.OutputSize
	W1 := n.model.W1
	W2 := n.model.W2
	W3 := n.model.W3
	b1 := n.model.B1
	b2 := n.model.B2
	b3 := n.model.B3

	// Layer 1: input -> hidden1
	hidden1 := make([]float64, hiddenSize)
	for j := 0; j < hiddenSize; j++ {
		sum := b1[j]
		for i := 0; i < inputSize; i++ {
			sum += input[i] * W1[j*inputSize+i]
		}
		hidden1[j] = relu(sum)
	}

	// Layer 2: hidden1 -> hidden2
	hidden2 := make([]float64, hiddenSize)
	for j := 0; j < hiddenSize; j++ {
		sum := b2[j]
		for i := 0; i < hiddenSize; i++ {
			sum += hidden1[i] * W2[j*hiddenSize+i]
		}
		hidden2[j] = relu(sum)
	}

	// Layer 3: hidden2 -> output
	output := make([]float64, outputSize)
	for j := 0; j < outputSize; j++ {
		sum := b3[j]
		for i := 0; i < hiddenSize; i++ {
			sum += hidden2[i] * W3[j*hiddenSize+i]
		}
		output[j] = sigmoid(sum)
	}

	return output
}

func relu(x float64) float64 {
	if x > 0 {
		return x
	}
	return 0
}

func sigmoid(x float64) float64 {
	if x > 500 {
		return 1
	}
	if x < -500 {
		return 0
	}
	return 1 / (1 + math.Exp(-x))
}

func (n *NeuralAI) Update(g *Game) {
	if !n.Enabled {
		return
	}

	features := n.extractFeatures(g)

	normalized := make([]float64, len(features))
	for i := range features {
		if i < len(n.model.Std) && n.model.Std[i] != 0 {
			normalized[i] = (features[i] - n.model.Mean[i]) / n.model.Std[i]
		} else {
			normalized[i] = features[i]
		}
	}

	output := n.forward(normalized)

	threshold := 0.35

	n.MoveForward = output[0] > threshold
	n.MoveBackward = output[1] > threshold

	turnStrength := 0.05
	if output[2] > threshold {
		g.Angle -= turnStrength
	}
	if output[3] > threshold {
		g.Angle += turnStrength
	}

	n.ShootRequested = output[4] > threshold
}

func (n *NeuralAI) extractFeatures(g *Game) []float64 {
	features := make([]float64, 22)

	features[0] = g.PlayerX / 32.0
	features[1] = g.PlayerY / 32.0
	features[2] = g.Angle / (2 * math.Pi)
	features[3] = float64(g.Health) / 100.0
	features[4] = float64(g.Ammo) / 50.0
	features[5] = float64(g.WeaponType) / 2.0

	enemyCount := 0
	var enemyDists []float64
	var enemyAngles []float64
	for _, e := range g.Entities {
		if !e.Dead {
			enemyCount++
			dx := e.X - g.PlayerX
			dy := e.Y - g.PlayerY
			dist := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx) - g.Angle
			for angle > math.Pi {
				angle -= 2 * math.Pi
			}
			for angle < -math.Pi {
				angle += 2 * math.Pi
			}
			enemyDists = append(enemyDists, dist)
			enemyAngles = append(enemyAngles, angle)
		}
	}

	features[6] = float64(enemyCount) / 20.0
	features[7] = float64(g.Wave) / 5.0
	features[8] = float64(g.CurrentMap) / 4.0

	var portalX, portalY float64
	if g.CurrentMap <= 1 {
		portalX, portalY = 13.0, 1.0
	} else {
		portalX, portalY = 28.0, 1.5
	}
	dx := portalX - g.PlayerX
	dy := portalY - g.PlayerY
	features[9] = math.Sqrt(dx*dx+dy*dy) / 20.0

	minAmmoDist := 99.0
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
	features[10] = minAmmoDist

	minHealthDist := 99.0
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
	features[11] = minHealthDist

	for i := 0; i < 5; i++ {
		if i < len(enemyDists) {
			features[12+i*2] = enemyDists[i] / 15.0
			features[12+i*2+1] = enemyAngles[i] / math.Pi
		} else {
			features[12+i*2] = 1.0
			features[12+i*2+1] = 0.0
		}
	}

	return features
}
