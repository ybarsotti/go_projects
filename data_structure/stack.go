package main

import "errors"

type Stack struct {
    storage []int
}

func (s *Stack) Push(value int) {
    s.storage = append(s.storage, value)
}

func (s *Stack) Pop() (int, error) {
    last := len(s.storage) - 1
    if last <= -1 { // check the size
        return 0, errors.New("Stack is empty") // and return error
    }

    value := s.storage[last]     // save the value
    s.storage = s.storage[:last] // remove the last element

    return value, nil // return saved value and nil error
}