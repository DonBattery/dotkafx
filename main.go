package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"dotkafx/client"
	"dotkafx/config"
	"dotkafx/log"
	"dotkafx/model"
	"dotkafx/scheduler"
	"dotkafx/server"
	"dotkafx/sound"
)

var (
	//go:embed dotkafx_config.yml
	defaultConfig []byte

	//go:embed embedded_sounds/*.mp3
	embeddedSounds embed.FS
)

func main() {
	command := model.RootCommand{}

	arg.MustParse(&command)

	if command.Debug {
		log.LoggingLevel = log.DebugLevel
	}

	log.Debug("Running with command: %+v", command)

	// if there is a positional argument, run the Client and pass the argument to it as the command.
	if len(command.Command) > 0 {
		log.Debug("Sending message: %s to DotkaFX Server via TCP Port: %d", command.Command, command.Port)
		response, err := client.NewClient(command.Port).SendRequest(command.Command)
		if err != nil {
			quit(err)
		}
		log.Info(response)
	} else {
		runServer(command)
	}
}

func runServer(cmd model.RootCommand) {
	// get the configuration
	confData, err := config.GetConfigData(defaultConfig)
	if err != nil {
		quit(err)
	}
	conf, err := config.CreateConfig(confData)
	if err != nil {
		quit(err)
	}
	log.Debug("Config object:\n%s", conf)
	log.Debug("Config Profile: %s", cmd.ConfigProfileName)
	profile, err := conf.CreateAndValidateProfile(cmd.ConfigProfileName)
	if err != nil {
		quit(err)
	}
	log.Debug("Loaded SoundEffects: %+v", profile.AllSoundEffect())

	// create the Sound Effect Player
	fx := sound.NewPlayer(embeddedSounds)
	if err := fx.LoadSoundsAndInitSpeaker(profile); err != nil {
		quit(err)
	}

	// create the Scheduler
	sch := scheduler.NewScheduler(profile)
	log.Debug("Scheduler Timeline:\n%s", sch.TimelineString())

	// create and run the Server
	srv := server.NewServer(fx, sch, cmd)
	quit(srv.Run())
}

func quit(errorMessage any) {
	if errorMessage != nil {
		log.Fatal(fmt.Sprintf("%s", errorMessage))
	}
	os.Exit(0)
}
