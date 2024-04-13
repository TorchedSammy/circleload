package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/TorchedSammy/circleload/mirror"
	"github.com/TorchedSammy/circleload/log"
	"github.com/TorchedSammy/circleload/beatmap"

	"github.com/cheggaaa/pb"
	flag "github.com/spf13/pflag"
)

const version = "0.4.0"

var (
	outDir string
	mirrorName string
	mirrorFallback bool
	noVideo bool
	versionFlag bool
	mirrorListFlag bool
	maxResults int
)

var kvRegex = regexp.MustCompile(`([\w]+)=([\w]+)`)
var mirrors = []string{"chimu", "osudirect"}
var dlmirror mirror.Mirror
var mirrorOpts mirror.Options

func main() {
	homedir, _ := os.UserHomeDir()
	flag.StringVarP(&outDir, "downloadDir", "d", filepath.Join(homedir, "Downloads"), "Directory Circeload will download maps into")
	flag.StringVarP(&mirrorName, "mirror", "m", "chimu", fmt.Sprintf("Mirror to download from (Options: %s)", strings.Join(mirrors, ", ")))
	flag.BoolVarP(&mirrorFallback, "fallback", "f", true, "Fallback to other mirrors if main mirror fails")
	flag.BoolVarP(&noVideo, "no-video", "n", false, "Download mapset without video")
	flag.BoolVarP(&versionFlag, "version", "v", false, "Print version and exit")
	flag.BoolVarP(&mirrorListFlag, "mirrors", "M", false, "List available mirrors and exit")
	flag.IntVarP(&maxResults, "max-results", "r", 5, "Amount of mapsets to return from a search")

	flag.Parse()

	if versionFlag {
		fmt.Println("Circleload v" + version)
		return
	}

	if mirrorListFlag {
		fmt.Println("Available mirrors:", strings.Join(mirrors, ", "))
		return
	}

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: Circeload [flags] <mapset> [mapset] ...")
		fmt.Println("mapset can be the url or just the id")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if maxResults < 2 && maxResults > 100 {
		log.Error("Search amount must be between 2 and 100.")
		os.Exit(1)
	}

	mirrorOpts = mirror.Options{
		NoVideo: noVideo,
		MaxResults: maxResults,
		Mode: beatmap.ModeStandard,
		Status: beatmap.StatusRanked,
	}
	dlmirror = getMirror(mirrorName, mirrorOpts)

	if dlmirror == nil {
		fmt.Println("Invalid mirror", mirrorName)
		fmt.Println("Valid mirrors are:", strings.Join(mirrors, ", "))
		os.Exit(1)
	}

	cycleMirrors()

	for _, v := range flag.Args() {
	start:
		var set beatmap.Mapset
		idInt, err := strconv.Atoi(v)
		if err != nil {
			//dlmirror.SetMode(beatmap.ModeStandard)
			//dlmirror.SetStatus(beatmap.StatusRanked)
			matches := kvRegex.FindAllStringSubmatch(v, -1)
			if len(matches) > 0 {
				// apply filters and remove them from search query
				v = applyFilters(matches, v)
			}

			// match all key value pairs in search query
			// strip spaces
			v = strings.TrimSpace(v)

			// will assume its a url
			// try to parse it
			u, err := url.ParseRequestURI(v)
			if err == nil {
				// check if it is a mapset url
				if u.Host == "osu.ppy.sh" {
					// we accept the path "beatmapsets/<id>" and "beatmaps/<id>"
					// a beatmapset has just a mapset, where for beatmaps we need to get the mapset id
					id := strings.Split(strings.Split(u.Path, "/")[2], "#")[0]
					idInt, err = strconv.Atoi(id)
					if strings.HasPrefix(u.Path, "/beatmapsets/") || strings.HasPrefix(u.Path, "/s/") {
						// err from above
						if err != nil {
							log.Error("Ignoring invalid mapset url: " + v)
							continue
						}

						set, err = dlmirror.GetMapset(idInt)
						goto download
					} else if strings.HasPrefix(u.Path, "/beatmaps/") || strings.HasPrefix(u.Path, "/b/") {
						// err from above
						if err != nil {
							log.Error("Ignoring invalid mapset url: " + v)
							continue
						}

						mapset, err := dlmirror.GetMapsetFromMap(idInt)
						if err != nil {
							log.Error(fmt.Sprintln("Could not get mapset from map:", err))
							continue
						}
						set = mapset
						goto download
					} else {
						log.Error("Invalid mapset url ", v)
						continue
					}
				} else {
					log.Error("Ignoring non-osu url: ", v)
					continue
				}
			}

			r, err := searchBeatmaps(v)
			if err == errNoResults {
				log.Error("No results found.")
				if mirrorFallback {
					fallbackMirror()
					goto start
				} else {
					continue
				}
			} else if err == errSearchAborted {
				log.Error("Search aborted.")
				continue
			} else if err != nil {
				log.Error("Error searching for ", v, ": ", err)
				continue
			}

			set = r
		} else {
			set, err = dlmirror.GetMapset(idInt)
			if err != nil {
				log.Error("Error getting mapset: ", err)
				if mirrorFallback {
					fallbackMirror()
					goto start
				}
				continue
			}
		}

	download:
		name := fmt.Sprintf("%d %s - %s", set.SetID, set.Artist, set.Title)
		err = downloadMapset(set.SetID, name, dlmirror, mirrorOpts)
		if err != nil {
			// i dont really like the repeating code here but i dont know how to do it better
			log.Error("Error downloading mapset: ", err)
			if mirrorFallback {
				fallbackMirror()
				goto start
			}
			continue
		}
	}
}

func downloadMapset(mapsetID int, name string, mirror mirror.Mirror, options mirror.Options) error {
	mapsetResp, err := mirror.GetMapsetData(mapsetID, options)
	if err != nil {
		return err
	}

	log.Info("Downloading " + name)

	// fName is the name without illegal characters
	fName := name
	illegalChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, c := range illegalChars {
		fName = strings.Replace(fName, c, "_", -1)
	}

	// write body to file
	dest := filepath.Join(outDir, fName + ".cdl")
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	contentLength := mapsetResp.ContentLength
	bar := pb.New64(contentLength)
	bar.SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = true
	barWriter := bar.NewProxyReader(mapsetResp.Body)

	bar.Start()
	// mapset is a ReadCloser
	_, err = io.Copy(file, barWriter)
	if err != nil {
		return err
	}
	mapsetResp.Body.Close()
	bar.Finish()

	// remove cdl extension and add .osz
	os.Rename(dest, dest[:len(dest) - 4] + ".osz")

	return nil
}

func getMirror(name string, opts mirror.Options) mirror.Mirror {
	switch name {
	case "chimu":
		return &mirror.Chimu{Options: opts}
	case "osudirect":
		return &mirror.OsuDirect{Options: opts}
	// perhaps in the future, copilot
	/*
	case "osu":
		return osuMirror{}
	*/
	}

	return nil
}

func fallbackMirror() {
	if len(mirrors) == 0 {
		log.Error("All mirrors tried, exiting..")
		os.Exit(1)
	}
	dlmirror = getMirror(mirrors[0], mirrorOpts)
	log.Info("Falling back to " + mirrors[0])
	cycleMirrors()
}

func cycleMirrors() {
	for i, m := range mirrors {
		if m == mirrors[0] {
			mirrors = append(mirrors[:i], mirrors[i + 1:]...)
			break
		}
	}
}
