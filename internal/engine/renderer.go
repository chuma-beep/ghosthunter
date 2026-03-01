package engine

import (
    "math"
    "github.com/hajimehoshi/ebiten/v2"
)

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
    dx := g.SpriteX - g.PlayerX
    dy := g.SpriteY - g.PlayerY
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
                for sy := yStart; sy < yEnd; sy++ {
                    idx := (sy*ScreenWidth + sx) * 4
                    g.Pixels[idx+0] = 255
                    g.Pixels[idx+1] = 0
                    g.Pixels[idx+2] = 0
                    g.Pixels[idx+3] = 255
                }
            }
        }
    }

    screen.ReplacePixels(g.Pixels)
}
