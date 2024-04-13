package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/TorchedSammy/circleload/beatmap"
	"github.com/TorchedSammy/circleload/log"
	"github.com/manifoldco/promptui"
)

var errNoResults = errors.New("no results")
var errSearchAborted = errors.New("search aborted")

func applyFilters(matches [][]string, query string) string {
	arg := query
	for _, m := range matches {
		// expected format: [group, key, value]
		if len(m) == 3 {
			group, key, val := m[0], m[1], m[2]
			// check key
			switch key {
				case "mode": // specific gamemode
				modes := []beatmap.Mode{beatmap.ModeStandard, beatmap.ModeTaiko, beatmap.ModeCatch, beatmap.ModeMania, beatmap.ModeAny}

				for i, mode := range modes {
					if val == "osu" {
						mirrorOptions := dlmirror.MirrorOptions()
						mirrorOptions.Mode = beatmap.ModeStandard

						dlmirror.SetOptions(mirrorOptions)
						//dlmirror.SetMode(beatmap.ModeStandard)
					} else if val == mode.String() {
						mirrorOptions := dlmirror.MirrorOptions()
						mirrorOptions.Mode = mode

						dlmirror.SetOptions(mirrorOptions)
						//dlmirror.SetMode(mode)
						break
					} else if i == len(modes) - 1 {
						log.Warn("Unknown gamemode ", val, ", filtering by ", mirrorOpts.Mode.String(), " instead.")
					}
				}
				case "status": // mapset status
					statuses := []beatmap.Status{
						beatmap.StatusGraveyard, beatmap.StatusWIP, beatmap.StatusPending,
						beatmap.StatusRanked, beatmap.StatusApproved, beatmap.StatusQualified,
						beatmap.StatusLoved, beatmap.StatusAny,
					}

					for i, status := range statuses {
						if val == status.String() {
							mirrorOptions := dlmirror.MirrorOptions()
							mirrorOptions.Status = status

							dlmirror.SetOptions(mirrorOptions)
							break
						} else if i == len(statuses) - 1 {
							log.Warn("Unknown status ", val, ", filtering by ", mirrorOpts.Status.String(), " maps instead.")
						}
					}
			}
			arg = strings.Replace(arg, group, "", -1)
		}
	}

	return arg
}

func searchBeatmaps(query string) (beatmap.Mapset, error) {
	escapedSearch := url.PathEscape(query)

	log.Info("Searching for query ", query)
	sets, _ := dlmirror.Search(escapedSearch, dlmirror.MirrorOptions())
	if len(sets) == 0 {
		return beatmap.Mapset{}, errNoResults
	}

	setTiles := make([]string, len(sets))
	for i, s := range sets {
		// artist - title by mapper
		setTiles[i] = fmt.Sprintf("%s - %s by %s", s.Artist, s.Title, s.Mapper)
	}

	prompt := promptui.Select{
		Label: "Select a mapset",
		Items: setTiles,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Aborting search...")
		return beatmap.Mapset{}, errSearchAborted
	}

	return sets[idx], nil
}
