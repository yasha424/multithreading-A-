package main

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	divideGraph(1000, 100)
	//return
	//concurrentPriorityEvaluation(1000)
}

func divideGraph(sizeOfMaze, threadsNum int) {
	start := make([]astar.Node, threadsNum)
	dest := make([]astar.Node, threadsNum)
	graphs := make([]astar.Graph, threadsNum)
	mazes := make([]mazeGenerator.Maze, threadsNum)

	destY := 1
	for i := 0; i < threadsNum; i++ {

		start[i] = astar.Node{X: 0, Y: destY}
		destY = rand.Intn(sizeOfMaze)
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

	serialStartTime := time.Now().UnixNano()

	for i := 0; i < threadsNum; i++ {
		serialPairs[i] = astar.FindPath(mazes[i], start[i], dest[i], mazeGenerator.Distance, mazeGenerator.Distance)
	}
	serialTime := float64(time.Now().UnixNano() - serialStartTime)
	fmt.Println("Serial execution time:", serialTime/1_000_000_000)

	concurrentStartTime := time.Now().UnixNano()
	concurrentPairs = astar.FindPaths(graphs, start, dest, mazeGenerator.Distance, mazeGenerator.Distance, threadsNum)
	concurrentTime := float64(time.Now().UnixNano() - concurrentStartTime)
	fmt.Println("Concurrent execution time:", concurrentTime/1_000_000_000)

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
	concurrentStart := time.Now().UnixNano()
	concurrentPair := astar.FindPathWithConcurrentPriorityEvaluation(concurrentMaze, start, dest, mazeGenerator.Distance, mazeGenerator.EuclidianDistance)
	concurrentTime := float64(time.Now().UnixNano() - concurrentStart)
	fmt.Println("Concurrent execution time:", concurrentTime/1_000_000_000)

	serialMaze := make(mazeGenerator.Maze, sizeOfMaze)

	for i, row := range concurrentMaze {
		serialMaze[i] = make([]rune, sizeOfMaze)
		for j, c := range row {
			serialMaze[i][j] = c
		}
	}

	serialStart := time.Now().UnixNano()
	serialPair := astar.FindPath(serialMaze, start, dest, mazeGenerator.Distance, mazeGenerator.EuclidianDistance)
	serialTime := float64(time.Now().UnixNano() - serialStart)
	fmt.Println("Serial execution time:", serialTime/1_000_000_000)

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
