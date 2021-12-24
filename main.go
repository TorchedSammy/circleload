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

	"github.com/TorchedSammy/circleload/beatmap"
	"github.com/TorchedSammy/circleload/log"
	"github.com/cheggaaa/pb"
	"github.com/manifoldco/promptui"
	flag "github.com/spf13/pflag"
)

var (
	outDir string
	mirrorName string
	mirrorFallback bool
	noVideo bool
	versionFlag bool
	mirrorListFlag bool
	maxResults int
)

const version = "0.2.0"
var kvRegex = regexp.MustCompile(`([\w]+)=([\w]+)`)

func main() {
	homedir, _ := os.UserHomeDir()
	flag.StringVarP(&outDir, "downloadDir", "d", filepath.Join(homedir, "Downloads"), "Directory Circeload will download maps into")
	flag.StringVarP(&mirrorName, "mirror", "m", "chimu", "Mirror to download from (kitsu or chimu)")
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

	mirrors := []string{"kitsu", "chimu"}
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

	mirrorOpts := beatmap.Options{
		NoVideo: noVideo,
		MaxResults: maxResults,
		Mode: beatmap.ModeAny,
	}
	dlmirror := getMirror(mirrorName, mirrorOpts)

	if dlmirror == nil {
		fmt.Println("Invalid mirror", mirrorName)
		fmt.Println("Valid mirrors are:", strings.Join(mirrors, ", "))
		os.Exit(1)
	}

	// remove mirror we are using from list
	for i, m := range mirrors {
		if m == mirrorName {
			mirrors = append(mirrors[:i], mirrors[i + 1:]...)
			break
		}
	}

	for _, v := range flag.Args() {
	start:
		var set beatmap.Mapset
		idInt, err := strconv.Atoi(v)
		if err != nil {
			mirrorOpts.Mode = beatmap.ModeAny
			// match all key value pairs in search query
			matches := kvRegex.FindAllStringSubmatch(v, -1)
			if len(matches) > 0 {
				for _, m := range matches {
					// expected format: [group, key, value]
					if len(m) == 3 {
						// check key
						switch m[1] {
							case "mode": // specific gamemode
							modes := []beatmap.Mode{beatmap.ModeStandard, beatmap.ModeTaiko, beatmap.ModeCatch, beatmap.ModeMania}

							for i, mode := range modes {
								m := m[2]
								if m == "osu" {
									dlmirror.SetMode(beatmap.ModeStandard)
								} else if m == mode.String() {
									dlmirror.SetMode(mode)
									break
								} else if i == len(modes) - 1 {
									log.Warn("Unknown gamemode ", m, ", getting maps for any gamemode.")
								}		
							}

							v = strings.Replace(v, m[0], "", -1)
						}
					}
				}
			}

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
					if strings.HasPrefix(u.Path, "/beatmapsets/") {
						idInt, err = strconv.Atoi(u.Path[len("/beatmapsets/"):])
						if err != nil {
							log.Error("Ignoring invalid mapset url: " + v)
						}

						set, err = dlmirror.GetMapset(idInt)
						goto download
					} else if strings.HasPrefix(u.Path, "/beatmaps/") {
						idInt, err = strconv.Atoi(u.Path[len("/beatmaps/"):])
						if err != nil {
							log.Error("Ignoring invalid mapset url: " + v)
						}

						mapset, err := dlmirror.GetMapsetFromMap(idInt)
						if err != nil {
							log.Error(fmt.Sprintln("Could not get mapset from map:", err))
						}
						set = mapset
						goto download
					} else {
						log.Error("Invalid mapset url")
						os.Exit(1)
					}
				} else {
					log.Error("Ignoring non-osu url: " + v)
					continue
				}
			}

			escapedSearch := url.PathEscape(v)
			log.Info("Searching for query ", v)
			sets, _ := dlmirror.Search(escapedSearch)
			if len(sets) == 0 {
				log.Error("No results found.")
				continue
			}
			setTiles := make([]string, len(sets))
			for i, s := range sets {
				// title - artist by mapper
				setTiles[i] = fmt.Sprintf("%s - %s by %s", s.Title, s.Artist, s.Mapper)
			}

			prompt := promptui.Select{
				Label: "Select a mapset",
				Items: setTiles,
			}

			idx, _, err := prompt.Run()
			if err != nil {
				fmt.Println("Aborting search...")
				continue
			}

			set = sets[idx]
			goto download
		}

		set, err = dlmirror.GetMapset(idInt)
		if err != nil {
			fmt.Println("Error getting mapset:", err)
			if mirrorFallback {
				// if no other mirrors, exit
				if len(mirrors) == 0 {
					fmt.Println("All mirrors tried, exiting..")
					os.Exit(1)
				}
				dlmirror = getMirror(mirrors[0], mirrorOpts)
				log.Info("Falling back to " + mirrors[0])
				// remove mirror we are using from list
				for i, m := range mirrors {
					if m == mirrors[0] {
						mirrors = append(mirrors[:i], mirrors[i + 1:]...)
						break
					}
				}
				goto start
			} else {
				continue
			}
		}
	download:
		name := fmt.Sprintf("%d %s - %s", set.SetID, set.Artist, set.Title)
		err = downloadMapset(set.SetID, name, dlmirror)
		if err != nil {
			// i dont really like the repeating code here but i dont know how to do it better
			log.Error(fmt.Sprint("Error downloading mapset:", err))
			if mirrorFallback {
				// if no other mirrors, exit
				if len(mirrors) == 0 {
					log.Error("All mirrors tried, exiting..")
					os.Exit(1)
				}
				log.Info("Falling back to " + mirrors[0])
				dlmirror = getMirror(mirrors[0], mirrorOpts)
				// remove mirror we are using from list
				for i, m := range mirrors {
					if m == mirrors[0] {
						mirrors = append(mirrors[:i], mirrors[i + 1:]...)
						break
					}
				}
				goto start
			} else {
				continue
			}
		}
	}
}

func downloadMapset(mapsetID int, name string, mirror beatmap.Mirror) error {
	mapsetResp, err := mirror.GetMapsetData(mapsetID)
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
	os.Rename(dest, dest[:len(dest)-4] + ".osz")

	return nil
}

func getMirror(name string, opts beatmap.Options) beatmap.Mirror {
	switch name {
	case "kitsu":
		return &beatmap.Kitsu{Options: opts}
	case "chimu":
		return &beatmap.Chimu{Options: opts}
	// perhaps in the future, copilot
	/*
	case "osu":
		return osuMirror{}
	*/
	}

	return nil
}

