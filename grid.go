package main

import (
	"container/heap"
	"math"
	"pathfinding/pair"
	"time"
)

var MS_COOLDOWN = 10

const BASE_WEIGHT = 1

type Status string

const (
	STATUS_IDLE        Status = "STATUS_IDLE"
	STATUS_PATHING     Status = "STATUS_PATHING"
	STATUS_END_SUCCESS Status = "STATUS_END_SUCCESS"
	STATUS_END_NOPATH  Status = "STATUS_END_NOPATH"
)

type Node struct {
	Coord  pair.Pair
	IsWall bool
	Prev   *Node

	Cost  float64 // Cost for Dijkstra, Fcost for A*
	Gcost float64 // For A*

	Visited bool
	Added   bool

	IsPath bool

	index int
}

type Grid struct {
	Cells [][]Node
	Start *Node
	End   *Node

	Status Status

	PathLength int
	Iterations int
}

func NewGrid(size int, start, end pair.Pair) Grid {
	cells := make([][]Node, size)
	for i := 0; i < size; i++ {
		cells[i] = make([]Node, size)
		for j := 0; j < size; j++ {
			cells[i][j] = Node{
				Coord: pair.New(i, j),
				Cost:  math.MaxInt}
		}
	}

	return Grid{Cells: cells, Start: &cells[start.I][start.J], End: &cells[end.I][end.J], Status: STATUS_IDLE}
}

func (g *Grid) Restart(keepLayout bool) {
	cells := make([][]Node, len(g.Cells))
	for i := 0; i < len(g.Cells); i++ {
		cells[i] = make([]Node, len(g.Cells))
		for j := 0; j < len(g.Cells); j++ {
			cells[i][j] = Node{Coord: pair.New(i, j), IsWall: keepLayout && g.Cells[i][j].IsWall}
		}
	}

	g.Start = &cells[g.Start.Coord.I][g.Start.Coord.J]
	g.End = &cells[g.End.Coord.I][g.End.Coord.J]

	g.PathLength = 0
	g.Iterations = 0
	g.Status = STATUS_IDLE
	g.Cells = cells
}

func (grid *Grid) DoDijkstra() {
	grid.Status = STATUS_PATHING
	pq := PriorityQueue{}

	grid.Start.Cost = 0
	heap.Push(&pq, grid.Start)

	for i, row := range grid.Cells {
		for j := range row {
			if &grid.Cells[i][j] != grid.Start {
				grid.Cells[i][j].Cost = math.MaxInt
			}
		}
	}

	heap.Init(&pq)

	for pq.Len() > 0 {
		select {
		case <-stopSignal:
			grid.Status = STATUS_IDLE
			return
		default:
		}

		grid.Iterations++
		time.Sleep(time.Millisecond * time.Duration(MS_COOLDOWN))

		u := heap.Pop(&pq).(*Node)

		if u.Visited {
			continue
		}

		if u == grid.End {
			break
		}

		u.Visited = true

		directions := []pair.Pair{pair.Up(), pair.Down(), pair.Left(), pair.Right()}
		for _, dir := range directions {
			neighborPos := u.Coord.Add(dir)

			if neighborPos.InBounds(0, 0, len(grid.Cells), len(grid.Cells[0])) {
				neighbor := &grid.Cells[neighborPos.I][neighborPos.J]

				if neighbor.IsWall || neighbor.Visited {
					continue
				}

				alt := u.Cost + 1
				if !neighbor.IsWall && alt < neighbor.Cost {
					neighbor.Cost = alt
					neighbor.Prev = u
					if !neighbor.Added {
						neighbor.Added = true
						heap.Push(&pq, neighbor)
					}
				}
			}
		}
	}

	grid.Status = STATUS_END_SUCCESS
	success := grid.constructPath()
	if !success {
		grid.Status = STATUS_END_NOPATH
	}
}

func (grid *Grid) DoAStar() {
	grid.Status = STATUS_PATHING
	for i := range grid.Cells {
		for j := range grid.Cells[i] {
			if grid.Start != &grid.Cells[i][j] {
				grid.Cells[i][j].Gcost = math.MaxFloat64
				grid.Cells[i][j].Cost = math.MaxFloat64
			}
		}
	}

	grid.Start.Gcost = 0
	grid.Start.Cost = grid.h(*grid.Start)

	pq := PriorityQueue{grid.Start}
	heap.Init(&pq)

	for pq.Len() > 0 {
		select {
		case <-stopSignal:
			grid.Status = STATUS_IDLE
			return
		default:
		}

		grid.Iterations++
		time.Sleep(time.Duration(MS_COOLDOWN) * time.Millisecond)

		current := heap.Pop(&pq).(*Node)
		current.Visited = true

		if current == grid.End {
			break
		}

		directions := []pair.Pair{pair.Up(), pair.Down(), pair.Left(), pair.Right()}
		for _, dir := range directions {
			neighborPos := current.Coord.Add(dir)

			if neighborPos.InBounds(0, 0, len(grid.Cells), len(grid.Cells[0])) {
				neighbor := &grid.Cells[neighborPos.I][neighborPos.J]

				if neighbor.IsWall {
					continue
				}

				gcost := current.Gcost + g(*current, *neighbor)
				if gcost < neighbor.Gcost {
					neighbor.Prev = current
					neighbor.Gcost = gcost
					neighbor.Cost = gcost + grid.h(*neighbor)

					if neighbor.Added {
						heap.Fix(&pq, neighbor.index)
					} else {
						neighbor.Added = true
						heap.Push(&pq, neighbor)
					}
				}
			}
		}
	}

	grid.Status = STATUS_END_SUCCESS
	success := grid.constructPath()
	if !success {
		grid.Status = STATUS_END_NOPATH
	}
}

func g(a, b Node) float64 {
	return BASE_WEIGHT // This could be changed to use diagonals (e.g 1 for horizontal & vertical, 1.4 for diagonals)
}

func (grid Grid) h(a Node) float64 {
	dy := math.Abs(float64(a.Coord.I - grid.End.Coord.I))
	dx := math.Abs(float64(a.Coord.J - grid.End.Coord.J))
	heuristic := BASE_WEIGHT * (dx + dy)
	tb := float64(1)/float64(len(grid.Cells)*len(grid.Cells[0])) + 1 // Tie breaker
	return heuristic * tb
}

func (grid *Grid) constructPath() bool {
	node := grid.End

	if node.Prev == nil {
		grid.PathLength = -1
		return false
	} else {
		for node != nil {
			if node != grid.Start && node != grid.End {
				grid.PathLength++
			}

			node.IsPath = true
			node = node.Prev
		}
	}

	return true
}
