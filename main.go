package main

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//serialSearch(1000)
	//divideGraphTest(1000, 4)
	//divideGraphTestAvg(800, 4, 10)
	//concurrentPriorityEvaluation(1000)
	bidirectionalSearch(500)
	//bidirectionalSearchAvg(600, 5)
}

func bidirectionalSearchAvg(sizeOfMaze, iterationsNum int) {
	concurrentTimeSum := 0.0
	serialTimeSum := 0.0

	for i := 0; i < iterationsNum; i++ {

		startNode := astar.Node{X: 0, Y: 1}
		endNode := astar.Node{X: sizeOfMaze - 1, Y: sizeOfMaze - 2}
		mg := mazeGenerator.NewMazeGenerator(sizeOfMaze, sizeOfMaze, startNode.X, startNode.Y, endNode.X, endNode.Y)
		maze := mg.GenerateMaze()

		fmt.Println("Starting search", i)

		concurrentStartTime := time.Now().UnixNano()
		concurrentPath := astar.FindPathWithBidirectionalSearch(maze, startNode, endNode, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
		concurrentEndTime := time.Now().UnixNano()

		serialStartTime := time.Now().UnixNano()
		serialPath := astar.FindPath(maze, startNode, endNode, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
		serialEndTime := time.Now().UnixNano()

		concurrentTime := float64(concurrentEndTime - concurrentStartTime)
		serialTime := float64(serialEndTime - serialStartTime)

		concurrentTimeSum += concurrentTime
		serialTimeSum += serialTime

		if serialPath.Cost != concurrentPath.Cost {
			fmt.Println("Paths costs aren't equal!")
		}
	}

	fmt.Println("Average concurrent execution time:", concurrentTimeSum/float64(iterationsNum), "ns")
	fmt.Println("Average serial execution time:", serialTimeSum/float64(iterationsNum), "ns")
	fmt.Println("Average speedup:", serialTimeSum/concurrentTimeSum)
}

func bidirectionalSearch(sizeOfMaze int) {
	startNode := astar.Node{X: 0, Y: 1}
	endNode := astar.Node{X: sizeOfMaze - 1, Y: sizeOfMaze - 2}
	mg := mazeGenerator.NewMazeGenerator(sizeOfMaze, sizeOfMaze, startNode.X, startNode.Y, endNode.X, endNode.Y)
	maze := mg.GenerateMaze()

	concurrentStartTime := time.Now().UnixNano()
	concurrentPath := astar.FindPathWithBidirectionalSearch(maze, startNode, endNode, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	concurrentEndTime := time.Now().UnixNano()

	serialStartTime := time.Now().UnixNano()
	serialPath := astar.FindPath(maze, startNode, endNode, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	serialEndTime := time.Now().UnixNano()

	fmt.Println("Serial Search cost:", serialPath.Cost)
	fmt.Println("Concurrent Search cost:", concurrentPath.Cost)

	concurrentTime := float64(concurrentEndTime - concurrentStartTime)
	fmt.Println("Concurrent execution time:", concurrentTime/1_000_000, "ms")

	serialTime := float64(serialEndTime - serialStartTime)
	fmt.Println("Serial execution time:", serialTime/1_000_000, "ms")

	fmt.Println("Speedup:", serialTime/concurrentTime)

	if serialPath.Cost != concurrentPath.Cost {
		fmt.Println("Paths costs aren't equal!")
	}

	serialMaze := make(mazeGenerator.Maze, sizeOfMaze)

	for i, row := range maze {
		serialMaze[i] = make([]rune, sizeOfMaze)
		for j, c := range row {
			serialMaze[i][j] = c
		}
	}

	for _, node := range concurrentPath.Path {
		maze.Put(node, '.')
	}

	for _, node := range serialPath.Path {
		serialMaze.Put(node, '.')
	}

	maze.WriteToFile("bidirectionalSearch.txt")
	serialMaze.WriteToFile("serialPath.txt")
}

func divideGraphTestAvg(sizeOfMaze, threadsNum, iterationsNum int) {
	concurrentTimeSum := 0.0
	serialTimeSum := 0.0

	for i := 0; i < iterationsNum; i++ {
		start := make([]astar.Node, threadsNum)
		dest := make([]astar.Node, threadsNum)
		graphs := make([]astar.Graph, threadsNum)
		mazes := make([]mazeGenerator.Maze, threadsNum)
		serialPaths := make([]astar.Pair, threadsNum)

		destY := 1
		for j := 0; j < threadsNum; j++ {

			start[j] = astar.Node{X: 0, Y: destY}
			destY = 1 + rand.Intn(sizeOfMaze-2)
			dest[j] = astar.Node{X: sizeOfMaze/threadsNum - 1, Y: destY}

			mg := mazeGenerator.NewMazeGenerator(sizeOfMaze/threadsNum, sizeOfMaze, start[j].X, start[j].Y, dest[j].X, dest[j].Y)
			maze := mg.GenerateMaze()
			graphs[j] = maze
			mazes[j] = maze
		}

		newMaze := make(mazeGenerator.Maze, sizeOfMaze)
		for j, maze := range mazes {
			for k, row := range maze {
				newMaze[j*len(maze)+k] = make([]rune, sizeOfMaze)
				for l, c := range row {
					newMaze[j*len(maze)+k][l] = c
				}
			}
		}

		fmt.Println("Starting search", i)

		serialStartTime := time.Now().UnixNano()
		for j := 0; j < threadsNum; j++ {
			serialPaths[j] = astar.FindPath(mazes[j], start[j], dest[j], mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
		}
		serialTimeSum += float64(time.Now().UnixNano() - serialStartTime)

		concurrentStartTime := time.Now().UnixNano()
		concurrentPaths := astar.FindPaths(graphs, start, dest, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance, threadsNum)
		concurrentTimeSum += float64(time.Now().UnixNano() - concurrentStartTime)
		for j := range concurrentPaths {
			if !concurrentPaths[j].Path.Equals(serialPaths[j].Path) {
				fmt.Println("Not equal!!!")
			}
		}
	}
	fmt.Println("Average concurrent execution time:", concurrentTimeSum/float64(iterationsNum), "ns")
	fmt.Println("Average serial execution time:", serialTimeSum/float64(iterationsNum), "ns")
	fmt.Println("Average speedup:", serialTimeSum/concurrentTimeSum)
}

func divideGraphTest(sizeOfMaze, threadsNum int) {
	start := make([]astar.Node, threadsNum)
	dest := make([]astar.Node, threadsNum)
	graphs := make([]astar.Graph, threadsNum)
	mazes := make([]mazeGenerator.Maze, threadsNum)

	destY := 1
	for i := 0; i < threadsNum; i++ {

		start[i] = astar.Node{X: 0, Y: destY}
		destY = 1 + rand.Intn(sizeOfMaze-2)
		dest[i] = astar.Node{X: sizeOfMaze/threadsNum - 1, Y: destY}

		mg := mazeGenerator.NewMazeGenerator(sizeOfMaze/threadsNum, sizeOfMaze, start[i].X, start[i].Y, dest[i].X, dest[i].Y)
		maze := mg.GenerateMaze()
		graphs[i] = maze
		mazes[i] = maze
	}

	serialPairs := make([]astar.Pair, threadsNum)
	concurrentPairs := make([]astar.Pair, threadsNum)

	newMaze := make(mazeGenerator.Maze, sizeOfMaze)
	for i, maze := range mazes {
		for j, row := range maze {
			newMaze[i*len(maze)+j] = make([]rune, sizeOfMaze)
			for k, c := range row {
				newMaze[i*len(maze)+j][k] = c
			}
		}
	}

	fmt.Println("Starting search...")

	serialStartTime := time.Now().UnixNano()
	for i := 0; i < threadsNum; i++ {
		serialPairs[i] = astar.FindPath(mazes[i], start[i], dest[i], mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	}
	serialTime := float64(time.Now().UnixNano() - serialStartTime)
	fmt.Println("Serial execution time:", serialTime/1_000_000, "ms")

	concurrentStartTime := time.Now().UnixNano()
	concurrentPairs = astar.FindPaths(graphs, start, dest, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance, threadsNum)
	concurrentTime := float64(time.Now().UnixNano() - concurrentStartTime)
	fmt.Println("Concurrent execution time:", concurrentTime/1_000_000, "ms")

	concurrentCost := 0
	serialCost := 0
	for i := 0; i < threadsNum; i++ {
		concurrentCost += concurrentPairs[i].Cost
		serialCost += serialPairs[i].Cost
	}
	fmt.Println("Concurrent cost:", concurrentCost)
	fmt.Println("Concurrent cost:", serialCost)

	fmt.Println("Speedup:", serialTime/concurrentTime)
	serialMazes := mazes

	for i, pair := range concurrentPairs {
		for _, node := range pair.Path {
			mazes[i].Put(node, '.')
		}
	}

	mazeGenerator.WriteToFile("concurrentPath.txt", mazes)

	for i, pair := range serialPairs {
		for _, node := range pair.Path {
			serialMazes[i].Put(node, '.')
		}
	}

	mazeGenerator.WriteToFile("serialPath.txt", serialMazes)
}

func concurrentPriorityEvaluation(sizeOfMaze int) {
	start := astar.Node{Y: 1}
	dest := astar.Node{X: sizeOfMaze - 1, Y: rand.Intn(sizeOfMaze)}
	mg := mazeGenerator.NewMazeGenerator(sizeOfMaze, sizeOfMaze, start.X, start.Y, dest.X, dest.Y)
	concurrentMaze := mg.GenerateMaze()

	fmt.Println("Starting search...")
	concurrentStart := time.Now().UnixNano()
	concurrentPair := astar.FindPathWithConcurrentPriorityEvaluation(concurrentMaze, start, dest, mazeGenerator.ManhattanDistance, mazeGenerator.EuclideanDistance)
	concurrentTime := float64(time.Now().UnixNano() - concurrentStart)
	fmt.Println("Concurrent execution time:", concurrentTime, "ns")

	serialMaze := make(mazeGenerator.Maze, sizeOfMaze)

	for i, row := range concurrentMaze {
		serialMaze[i] = make([]rune, sizeOfMaze)
		for j, c := range row {
			serialMaze[i][j] = c
		}
	}

	serialStart := time.Now().UnixNano()
	serialPair := astar.FindPath(serialMaze, start, dest, mazeGenerator.ManhattanDistance, mazeGenerator.EuclideanDistance)
	serialTime := float64(time.Now().UnixNano() - serialStart)
	fmt.Println("Serial execution time:", serialTime, "ns")

	fmt.Println("Speedup:", serialTime/concurrentTime)

	for _, node := range concurrentPair.Path {
		concurrentMaze.Put(node, '.')
	}
	fmt.Println("Concurrent cost:", concurrentPair.Cost)

	concurrentMaze.WriteToFile("concurrentNodePath.txt")

	for _, node := range serialPair.Path {
		serialMaze.Put(node, '.')
	}
	fmt.Println("Serial cost:", serialPair.Cost)

	serialMaze.WriteToFile("serialPath.txt")
}

func serialSearch(sizeOfMaze int) {
	start := astar.Node{Y: 1}
	end := astar.Node{X: sizeOfMaze - 1, Y: 1 + rand.Intn(sizeOfMaze-2)}
	mg := mazeGenerator.NewMazeGenerator(sizeOfMaze, sizeOfMaze, start.X, start.Y, end.X, end.Y)
	maze := mg.GenerateMaze()

	fmt.Println("Starting search...")
	startTime := time.Now().UnixNano()
	pair := astar.FindPath(maze, start, end, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	endTime := time.Now().UnixNano()
	fmt.Println("Execution time:", float64(endTime-startTime), "ns")

	for _, n := range pair.Path {
		maze.Put(n, '.')
	}
	maze.WriteToFile("serialPath.txt")
}
