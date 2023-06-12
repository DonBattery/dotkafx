package model

import (
	"fmt"

	"dotkafx/tools"
)

type Event struct {
	Offset         int
	FirstHappensAt int
	Interval       int
	Repeats        int
	SoundEffect    string
}

type EventInput struct {
	Offset         string `yaml:"Offset"`
	FirstHappensAt string `yaml:"FirstHappensAt"`
	Interval       string `yaml:"Interval"`
	Repeats        int    `yaml:"Repeats"`
	SoundEffect    string `yaml:"SoundEffect"`
}

func (ei EventInput) Parse() (Event, error) {
	ev := Event{}

	val, err := tools.StringToSeconds(ei.Offset)
	if err != nil {
		return ev, err
	}
	ev.Offset = val

	val, err = tools.StringToSeconds(ei.FirstHappensAt)
	if err != nil {
		return ev, err
	}
	ev.FirstHappensAt = val

	val, err = tools.StringToSeconds(ei.Interval)
	if err != nil {
		return ev, err
	}
	ev.Interval = val

	ev.Repeats = ei.Repeats

	ev.SoundEffect = ei.SoundEffect

	return ev, nil
}

type ConfigProfile struct {
	GlobalOffset int
	MatchLength  int
	Countdown    int
	Events       map[string]Event
}

type ConfigProfileInput struct {
	GlobalOffset string                `yaml:"GlobalOffset"`
	MatchLength  string                `yaml:"MatchLength"`
	Countdown    string                `yaml:"Countdown"`
	Events       map[string]EventInput `yaml:"Events"`
}

func (cpi ConfigProfileInput) Parse() (ConfigProfile, error) {
	cp := ConfigProfile{
		Events: make(map[string]Event),
	}

	val, err := tools.StringToSeconds(cpi.GlobalOffset)
	if err != nil {
		return cp, err
	}
	cp.GlobalOffset = val

	val, err = tools.StringToSeconds(cpi.MatchLength)
	if err != nil {
		return cp, err
	}
	cp.MatchLength = val

	val, err = tools.StringToSeconds(cpi.Countdown)
	if err != nil {
		return cp, err
	}
	cp.Countdown = val

	if cpi.Events == nil {
		return cp, fmt.Errorf("The config profile must have an Events map")
	}

	if len(cpi.Events) == 0 {
		return cp, fmt.Errorf("The Events map must have at least one Event")
	}

	for eventName, event := range cpi.Events {
		val, err := event.Parse()
		if err != nil {
			return cp, err
		}
		cp.Events[eventName] = val
	}

	return cp, nil
}

// AllSoundEffect returns a map where the keys are the used SoundEffects across the Profile.
func (cp ConfigProfile) AllSoundEffect() (soundEffects map[string]bool) {
	soundEffects = map[string]bool{}
	for _, event := range cp.Events {
		soundEffects[event.SoundEffect] = true
	}
	return
}

type Config struct {
	Profiles map[string]ConfigProfile
}

type ConfigInput struct {
	Profiles map[string]ConfigProfileInput `yaml:"Profiles"`
}

func (ci ConfigInput) Parse() (Config, error) {
	c := Config{
		Profiles: make(map[string]ConfigProfile),
	}

	if ci.Profiles == nil {
		return c, fmt.Errorf("The config must have a Profiles map")
	}

	if len(ci.Profiles) == 0 {
		return c, fmt.Errorf("The Profiles map must have at least one Profile")
	}

	for profileName, profile := range ci.Profiles {
		val, err := profile.Parse()
		if err != nil {
			return c, err
		}
		c.Profiles[profileName] = val
	}

	return c, nil
}

func (conf Config) CreateAndValidateProfile(profileName string) (profile ConfigProfile, err error) {
	profile, ok := conf.Profiles[profileName]
	if !ok {
		err = fmt.Errorf("The profile with the name: %s cannot be found in the configuration.", profileName)
		return
	}
	return
}

func (conf Config) String() string {
	out := "Profiles:\n"

	for profileName, profile := range conf.Profiles {
		out += "  " + profileName + ":\n"
		out += fmt.Sprintf(
			`    GlobalOffset: %s
    MatchLength : %s
    Countdown   : %s
`,
			tools.SecondsToString(profile.GlobalOffset),
			tools.SecondsToString(profile.MatchLength),
			tools.SecondsToString(profile.Countdown))
		out += "    Events:\n"
		for eventName, event := range profile.Events {
			out += "      " + eventName + ":\n"
			out += fmt.Sprintf(
				`        Offset        : %s
        FirstHappensAt: %s
        Interval      : %s
        Repeats       : %d
        SoundEffect   : %s
`,
				tools.SecondsToString(event.Offset),
				tools.SecondsToString(event.FirstHappensAt),
				tools.SecondsToString(event.Interval),
				event.Repeats,
				event.SoundEffect,
			)
		}
	}

	return out
}
