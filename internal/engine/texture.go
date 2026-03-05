package engine

import (
    "image"
    _ "image/png"
    "os"
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
