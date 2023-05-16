package main

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"fmt"
	"strconv"
	"time"
)

func main() {

	isConcurrent := true

	size := 100
	start := make([]astar.Node, size)
	dest := make([]astar.Node, size)
	graphs := make([]astar.Graph, size)
	mazes := make([]mazeGenerator.Maze, size)

	for i := 0; i < size; i++ {
		start[i] = astar.Node{1, 0}
		dest[i] = astar.Node{98, 99}

		mg := mazeGenerator.NewMazeGenerator(100, 100, start[0].X, start[0].Y, dest[0].X, dest[0].Y)
		maze := mg.GenerateMaze()
		graphs[i] = maze
		mazes[i] = maze
	}

	pairs := make([]astar.Pair, size)

	startTime := time.Now().UnixNano()
	if !isConcurrent {
		for i := 0; i < size; i++ {
			pairs[i] = astar.FindPath(graphs[i], start[i], dest[i], mazeGenerator.Distance, mazeGenerator.Distance)
		}
	} else {
		pairs = astar.FindPaths(graphs, start, dest, mazeGenerator.Distance, mazeGenerator.Distance, 8)
	}
	endTime := time.Now().UnixNano()
	fmt.Println(float64(endTime-startTime) / 100000000)

	for i := 0; i < size; i++ {
		for _, c := range pairs[i].Path {
			mazes[i].Put(c, '.')
		}

		fileName := "path" + strconv.FormatInt(int64(i), 10) + ".txt"
		//fmt.Println(fileName)
		mazes[i].WriteToFile(fileName)
	}

	//maze.Print()
	//fmt.Println("Cost is", pair.Cost)
}
