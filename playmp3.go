package lotte

import (
	"log"
	"os/user"
	"path/filepath"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func musicStopped() {
	log.Println("Done playing audio file...")
}

// PlayOgg Play's the ogg file
// (omit the path this is set to usr.HomeDir/music)
func PlayOgg(fileToPlay string) {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	if err := mix.Init(mix.INIT_OGG); err != nil {
		log.Println(err)
		return
	}
	defer mix.Quit()

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		log.Println(err)
		return
	}
	defer mix.CloseAudio()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fileToPlayfolder := filepath.Join(usr.HomeDir, "Music")
	fullfileNameToPlay := filepath.Join(fileToPlayfolder, fileToPlay)
	log.Println("fileToPlay:", fullfileNameToPlay)

	if music, err := mix.LoadMUS(fullfileNameToPlay); err != nil {
		log.Println(err)
	} else if err = music.Play(1); err != nil {
		log.Println(err)
	} else {
		mix.HookMusicFinished(musicStopped)
		var running bool
		running = true
		for running {
			sdl.Delay(100)
			running = mix.PlayingMusic()
		}
		music.Free()
	}
}
