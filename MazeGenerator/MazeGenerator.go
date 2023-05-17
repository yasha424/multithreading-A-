package mazeGenerator

import (
	astar "course-work/AStar"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
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
	for i := range mg.Maze {
		for j := range mg.Maze[i] {
			mg.Maze[i][j] = '#'
			mg.visited[i][j] = false
		}
	}

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

	start := astar.Node{X: mg.startRow, Y: mg.startCol}
	dest := astar.Node{X: mg.finishRow, Y: mg.finishCol}
	pair := astar.FindPath(mg.Maze, start, dest, ManhattanDistance, ManhattanDistance)

	if pair.Path != nil {
		return mg.Maze
	} else {
		return mg.GenerateMaze()
	}
}

func (mg *MazeGenerator) recursiveBacktracking(row, col int) {
	directions := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	mg.visited[row][col] = true

	for i, dir := range directions {
		nextRow := row + dir[0]
		nextCol := col + dir[1]

		if mg.isValidCell(nextRow, nextCol) && !mg.visited[nextRow][nextCol] {
			if i < len(directions)/2 ||
				(nextRow == mg.finishRow-1 && nextCol == mg.finishCol) ||
				(nextRow == mg.startRow+1 && nextCol == mg.startCol) {
				mg.Maze[nextRow][nextCol] = ' '
			} else {
				mg.Maze[nextRow][nextCol] = '#'
			}

			mg.recursiveBacktracking(nextRow, nextCol)
		}
	}
}

func (mg *MazeGenerator) isValidCell(row, col int) bool {
	return row >= 1 && row < mg.rows-1 && col >= 1 && col < mg.columns-1
}

func ManhattanDistance(p, q astar.Node) float64 {
	return math.Abs(float64(p.X-q.X)) + math.Abs(float64(p.Y-q.Y))
}

func EuclidianDistance(p, q astar.Node) float64 {
	return math.Sqrt(math.Pow(float64(p.X-q.X), 2) + math.Pow(float64(p.Y-q.Y), 2))
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

func (f Maze) Equals(m Maze) bool {
	if len(m) == len(f) {
		if len(m) == 0 {
			return true
		} else {
			for i := 0; i < len(m); i++ {
				if len(m[i]) != len(f[i]) {
					return false
				}
				for j := 0; j < len(m[i]); j++ {
					if m[i][j] != f[i][j] {
						return false
					}
				}
			}
			return true
		}
	} else {
		return false
	}
}

func (f Maze) Print() {
	for _, row := range f {
		for _, c := range row {
			fmt.Print(string(c))
		}
		fmt.Println()
	}
}

func (f Maze) WriteToFile(fileName string) {
	file, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	for _, row := range f {
		for _, c := range row {
			_, err := file.WriteString(string(c))
			if err != nil {
				log.Fatal(err)
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WriteToFile(fileName string, mazes []Maze) {
	file, _ := os.Create(fileName)
	defer file.Close()

	for _, maze := range mazes {
		for _, row := range maze {
			for _, c := range row {
				file.WriteString(string(c))
			}
			file.WriteString("\n")
		}
	}
}
