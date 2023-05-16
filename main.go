package main

import (
	astar "course-work/AStar"
	mazeGenerator "course-work/MazeGenerator"
	"fmt"
)

//func main() {
//	maze := Maze{
//		"#############  ",
//		"#             #",
//		"# # # #########",
//		"#           # #",
//		"### ##### ### #",
//		"#   #         #",
//		"# # ### ### ###",
//		"#   # # # #   #",
//		"# # # # # # ###",
//		"# #       # # #",
//		"# # ######### #",
//		"#         #   #",
//		"# ### # # ### #",
//		"#   # # #     #",
//		"  #############",
//	}
//
//	start := astar.Node{0, 14}
//	dest := astar.Node{14, 0}
//
//	pair := astar.FindPath(maze, start, dest, distance, distance)
//
//	for _, p := range pair.Path {
//		maze.put(p, '.')
//	}
//	maze.print()
//	fmt.Println("Cost is", pair.Cost)
//}

func main() {
	mg := mazeGenerator.NewMazeGenerator(10, 10, 0, 1, 9, 8)
	maze := mg.GenerateMaze()
	//maze.Print()

	start := astar.Node{0, 1}
	dest := astar.Node{9, 8}
	pair := astar.FindPath(maze, start, dest, mazeGenerator.Distance, mazeGenerator.Distance)

	for _, p := range pair.Path {
		maze.Put(p, '.')
	}
	maze.Print()
	fmt.Println("Cost is", pair.Cost)
}
