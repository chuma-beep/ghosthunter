package engine

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var audioContext *audio.Context

func InitAudio() {
	audioContext = audio.NewContext(44100)
}

func PlaySound(path string) {
	if audioContext == nil {
		return
	}
	data, err := assets.ReadFile(path)
	if err != nil {
		log.Println("sound error:", err)
		return
	}

	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Println("decode error:", err)
		return
	}

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		log.Println("player error:", err)
		return
	}

	player.Play()
}

var musicPlayer *audio.Player

func PlayMusic(path string) {
	data, err := assets.ReadFile(path)
	if err != nil {
		log.Println("music error:", err)
		return
	}

	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
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
