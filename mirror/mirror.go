package mirror

import (
	"errors"
	"net/http"

	"github.com/TorchedSammy/circleload/beatmap"
)

var (
	ErrMapsetNotFound = errors.New("mapset not found")
	ErrBeatmapNotFound = errors.New("map not found")
)

type Mirror interface {
	GetMapset(id int) (beatmap.Mapset, error)
	GetMapsetFromMap(id int) (beatmap.Mapset, error)
	Search(query string, options Options) ([]beatmap.Mapset, error)
	GetMapsetData(id int, options Options) (*http.Response, error)
	MirrorOptions() Options
	SetOptions(options Options)
}

type Options struct {
	NoVideo bool
	MaxResults int
	Mode beatmap.Mode
	Status beatmap.Status
}

