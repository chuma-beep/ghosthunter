package main 

import (
	"log"
	"math"
    // "image"
	// "image/color"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth = 320
	screenHeight = 200 
)

//types for game sruct here 
type Game struct{
	pixels []byte
	playerX float64
	playerY float64
	angle float64
}

//world map 
var worldMap = [8][8]int{
	{1, 1, 1, 1, 1, 1, 1, 1,},
	{1, 0, 0, 0, 0, 0, 0, 1,},
	{1, 0, 0, 0, 0, 0, 0, 1,},
	{1, 0, 0, 0, 0, 0, 0, 1,},
	{1, 0, 0, 0, 0, 0, 0, 1,},
	{1, 0, 0, 0, 0, 0, 0, 1,},
	{1, 0, 0, 0, 0, 0, 0, 1,},
	{1, 1, 1, 1, 1, 1, 1, 1,},
}


//Update is called every tick 
func (g *Game) Update() error {
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft){
		g.angle -= 0.03
	} 
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight){
		g.angle += 0.03
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp){
		g.playerX += math.Cos(g.angle) * 0.05
		g.playerY += math.Sin(g.angle) * 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown){
		g.playerX -= math.Cos(g.angle) * 0.05
		g.playerY -= math.Sin(g.angle) * 0.05
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image){
	//clear to black 
	for i := range g.pixels{
		g.pixels[i] = 0
	}

	fov := math.Pi / 3 

    for x := 0; x < screenWidth; x++ {
		//angle of this ray 
		rayAngle := g.angle - fov/2 + fov+float64(x)/float64(screenWidth)
	
        //cast the ray
		var distance float64
		for distance = 0; distance < 20; distance += 0.01{
             rayX := g.playerX + math.Cos(rayAngle)*distance
			 rayY := g.playerY + math.Sin(rayAngle)*distance


			 if worldMap[int(rayY)][int(rayX)] == 1 {
				 break
			 }
		}


		//draw column 
		height := int(float64(screenHeight) / distance)
		yStart := (screenHeight - height) / 2
		yEnd := (screenHeight + height) / 2

		if yStart < 0{
			yStart = 0
		}
		if yEnd > screenHeight{
			yEnd  = screenHeight 
		}

		for y := yStart; y < yEnd; y++{
			idx := (y*screenWidth + x) * 4 
			g.pixels[idx+0] = 255
			g.pixels[idx+1] = 255 
			g.pixels[idx+2] = 255
			g.pixels[idx+3] = 255
		}

       for y := 0; y < yStart; y++{
		   idx := (y*screenWidth + x) * 4 
		   g.pixels[idx+0] = 50
		   g.pixels[idx+1] = 50
		   g.pixels[idx+2] = 139
		   g.pixels[idx+3] = 255
	   }
          
	   for y := yEnd; y < screenHeight; y++{
		   idx := (y*screenWidth + x) * 4 
	       g.pixels[idx+0] = 139
		   g.pixels[idx+1] = 50
		   g.pixels[idx+2] = 50
		   g.pixels[idx+3] = 255
	   }
	}
	
      screen.ReplacePixels(g.pixels)
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight 
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("doom-go")
     
	game := &Game{
		pixels: make([]byte, screenWidth*screenHeight*4),
	    playerX: 2.0,
		playerY: 2.0,
        angle: 0.0,
	}

	if err := ebiten.RunGame(game); err != nil{
		log.Fatal(err)
	}
}
