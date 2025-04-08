package api

import (
	"github.com/ubombar/soa/internal/buffer"
	"github.com/ubombar/soa/internal/datetime"
)

type QuestionHeader struct {
	Kind     string        `buffer:"kind"`     // kind, always question
	Created  datetime.Date `buffer:"created"`  // creation date
	Question string        `buffer:"question"` // actual question string
	From     string        `buffer:"from"`     // this denotes where this note is spawed from
	Tags     []string      `buffer:"tags"`     // tags of the note
}

type LiteratureHeader struct {
	Kind    string        `buffer:"kind"`     // kind, always literature
	Created datetime.Date `buffer:"created"`  // creation date
	Source  string        `buffer:"question"` // actual question string
	From    string        `buffer:"from"`     // this denotes where this note is spawed from
	Tags    []string      `buffer:"tags"`     // tags of the note
}

type MeetingHeader struct {
	Kind     string        `buffer:"kind"`     // kind, always question
	Created  datetime.Date `buffer:"created"`  // creation date
	Question string        `buffer:"question"` // actual question string
	From     string        `buffer:"from"`     // this denotes where this note is spawed from
	Tags     []string      `buffer:"tags"`     // tags of the note
}

type PermanentHeader struct {
	Kind    string        `buffer:"kind"`     // kind, always literature
	Created datetime.Date `buffer:"created"`  // creation date
	Source  string        `buffer:"question"` // actual question string
	From    string        `buffer:"from"`     // this denotes where this note is spawed from
	Tags    []string      `buffer:"tags"`     // tags of the note
}

func QuestionFromBuffer(b *buffer.Buffer) (*QuestionHeader, error) {
	var header QuestionHeader
	if err := b.ReadHeader(&header, false); err != nil {
		return nil, err
	}
	header.Kind = "question"
	return &header, nil
}

func LiteratureFromBuffer(b *buffer.Buffer) (*LiteratureHeader, error) {
	var header LiteratureHeader
	if err := b.ReadHeader(&header, false); err != nil {
		return nil, err
	}
	header.Kind = "literature"
	return &header, nil
}

func MeetingFromBuffer(b *buffer.Buffer) (*MeetingHeader, error) {
	var header MeetingHeader
	if err := b.ReadHeader(&header, false); err != nil {
		return nil, err
	}
	header.Kind = "meeting"
	return &header, nil
}

func PermanentFromBuffer(b *buffer.Buffer) (*PermanentHeader, error) {
	var header PermanentHeader
	if err := b.ReadHeader(&header, false); err != nil {
		return nil, err
	}
	header.Kind = "permanent"
	return &header, nil
}
