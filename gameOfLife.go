package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	r "github.com/gen2brain/raylib-go/raylib"
)

const (
	alive = 1
	dead  = 0
)

var (
	totalRows int
	totalCols int

	totalCells int
	deadCells  int
	aliveCells *int

	shutdownIndex int
	shutdownMode  bool
)

func printMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func main() {
	cellSizeVal := 8
	cellSize := &cellSizeVal
	monitorIndex := 0

	windowWidth := 1920
	windowHeight := 1080
	// small windows to check for HIGH density pixel behaviour
	// windowWidth := 900
	// windowHeight := 500

	fmt.Printf("Height: %v, Width: %v\n", windowHeight, windowWidth)

	totalRows = windowHeight / *cellSize
	totalCols = windowWidth / *cellSize

	totalCells = totalRows * totalCols
	aliveCount := 0
	aliveCells = &aliveCount

	r.InitWindow(int32(windowWidth), int32(windowHeight), "Game Of Life")
	defer r.CloseWindow()

	refreshRate := r.GetMonitorRefreshRate(monitorIndex)
	if refreshRate == 0 {
		// default FPS if GetMonitorRefreshRate fails
		refreshRate = 60
	}

	fmt.Println("Current refresh rate is: ", refreshRate)

	mainGrid := MakeGrid()
	GenGrid(mainGrid)
	// initial calculation of the ALIVE cells
	for row := 0; row < totalRows; row++ {
		for col := 0; col < totalCols; col++ {
			if mainGrid[row][col] == alive {
				aliveCount++
			}
		}
	}

	ticker := time.NewTicker(time.Duration(refreshRate))
	defer ticker.Stop()

	genCounter := 0

	printMemoryUsage()
	for !r.WindowShouldClose() {
		r.BeginDrawing()
		select {
		case <-ticker.C:
			r.ClearBackground(r.Black)
			printGrid(mainGrid, *cellSize)

			if r.IsKeyPressed(r.KeyW) {
				*cellSize += 1
				totalRows = windowHeight / *cellSize
				totalCols = windowWidth / *cellSize
				totalCells = totalRows * totalCols
				mainGrid = MakeGrid()
				GenGrid(mainGrid)
			}
			if r.IsKeyPressed(r.KeyE) {
				fmt.Println("*cellSize is: ", *cellSize)
				if *cellSize > 3 {
					*cellSize -= 1
					totalRows = windowHeight / *cellSize
					totalCols = windowWidth / *cellSize
					totalCells = totalRows * totalCols
					mainGrid = MakeGrid()
					GenGrid(mainGrid)
				}
			}

			if r.IsKeyPressed(r.KeyG) {
				GenGrid(mainGrid)
			}

			nextGen := MakeGrid()
			ComputeNextGen(mainGrid, nextGen, aliveCells)
			mainGrid, nextGen = nextGen, mainGrid
			genCounter++

			deadCells = totalCells - *(aliveCells)

			r.DrawFPS(50, 300)
			r.DrawText(fmt.Sprintf("Generation: %d", genCounter), 50, 30, 50, r.White)
			r.DrawText(fmt.Sprintf("TotalCells: %d", totalCells), 50, 90, 10, r.Gray)
			// r.DrawText(fmt.Sprintf("Alive: %d", *aliveCells), 50, 150, 50, r.Green)
			// r.DrawText(fmt.Sprintf("Dead: %d", deadCells), 50, 200, 50, r.Red)
			// r.DrawText(fmt.Sprintf("CellSize: %d", cellSize), 50, 400, 50, r.RayWhite)
		}
		r.EndDrawing()
	}
	printMemoryUsage()
}

func printGrid(mainGrid [][]int, cellSize int) {
	for row := range mainGrid {
		for col := range mainGrid[row] {
			x := int32(10 + col*cellSize)
			y := int32(10 + row*cellSize)
			if mainGrid[row][col] == alive {
				r.DrawRectangle(x, y, int32(cellSize)-2, int32(cellSize)-2, r.DarkGray)
			} else {
				r.DrawRectangle(x, y, int32(cellSize)-2, int32(cellSize)-2, r.Black)
			}
		}
	}
}

func MakeGrid() [][]int {
	grid := make([][]int, totalRows)
	for i := range grid {
		grid[i] = make([]int, totalCols)
	}
	return grid
}

func GenGrid(grid [][]int) {
	for row := range grid {
		for col := range grid[row] {
			grid[row][col] = rand.Intn(2)
		}
	}
}

func ComputeNextGen(grid, nextGen [][]int, aliveCells *int) {
	for row := 0; row < totalRows; row++ {
		for col := 0; col < totalCols; col++ {
			state := grid[row][col]
			neighbourCount := CountNeighbours(grid, row, col)

			if state == alive && (neighbourCount < 2 || neighbourCount > 3) {
				nextGen[row][col] = dead
				*aliveCells--
			} else if state == dead && neighbourCount == 3 {
				nextGen[row][col] = alive
				*aliveCells--
			} else {
				nextGen[row][col] = state
			}
		}
	}
}

func CountNeighbours(grid [][]int, x, y int) int {
	count := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			row := (x + i + totalRows) % totalRows
			col := (y + j + totalCols) % totalCols
			count += grid[row][col]
		}
	}
	return count - grid[x][y]
}
