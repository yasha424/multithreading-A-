package main

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"fmt"
)

func main() {
	for i := 0; i < 100; i++ {
		start := astar.Node{1, 0}
		dest := astar.Node{98, 99}
		mg := mazeGenerator.NewMazeGenerator(100, 100, start.X, start.Y, dest.X, dest.Y)
		maze := mg.GenerateMaze()
		//maze.Print()

		//start := astar.Node{0, 1}
		//dest := astar.Node{9, 18}
		pair := astar.FindPath(maze, start, dest, mazeGenerator.Distance, mazeGenerator.Distance)

		for _, p := range pair.Path {
			maze.Put(p, '.')
		}
		if pair.Path == nil {
			fmt.Println("path not found")
		}

		maze.WriteToFile("path.txt")
		//maze.Print()
		//fmt.Println("Cost is", pair.Cost)
	}
}
