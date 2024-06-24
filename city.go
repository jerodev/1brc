package main

import (
	"fmt"
	"math"
)

type City struct {
	Name  string
	Count int
	Sum   int
	Min   int
	Max   int
}

func (r *City) Merge(temperatures []int) {
	for _, t := range temperatures {
		r.Sum += t
		r.Count++

		if t < r.Min {
			r.Min = t
		}
		if t > r.Max {
			r.Max = t
		}
	}
}

func (r *City) ToString() string {
	minVal := math.Round(float64(r.Min)) / 10
	avgVal := math.Round(float64(r.Sum)/float64(r.Count)) / 10
	maxVal := math.Round(float64(r.Max)) / 10

	return r.Name + "=" + fmt.Sprintf("%.1f", minVal) + "/" + fmt.Sprintf("%.1f", avgVal) + "/" + fmt.Sprintf("%.1f", maxVal)
}

func NewCity(name string, temperatures []int) *City {
	c := City{
		Name: name,
		Min:  100,
		Max:  -100,
	}
	c.Merge(temperatures)

	return &c
}
