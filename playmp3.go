package main

import (
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os/user"
	"path/filepath"
)

func main() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	if err := mix.Init(mix.INIT_MP3); err != nil {
		log.Println(err)
		return
	}
	defer mix.Quit()

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		log.Println(err)
		return
	}
	defer mix.CloseAudio()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fileToPlay := filepath.Join(usr.HomeDir, "Music/test.mp3")
	log.Println("fileToPlay:", fileToPlay)

	if music, err := mix.LoadMUS(fileToPlay); err != nil {
		log.Println(err)
	} else if err = music.Play(1); err != nil {
		log.Println(err)
	} else {
		/*		if sdl.GetAudioStatus() == sdl.AUDIO_PLAYING {
				for sdl.GetAudioStatus() == sdl.AUDIO_PLAYING {
					sdl.Delay(5000)
					log.Println("now:", sdl.GetAudioStatus())
				}
			}*/
		var running bool
		var event sdl.Event
		running = true
		for running {
			event = sdl.WaitEvent()
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
		music.Free()
	}
}
