package engine 

import(
	"bytes" 
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
  "github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var audioContext *audio.Context 

func InitAudio() {
	audioContext = audio.NewContext(44100)
}


func PlaySound(path string) {
	if audioContext == nil{
		return 
	}
	f, err := os.ReadFile(path)
	if  err != nil {
		log.Println("sound error:", err)
		return 
	}

	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(f))
	if err != nil{
		log.Println("decode error:", err)
		return
	}

    player, err := audioContext.NewPlayer(stream) 
	if err != nil{
		log.Println("player error:", err)
		return 
	}

	player.Play()
}

var musicPlayer *audio.Player

func PlayMusic(path string) {
    f, err := os.ReadFile(path)
    if err != nil {
        log.Println("music error:", err)
        return
    }

    stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(f))
    if err != nil {
        log.Println("music decode error:", err)
        return
    }

    loop := audio.NewInfiniteLoop(stream, stream.Length())

    musicPlayer, err = audioContext.NewPlayer(loop)
    if err != nil {
        log.Println("music player error:", err)
        return
    }

    musicPlayer.Play()
}
