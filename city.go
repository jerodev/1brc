package main

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

func NewCity(name string, temperatures []int) *City {
	c := City{
		Name: name,
		Min:  100,
		Max:  -100,
	}
	c.Merge(temperatures)

	return &c
}
