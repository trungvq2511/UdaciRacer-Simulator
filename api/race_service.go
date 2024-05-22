package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// RaceService is responsible for running races, sharing updates
// and returning the results
type RaceService struct {
	Cars   []*Car
	Tracks []*Track
	Races  []*Race
}

// NewRaceService takes the options and returns a new RaceService
func NewRaceService(opts ...RaceServiceOpt) (*RaceService, error) {
	var s RaceService

	for _, opt := range opts {
		if err := opt.Apply(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

// CreateRace makes a new race
func (s *RaceService) CreateRace(playerID int, trackID int) (*Race, error) {
	id := len(s.Races) + 1

	race, err := NewRace(
		WithID(id),
		WithPlayerID(playerID),
		WithCars(s.Cars),
		WithTrack(s.Tracks[trackID - 1]),
	)
	if err != nil {
		return nil, err
	}

	s.Races = append(s.Races, race)

	return race, nil
}

// GetRace fetches an existing race
func (s *RaceService) GetRace(raceID uint) (*Race, error) {
	if raceID-1 < 0 || raceID > uint(len(s.Races)) {
		return nil, fmt.Errorf("invalid race id")
	}

	return s.Races[raceID-1], nil
}

// Accelerate exposes the Race's Accelerate method
func (s *RaceService) Accelerate(raceID uint) error {
	race, err := s.GetRace(raceID)
	if err != nil {
		return err
	}

	err = race.Accelerate()
	return nil
}

// RaceServiceOpt are options for building a RaceService
type RaceServiceOpt interface {
	Apply(*RaceService) error
}

// DataFromJSONFile tells the RaceService where the data lives
func DataFromJSONFile(path string) RaceServiceOpt {
	return &dataFromJSONFile{path}
}

type dataFromJSONFile struct {
	path string
}

func (opt *dataFromJSONFile) Apply(s *RaceService) error {
	f, err := os.Open(opt.path)
	if err != nil {
		return err
	}

	if err := DataFromReader(f).Apply(s); err != nil {
		return err
	}

	return nil
}

// DataFromReader ...
func DataFromReader(r io.Reader) RaceServiceOpt {
	return &dataFromReader{r}
}

type dataFromReader struct {
	r io.Reader
}

func (opt *dataFromReader) Apply(s *RaceService) error {
	var d data

	if err := json.NewDecoder(opt.r).Decode(&d); err != nil {
		return err
	}

	s.Cars = d.Cars
	s.Tracks = d.Tracks

	return nil
}

type data struct {
	Cars   []*Car   `json:"cars"`
	Tracks []*Track `json:"tracks"`
}
