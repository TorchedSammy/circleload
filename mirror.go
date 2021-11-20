package main

import (
	"errors"
	"net/http"
)

var (
	ErrMapsetNotFound = errors.New("mapset not found")
)

type mapsetMirror interface {
	GetMapset(id int) (osuMapset, error)
	Search(query string) ([]osuMapset, error)
	GetMapsetData(id int, opts mirrorOptions) (*http.Response, error)
}

type mirrorOptions struct {
	noVideo bool
}

