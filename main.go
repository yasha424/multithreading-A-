package main

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	threadsNum := 2
	sizeOfMaze := 1000
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

	pairs := make([]astar.Pair, threadsNum)

	newMaze := make(mazeGenerator.Maze, sizeOfMaze)
	for i, maze := range mazes {
		for j, row := range maze {
			newMaze[i*len(maze)+j] = make([]rune, sizeOfMaze)
			newMaze[i*len(maze)+j] = row
		}
	}

	serialStartTime := time.Now().UnixNano()
	//newMaze.Print()
	//for i := 0; i < threadsNum; i++ {
	fmt.Println("start")
	serialPair := astar.FindPath(newMaze, start[0], dest[threadsNum-1], mazeGenerator.Distance, mazeGenerator.Distance)
	//}
	serialEndTime := time.Now().UnixNano()
	fmt.Println("Serial execution time:", float64(serialEndTime-serialStartTime)/1_000_000_000)

	concurrentStartTime := time.Now().UnixNano()
	pairs = astar.FindPaths(graphs, start, dest, mazeGenerator.Distance, mazeGenerator.Distance, threadsNum)
	concurrentEndTime := time.Now().UnixNano()
	fmt.Println("Concurrent execution time:", float64(concurrentEndTime-concurrentStartTime)/1_000_000_000)

	fileName := "path.txt"
	file, _ := os.Create(fileName)
	defer file.Close()

	for i := 0; i < threadsNum; i++ {
		for _, c := range pairs[i].Path {
			mazes[i].Put(c, '.')
		}

		for _, row := range mazes[i] {
			for _, c := range row {
				file.WriteString(string(c))
			}
			file.WriteString("\n")
		}
	}
	file.Close()

	fileName = "serialPath.txt"
	file, _ = os.Create(fileName)
	defer file.Close()

	//for i := 0; i < threadsNum; i++ {
	for _, c := range serialPair.Path {
		newMaze.Put(c, '.')
	}

	for _, row := range newMaze {
		for _, c := range row {
			file.WriteString(string(c))
		}
		file.WriteString("\n")
	}
	//}

}
