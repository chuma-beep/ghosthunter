package engine

import (
	"fmt"
	"image/color"
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
	cellSize := 2
	mapW := GetMapWidth(g.CurrentMap)
	mapH := GetMapHeight(g.CurrentMap)
	for row := 0; row < mapH; row++ {
		for col := 0; col < mapW; col++ {
			var r, gr, b uint8
			if GetMap(g.CurrentMap)[row][col] == 1 {
				r, gr, b = 255, 255, 255
			} else {
				r, gr, b = 50, 50, 50
			}
			for py := 0; py < cellSize; py++ {
				for px := 0; px < cellSize; px++ {
					x := col*cellSize + px
					y := row*cellSize + py
					if x < ScreenWidth && y < ScreenHeight {
						idx := (y*ScreenWidth + x) * 4
						g.Pixels[idx+0] = r
						g.Pixels[idx+1] = gr
						g.Pixels[idx+2] = b
						g.Pixels[idx+3] = 255
					}
				}
			}
		}
	}
	px := int(g.PlayerX * float64(cellSize))
	py := int(g.PlayerY * float64(cellSize))
	if px >= 0 && px < ScreenWidth && py >= 0 && py < ScreenHeight {
		idx := (py*ScreenWidth + px) * 4
		g.Pixels[idx+0] = 255
		g.Pixels[idx+1] = 255
		g.Pixels[idx+2] = 0
		g.Pixels[idx+3] = 255
	}
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

func (g *Game) DrawGun() {
	gunWidth := 64
	gunHeight := 64
	startX := (ScreenWidth - gunWidth) / 2
	startY := ScreenHeight - gunHeight + g.GunKick
	for y := 0; y < gunHeight; y++ {
		for x := 0; x < gunWidth; x++ {
			texIdx := (y*gunWidth + x) * 4
			a := gunTexture[texIdx+3]
			if a > 128 {
				px := startX + x
				py := startY + y
				if px >= 0 && px < ScreenWidth && py >= 0 && py < ScreenHeight {
					idx := (py*ScreenWidth + px) * 4
					g.Pixels[idx+0] = gunTexture[texIdx+0]
					g.Pixels[idx+1] = gunTexture[texIdx+1]
					g.Pixels[idx+2] = gunTexture[texIdx+2]
					g.Pixels[idx+3] = 255
				}
			}
		}
	}
}

func (g *Game) DrawHUD() {
	barWidth := 100
	barHeight := 2
	barX := 5
	barY := ScreenHeight - 20
	for y := barY; y < barY+barHeight; y++ {
		for x := barX; x < barX+barWidth; x++ {
			idx := (y*ScreenWidth + x) * 4
			g.Pixels[idx+0] = 100
			g.Pixels[idx+1] = 0
			g.Pixels[idx+2] = 0
			g.Pixels[idx+3] = 255
		}
	}
	fillWidth := barWidth * g.Health / 100
	if fillWidth < 0 {
		fillWidth = 0
	}
	for y := barY; y < barY+barHeight; y++ {
		for x := barX; x < barX+fillWidth; x++ {
			idx := (y*ScreenWidth + x) * 4
			g.Pixels[idx+0] = 0
			g.Pixels[idx+1] = 255
			g.Pixels[idx+2] = 0
			g.Pixels[idx+3] = 255
		}
	}
	ammoBarWidth := 100
	ammoBarHeight := 2
	ammoBarX := ScreenWidth - ammoBarWidth - 5
	ammoBarY := ScreenHeight - 20
	for y := ammoBarY; y < ammoBarY+ammoBarHeight; y++ {
		for x := ammoBarX; x < ammoBarX+ammoBarWidth; x++ {
			idx := (y*ScreenWidth + x) * 4
			g.Pixels[idx+0] = 50
			g.Pixels[idx+1] = 50
			g.Pixels[idx+2] = 0
			g.Pixels[idx+3] = 255
		}
	}
	maxAmmo := 50
	ammoFill := ammoBarWidth * g.Ammo / maxAmmo
	if ammoFill > ammoBarWidth {
		ammoFill = ammoBarWidth
	}
	if ammoFill < 0 {
		ammoFill = 0
	}
	for y := ammoBarY; y < ammoBarY+ammoBarHeight; y++ {
		for x := ammoBarX; x < ammoBarX+ammoFill; x++ {
			idx := (y*ScreenWidth + x) * 4
			g.Pixels[idx+0] = 255
			g.Pixels[idx+1] = 165
			g.Pixels[idx+2] = 0
			g.Pixels[idx+3] = 255
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// --- Start screen ---
	if g.GameState == 0 {
		screen.Fill(color.Black)
		ebitenutil.DebugPrint(screen, "GHOST HUNTER\n\nPress SPACE to start\n\nArrow keys to move\nSPACE to shoot\nESC to pause")
		return
	}

	// --- Game over screen ---
	if g.GameState == 2 {
		screen.Fill(color.Black)
		ebitenutil.DebugPrint(screen, fmt.Sprintf(
			"\n\n\n\n\n\n\n          GAME OVER\n\n          Score:      %d\n          Best Score: %d\n          Wave:       %d\n          Map:        %d/5\n\n          Press R to restart",
			g.Score, g.HighScore, g.Wave, g.CurrentMap+1,
		))
		return
	}

	// --- Clear pixels ---
	for i := range g.Pixels {
		g.Pixels[i] = 0
	}

	zBuffer := make([]float64, ScreenWidth)
	fov := math.Pi / 3

	// --- Ray loop ---
	for x := 0; x < ScreenWidth; x++ {
		rayAngle := g.Angle - fov/2 + fov*float64(x)/float64(ScreenWidth)
		var distance float64
		for distance = 0; distance < 32; distance += 0.01 {
			rayX := g.PlayerX + math.Cos(rayAngle)*distance
			rayY := g.PlayerY + math.Sin(rayAngle)*distance
			if int(rayY) < 0 || int(rayY) >= GetMapHeight(g.CurrentMap) ||
				int(rayX) < 0 || int(rayX) >= GetMapWidth(g.CurrentMap) {
				break
			}
			if GetMap(g.CurrentMap)[int(rayY)][int(rayX)] == 1 {
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
		if height == 0 {
			continue
		}
		yStart := (ScreenHeight - height) / 2
		yEnd := (ScreenHeight + height) / 2
		if yStart < 0 {
			yStart = 0
		}
		if yEnd > ScreenHeight {
			yEnd = ScreenHeight
		}

		// Wall
		for y := yStart; y < yEnd; y++ {
			texY := (y - yStart) * TexSize / height
			if texY >= TexSize {
				texY = TexSize - 1
			}
			texIdx := (texY*TexSize + texX) * 4
			idx := (y*ScreenWidth + x) * 4
			var tex []byte
			if g.CurrentMap == 0 {
				tex = WallTexture[:]
			} else {
				tex = WallTexture2[:]
			}
			if idx+3 < len(g.Pixels) && texIdx+3 < len(tex) {
				g.Pixels[idx+0] = uint8(float64(tex[texIdx+0]) / distance)
				g.Pixels[idx+1] = uint8(float64(tex[texIdx+1]) / distance)
				g.Pixels[idx+2] = uint8(float64(tex[texIdx+2]) / distance)
				g.Pixels[idx+3] = 255
			}
		}

		// Ceiling
		var rowDist float64
		for y := 0; y < yStart; y++ {
			denom := ScreenHeight - 2*y
			if denom == 0 {
				continue
			}
			rowDist = float64(ScreenHeight) / float64(denom)
			var ceilTex []byte
			if g.CurrentMap == 0 {
				ceilTex = FloorTexture2[:]
			} else {
				ceilTex = FloorTexture[:]
			}
			floorX := g.PlayerX + rowDist*math.Cos(rayAngle)
			floorY := g.PlayerY + rowDist*math.Sin(rayAngle)
			tX := int(floorX*float64(TexSize)) & (TexSize - 1)
			tY := int(floorY*float64(TexSize)) & (TexSize - 1)
			texIdx := (tY*TexSize + tX) * 4
			idx := (y*ScreenWidth + x) * 4
			if idx+3 >= len(g.Pixels) || texIdx+3 >= len(ceilTex) {
				continue
			}
			g.Pixels[idx+0] = uint8(float64(ceilTex[texIdx+0]) * 0.5)
			g.Pixels[idx+1] = uint8(float64(ceilTex[texIdx+1]) * 0.5)
			g.Pixels[idx+2] = uint8(float64(ceilTex[texIdx+2]) * 0.5)
			g.Pixels[idx+3] = 255
		}

		// Floor
		for y := yEnd; y < ScreenHeight; y++ {
			denom := 2*y - ScreenHeight
			if denom == 0 {
				continue
			}
			rowDist = float64(ScreenHeight) / float64(denom)
			var floorTex []byte
			if g.CurrentMap == 0 {
				floorTex = FloorTexture[:]
			} else {
				floorTex = FloorTexture2[:]
			}
			floorX := g.PlayerX + rowDist*math.Cos(rayAngle)
			floorY := g.PlayerY + rowDist*math.Sin(rayAngle)
			tX := int(floorX*float64(TexSize)) & (TexSize - 1)
			tY := int(floorY*float64(TexSize)) & (TexSize - 1)
			texIdx := (tY*TexSize + tX) * 4
			idx := (y*ScreenWidth + x) * 4
			if idx+3 >= len(g.Pixels) || texIdx+3 >= len(floorTex) {
				continue
			}
			g.Pixels[idx+0] = uint8(float64(floorTex[texIdx+0]) * 0.6)
			g.Pixels[idx+1] = uint8(float64(floorTex[texIdx+1]) * 0.6)
			g.Pixels[idx+2] = uint8(float64(floorTex[texIdx+2]) * 0.6)
			g.Pixels[idx+3] = 255
		}
	} // end ray loop

	// --- Sprites ---
	for _, sprite := range g.Entities {
		dx := sprite.X - g.PlayerX
		dy := sprite.Y - g.PlayerY
		spriteDist := math.Sqrt(dx*dx + dy*dy)
		spriteAngle := math.Atan2(dy, dx) - g.Angle
		for spriteAngle > math.Pi {
			spriteAngle -= 2 * math.Pi
		}
		for spriteAngle < -math.Pi {
			spriteAngle += 2 * math.Pi
		}
		if math.Abs(spriteAngle) >= fov/2 {
			continue
		}
		var tex []byte
		var texSize int
		if sprite.Type == EntityWizard {
			tex = wizardTexture[:]
			texSize = wizardTexSize
		} else {
			tex = spriteTexture[:]
			texSize = spriteTexSize
		}
		spriteScreenX := int((0.5 + spriteAngle/fov) * float64(ScreenWidth))
		spriteHeight := int(float64(ScreenHeight) / spriteDist)
		if spriteHeight == 0 {
			continue
		}
		spriteWidth := spriteHeight
		yStart := (ScreenHeight - spriteHeight) / 2
		yEnd := (ScreenHeight + spriteHeight) / 2
		xStart := spriteScreenX - spriteWidth/2
		xEnd := spriteScreenX + spriteWidth/2
		if yStart < 0 {
			yStart = 0
		}
		if yEnd > ScreenHeight {
			yEnd = ScreenHeight
		}
		if xStart < 0 {
			xStart = 0
		}
		if xEnd > ScreenWidth {
			xEnd = ScreenWidth
		}
		for sx := xStart; sx < xEnd; sx++ {
			if spriteDist >= zBuffer[sx] {
				continue
			}
			tX := (sx - xStart) * texSize / spriteWidth
			for sy := yStart; sy < yEnd; sy++ {
				tY := (sy - yStart) * texSize / spriteHeight
				texIdx := (tY*texSize + tX) * 4
				if texIdx+3 >= len(tex) {
					continue
				}
				a := tex[texIdx+3]
				if a > 128 {
					fade := 1.0
					if sprite.FadeTimer > 0 {
						fade = float64(sprite.FadeTimer) / 20.0
					}
					idx := (sy*ScreenWidth + sx) * 4
					if idx+3 < len(g.Pixels) {
						g.Pixels[idx+0] = uint8(float64(tex[texIdx+0]) * fade)
						g.Pixels[idx+1] = uint8(float64(tex[texIdx+1]) * fade)
						g.Pixels[idx+2] = uint8(float64(tex[texIdx+2]) * fade)
						g.Pixels[idx+3] = 255
					}
				}
			}
		}
	}

	// --- Ammo pickups ---
	for _, pickup := range g.AmmoPickups {
		if !pickup.Active {
			continue
		}
		dx := pickup.X - g.PlayerX
		dy := pickup.Y - g.PlayerY
		spriteDist := math.Sqrt(dx*dx + dy*dy)
		spriteAngle := math.Atan2(dy, dx) - g.Angle
		for spriteAngle > math.Pi {
			spriteAngle -= 2 * math.Pi
		}
		for spriteAngle < -math.Pi {
			spriteAngle += 2 * math.Pi
		}
		if math.Abs(spriteAngle) >= fov/2 {
			continue
		}
		spriteScreenX := int((0.5 + spriteAngle/fov) * float64(ScreenWidth))
		spriteHeight := int(float64(ScreenHeight)/spriteDist) / 2
		if spriteHeight == 0 {
			continue
		}
		spriteWidth := spriteHeight
		yStart := (ScreenHeight - spriteHeight) / 2
		yEnd := (ScreenHeight + spriteHeight) / 2
		xStart := spriteScreenX - spriteWidth/2
		xEnd := spriteScreenX + spriteWidth/2
		if yStart < 0 {
			yStart = 0
		}
		if yEnd > ScreenHeight {
			yEnd = ScreenHeight
		}
		if xStart < 0 {
			xStart = 0
		}
		if xEnd > ScreenWidth {
			xEnd = ScreenWidth
		}
		for sx := xStart; sx < xEnd; sx++ {
			if spriteDist >= zBuffer[sx] {
				continue
			}
			for sy := yStart; sy < yEnd; sy++ {
				idx := (sy*ScreenWidth + sx) * 4
				if idx+3 < len(g.Pixels) {
					g.Pixels[idx+0] = 255
					g.Pixels[idx+1] = 255
					g.Pixels[idx+2] = 0
					g.Pixels[idx+3] = 255
				}
			}
		}
	}

	// --- Portal ---
	var portalX, portalY float64
	if g.CurrentMap <= 1 {
		portalX, portalY = 13.0, 1.0
	} else {
		portalX, portalY = 28.0, 1.5
	}
	{
		dx := portalX - g.PlayerX
		dy := portalY - g.PlayerY
		portalDist := math.Sqrt(dx*dx + dy*dy)
		portalAngle := math.Atan2(dy, dx) - g.Angle
		for portalAngle > math.Pi {
			portalAngle -= 2 * math.Pi
		}
		for portalAngle < -math.Pi {
			portalAngle += 2 * math.Pi
		}
		if math.Abs(portalAngle) < fov/2 && portalDist > 0 {
			spriteScreenX := int((0.5 + portalAngle/fov) * float64(ScreenWidth))
			portalHeight := int(float64(ScreenHeight) / portalDist)
			if portalHeight > 0 {
				portalWidth := portalHeight / 2
				yStart := (ScreenHeight - portalHeight) / 2
				yEnd := (ScreenHeight + portalHeight) / 2
				xStart := spriteScreenX - portalWidth/2
				xEnd := spriteScreenX + portalWidth/2
				if yStart < 0 {
					yStart = 0
				}
				if yEnd > ScreenHeight {
					yEnd = ScreenHeight
				}
				if xStart < 0 {
					xStart = 0
				}
				if xEnd > ScreenWidth {
					xEnd = ScreenWidth
				}
				for sx := xStart; sx < xEnd; sx++ {
					if portalDist >= zBuffer[sx] {
						continue
					}
					for sy := yStart; sy < yEnd; sy++ {
						t := float64(sy-yStart) / float64(yEnd-yStart+1)
						idx := (sy*ScreenWidth + sx) * 4
						if idx+3 < len(g.Pixels) {
							g.Pixels[idx+0] = uint8(150 + 50*t)
							g.Pixels[idx+1] = 0
							g.Pixels[idx+2] = uint8(200 + 55*t)
							g.Pixels[idx+3] = 255
						}
					}
				}
			}
		}
	}

	// --- Crosshair ---
	cx := ScreenWidth / 2
	cy := ScreenHeight / 2
	for i := -5; i <= 5; i++ {
		idx := (cy*ScreenWidth + (cx + i)) * 4
		if idx+3 < len(g.Pixels) {
			g.Pixels[idx+0] = 255
			g.Pixels[idx+1] = 255
			g.Pixels[idx+2] = 255
			g.Pixels[idx+3] = 255
		}
		idx = ((cy+i)*ScreenWidth + cx) * 4
		if idx+3 < len(g.Pixels) {
			g.Pixels[idx+0] = 255
			g.Pixels[idx+1] = 255
			g.Pixels[idx+2] = 255
			g.Pixels[idx+3] = 255
		}
	}

	// --- Damage flash ---
	if g.DamageFlash > 0 {
		for i := 0; i < len(g.Pixels); i += 4 {
			g.Pixels[i+0] = uint8(min(int(g.Pixels[i+0])+100, 255))
			g.Pixels[i+3] = 255
		}
		g.DamageFlash--
	}

	g.DrawGun()
	g.DrawHUD()
	g.DrawMinimap()
	screen.ReplacePixels(g.Pixels)

	// --- Overlays ---
	if g.Paused {
		ebitenutil.DebugPrint(screen, "\n\n\n\n\n\n\n          PAUSED\n\n          Press ESC to resume\n          Press R to restart")
		return
	}
	if g.WaveTransition > 0 {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n\n\n\n          Wave %d incoming!", g.Wave))
	}
	if g.LevelNameTimer > 0 {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n          %s", MapNames[g.CurrentMap]))
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Wave: %d  Score: %d  Best: %d  Health: %d  Ammo: %d", g.Wave, g.Score, g.HighScore, g.Health, g.Ammo))
}
