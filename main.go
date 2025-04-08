package main

import (
	"fmt"

	"github.com/ubombar/soa/internal/buffer"
	"github.com/ubombar/soa/internal/datetime"
)

type QuestionHeader struct {
	Kind      string        `buffer:"kind"`
	Question  string        `buffer:"question"`
	Attendees []string      `buffer:"attendees"`
	Test      datetime.Date `buffer:"test"`
}

func main() {
	b := buffer.NewBuffer()
	if err := b.FromFile("example.md"); err != nil {
		panic(err)
	}

	q := &QuestionHeader{
		Question:  "lol",
		Kind:      "question",
		Attendees: []string{"hello", "hi"},
	}

	if err := b.ReadHeader(q, false); err != nil {
		panic(err)
	}

	q.Attendees = append(q.Attendees, "one")
	q.Test = datetime.CurrentDate()
	fmt.Printf("q.Attendees: %v\n", q.Attendees)

	if err := b.WriteHeader(q, false); err != nil {
		panic(err)
	}

	fmt.Printf("q.Kind: %v\n", q.Kind)
	fmt.Printf("q.Question: %v\n", q.Question)

	fmt.Printf("b.Header: %v\n", b.Header)
	fmt.Printf("b.Content: %v\n", b.Content)

	if err := b.ToFile("example.md"); err != nil {
		panic(err)
	}
}
