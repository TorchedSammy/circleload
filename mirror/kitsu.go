package mirror

// This defines the kitsu.moe mirror.
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TorchedSammy/circleload/log"
)

type Kitsu struct {
	Options Options
}

// kitsu is NOT a proper restful api
// on success, it doesnt send code
// but on failure it does (and the actual response code is still 200)
// so to check for errors we just check is Code isnt 0
type kitsuResponse struct {
	Mapset
	Code int
	Message string
}

type kitsuMap struct {
	BeatmapID int
	ParentSetID int
}

type kitsuMapResponse struct {
	kitsuMap
	Code int
	Message string
}

func (k Kitsu) GetMapset(id int) (Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/api/s/%d", id))
	if err != nil {
		return Mapset{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Mapset{}, err
	}

	var apiResp kitsuResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return Mapset{}, err
	}

	if apiResp.Code != 0 {
		switch apiResp.Code {
		case 404:
			return Mapset{}, ErrMapsetNotFound
		}
	}

	set := Mapset{
		SetID: apiResp.SetID,
		Title: apiResp.Title,
		Artist: apiResp.Artist,
	}

	return set, nil
}

func (k Kitsu) GetMapsetFromMap(id int) (Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/api/b/%d", id))
	if err != nil {
		return Mapset{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Mapset{}, err
	}

	var apiResp kitsuMapResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return Mapset{}, err
	}

	if apiResp.Code != 0 {
		switch apiResp.Code {
		case 404:
			return Mapset{}, ErrBeatmapNotFound
		}
	}

	beatmap, err := k.GetMapset(apiResp.ParentSetID)
	if err != nil {
		return Mapset{}, err
	}

	return k.GetMapset(beatmap.SetID)
}

// thanks copilot
func (k Kitsu) Search(query string) ([]Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/api/search?query=%s&amount=%d", query, k.Options.MaxResults))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sets []Mapset
	err = json.Unmarshal(body, &sets)
	if err != nil {
		return nil, err
	}

	return sets, nil
}

// get beatmap from kitsu
func (k Kitsu) GetMapsetData(id int) (*http.Response, error) {
	// kitsu doesnt have a noVideo option
	// log that it doesnt
	if k.Options.NoVideo {
		log.Warn("kitsu mirror doesnt support noVideo")
	}

	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/d/%d", id))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrMapsetNotFound
	}

	return resp, nil
}

