package engine

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
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

// draw guns??
func (g *Game) DrawGun(screen *ebiten.Image) {
	state := WeaponStates[g.WeaponStateID]
	weaponIdx := state.Weapon
	frameIdx := state.Frame

	if weaponIdx >= len(weaponAnimations) {
		weaponIdx = 0
	}
	frames := weaponAnimations[weaponIdx]
	if len(frames) == 0 || frameIdx >= len(frames) {
		return
	}
	img := frames[frameIdx]
	bounds := img.Bounds()
	dstW := ScreenWidth / 2
	dstH := ScreenHeight / 2
	startX := float64(ScreenWidth-dstW) / 2
	startY := float64(ScreenHeight - dstH + g.GunKick)
	scaleX := float64(dstW) / float64(bounds.Max.X)
	scaleY := float64(dstH) / float64(bounds.Max.Y)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(startX, startY)
	screen.DrawImage(img, op)
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

	if len(g.Pixels) != ScreenWidth*ScreenHeight*4 {
		g.Pixels = make([]byte, ScreenWidth*ScreenHeight*4)
	}

	// --- Start screen ---
	if g.GameState == 0 {
		screen.Fill(color.Black)
		fontFace := basicfont.Face7x13
		startLines := []string{
			"GHOST HUNTER",
			"",
			"Press SPACE to start",
			"",
			"Arrow keys to move",
			"SPACE to shoot",
			"ESC to pause",
		}
		startY := ScreenHeight/2 - len(startLines)*fontFace.Height/2
		for i, line := range startLines {
			x := (ScreenWidth - len(line)*fontFace.Width) / 2
			y := startY + i*fontFace.Height
			text.Draw(screen, line, fontFace, x, y, color.White)
		}
		return
	}

	// --- Game over screen ---
	if g.GameState == 2 {
		screen.Fill(color.Black)
		fontFace := basicfont.Face7x13
		gameOverLines := []string{
			"GAME OVER",
			fmt.Sprintf("Score:     %d", g.Score),
			fmt.Sprintf("Best:     %d", g.HighScore),
			fmt.Sprintf("Wave:      %d", g.Wave),
			fmt.Sprintf("Map:       %d/5", g.CurrentMap+1),
			"",
			"Press R to restart",
		}
		startY := ScreenHeight/2 - len(gameOverLines)*fontFace.Height/2
		for i, line := range gameOverLines {
			x := (ScreenWidth - len(line)*fontFace.Width) / 2
			y := startY + i*fontFace.Height
			text.Draw(screen, line, fontFace, x, y, color.White)
		}
		return
	}

	// --- Clear pixels ---
	for i := range g.Pixels {
		g.Pixels[i] = 0
	}

	// shake offset
	shakeX := 0
	shakeY := 0
	if g.ScreenShake > 0 {
		shakeX = (g.ScreenShake % 3) - 1
		shakeY = (g.ScreenShake % 2) - 1
		g.ScreenShake--
	}

	zBuffer := make([]float64, ScreenWidth)
	fov := math.Pi / 3

	// --- Ray loop ---
	for x := 0; x < ScreenWidth; x++ {
		rayAngle := g.Angle + float64(shakeX)*0.003 - fov/2 + fov*float64(x)/float64(ScreenWidth)
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
			// idx := (y*ScreenWidth + x) * 4

			sy := y + shakeY
			if sy < 0 || sy >= ScreenHeight {
				continue
			}
			idx := (sy*ScreenWidth + x) * 4

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
		var frameCount int
		switch sprite.Type {
		case EntityWizard:
			tex = wizardTexture[:]
			texSize = wizardTexSize
			frameCount = 1
		case EntityDemon:
			tex = demonTexture
			texSize = 64
			frameCount = demonFrames
		case EntityWraith:
			tex = wraithTexture
			texSize = 64
			frameCount = wraithFrames
		case EntityReaper:
			tex = reaperTexture
			texSize = 64
			frameCount = reaperFrames
		default:
			tex = spriteTexture[:]
			texSize = spriteTexSize
			frameCount = 1
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
			tX := (sx-xStart)*texSize/spriteWidth + sprite.Frame*texSize
			for sy := yStart; sy < yEnd; sy++ {
				tY := (sy - yStart) * texSize / spriteHeight
				texIdx := (tY*texSize*frameCount + tX) * 4
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

	// health pickups
	for _, pickup := range g.HealthPickups {
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
				if idx+3 >= len(g.Pixels) {
					continue
				}
				// red cross pattern
				cx := (xStart + xEnd) / 2
				cy := (yStart + yEnd) / 2
				onCross := (sx == cx) || (sy == cy)
				if onCross {
					g.Pixels[idx+0] = 255
					g.Pixels[idx+1] = 0
					g.Pixels[idx+2] = 0
					g.Pixels[idx+3] = 255
				} else {
					g.Pixels[idx+0] = 200
					g.Pixels[idx+1] = 200
					g.Pixels[idx+2] = 200
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

	g.DrawHUD()
	g.DrawMinimap()
	screen.ReplacePixels(g.Pixels)
	g.DrawGun(screen)

	// --- Overlays ---
	if g.Paused {
		// Darken the existing pixel buffer (multiply RGB by 0.3, alpha unchanged)
		for i := 0; i < len(g.Pixels); i += 4 {
			r := float64(g.Pixels[i]) * 0.3
			gn := float64(g.Pixels[i+1]) * 0.3
			b := float64(g.Pixels[i+2]) * 0.3
			g.Pixels[i] = uint8(r)
			g.Pixels[i+1] = uint8(gn)
			g.Pixels[i+2] = uint8(b)
		}
		screen.ReplacePixels(g.Pixels)

		fontFace := basicfont.Face7x13
		menuW, menuH := 200, 150
		menuImg := ebiten.NewImage(menuW, menuH)
		menuImg.Fill(color.RGBA{0, 0, 0, 150})

		originX := (ScreenWidth - menuW) / 2
		originY := (ScreenHeight - menuH) / 2

		if g.ShowControls {
			controls := []string{
				"CONTROLS",
				"",
				"Arrow Keys : Move",
				"Space      : Shoot",
				"Escape     : Pause",
				"A          : Toggle AI",
				"",
				"Press ESC to go back",
			}
			startY := 20
			for i, line := range controls {
				x := (menuW - len(line)*fontFace.Width) / 2
				y := startY + i*fontFace.Height
				text.Draw(menuImg, line, fontFace, x, y, color.White)
			}
		} else {
			menuItems := []string{"Resume", "Controls", "Quit"}
			startY := 30
			for i, item := range menuItems {
				y := startY + i*fontFace.Height*2

				if i == g.PauseMenuSelection {
					text.Draw(menuImg, ">", fontFace, 30, y, color.RGBA{255, 255, 0, 255})
					text.Draw(menuImg, item, fontFace, 50, y, color.RGBA{255, 255, 0, 255})
				} else {
					text.Draw(menuImg, item, fontFace, 50, y, color.White)
				}
			}
			hint := "Arrow Keys + Enter | ESC"
			hintX := (menuW - len(hint)*fontFace.Width) / 2
			text.Draw(menuImg, hint, fontFace, hintX, menuH-10, color.RGBA{180, 180, 180, 255})
		}

		ops := &ebiten.DrawImageOptions{}
		ops.GeoM.Translate(float64(originX), float64(originY))
		screen.DrawImage(menuImg, ops)
		return
	}
	if g.WaveTransition > 0 {
		fontFace := basicfont.Face7x13
		waveText := fmt.Sprintf("Wave %d incoming!", g.Wave)
		x := (ScreenWidth - len(waveText)*fontFace.Width) / 2
		y := ScreenHeight / 2
		text.Draw(screen, waveText, fontFace, x, y, color.RGBA{255, 255, 0, 255})
	}
	if g.LevelNameTimer > 0 {
		fontFace := basicfont.Face7x13
		mapName := MapNames[g.CurrentMap]
		x := (ScreenWidth - len(mapName)*fontFace.Width) / 2
		y := ScreenHeight/2 - 40
		text.Draw(screen, mapName, fontFace, x, y, color.White)
	}

	// weapon name
	var weaponName string
	switch g.WeaponType {
	case 0:
		weaponName = "PST"
	case 1:
		weaponName = "SG"
	case 2:
		weaponName = "MG"
	}
	hudText := fmt.Sprintf("W%d S%d H%d A%d+%s", g.Wave, g.Score, g.Health, g.Ammo, weaponName)
	x := ScreenWidth - len(hudText)*basicfont.Face7x13.Width - 20
	text.Draw(screen, hudText, basicfont.Face7x13, x, 15, color.RGBA{120, 120, 120, 255})
}
