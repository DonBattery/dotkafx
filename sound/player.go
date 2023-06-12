package sound

import (
	"dotkafx/log"
	"dotkafx/model"
	"embed"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

const (
	embeddedSoundsFolder        = "embedded_sounds"
	ChaosDunk                   = "chaos_dunk"
	DotkaFXSercerIsOnline       = "dotkafx_server_is_online"
	DotkaFXServerIsShuttingDown = "dotkafx_server_is_shutting_down"
	SchedulerPaused             = "scheduler_paused"
	SchedulerRestarted          = "scheduler_restarted"
	SchedulerResumed            = "scheduler_resumed"
	SchedulerStarted            = "scheduler_started"
	SchedulerStopped            = "scheduler_stopped"
	SchedulerRolledBackward     = "scheduler_rolled_backward"
	SchedulerRolledForward      = "scheduler_rolled_forward"
)

type Player struct {
	embedded embed.FS
	format   beep.Format
	sounds   map[string]*beep.Buffer
}

func NewPlayer(embedded embed.FS) *Player {
	return &Player{
		embedded: embedded,
		sounds:   make(map[string]*beep.Buffer),
	}
}

func (player *Player) loadFromReadCloser(name string, rc io.ReadCloser) error {
	streamer, format, err := mp3.Decode(rc)
	if err != nil {
		return err
	}

	player.format = format

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	player.sounds[name] = buffer

	return nil
}

// loadSound loads an mp3 file into the Player's memory. If it ends with .mp3 we try to load it from the filesystem path,
// if not we will try to load it from the embedded sounds.
func (player *Player) loadSound(name string) error {
	if strings.HasSuffix(name, ".mp3") {
		file, err := os.Open(name)
		if err != nil {
			return err
		}
		defer file.Close()
		return player.loadFromReadCloser(name, file)
	}

	file, err := player.embedded.Open(fmt.Sprintf("%s/%s.mp3", embeddedSoundsFolder, name))
	if err != nil {
		return err
	}
	defer file.Close()
	return player.loadFromReadCloser(name, file)
}

func (player *Player) LoadSoundsAndInitSpeaker(profile model.ConfigProfile) error {
	for soundEffect := range profile.AllSoundEffect() {
		if err := player.loadSound(soundEffect); err != nil {
			return err
		}
	}

	for _, soundEffect := range []string{
		ChaosDunk,
		DotkaFXSercerIsOnline,
		DotkaFXServerIsShuttingDown,
		SchedulerPaused,
		SchedulerRestarted,
		SchedulerResumed,
		SchedulerStarted,
		SchedulerStopped,
		SchedulerRolledBackward,
		SchedulerRolledForward,
	} {
		if err := player.loadSound(soundEffect); err != nil {
			return err
		}
	}

	return speaker.Init(player.format.SampleRate, player.format.SampleRate.N(time.Second/10))
}

func (player *Player) Play(name string) {
	fx, ok := player.sounds[name]
	if !ok {
		return
	}
	sound := fx.Streamer(0, fx.Len())
	log.Debug("SoundPlayer is now playing: %s", name)
	speaker.Play(sound)
}

func (player *Player) Names() (names []string) {
	for key := range player.sounds {
		names = append(names, key)
	}
	return
}
