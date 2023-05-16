package mazeGenerator

import (
	astar "course-work/AStar"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type MazeGenerator struct {
	rows      int
	columns   int
	Maze      Maze
	visited   [][]bool
	startRow  int
	startCol  int
	finishRow int
	finishCol int
}

func NewMazeGenerator(rows, columns, startRow, startCol, finishRow, finishCol int) *MazeGenerator {
	maze := make([][]rune, rows)
	visited := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		maze[i] = make([]rune, columns)
		visited[i] = make([]bool, columns)
		for j := 0; j < columns; j++ {
			maze[i][j] = '#'
		}
	}
	return &MazeGenerator{
		rows:      rows,
		columns:   columns,
		Maze:      maze,
		visited:   visited,
		startRow:  startRow,
		startCol:  startCol,
		finishRow: finishRow,
		finishCol: finishCol,
	}
}

func (mg *MazeGenerator) GenerateMaze() Maze {
	mg.Maze[mg.startRow][mg.startCol] = ' '
	mg.visited[mg.startRow][mg.startCol] = true

	rand.Seed(time.Now().UnixNano())
	mg.recursiveBacktracking(mg.startRow, mg.startCol)

	for row := 0; row < mg.rows; row++ {
		mg.Maze[row][0] = '#'
		mg.Maze[row][mg.columns-1] = '#'
	}

	for col := 0; col < mg.columns; col++ {
		mg.Maze[0][col] = '#'
		mg.Maze[mg.rows-1][col] = '#'
	}

	mg.Maze[mg.startRow][mg.startCol] = ' '
	mg.Maze[mg.finishRow][mg.finishCol] = ' '

	return mg.Maze
}

func (mg *MazeGenerator) recursiveBacktracking(row, col int) {
	directions := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	mg.visited[row][col] = true

	for _, dir := range directions {
		newRow := row + dir[0]*2
		newCol := col + dir[1]*2

		if mg.isValidCell(newRow, newCol) && !mg.visited[newRow][newCol] {
			mg.Maze[newRow][newCol] = ' '
			mg.Maze[row+dir[0]][col+dir[1]] = ' '
			if newRow != mg.rows-1 && newCol != mg.columns-1 {
				mg.recursiveBacktracking(newRow, newCol)
			}
		}
	}
}

func (mg *MazeGenerator) isValidCell(row, col int) bool {
	return row >= 0 && row < mg.rows && col >= 0 && col < mg.columns
}

//func (mg *MazeGenerator) PrintMaze() {
//	for _, row := range mg.maze {
//		for _, cell := range row {
//			if cell ==  {
//				fmt.Print(" ") // Open path
//			} else {
//				fmt.Print("#") // Blocked wall
//			}
//		}
//		fmt.Println()
//	}
//}

func Distance(p, q astar.Node) float64 {
	return math.Abs(float64(p.X-q.X)) + math.Abs(float64(p.Y-q.Y))
}

type Maze [][]rune

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
	return f.isInBounds(node) && f[node.X][node.Y] == ' '
}

func (f Maze) isInBounds(node astar.Node) bool {
	return node.Y >= 0 && node.X >= 0 && node.X < len(f) && node.Y < len(f[node.X])
}

func (f Maze) Put(node astar.Node, c rune) {
	f[node.X][node.Y] = c
}

func (f Maze) Print() {
	for _, row := range f {
		for _, c := range row {
			fmt.Print(string(c))
		}
		fmt.Println()
	}
}
