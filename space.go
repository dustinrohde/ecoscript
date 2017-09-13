package main

import "math/rand"

// Space is a container that has width and height and contains Cells.
type Space interface {
	// Width returns the length of the x-axis.
	Width() int

	// Height returns the length of the y-axis.
	Height() int

	// InBounds returns true if the given Vector is in bounds.
	InBounds(vec Vector) bool

	// Walkable returns true if the given Vector is walkable.
	Walkable(vec Vector) bool

	// Cell returns the Cell at the given vector.
	Cell(vec Vector) *Cell

	// View returns all Vectors that are in-bounds within a radius.
	View(origin Vector, radius int) []Vector

	// ViewR is like View but randomizes the returned Vectors.
	ViewR(origin Vector, radius int) []Vector

	// ViewWalkable returns all walkable Vectors within a radius.
	ViewWalkable(origin Vector, radius int) []Vector

	// ViewWalkableR is like ViewWalkable but randomizes the returned Vectors.
	ViewWalkableR(origin Vector, radius int) []Vector

	// RandWalkable finds a random walkable Vector within a radius.
	RandWalkable(origin Vector, radius int) Vector

	// Add attempts to add an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Add(org *Organism, vec Vector) (func(), bool)

	// Remove attempts to remove an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Remove(org *Organism, vec Vector) (func(), bool)

	// Remove attempts to remove and kill an Organism at the given Vector.
	// It returns true if it succeeded or false if it wasn't found.
	Kill(org *Organism, vec Vector) (func(), bool)

	// Move attempts to move an Organism from one Vector to another
	// It returns true if it succeeded or false if it wasn't found.
	Move(org *Organism, src Vector, dst Vector) (func(), bool)
}

func InBounds(s Space, vec Vector) bool {
	return vec.Flatten(s.Height()) < s.Width()*s.Height()
}

func Walkable(s Space, vec Vector) bool {
	return s.InBounds(vec) && !s.Cell(vec).Occupied()
}

func View(s Space, origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, s.InBounds)
}

func ViewR(s Space, origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, s.InBounds)
}

func ViewWalkable(s Space, origin Vector, radius int) []Vector {
	vectors := origin.Radius(radius)
	return VecFilter(vectors, s.Walkable)
}

func ViewWalkableR(s Space, origin Vector, radius int) []Vector {
	vectors := origin.RadiusR(radius)
	return VecFilter(vectors, s.Walkable)
}

func RandWalkable(s Space, origin Vector, radius int) Vector {
	vectors := s.ViewWalkable(origin, radius)
	index := rand.Intn(len(vectors))
	return vectors[index]
}

func Add(s Space, organism *Organism, vec Vector) (exec func(), ok bool) {
	cell := s.Cell(vec)
	return cell.Add(organism)
}

func Remove(s Space, organism *Organism, vec Vector) (exec func(), ok bool) {
	cell := s.Cell(vec)
	return cell.Remove(organism)
}

func Move(s Space, organism *Organism, src Vector, dst Vector) (exec func(), ok bool) {
	oldCell := s.Cell(src)
	newCell := s.Cell(dst)
	execAdd, okAdd := newCell.Add(organism)
	execRm, okRm := oldCell.Remove(organism)

	ok = okAdd && okRm
	if ok {
		exec = chain(execAdd, execRm)
	}
	return
}

func Kill(s Space, organism *Organism, vec Vector) (exec func(), ok bool) {
	// TODO: implement corpses
	execRm, ok := s.Remove(organism, vec)
	if ok {
		exec = chain(execRm, organism.EndLife)
	}
	return
}
