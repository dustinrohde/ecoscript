package main

import "math/rand"

type World struct {
	width     int
	height    int
	organisms []*Organism
}

func (w *World) index(v Vector) (i int, ok bool) {
	i = v.X + (v.Y * w.height)
	if i <= len(w.organisms) {
		ok = true
	}
	return
}

func (w *World) Get(vector Vector) (organism *Organism, ok bool) {
	i, ok := w.index(vector)
	if !ok {
		// TODO: figure out if we should crash the program here
		panic("seeing if this triggers; may just need to fail silently")
		return
	}
	organism = w.organisms[i]
	if organism != nil {
		ok = true
	}
	return
}

func (w *World) View(origin Vector, distance int) []Vector {
	n := (2*distance + 1) ^ 2 - 1
	vectors := make([]Vector, n)

	i := 0
	for y := -distance; y < distance; y++ {
		for x := -distance; x < distance; x++ {
			vec := origin.Plus(Vector{x, y})
			if !vec.Equals(origin) {
				vectors[i] = vec
				i++
			}
		}
	}
	return vectors
}

func (w *World) ViewShuffled(origin Vector, distance int) []Vector {
	vectors := w.View(origin, distance)
	n := len(vectors)

	shuffled := make([]Vector, n)
	for i, j := range rand.Perm(n) {
		shuffled[i] = vectors[j]
	}
	return shuffled
}

func (w *World) Remove(vector Vector) (ok bool) {
	i, ok := w.index(vector)
	if ok {
		w.organisms[i] = nil
	}
	return
}

func (w *World) EndLifeAt(vector Vector) (ok bool) {
	organism, ok := w.Get(vector)
	if ok {
		organism.EndLife()
		w.Remove(vector)
	}
	return
}
