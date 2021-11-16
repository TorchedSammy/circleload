package main

import (
	"errors"
	"io"
)

var (
	ErrMapsetNotFound = errors.New("mapset not found")
)

type mapsetMirror interface {
	GetMapset(id int) (osuMapset, error)
	GetMapsetData(id int) (io.ReadCloser, error)
}

