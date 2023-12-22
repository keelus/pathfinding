package main

import (
	"container/heap"
	"fmt"
	"image"
	"math"
	"time"
)

type Status string

const (
	STATUS_IDLE        Status = "STATUS_IDLE"
	STATUS_PATHING     Status = "STATUS_PATHING"
	STATUS_END_SUCCESS Status = "STATUS_END_SUCCESS"
	STATUS_END_NOPATH  Status = "STATUS_END_NOPATH"
)

type Node struct {
	Coord    image.Point
	Distance int
	Visited  bool
	Added    bool

	IsWall bool

	IsPath bool
	Prev   *Node
}

type Grid struct {
	Cells [][]Node
	Start image.Point
	End   image.Point

	Status Status

	PathLength int
	Iterations int
}

func NewGrid(size int, start, end image.Point) Grid {
	cells := make([][]Node, size)
	for i := 0; i < size; i++ {
		cells[i] = make([]Node, size)
		for j := 0; j < size; j++ {
			cells[i][j] = Node{Coord: image.Pt(i, j), Distance: math.MaxInt}
		}
	}

	return Grid{Cells: cells, Start: start, End: end, Status: STATUS_IDLE}
}

func (g *Grid) Restart(keepLayout bool) {
	cells := make([][]Node, len(g.Cells))
	for i := 0; i < len(g.Cells); i++ {
		cells[i] = make([]Node, len(g.Cells))
		for j := 0; j < len(g.Cells); j++ {
			if !keepLayout {
				cells[i][j] = Node{Coord: image.Pt(i, j), Distance: math.MaxInt}
			} else {
				cells[i][j] = Node{Coord: image.Pt(i, j), Distance: math.MaxInt, IsWall: g.Cells[i][j].IsWall}
			}
		}
	}
	g.PathLength = 0
	g.Iterations = 0
	g.Status = STATUS_IDLE
	g.Cells = cells
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Distance < pq[j].Distance }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*Node)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

const MS_COOLDOWN = 10

func (grid *Grid) DoDijkstra() {
	grid.Status = STATUS_PATHING
	pq := PriorityQueue{}

	grid.Cells[grid.Start.X][grid.Start.Y].Distance = 0
	heap.Push(&pq, &grid.Cells[grid.Start.X][grid.Start.Y])

	for i, row := range grid.Cells {
		for j := range row {
			if i != grid.Start.X && j != grid.Start.Y {
				grid.Cells[i][j].Distance = math.MaxInt
			}
		}
	}

	heap.Init(&pq)

	maxVal := -1

	for pq.Len() > 0 {
		grid.Iterations++
		if pq.Len() > maxVal {
			maxVal = pq.Len()
		}
		time.Sleep(MS_COOLDOWN * time.Millisecond)
		u := heap.Pop(&pq).(*Node)

		if u.Visited {
			continue
		}

		if u.Coord == grid.End {
			break
		}

		u.Visited = true

		directions := []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

		for _, dir := range directions {
			newPos := u.Coord.Add(dir)
			if newPos.X >= 0 && newPos.X < len(grid.Cells) && newPos.Y >= 0 && newPos.Y < len(grid.Cells[0]) {
				neighbor := &grid.Cells[newPos.X][newPos.Y]

				if neighbor.IsWall || neighbor.Visited {
					continue
				}

				alt := u.Distance + 1

				if !neighbor.IsWall && alt < neighbor.Distance {
					neighbor.Distance = alt
					neighbor.Prev = u
					if !neighbor.Added {
						neighbor.Added = true
						heap.Push(&pq, neighbor)
					}
				}
			}
		}
	}

	node := &grid.Cells[grid.End.X][grid.End.Y]
	if node.Prev == nil {
		grid.PathLength = -1
		grid.Status = STATUS_END_NOPATH
	} else {
		grid.Status = STATUS_END_SUCCESS
		grid.PathLength = node.Distance - 1

		for node != nil {
			node.IsPath = true
			node = node.Prev
		}
	}

	fmt.Printf("Dijkstra Finished\nDistance:%d", grid.PathLength)
}

func (grid *Grid) DoAStar() {
	grid.Status = STATUS_PATHING
	pq := PriorityQueue{}

	grid.Cells[grid.Start.X][grid.Start.Y].Distance = 0
	heap.Push(&pq, &grid.Cells[grid.Start.X][grid.Start.Y])

	itemIndex := 1
	for i, row := range grid.Cells {
		for j := range row {
			if i != grid.Start.X && j != grid.Start.Y {
				grid.Cells[i][j].Distance = math.MaxInt
				itemIndex++
			}
		}
	}

	heap.Init(&pq)

	for pq.Len() > 0 {
		grid.Iterations++
		time.Sleep(MS_COOLDOWN * time.Millisecond)
		u := heap.Pop(&pq).(*Node)

		if u.Visited {
			continue
		}

		if u.Coord == grid.End {
			break
		}

		u.Visited = true

		directions := []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

		for _, dir := range directions {
			newPos := u.Coord.Add(dir)
			if newPos.X >= 0 && newPos.X < len(grid.Cells) && newPos.Y >= 0 && newPos.Y < len(grid.Cells[0]) {
				neighbor := &grid.Cells[newPos.X][newPos.Y]

				if neighbor.IsWall || neighbor.Visited {
					continue
				}

				alt := u.Distance + manhattanDistance(neighbor.Coord, grid.End)

				if alt < neighbor.Distance {
					neighbor.Distance = alt
					neighbor.Prev = u
					neighbor.Added = true
					heap.Push(&pq, &grid.Cells[newPos.X][newPos.Y])
				}
			}
		}
	}

	node := &grid.Cells[grid.End.X][grid.End.Y]

	if node.Prev == nil {
		grid.PathLength = -1
		grid.Status = STATUS_END_NOPATH
	} else {
		grid.Status = STATUS_END_SUCCESS
		for node != nil {
			node.IsPath = true
			node = node.Prev
			grid.PathLength++
		}
		grid.PathLength -= 2 // Remove start and end
	}

	fmt.Print("A* Finished\n")
}

func manhattanDistance(a, b image.Point) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

func factor() {

}
