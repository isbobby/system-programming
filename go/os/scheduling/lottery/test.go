package main

import (
	"errors"
	"math"
)

func acceptableProbability(expected, actual, tolerance float64) error {
	if math.Abs(actual-expected) > tolerance {
		return errors.New("probability not in acceptable tolerance")
	}
	return nil
}
