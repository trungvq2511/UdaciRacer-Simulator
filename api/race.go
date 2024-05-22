package main

import (
	"fmt"
	"time"
)

// Race holds all of the information about the race
type Race struct {
	ID                   int
	Track                *Track
	PlayerID             int
	clicksSinceLastCheck int
	finishPosition       int
	Cars                 []*Car
	Results              *RaceResults
	actionCh             chan int
	tickerCh             <-chan time.Time
	closeCh              chan struct{}
}

// NewRace constructs a Race from the provided options
func NewRace(opts ...RaceOpt) (*Race, error) {
	results := &RaceResults{
		Status: Unstarted,
	}

	r := Race{
		actionCh: make(chan int),
		closeCh:  make(chan struct{}),
		Results:  results,
	}

	for _, opt := range opts {
		if err := opt.Apply(&r); err != nil {
			return nil, err
		}
	}

	return &r, nil
}

// Start starts a race
func (r *Race) Start() error {
	if r.Results.Status != Unstarted {
		return fmt.Errorf("race has already started")
	}

	r.Results.Status = InProgress

	r.tickerCh = time.NewTicker(time.Second).C

	go r.loop()

	return nil
}

// Finish marks the end of a race
func (r *Race) Finish() error {
	close(r.closeCh)
	close(r.actionCh)

	return nil
}

// Refresh computes the location of cars on the track
// based on speed and acceleration
func (r *Race) Refresh() (*RaceResults, error) {
	raceFinished := true

	for _, car := range r.Results.Positions {
		if car.FinalPosition > 0 {
			continue
		}

		if car.ID == r.PlayerID {
			if r.clicksSinceLastCheck > 0 {
				car.Speed = car.Speed + (car.Acceleration / 1) + r.clicksSinceLastCheck
				r.clicksSinceLastCheck = 0
			} else {
				car.Speed = car.Speed - 5
			}
		} else {
			car.Speed = car.Speed + car.Acceleration
		}

		if car.Speed > car.TopSpeed {
			car.Speed = car.TopSpeed
		}

		if car.Speed <= 0 {
			car.Speed = 0
		}

		car.Segment = car.Segment + (car.Speed / 30)

		if car.Segment >= len(r.Track.Segments) {
			car.Segment = len(r.Track.Segments)
			r.finishPosition++
			car.FinalPosition = r.finishPosition
		} else {
			raceFinished = false
		}
	}

	if raceFinished {
		r.Finish()
	}

	return r.Results, nil
}

func (r *Race) loop() {
	for {
		select {
		case <-r.tickerCh:
			r.Refresh()
		case <-r.actionCh:
			r.clicksSinceLastCheck++
		case <-r.closeCh:
			r.Results.Status = Finished
			return
		}

		// time.Sleep(time.Second * 1)
	}
}

// Accelerate handles the action from the user to increase the speed
// Acceleration sets the cap for this.
func (r *Race) Accelerate() error {
	if r.actionCh == nil {
		return nil
	}

	r.actionCh <- 0

	return nil
}

// RaceOpt allows the configuration of a new Race
type RaceOpt interface {
	Apply(*Race) error
}

// WithID sets the id of the Race
func WithID(id int) RaceOpt {
	return &withID{id}
}

type withID struct {
	id int
}

func (opt *withID) Apply(r *Race) error {
	r.ID = opt.id

	return nil
}

// WithCars sets the cars for the race
func WithCars(cars []*Car) RaceOpt {
	return &withCars{cars}
}

type withCars struct {
	cars []*Car
}

func (opt *withCars) Apply(r *Race) error {
	r.Cars = opt.cars
	for _, car := range opt.cars {
		r.Results.Positions = append(r.Results.Positions, &CarPosition{Car: *car})
	}

	return nil
}

// WithTrack sets the track
func WithTrack(t *Track) RaceOpt {
	return &withTrack{t}
}

type withTrack struct {
	t *Track
}

func (opt *withTrack) Apply(r *Race) error {
	r.Track = opt.t

	return nil
}

// WithPlayerID sets which car a user wants to use
func WithPlayerID(playerID int) RaceOpt {
	return &withPlayerID{playerID}
}

type withPlayerID struct {
	playerID int
}

func (opt *withPlayerID) Apply(r *Race) error {
	r.PlayerID = opt.playerID

	return nil
}

// RaceResults represent the results of a race
// These results can be a finised race, an in-progress race
// or a race that has yet to be started.
type RaceResults struct {
	Status    RaceStatus     `json:"status"`
	Positions []*CarPosition `json:"positions"`
}

// RaceStatus is the status of the race.
// It is eaither unstarted, in-progress, or finished
type RaceStatus string

// Values for RaceStatus
const (
	Unstarted  = "unstarted"
	InProgress = "in-progress"
	Finished   = "finished"
)

// CarPosition wraps a car and determines the position
// of a car in a race.
type CarPosition struct {
	Car
	FinalPosition int `json:"final_position,omitempty"`
	Speed         int `json:"speed"`
	Segment       int `json:"segment"`
}
