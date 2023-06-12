package model

type RootCommand struct {
	ConfigFile        string `arg:"-f,--config-file" default:"C:\\Users\\your_username\\dotkafx_config.yml"`
	ConfigProfileName string `arg:"-n,--config-profile-name" default:"default"`
	Port              int    `arg:"-p,--port" default:"38383"`
	Command           string `arg:"positional"`
	Debug             bool
}

func (rc RootCommand) Description() string {
	return `DotkaFX is a sound effect scheduler for Dota2.

Run it once without a command to spin up the server.
Run it again with a command argument which can be: start, stop, pause, back, forward or shutdown
Use the dotkafx_config.yml file in your home folder to adjust the timeline or create a personal configuration.
`
}
