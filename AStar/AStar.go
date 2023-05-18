package astar

import (
	"container/heap"
	pqueue "course-work/PriorityQueue"
	"sync"
)

type Node struct {
	X int
	Y int
}

func (node Node) Equals(tn Node) bool {
	return node.X == tn.X && node.Y == tn.Y
}

func (node Node) Add(n Node) Node {
	return Node{node.X + n.X, node.Y + n.Y}
}

type Graph interface {
	Neighbours(node Node) []Node
}

type CostFunc func(a, b Node) float64

type Pair struct {
	Path Path
	Cost int
}

type Path []Node

func (p Path) Equals(tp Path) bool {
	if len(p) != len(tp) {
		return false
	}
	for i := range p {
		if !p[i].Equals(tp[i]) {
			return false
		}
	}

	return true
}

func NewPath(start Node) Path {
	return []Node{start}
}

func (p Path) Last() Node {
	return p[len(p)-1]
}

func (p Path) Cont(n Node) Path {
	newPath := make([]Node, len(p), len(p)+1)
	copy(newPath, p)
	newPath = append(newPath, n)
	return newPath
}

func (p Path) Cost(d CostFunc) (c float64) {
	for i := 1; i < len(p); i++ {
		c += d(p[i-1], p[i])
	}
	return c
}

func FindPath(g Graph, start, dest Node, d, h CostFunc) Pair {
	closed := make(map[Node]bool)

	pq := &pqueue.PriorityQueue[Path]{}
	heap.Init(pq)
	heap.Push(pq, &pqueue.Item[Path]{Value: NewPath(start)})

	for pq.Len() > 0 {
		p := heap.Pop(pq).(*pqueue.Item[Path])
		n := p.Value.Last()
		if closed[n] {
			continue
		}
		if n == dest {
			return Pair{p.Value, int(p.Priority)}
		}
		closed[n] = true

		for _, nb := range g.Neighbours(n) {
			newPath := p.Value.Cont(nb)
			heap.Push(pq, &pqueue.Item[Path]{
				Value:    newPath,
				Priority: newPath.Cost(d) + h(nb, dest),
			})
		}
	}

	return Pair{nil, 0}
}

func FindPaths(g []Graph, start, dest []Node, d, h CostFunc, threadsNum int) []Pair {
	paths := make([]Pair, len(g))

	var wg sync.WaitGroup
	for i := 0; i < len(g)/threadsNum; i++ {
		if i == len(g)/threadsNum-1 && float32(len(g))/float32(threadsNum) != float32(len(g)/threadsNum) {
			wg.Add(threadsNum)
			for j := 0; j < threadsNum; j++ {
				go func(i int) {
					paths[i] = FindPath(g[i], start[i], dest[i], d, h)
					wg.Done()
				}(i*threadsNum + j)
			}

			wg.Add(len(g) - (i*threadsNum + threadsNum))
			for j := i*threadsNum + threadsNum; j < len(g); j++ {
				go func(i int) {
					paths[i] = FindPath(g[i], start[i], dest[i], d, h)
					wg.Done()
				}(j)
			}

		} else {
			wg.Add(threadsNum)
			for j := 0; j < threadsNum; j++ {
				go func(i int) {
					paths[i] = FindPath(g[i], start[i], dest[i], d, h)
					wg.Done()
				}(i*threadsNum + j)
			}
		}
	}
	wg.Wait()
	return paths
}

func FindPathWithConcurrentPriorityEvaluation(g Graph, start, dest Node, d, h CostFunc) Pair {
	closed := make(map[Node]bool)

	pq := &pqueue.PriorityQueue[Path]{}
	heap.Init(pq)
	heap.Push(pq, &pqueue.Item[Path]{Value: NewPath(start)})

	for pq.Len() > 0 {
		p := heap.Pop(pq).(*pqueue.Item[Path])
		n := p.Value.Last()
		if closed[n] {
			continue
		}
		if n == dest {
			return Pair{p.Value, int(p.Priority)}
		}
		closed[n] = true

		neighbours := g.Neighbours(n)

		type Pair struct {
			priority float64
			path     Path
		}

		pairChan := make(chan Pair, len(neighbours))

		for _, nb := range neighbours {
			go func(n Node) {
				newPath := p.Value.Cont(n)

				pairChan <- Pair{
					priority: newPath.Cost(d) + h(n, dest),
					path:     newPath,
				}
			}(nb)
		}

		for range neighbours {
			pair := <-pairChan
			heap.Push(pq, &pqueue.Item[Path]{
				Value:    pair.path,
				Priority: pair.priority,
			})
		}
	}

	return Pair{nil, 0}
}
