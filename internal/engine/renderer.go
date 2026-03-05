package engine

import (
     "fmt"
	"math"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)


func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}


func (g *Game) DrawMinimap() {
    cellSize := 4

    for row := 0; row < 16; row++ {
        for col := 0; col < 16; col++ {
            var r, gr, b uint8
            if WorldMap[row][col] == 1 {
                r, gr, b = 255, 255, 255 // wall
            } else {
                r, gr, b = 50, 50, 50 // empty
            }

            for py := 0; py < cellSize; py++ {
                for px := 0; px < cellSize; px++ {
                    x := col*cellSize + px
                    y := row*cellSize + py
                    idx := (y*ScreenWidth + x) * 4
                    g.Pixels[idx+0] = r
                    g.Pixels[idx+1] = gr
                    g.Pixels[idx+2] = b
                    g.Pixels[idx+3] = 255
                }
            }
        }
    }


	// draw player
px := int(g.PlayerX * float64(cellSize))
py := int(g.PlayerY * float64(cellSize))

idx := (py*ScreenWidth + px) * 4
g.Pixels[idx+0] = 255
g.Pixels[idx+1] = 255
g.Pixels[idx+2] = 0
g.Pixels[idx+3] = 255

// draw direction line
for i := 0; i < 10; i++ {
    lx := int(g.PlayerX*float64(cellSize) + math.Cos(g.Angle)*float64(i))
    ly := int(g.PlayerY*float64(cellSize) + math.Sin(g.Angle)*float64(i))
    if lx >= 0 && lx < ScreenWidth && ly >= 0 && ly < ScreenHeight {
        idx := (ly*ScreenWidth + lx) * 4
        g.Pixels[idx+0] = 255
        g.Pixels[idx+1] = 255
        g.Pixels[idx+2] = 0
        g.Pixels[idx+3] = 255
    }
}




}




func (g *Game) Draw(screen *ebiten.Image) {
    // clear to black
    for i := range g.Pixels {
        g.Pixels[i] = 0
    }

    zBuffer := make([]float64, ScreenWidth)
    fov := math.Pi / 3

    for x := 0; x < ScreenWidth; x++ {
        rayAngle := g.Angle - fov/2 + fov*float64(x)/float64(ScreenWidth)

        var distance float64
        for distance = 0; distance < 20; distance += 0.01 {
            rayX := g.PlayerX + math.Cos(rayAngle)*distance
            rayY := g.PlayerY + math.Sin(rayAngle)*distance
            if WorldMap[int(rayY)][int(rayX)] == 1 {
                break
            }
        }

        zBuffer[x] = distance

        hitX := g.PlayerX + math.Cos(rayAngle)*distance
        hitY := g.PlayerY + math.Sin(rayAngle)*distance

        var wallX float64
        if math.Abs(math.Cos(rayAngle)) > math.Abs(math.Sin(rayAngle)) {
            wallX = hitY - math.Floor(hitY)
        } else {
            wallX = hitX - math.Floor(hitX)
        }

        texX := int(wallX * float64(TexSize))

        height := int(float64(ScreenHeight) / distance)
        yStart := (ScreenHeight - height) / 2
        yEnd := (ScreenHeight + height) / 2

        if yStart < 0 { yStart = 0 }
        if yEnd > ScreenHeight { yEnd = ScreenHeight }

        for y := yStart; y < yEnd; y++ {
            texY := (y - yStart) * TexSize / height
            if texY >= TexSize { texY = TexSize - 1 }
            texIdx := (texY*TexSize + texX) * 4
            idx := (y*ScreenWidth + x) * 4
            g.Pixels[idx+0] = uint8(float64(WallTexture[texIdx+0]) / distance)
            g.Pixels[idx+1] = uint8(float64(WallTexture[texIdx+1]) / distance)
            g.Pixels[idx+2] = uint8(float64(WallTexture[texIdx+2]) / distance)
            g.Pixels[idx+3] = 255
        }

        for y := 0; y < yStart; y++ {
            idx := (y*ScreenWidth + x) * 4
            g.Pixels[idx+0] = 50
            g.Pixels[idx+1] = 50
            g.Pixels[idx+2] = 139
            g.Pixels[idx+3] = 255
        }

        for y := yEnd; y < ScreenHeight; y++ {
            idx := (y*ScreenWidth + x) * 4
            g.Pixels[idx+0] = 139
            g.Pixels[idx+1] = 50
            g.Pixels[idx+2] = 50
            g.Pixels[idx+3] = 255
        }
    }

// sprite rendering
for _, sprite := range g.Sprites {
    dx := sprite.X - g.PlayerX
    dy := sprite.Y - g.PlayerY
    spriteDist := math.Sqrt(dx*dx + dy*dy)
    spriteAngle := math.Atan2(dy, dx) - g.Angle

    for spriteAngle > math.Pi { spriteAngle -= 2 * math.Pi }
    for spriteAngle < -math.Pi { spriteAngle += 2 * math.Pi }

    if math.Abs(spriteAngle) < fov/2 {
        spriteScreenX := int((0.5 + spriteAngle/fov) * float64(ScreenWidth))
        spriteHeight := int(float64(ScreenHeight) / spriteDist)
        spriteWidth := spriteHeight

        yStart := (ScreenHeight - spriteHeight) / 2
        yEnd := (ScreenHeight + spriteHeight) / 2
        xStart := spriteScreenX - spriteWidth/2
        xEnd := spriteScreenX + spriteWidth/2

        if yStart < 0 { yStart = 0 }
        if yEnd > ScreenHeight { yEnd = ScreenHeight }
        if xStart < 0 { xStart = 0 }
        if xEnd > ScreenWidth { xEnd = ScreenWidth }

        for sx := xStart; sx < xEnd; sx++ {
            if spriteDist < zBuffer[sx] {
                texX := (sx - xStart) * spriteTexSize / spriteWidth
                for sy := yStart; sy < yEnd; sy++ {
                    texY := (sy - yStart) * spriteTexSize / spriteHeight
                    texIdx := (texY*spriteTexSize + texX) * 4
                    a := spriteTexture[texIdx+3]
                    if a > 128 {
                        idx := (sy*ScreenWidth + sx) * 4
                        g.Pixels[idx+0] = spriteTexture[texIdx+0]
                        g.Pixels[idx+1] = spriteTexture[texIdx+1]
                        g.Pixels[idx+2] = spriteTexture[texIdx+2]
                        g.Pixels[idx+3] = 255
                    }
                }
            }
        }
    }
}

   // draw crosshair
cx := ScreenWidth / 2
cy := ScreenHeight / 2
for i := -5; i <= 5; i++ {
    idx := (cy*ScreenWidth + (cx + i)) * 4
    g.Pixels[idx+0] = 255
    g.Pixels[idx+1] = 255
    g.Pixels[idx+2] = 255
    g.Pixels[idx+3] = 255

    idx = ((cy+i)*ScreenWidth + cx) * 4
    g.Pixels[idx+0] = 255
    g.Pixels[idx+1] = 255
    g.Pixels[idx+2] = 255
    g.Pixels[idx+3] = 255
}




// screen flash when taking damage
if g.DamageFlash > 0 {
    for i := 0; i < len(g.Pixels); i += 4 {
        g.Pixels[i+0] = uint8(min(int(g.Pixels[i+0])+100, 255))
        g.Pixels[i+3] = 255
    }
    g.DamageFlash--
}

  
    g.DrawMinimap()
    screen.ReplacePixels(g.Pixels)
    
	if g.Health <= 0{
		ebitenutil.DebugPrint(screen, "GAME OVER") 
	}else{

		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d Health: %d", g.Score, g.Health))
  }
}
