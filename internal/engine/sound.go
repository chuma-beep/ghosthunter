package engine 

import(
	"bytes" 
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var audioContext *audio.Context 

func InitAudio() {
	audioContext = audio.NewContext(44100)
}


func PlaySound(path string) {
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


