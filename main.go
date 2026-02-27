package main 

import (
	"log"
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
}

//Update is called every tick 
func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//clear to black 
	for i := range g.pixels{
		g.pixels[i] = 0
	}
	

    x := screenWidth/2
	for y := 0; y < screenHeight; y++{

	 	 idx := (y*screenWidth + x) * 4
     g.pixels[idx+0] = 255 
     g.pixels[idx+1] = 0
     g.pixels[idx+2] = 0 
     g.pixels[idx+3] = 255 
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
	}

	if err := ebiten.RunGame(game); err != nil{
		log.Fatal(err)
	}
}
