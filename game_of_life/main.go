package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
    var N int
    fmt.Print("Digite o tamanho do universo (N): ")
    fmt.Scanf("%d", &N)

    rand.Seed(time.Now().UnixNano())

    currentGen := make([][]rune, N)
    nextGen := make([][]rune, N)
    for i := range currentGen {
        currentGen[i] = make([]rune, N)
        nextGen[i] = make([]rune, N)
        for j := range currentGen[i] {
            if rand.Intn(2) == 1 {
                currentGen[i][j] = 'O'
            } else {
                currentGen[i][j] = ' '
            }
        }
    }

    countLiveNeighbors := func(grid [][]rune, x, y int) int {
        directions := [8][2]int{
            {-1, -1}, {-1, 0}, {-1, 1},
            {0, -1},          {0, 1},
            {1, -1}, {1, 0}, {1, 1},
        }
        count := 0
        for _, d := range directions {
            nx, ny := (x + d[0] + N) % N, (y + d[1] + N) % N
            if grid[nx][ny] == 'O' {
                count++
            }
        }
        return count
    }

    calculateNextGen := func(current, next [][]rune) {
        for i := 0; i < N; i++ {
            for j := 0; j < N; j++ {
                liveNeighbors := countLiveNeighbors(current, i, j)
                if current[i][j] == 'O' {
                    if liveNeighbors == 2 || liveNeighbors == 3 {
                        next[i][j] = 'O'
                    } else {
                        next[i][j] = ' '
                    }
                } else {
                    if liveNeighbors == 3 {
                        next[i][j] = 'O'
                    } else {
                        next[i][j] = ' '
                    }
                }
            }
        }
    }

    printGeneration := func(grid [][]rune, gen int) {
        aliveCount := 0
        fmt.Printf("Generation #%d\n", gen)
        for i := 0; i < N; i++ {
            for j := 0; j < N; j++ {
                if grid[i][j] == 'O' {
                    aliveCount++
                }
            }
        }
        fmt.Printf("Alive: %d\n", aliveCount)
        for i := 0; i < N; i++ {
            for j := 0; j < N; j++ {
                fmt.Print(string(grid[i][j]))
            }
            fmt.Println()
        }
    }

    generations := 20
    for g := 1; g <= generations; g++ {
        printGeneration(currentGen, g)
        calculateNextGen(currentGen, nextGen)
        currentGen, nextGen = nextGen, currentGen
        time.Sleep(500 * time.Millisecond)
        fmt.Print("\033[H\033[2J") 
    }
}
