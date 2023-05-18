package Tests

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"math/rand"
	"testing"
)

func TestCompareConcurrentToSerial(t *testing.T) {
	threadsNum := 8
	mazes := make([]mazeGenerator.Maze, threadsNum)
	graphs := make([]astar.Graph, threadsNum)
	starts := make([]astar.Node, threadsNum)
	ends := make([]astar.Node, threadsNum)

	for i := 0; i < threadsNum; i++ {
		sizeOfMaze := 100
		if i != 0 {
			starts[i] = astar.Node{X: 0, Y: ends[i-1].Y}
		} else {
			starts[i] = astar.Node{X: 0, Y: 1}
		}
		ends[i] = astar.Node{X: sizeOfMaze - 1, Y: 1 + rand.Intn(sizeOfMaze-2)}
		mg := mazeGenerator.NewMazeGenerator(sizeOfMaze, sizeOfMaze, starts[i].X, starts[i].Y, ends[i].X, ends[i].Y)
		mazes[i] = mg.GenerateMaze()
		graphs[i] = mazes[i]
	}

	pairs := astar.FindPaths(graphs, starts, ends, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance, threadsNum)

	serialPairs := make([]astar.Pair, threadsNum)
	for i, graph := range graphs {
		serialPairs[i] = astar.FindPath(graph, starts[i], ends[i], mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	}

	for i, pair := range pairs {
		if !pair.Path.Equals(serialPairs[i].Path) {
			t.Error("Concurrent path is not equal to serial")
		}
	}
}

func TestConcurrentNodeEval(t *testing.T) {
	sizeOfMaze := 500
	start := astar.Node{Y: 1}
	end := astar.Node{X: sizeOfMaze - 1, Y: 1 + rand.Intn(sizeOfMaze-2)}
	mg := mazeGenerator.NewMazeGenerator(sizeOfMaze, sizeOfMaze, start.X, start.Y, end.X, end.Y)
	maze := mg.GenerateMaze()

	pair := astar.FindPathWithConcurrentPriorityEvaluation(maze, start, end, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	serialPair := astar.FindPath(maze, start, end, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)

	if serialPair.Cost != pair.Cost {
		t.Error("Different costs")
	}
}