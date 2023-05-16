package main

import (
	"container/heap"
	astar "course-work/AStar"
	pqueue "course-work/PriorityQueue"
	"fmt"
	"math"
)

//func main() {
//	maze := Maze{
//		"#############  ",
//		"#             #",
//		"# # # #########",
//		"#           # #",
//		"### ##### ### #",
//		"#   #         #",
//		"# # ### ### ###",
//		"#   # # # #   #",
//		"# # # # # # ###",
//		"# #       # # #",
//		"# # ######### #",
//		"#         #   #",
//		"# ### # # ### #",
//		"#   # # #     #",
//		"  #############",
//	}
//
//	start := astar.Node{0, 14}
//	dest := astar.Node{14, 0}
//
//	pair := astar.FindPath(maze, start, dest, distance, distance)
//
//	for _, p := range pair.Path {
//		maze.put(p, '.')
//	}
//	maze.print()
//	fmt.Println("Cost is", pair.Cost)
//}

func distance(p, q astar.Node) float64 {
	return math.Abs(float64((q.X - p.X) + (q.Y - p.Y)))
}

type Maze []string

func (f Maze) Neighbours(node astar.Node) []astar.Node {
	offsets := []astar.Node{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}

	res := make([]astar.Node, 0, 4)
	for _, off := range offsets {
		q := node.Add(off)
		if f.isFreeAt(q) {
			res = append(res, q)
		}
	}
	return res
}

func (f Maze) isFreeAt(node astar.Node) bool {
	return f.isInBounds(node) && f[node.Y][node.X] == ' '
}

func (f Maze) isInBounds(node astar.Node) bool {
	return node.Y >= 0 && node.X >= 0 && node.Y < len(f) && node.X < len(f[node.Y])
}

func (f Maze) put(node astar.Node, c rune) {
	f[node.Y] = f[node.Y][:node.X] + string(c) + f[node.Y][node.X+1:]
}

func (f Maze) print() {
	for _, row := range f {
		fmt.Println(row)
	}
}

func main() {
	pq := &pqueue.PriorityQueue[astar.Path]{}
	heap.Init(pq)

	start := astar.Node{14, 0}
	dest := astar.Node{0, 14}
	heap.Push(pq, &pqueue.Item[astar.Path]{Value: astar.NewPath(start)})
	p := heap.Pop(pq).(*pqueue.Item[astar.Path])
	println(p.Value.Last().X, p.Value.Last().Y)

	nb := astar.Node{13, 0}
	newPath := p.Value.Cont(nb)
	heap.Push(pq, &pqueue.Item[astar.Path]{
		Value:    newPath,
		Priority: newPath.Cost(distance) + distance(nb, dest),
	})

	nb = astar.Node{15, 0}
	newPath = p.Value.Cont(nb)
	heap.Push(pq, &pqueue.Item[astar.Path]{
		Value:    newPath,
		Priority: newPath.Cost(distance) + distance(nb, dest),
	})

	nb = astar.Node{13, 1}
	newPath = p.Value.Cont(nb)
	heap.Push(pq, &pqueue.Item[astar.Path]{
		Value:    newPath,
		Priority: newPath.Cost(distance) + distance(nb, dest),
	})

	p = heap.Pop(pq).(*pqueue.Item[astar.Path])
	println(p.Value.Last().X, p.Value.Last().Y)
}
