package scheduler

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"dotkafx/log"
	"dotkafx/model"
	"dotkafx/sound"
	"dotkafx/tools"
)

type timeLineEvent struct {
	name        string
	happensAt   int
	soundEffect string
}

type Scheduler struct {
	profile          model.ConfigProfile
	state            string
	secondsFromStart int
	timeline         []*timeLineEvent
	EventChan        chan string
	mu               sync.Mutex
}

// NewScheduler  creates a new Scheduler initialized with the ConfigProfile in the "stopped" state.
func NewScheduler(profile model.ConfigProfile) *Scheduler {
	sch := &Scheduler{
		profile:   profile,
		EventChan: make(chan string),
		state:     "stopped",
	}

	sch.buildTimeline()

	go sch.initTicker()

	return sch
}

// buildTimeline builds up the timeline based on the ConfigProfile
func (sc *Scheduler) buildTimeline() {
	// for every occurrence of every Event in the ConfigProfile we put a timelineEvent into the timeline
	for eventName, event := range sc.profile.Events {
		willHappen := true
		occurred := 0
		for willHappen {
			if event.Repeats > 0 && occurred == event.Repeats {
				willHappen = false
				continue
			}

			nextOccurrenceAt := sc.secondsFromStart +
				sc.profile.Countdown +
				sc.profile.GlobalOffset +
				event.Offset +
				event.FirstHappensAt +
				event.Interval*occurred

			if nextOccurrenceAt > sc.profile.MatchLength {
				willHappen = false
				continue
			}

			sc.timeline = append(sc.timeline, &timeLineEvent{
				name:        eventName,
				soundEffect: event.SoundEffect,
				happensAt:   nextOccurrenceAt,
			})

			occurred += 1
		}
	}

	// we sort the timeline
	sort.Slice(sc.timeline, func(i, j int) bool {
		return sc.timeline[i].happensAt < (sc.timeline[j].happensAt)
	})

	// then we adjust the timeline so there will be no conflicting timelineEvents
	sc.adjustTimeline(0)

	// re-sort the timeline after adjustment
	sort.Slice(sc.timeline, func(i, j int) bool {
		return sc.timeline[i].happensAt < (sc.timeline[j].happensAt)
	})
}

// adjustTimeline shifts the overlapping events on the timeline by two seconds recursively (preferring negative adjustment)
func (sch *Scheduler) adjustTimeline(index int) {
	if index >= len(sch.timeline) {
		// We've checked all events, so we're done
		return
	}

	// check for conflicts with previous events
	for i := 0; i < index; i++ {
		if math.Abs(float64(sch.timeline[i].happensAt-sch.timeline[index].happensAt)) < 2 {
			// Conflict! Move the current event by 2 seconds (to the left if possible)
			if sch.timeline[index].happensAt > 2 {
				sch.timeline[index].happensAt -= 2
			} else {
				sch.timeline[index].happensAt += 2
			}

			// now that we've moved this event, we need to re-check all previous events
			sch.adjustTimeline(i)
			break
		}
	}

	// check next event
	sch.adjustTimeline(index + 1)
}

// TimelineString returns the timeline as a string
func (sc *Scheduler) TimelineString() string {
	out := ""

	for _, timelineEvent := range sc.timeline {
		out += fmt.Sprintf("Happens at: %s Name: %s SoundEffect: %s\n", tools.SecondsToString(timelineEvent.happensAt-sc.profile.Countdown), timelineEvent.name, timelineEvent.soundEffect)
	}

	return out
}

// gameTime returns the time as represented in the game (00:03:59), considering the seconds elapsed from Start
// minus the countdown seconds
func (sch *Scheduler) gameTime() string {
	return "GameTime: " + tools.SecondsToString(sch.secondsFromStart-sch.profile.Countdown)
}

// nextEvent returns the next timelineEvent from the timeline (according to secondsFromStart) and two boolean values
// the first indicates if this timelineEvent occurs right now,
// and the second indicates if we reached the end of the whole timeline.
func (sch *Scheduler) nextEvent() (nextEvent *timeLineEvent, happensNow bool, endOfMatch bool) {
	endOfMatch = sch.secondsFromStart-sch.profile.Countdown >= sch.profile.MatchLength
	for _, ev := range sch.timeline {
		nextEvent = ev
		happensNow = sch.secondsFromStart == ev.happensAt
		if ev.happensAt >= sch.secondsFromStart {
			return
		}
	}
	return
}

// initTicker creates a ticker that ticks every Second. If the Scheduler is in the "running" state
// it checks for the next event in the timeline and if we reached the end of the match. If the next event
// is happening in the current second we send the name of the correlating SoundEffect to the EventChan.
// if we have reached the end of the match we are stopping the scheduler.
func (sch *Scheduler) initTicker() {
	for {
		sch.mu.Lock()
		if sch.state == "running" {
			if sch.secondsFromStart%5 == 0 {
				log.Debug(sch.gameTime())
			}

			nextEvent, happensNow, endOfMatch := sch.nextEvent()

			if happensNow {
				log.Info("Timeline Event: %s %s", nextEvent.name, sch.gameTime())
				sch.EventChan <- nextEvent.soundEffect
			}

			if endOfMatch {
				sch.state = "stopped"
				log.Info("Maximum match length exceeded, Scheduler stopped automatically. %s", sch.gameTime())
			}

			sch.secondsFromStart += 1
		}
		sch.mu.Unlock()

		time.Sleep(time.Second)
	}
}

// Start resets the secondsFromStart and sets the state to "running"
func (sch *Scheduler) Start() string {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	message := "Scheduler restarted "
	if sch.state == "stopped" {
		message = "Scheduler started "
		sch.EventChan <- sound.SchedulerStarted
	} else {
		sch.EventChan <- sound.SchedulerRestarted
	}

	sch.state = "running"
	sch.secondsFromStart = 0

	return message + sch.gameTime()
}

// Stop stops the Scheduler without the possibility of resuming
func (sch *Scheduler) Stop() string {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	if sch.state == "stopped" {
		return "Scheduler is already stopped"
	}

	sch.state = "stopped"

	sch.EventChan <- sound.SchedulerStopped

	return "Scheduler stopped. " + sch.gameTime()
}

// Pause sets the state to "paused" if it was in "running" and vice versa.
func (sch *Scheduler) Pause() string {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	switch sch.state {
	case "running":
		sch.state = "paused"
		sch.EventChan <- sound.SchedulerPaused
		return "Scheduler paused. " + sch.gameTime()
	case "paused":
		sch.state = "running"
		sch.EventChan <- sound.SchedulerResumed
		return "Scheduler resumed. " + sch.gameTime()
	default:
		return fmt.Sprintf("Scheduler cannot be paused/unpaused in the %s state.", sch.state)
	}
}

// Back rolls the the Scheduler's secondsFromStart back by the input seconds (if it is running).
func (sch *Scheduler) Back(seconds int) string {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	if sch.state == "running" {
		movedBackwards := seconds
		newSeconds := sch.secondsFromStart - seconds
		if newSeconds < 0 {
			newSeconds = 0
			movedBackwards = sch.secondsFromStart
		}
		sch.secondsFromStart = newSeconds
		sch.EventChan <- sound.SchedulerRolledBackward
		return fmt.Sprintf("Scheduler rolled backwards by %d seconds. %s", movedBackwards, sch.gameTime())
	}

	return fmt.Sprintf("The Scheduler cannot be rolled backwards in the %s state", sch.state)
}

// Forward rolls the the Scheduler's secondsFromStart forward by the input seconds (if ti is running).
func (sch *Scheduler) Forward(seconds int) string {
	sch.mu.Lock()
	defer sch.mu.Unlock()

	if sch.state == "running" {
		sch.secondsFromStart = sch.secondsFromStart + seconds
		sch.EventChan <- sound.SchedulerRolledForward
		return fmt.Sprintf("Scheduler rolled forward by %d seconds. %s", seconds, sch.gameTime())
	}

	return fmt.Sprintf("The Scheduler cannot be rolled forward in the %s state", sch.state)
}
