package main

// MaxTopSpeed represents the absolute fastest a car can go
const MaxTopSpeed = 200

// MaxAcceleration represents how quickly a car reaches it's top speed.
// because the service updates it's state 60 times/second this is represented
// as an integer where each frame of the race increments the current speed by
// this ammount.
//
// Deceleration is the same for all cars.
const MaxAcceleration = 10

// MaxHandling represents the best handling a car can have.
// Handling dicates the max speed a car can go around a turn.
// It is represented as a float64 with a value between 0 and 1
const MaxHandling = 1

// A Car represents a race car
type Car struct {
	ID           int    `json:"id"`
	DriverName   string `json:"driver_name"`
	TopSpeed     int    `json:"top_speed"`
	Acceleration int    `json:"acceleration"`
	Handling     int    `json:"handling"`
}
