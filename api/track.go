package main

import (
	"encoding/json"
	"math/rand"
	"os"
)

// MaxTurnDegree represents the absolute value of the maximum
// degree that a single segment of a track can be.
const MaxTurnDegree = 90

// A Track is represented as an array of segments
type Track struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Segments []int  `json:"segments"`
}

func generateTrackSegments(dataFileName string) error {
	f, err := os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	var d data

	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return err
	}

	f.Close()
	os.Remove(dataFileName)

	for _, track := range d.Tracks {
		if len(track.Segments) == 0 {
			track.Segments = generateTrack(rand.Intn(100) + 120)
		}
	}

	f, err = os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(d); err != nil {
		return err
	}

	return nil
}

func generateTrack(totalSegments int) []int {
	results := make([]int, totalSegments)

	for i := 0; i < totalSegments; i++ {
		results[i] = rand.Intn(90)
	}

	return results
}
