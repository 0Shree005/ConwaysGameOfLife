package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	totalRows = 40
	totalCols = 70
	alive     = 1
	dead      = 0
)

func main() {
	grid := makeGrid()
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	genGrid(grid)
	genCounter := 0

	for {
		select {
		case <-ticker.C:
			clearScreen()
			fmt.Printf("Generation: %d\n", genCounter)
			printGrid(grid)

			nextGen := makeGrid()

			computeNextGen(grid, nextGen)
			grid, nextGen = nextGen, grid
			genCounter++
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func makeGrid() [][]int {
	grid := make([][]int, totalRows)
	for i := range grid {
		grid[i] = make([]int, totalCols)
	}
	return grid
}

func genGrid(grid [][]int) {
	for row := range grid {
		for col := range grid[row] {
			grid[row][col] = rand.Intn(2)
		}
	}
}

func computeNextGen(grid, nextGen [][]int) {
	for row := 0; row < totalRows; row++ {
		for col := 0; col < totalCols; col++ {
			state := grid[row][col]
			neighbourCount := countNeighbours(grid, row, col)

			if state == alive && (neighbourCount < 2 || neighbourCount > 3) {
				nextGen[row][col] = dead
			} else if state == dead && neighbourCount == 3 {
				nextGen[row][col] = alive
			} else {
				nextGen[row][col] = state
			}
		}
	}
}

func countNeighbours(grid [][]int, x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			row := (x + i + totalRows) % totalRows
			col := (y + j + totalCols) % totalCols
			count += grid[row][col]
		}
	}
	return count - grid[x][y]
}

func printGrid(grid [][]int) {
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] == dead {
				fmt.Print("  ")
			} else {
				fmt.Print("▄ ")
			}
		}
		fmt.Println()
	}
}
