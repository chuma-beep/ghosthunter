package engine

import (
	  "fmt"
    "image"
    _ "image/png"
    "os"
    "path/filepath"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)



const TexSize = 64

var WallTexture [TexSize * TexSize * 4]byte


func LoadTexture(path string) {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    for y := 0; y < TexSize; y++ {
        for x := 0; x < TexSize; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            idx := (y*TexSize + x) * 4
            WallTexture[idx+0] = byte(r >> 8)
            WallTexture[idx+1] = byte(g >> 8)
            WallTexture[idx+2] = byte(b >> 8)
            WallTexture[idx+3] = byte(a >> 8)
        }
    }
}



//WallTexture  second room 
var WallTexture2 [spriteTexSize * spriteTexSize * 4]byte 


   
func LoadTexture2(path string){
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()


	img, _, err := image.Decode(f)
	if err != nil{
		panic(err)
	}
     
	for y := 0; y < TexSize; y++ {
		for x := 0; x < TexSize; x++ {
			r, g, b, a := img.At(x,y).RGBA()
			idx := (y*TexSize + x) * 4
			WallTexture2[idx+0] = byte(r >> 8)
			WallTexture2[idx+1] = byte(g >> 8)
			WallTexture2[idx+2] = byte(b >> 8)
			WallTexture2[idx+3] = byte(a >> 8)
		}
	}


}




// floor FloorTexture 
var FloorTexture [TexSize * TexSize * 4]byte
var FloorTexture2 [TexSize * TexSize * 4]byte

func LoadFloor(path string, tex *[TexSize * TexSize * 4]byte) {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    for y := 0; y < TexSize; y++ {
        for x := 0; x < TexSize; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            idx := (y*TexSize + x) * 4
            tex[idx+0] = byte(r >> 8)
            tex[idx+1] = byte(g >> 8)
            tex[idx+2] = byte(b >> 8)
            tex[idx+3] = byte(a >> 8)
        }
    }
}


//loadImag func 
func loadImage(path string) (image.Image, string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, "", err
    }
    defer f.Close()
    img, format, err := image.Decode(f)
    return img, format, err
}


//loaded wea
var weaponAnimations [3][]*ebiten.Image

func LoadWeaponAnimations() {
    folders := []string{
        "assets/gun_pistol",
        "assets/gun_pistol",
        "assets/gun_machinegun",
    }
    for i, folder := range folders {
        files, err := os.ReadDir(folder)
        if err != nil {
            continue
        }
        for _, file := range files {
            if filepath.Ext(file.Name()) == ".png" {
                img, _, err := ebitenutil.NewImageFromFile(filepath.Join(folder, file.Name()))
                if err != nil {
                    continue
                }
                weaponAnimations[i] = append(weaponAnimations[i], img)
            }
        }
        fmt.Println("weapon", i, "loaded frames:", len(weaponAnimations[i]))
    }
}





const spriteTexSize = 64 

var spriteTexture [spriteTexSize * spriteTexSize * 4]byte 

func LoadSpriteTexture(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
    defer f.Close()

    img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

     
	for y := 0; y < spriteTexSize; y++ {
		for x := 0; x < spriteTexSize; x++{
			r, g, b, a := img.At(x, y).RGBA()
			idx := (y*spriteTexSize + x) * 4
			spriteTexture[idx+0] = byte(r >> 8)
			spriteTexture[idx+1] = byte(g >> 8)
			spriteTexture[idx+2] = byte(b >> 8)
			spriteTexture[idx+3] = byte(a >> 8)
		}
	}

}



//wizard entity 
var wizardTexture [64 * 64 * 4]byte
var wizardTexSize = 64

func LoadWizard(path string) {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    for y := 0; y < 64; y++ {
        for x := 0; x < 64; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            idx := (y*64 + x) * 4
            wizardTexture[idx+0] = byte(r >> 8)
            wizardTexture[idx+1] = byte(g >> 8)
            wizardTexture[idx+2] = byte(b >> 8)
            wizardTexture[idx+3] = byte(a >> 8)
        }
    }
}


const gunTexWidth = 64
const gunTexHeight = 64

var gunTexture [gunTexWidth * gunTexHeight * 4]byte

func LoadGun(path string) {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    for y := 0; y < gunTexHeight; y++ {
        for x := 0; x < gunTexWidth; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            idx := (y*gunTexWidth + x) * 4
            gunTexture[idx+0] = byte(r >> 8)
            gunTexture[idx+1] = byte(g >> 8)
            gunTexture[idx+2] = byte(b >> 8)
            gunTexture[idx+3] = byte(a >> 8)
        }
    }
}

//enemies 
var demonTexture []byte
var demonFrames int
var wraithTexture []byte
var wraithFrames int
var reaperTexture []byte
var reaperFrames int

func LoadEnemySprites() {
    demonTexture, demonFrames = loadSpriteSheet("assets/demon_final.png", 64)
    wraithTexture, wraithFrames = loadSpriteSheet("assets/wraith_final.png", 64)
    reaperTexture, reaperFrames = loadSpriteSheet("assets/reaper_final.png", 64)
    fmt.Println("demon frames:", demonFrames, "wraith frames:", wraithFrames, "reaper frames:", reaperFrames)
}

func loadSpriteSheet(path string, frameSize int) ([]byte, int) {
    img, _, err := loadImage(path)
    if err != nil {
        return nil, 0
    }
    bounds := img.Bounds()
    frames := bounds.Max.X / frameSize
    pixels := make([]byte, bounds.Max.X*frameSize*4)
    for y := 0; y < frameSize; y++ {
        for x := 0; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            idx := (y*bounds.Max.X + x) * 4
            pixels[idx+0] = uint8(r >> 8)
            pixels[idx+1] = uint8(g >> 8)
            pixels[idx+2] = uint8(b >> 8)
            pixels[idx+3] = uint8(a >> 8)
        }
    }
    return pixels, frames
}





var pistolTexture []byte
var pistolW, pistolH int
var shotgunTexture []byte
var shotgunW, shotgunH int
var machinegunFrames int
var machinegunTexture []byte
var machinegunW, machinegunH int
var machinegunSheet []byte
var machinegunFrameW, machinegunFrameH int


func LoadWeapons() {
    machinegunSheet, machinegunFrames = loadSpriteSheet("assets/machinegun_sheet.png", 512)
    machinegunFrameW = 512
    machinegunFrameH = 512
    machinegunTexture, machinegunW, machinegunH = loadWeaponImage("assets/machinegun_idle.png")
}

func loadWeaponImage(path string) ([]byte, int, int) {
    img, _, err := loadImage(path)
    if err != nil {
        return nil, 0, 0
    }
    bounds := img.Bounds()
    w := bounds.Max.X
    h := bounds.Max.Y
    pixels := make([]byte, w*h*4)
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            idx := (y*w + x) * 4
            pixels[idx+0] = uint8(r >> 8)
            pixels[idx+1] = uint8(g >> 8)
            pixels[idx+2] = uint8(b >> 8)
            pixels[idx+3] = uint8(a >> 8)
        }
    }
    return pixels, w, h
}




