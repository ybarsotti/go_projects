package main

import "fmt"

func main() {
    var word string
    var stackSolver Stack
    var queueSolver Queue

    fmt.Scan(&word)

    for _, char := range word {
        stackSolver.Push(int(char))
        queueSolver.Push(int(char))
    }

    for {
        fromStack, err := stackSolver.Pop()
        fromQueue, err1 := queueSolver.Pop()

        if fromStack == fromQueue {
            fmt.Println("Palindrome")
            break
        }

        if err != nil || err1 != nil || fromStack != fromQueue {
			fmt.Println("Not palindrome")
			break
		}
    }
}