package main

import "errors"

type Queue struct {
    input  Stack
    output Stack
}

func (q *Queue) Push(value int) {
    q.input.Push(value)
}

func (q *Queue) Pop() (int, error) {
    outputVal, outputErr := q.output.Pop()
    if outputErr == nil { // if output stack is not empty
        return outputVal, nil // just return value
    }

    inputVal, inputErr := q.input.Pop()
    if inputErr != nil { // if input stack is empty
        return 0, errors.New("Queue is empty") // return the error
    }

    // if the output stack is empty but the input is not empty
    for inputErr == nil { // while input stack not empty...
        q.output.Push(inputVal)            // rearrange input to output
        inputVal, inputErr = q.input.Pop() // and read again
    }

    return q.output.Pop() // and Pop the output
}