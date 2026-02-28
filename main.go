package main 

import (
	"log"
	"math"
    "github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth = 320
	screenHeight = 200 
)

const texSize = 64 

 var wallTexture [texSize * texSize * 4]byte 

func init(){
	for y := 0; y < texSize;y++ {
		for x := 0; x < texSize; x++ {
			idx := (y*texSize + x) * 4
			if (x/8+y/8)%2 == 0{
				wallTexture[idx+0] = 200    
				wallTexture[idx+1] = 100
				wallTexture[idx+2] =50
				wallTexture[idx+3] = 255 
			}else {
				wallTexture[idx+0] = 100 
				wallTexture[idx+1] = 50
				wallTexture[idx+2] = 25 
				wallTexture[idx+3] = 255
			}
		}
	}
}





//types for game sruct here 
type Game struct{
	pixels []byte
	playerX float64
	playerY float64
	angle float64
}


//worldMap
var worldMap = [16][16]int{
    {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1},
    {1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
    {1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
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
		newX := g.playerX + math.Cos(g.angle) * 0.05
		newY := g.playerY + math.Sin(g.angle) * 0.05
		if worldMap[int(newY)][int(newX)] == 0{
			g.playerX = newX 
			g.playerY = newY 
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown){
		newX := g.playerX - math.Cos(g.angle) * 0.05
		newY := g.playerY - math.Sin(g.angle) * 0.05
	    if worldMap[int(newY)][int(newX)] == 0 {
			g.playerX = newX 
			g.playerY = newY 
		}
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

		brightness := 255.0 / distance
		if brightness > 255 {
			brightness = 255
		}

           //calculate where on the wall the ray hit 
        hitX := g.playerX + math.Cos(rayAngle)*distance
		hitY := g.playerY + math.Sin(rayAngle)*distance

		var wallX float64
		if math.Abs(math.Cos(rayAngle)) > (math.Sin(rayAngle)) {
			wallX = hitY - math.Floor(hitY)
		}else {
			wallX = hitX - math.Floor(hitX)
		}
       
        texX := int(wallX * float64(texSize)) 
      

		//draw255
		height := int(float64(screenHeight) / distance)
		yStart := (screenHeight - height) / 2
		yEnd := (screenHeight + height) / 2

		if yStart < 0{
			yStart = 0
		}
		if yEnd > screenHeight{
			yEnd  = screenHeight 
		}


             // wall column 
 for y := yStart; y < yEnd; y++ {
    texY := (y - yStart) * texSize / height
    if texY >= texSize {
        texY = texSize - 1
    }
    texIdx := (texY*texSize + texX) * 4
    idx := (y*screenWidth + x) * 4
    g.pixels[idx+0] = uint8(float64(wallTexture[texIdx+0]) / distance)
    g.pixels[idx+1] = uint8(float64(wallTexture[texIdx+1]) / distance)
    g.pixels[idx+2] = uint8(float64(wallTexture[texIdx+2]) / distance)
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
	    playerX: 8.0,
		playerY: 8.0,
        angle: 0.0,
	}

	if err := ebiten.RunGame(game); err != nil{
		log.Fatal(err)
	}
}
