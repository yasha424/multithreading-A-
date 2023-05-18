package Tests

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"math/rand"
	"testing"
)

func TestSerialSearch(t *testing.T) {
	sizeOfMaze := 500
	start := astar.Node{Y: 1}
	end := astar.Node{X: sizeOfMaze - 1, Y: 1 + rand.Intn(sizeOfMaze-2)}
	mg := mazeGenerator.NewMazeGenerator(sizeOfMaze, sizeOfMaze, start.X, start.Y, end.X, end.Y)
	maze := mg.GenerateMaze()

	pair := astar.FindPath(maze, start, end, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)

	if pair.Path == nil {
		t.Error("Path not found")
	}
}

func TestNoPathSearch(t *testing.T) {
	maze := mazeGenerator.Maze{
		{'#', ' ', '#', '#', '#', '#'},
		{'#', ' ', ' ', ' ', ' ', '#'},
		{'#', ' ', ' ', ' ', ' ', '#'},
		{'#', ' ', ' ', ' ', ' ', '#'},
		{'#', ' ', ' ', ' ', '#', '#'},
		{'#', '#', '#', '#', ' ', '#'},
	}
	start := astar.Node{Y: 1}
	end := astar.Node{X: 5, Y: 4}

	pair := astar.FindPath(maze, start, end, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	if pair.Path != nil {
		t.Error("Path found, but not exists")
	}
}

func TestSearch(t *testing.T) {
	maze := mazeGenerator.Maze{
		{'#', ' ', '#', '#', '#', '#', '#'},
		{'#', ' ', ' ', '#', ' ', ' ', '#'},
		{'#', ' ', '#', ' ', ' ', ' ', '#'},
		{'#', ' ', ' ', ' ', '#', ' ', '#'},
		{'#', ' ', ' ', ' ', '#', ' ', '#'},
		{'#', '#', '#', '#', '#', ' ', '#'},
	}
	start := astar.Node{Y: 1}
	end := astar.Node{X: 5, Y: 5}

	pair := astar.FindPath(maze, start, end, mazeGenerator.ManhattanDistance, mazeGenerator.ManhattanDistance)
	if pair.Path == nil {
		t.Error("Path not found")
	}

	if pair.Cost != 11 {
		t.Error("Cost is not optimal")
	}
}
