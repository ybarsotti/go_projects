package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/yaml.v3"
)

type Movie struct {
	Title     string    `yaml:"title"`
	Genres    []string  `yaml:"genres"`
	Year      int       `yaml:"year"`
	Runtime   int32     `yaml:"movie_length_in_minutes"`
	Rating    float32   `yaml:"rating"`
	CreatedAt time.Time `yaml:"created_at"`
}

func main() {
	movie := Movie{
		Title:     "Titanic",
		Genres:    []string{"drama", "romance"},
		Year:      1997,
		Runtime:   197,
		Rating:    7.9,
		CreatedAt: time.Now(),
	}

	movieYAML, err := yaml.Marshal(movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(movieYAML))
}