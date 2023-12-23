package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"pathfinding/pair"
	"slices"
	"time"

	"golang.org/x/exp/maps"
)

type Status string

const (
	STATUS_IDLE        Status = "STATUS_IDLE"
	STATUS_PATHING     Status = "STATUS_PATHING"
	STATUS_END_SUCCESS Status = "STATUS_END_SUCCESS"
	STATUS_END_NOPATH  Status = "STATUS_END_NOPATH"
)

type Node struct {
	Coord   pair.Pair
	Visited bool
	Added   bool

	IsWall bool

	Gcost int // For A*
	Hcost int // For A*
	Cost  int // Fcost in A*

	IsPath bool
	Prev   *Node

	index int
}

type Grid struct {
	Cells [][]Node
	Start pair.Pair
	End   pair.Pair

	Status Status

	PathLength int
	Iterations int
}

func NewGrid(size int, start, end pair.Pair) Grid {
	cells := make([][]Node, size)
	for i := 0; i < size; i++ {
		cells[i] = make([]Node, size)
		for j := 0; j < size; j++ {
			cells[i][j] = Node{Coord: pair.New(i, j), Cost: math.MaxInt}
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
				cells[i][j] = Node{Coord: pair.New(i, j), Cost: math.MaxInt}
			} else {
				cells[i][j] = Node{Coord: pair.New(i, j), Cost: math.MaxInt, IsWall: g.Cells[i][j].IsWall}
			}
		}
	}
	g.PathLength = 0
	g.Iterations = 0
	g.Status = STATUS_IDLE
	g.Cells = cells
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil  // avoid memory leak
	node.index = -1 // for safety
	*pq = old[0 : n-1]
	return node
}

const MS_COOLDOWN = 10

func (grid *Grid) DoDijkstra() {
	grid.Status = STATUS_PATHING
	pq := PriorityQueue{}

	grid.Cells[grid.Start.I][grid.Start.J].Cost = 0
	heap.Push(&pq, &grid.Cells[grid.Start.I][grid.Start.J])

	for i, row := range grid.Cells {
		for j := range row {
			if i != grid.Start.I && j != grid.Start.J {
				grid.Cells[i][j].Cost = math.MaxInt
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

		directions := []pair.Pair{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

		for _, dir := range directions {
			newPos := u.Coord.Add(dir)
			if newPos.I >= 0 && newPos.I < len(grid.Cells) && newPos.J >= 0 && newPos.J < len(grid.Cells[0]) {
				neighbor := &grid.Cells[newPos.I][newPos.J]

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

	node := &grid.Cells[grid.End.I][grid.End.J]
	if node.Prev == nil {
		grid.PathLength = -1
		grid.Status = STATUS_END_NOPATH
	} else {
		grid.Status = STATUS_END_SUCCESS
		grid.PathLength = node.Cost - 1

		for node != nil {
			node.IsPath = true
			node = node.Prev
		}
	}

	fmt.Printf("Dijkstra Finished\nDistance:%d", grid.PathLength)
}

func (grid *Grid) DoAStar() {
	openSet := make(map[pair.Pair]struct{})

	cameFrom := make(map[pair.Pair]pair.Pair)

	gScore := make(map[pair.Pair]float64)
	fScore := make(map[pair.Pair]float64)

	for i := range grid.Cells {
		for j := range grid.Cells[i] {
			gScore[pair.New(i, j)] = math.MaxFloat64
			fScore[pair.New(i, j)] = math.MaxFloat64
		}
	}

	gScore[grid.Start] = 0
	fScore[grid.Start] = grid.h(grid.Start)

	openSet[grid.Start] = struct{}{}

	for len(openSet) > 0 {
		time.Sleep(10 * time.Millisecond)
		var minFscore float64 = 0
		minPair := pair.New(0, 0)

		p := 0
		for pair := range openSet {
			fs := fScore[pair]
			if p == 0 || fs < minFscore {
				minFscore = fs
				minPair = pair
			}
			p++
		}

		current := minPair
		grid.Cells[current.I][current.J].Visited = true
		if current.Eq(grid.End) {
			log.Print("We ended!")
			break
		}

		delete(openSet, current)

		directions := []pair.Pair{pair.Up(), pair.Down(), pair.Left(), pair.Right()}
		for _, dir := range directions {
			neighbor := current.Add(dir)
			if neighbor.InBounds(0, 0, len(grid.Cells), len(grid.Cells[0])) {
				if !grid.Cells[neighbor.I][neighbor.J].IsWall {
					tentative_gScore := gScore[current] + g(current, neighbor)
					if tentative_gScore < gScore[neighbor] {
						cameFrom[neighbor] = current
						gScore[neighbor] = tentative_gScore
						fScore[neighbor] = tentative_gScore + grid.h(neighbor)

						if !slices.Contains(maps.Keys(openSet), neighbor) {
							grid.Cells[neighbor.I][neighbor.J].Added = true
							openSet[neighbor] = struct{}{}
						}
					}
				}
			}
		}
	}

	log.Print(cameFrom)

	cur := grid.End
	for {
		if !cur.Eq(grid.Start) && !cur.Eq(grid.End) {
			grid.Cells[cur.I][cur.J].IsPath = true
			grid.PathLength++
		}

		if _, ok := cameFrom[cur]; ok {
			cur = cameFrom[cur]
		} else {
			if cur == grid.End {
				grid.Status = STATUS_END_NOPATH
			} else {
				grid.Status = STATUS_END_SUCCESS
			}
			break
		}
	}

	fmt.Print("A* Finished\n")
}

func g(a, b pair.Pair) float64 {
	return 1
}
func (grid Grid) h(a pair.Pair) float64 {
	// return 10 * int(math.Abs(float64(a.I-grid.End.I))+math.Abs(float64(a.J-grid.End.J)))
	dy := math.Abs(float64(a.I - grid.End.I))
	dx := math.Abs(float64(a.J - grid.End.J))
	heuristic := 1 * (dx + dy)
	tb := float64(1)/float64(len(grid.Cells)*len(grid.Cells[0])) + 1 // Tie breaker
	return heuristic * tb
}
