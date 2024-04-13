package mirror

import (
	"github.com/TorchedSammy/circleload/beatmap"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var baseUrl = "https://api.osu.direct"
type OsuDirect struct {
	Options Options
}

func (c OsuDirect) GetMapset(id int) (beatmap.Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/s/%d", baseUrl, id))
	if err != nil {
		return beatmap.Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return beatmap.Mapset{}, ErrMapsetNotFound
	}

	jsondata := struct{
		SetID int `json:"id"`
		Artist string `json:"artist"`
		Title string `json:"title"`
		Creator string `json:"creator"`
		Status string `json:"status"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	set := beatmap.Mapset{
		SetID: jsondata.SetID,
		Artist: jsondata.Artist,
		Title: jsondata.Title,
		Mapper: jsondata.Creator,
		Status: beatmap.StatusFromString(jsondata.Status),
	}

	return set, nil
}

func (c OsuDirect) GetMapsetFromMap(id int) (beatmap.Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/b/%d", baseUrl, id))
	if err != nil {
		return beatmap.Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return beatmap.Mapset{}, ErrMapsetNotFound
	}

	jsondata := struct{
		SetID int `json:"beatmapset_id"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	return c.GetMapset(jsondata.SetID)
}

func (c OsuDirect) Search(query string, options Options) ([]beatmap.Mapset, error) {
	reqUrl, _ := url.Parse(fmt.Sprintf("%s/v2/search?query=%s&amount=%d", baseUrl, query, options.MaxResults))
	// we should only add the mode query if we dont want all gamemodes
	if options.Mode != beatmap.ModeAny {
		q := reqUrl.Query()
		q.Add("mode", fmt.Sprintf("%d", options.Mode))
		reqUrl.RawQuery = q.Encode()
	}

	if options.Status != beatmap.StatusAny {
		q := reqUrl.Query()
		q.Add("status", fmt.Sprintf("%d", options.Status))
		reqUrl.RawQuery = q.Encode()
	}

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return nil, err
	}

	jsondata := []struct{
		SetID int `json:"id"`
		Artist string `json:"artist"`
		Title string `json:"title"`
		Creator string `json:"creator"`
		Status string `json:"status"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	sets := []beatmap.Mapset{}
	for _, set := range jsondata {
		sets = append(sets, beatmap.Mapset{
			SetID: set.SetID,
			Artist: set.Artist,
			Title: set.Title,
			Mapper: set.Creator,
			Status: beatmap.StatusFromString(set.Status),
		})
	}

	return sets, nil
}

func (c OsuDirect) GetMapsetData(id int, options Options) (*http.Response, error) {
	url := fmt.Sprintf("%s/d/%d", baseUrl, id)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrMapsetNotFound
	}

	return resp, nil
}

func (c OsuDirect) MirrorOptions() Options {
	return c.Options
}

func (c OsuDirect) SetOptions(options Options) {
	c.Options = options
}
