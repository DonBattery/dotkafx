# DotkaFX

DotkaFX is a simple scheduler for our favorite game written in Golang. It can run in Server and Client mode depending on the command line arguments. In Server mode it reads a profile form a configuration file. Based on this profile it builds up the timeline of one-shot and/or recurring events. In Client mode we can pass in control commands to the Server (e.g.: start, pause, back1m24s). Once started, the Server will play the appropriate sound effects when certain events are happening.  

## Installation

As the application is written in Golang, and does not use any platform specific package, we should be able to build and run it on both Windows, Linux and Mac. If Go is installed on your machine, just issue the `go install` command in the root of the repository and you should be good to go (issue the `dotkafx -h` command to check if the app is installed correctly).

For convenience I have built a binary for Windows (amd64) and put it in the **dist/windows** folder, this should work on both Intel and AMD based Windows systems. Once downloaded you need to add the folder where it is [to the PATH](https://stackoverflow.com/questions/44272416/how-to-add-a-folder-to-path-environment-variable-in-windows-10-with-screensho), so when you issue the **dotkafx** command in your terminal Windows will know where to find the exe file.

## The ConfigFile

Running the application for the first time will create the [dotkafx_config.yml](dotkafx_config.yml) file in your Home folder (C:\Users\YourUsername\dotkafx_config.yml). You may create new UserProfiles in this file (just copy the default one, and rename it). The config file must follow the following format:  

```YAML
---
Profiles:

  # There is a "default" Profile in the config file, but you can create as many custom Profiles as you want
  profile_name:

    # GlobalOffset is the duration added to all Event (basically if you set this to a negative number you will hear the sound effects before the events actually happens)
    GlobalOffset: -10s
    # Countdown is the duration before the match is actually started, the game clock is counting down until 00:00:00 and the it starts to count up.
    Countdown: 1m
    # This is the predicted maximum length of a match, the scheduler won't schedule any Event happening after this time
    MatchLength: 2h

    # Events is a map of objects where every key is the name of an event
    Events:

    "First Bounty Runes":
      # SoundEffect of the Event, if ends in .mp3 the app tries to load it from the local file system (so you either use absolute path, or place the mp3 file to the folder where you run the app).
      SoundEffect: 'C:\Users\myuser\Documents\my_favorite_sound_effect.mp3'
      # You can set individual offset to Events, so you can tune them even further
      Offset: 0
      # When should this event first happen
      FirstHappensAt: 0
      # Interval is the time duration between occurrences
      Interval: 0
      # Repeats tells the scheduler how many times this Event occurs, less than 1 means repeat infinitely
      Repeats: 1

    # We can continue adding more events...
    "Bounty Runes":
      # In this case we are using an embedded sound effect (which is built into the application)
      SoundEffect: "bounty_runes_appeared"
      Offset: 0
      FirstHappensAt: 3m
      Interval: 3m
      Repeats: 0
...
```  

**Durations** can be in the following formats: a simple integer number means seconds, "1h23m48s" will be translated to seconds.  

## CLI Usage  

```TEXT
Usage: dotkafx.exe [--config-file CONFIG-FILE] [--config-profile-name CONFIG-PROFILE-NAME] [--port PORT] [--debug] [COMMAND]

Positional arguments:
  COMMAND

Options:
  --config-file CONFIG-FILE, -f CONFIG-FILE [default: C:\Users\your_username\dotkafx_config.yml] 
  --config-profile-name CONFIG-PROFILE-NAME, -n CONFIG-PROFILE-NAME [default: default]
  --port PORT, -p PORT [default: 38383]
  --debug
  --help, -h             display this help and exit
```  

Running the application without any argument will spin up the Server with the default config profile. Running the application with these flags:  
```TEXT
dotkafx.exe --config-file myconfig.yml --config-profile-name myprofile --port 8080 --debug
```
will run the Server with the **myconfig.yml** config wile, and the **myprofile** profile (from this config file). Also the Server will listen on TCP Port **8080** and it will be running in **debug mode**. The same can be achieved with shorthands:  
```TEXT
dotkafx.exe -f myconfig.yml -n myprofile -p 8080 --debug
```  

Once the Server is running it is time to open another terminal and issue the following command:  
```TEXT
dotkafx.exe start
```  
this will dial into the Server and issue a start command, so the scheduler will be started and appropriate sound effects will be played on timeline events. Issue the
```TEXT
dotkafx.exe pause
```  
command to tell the Server to pause the scheduler, Issue the  
```TEXT
dotkafx.exe pause --debug
```  
again to unpause the scheduler, this time in debug mode. Issue the  
```TEXT
dotkafx.exe shutdown
```  
command to shut down the Server.  

## Auto Hotkey

Using [Auto Hotkey](https://www.autohotkey.com/) we can have a script [like this](dotkafx.ahk) so we can control DotkaFX with key combinations.
