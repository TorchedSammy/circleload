package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/manifoldco/promptui"
	"github.com/cheggaaa/pb"
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

type osuMapset struct {
	SetID   int
	Artist string
	Title string
	Mapper string `json:"Creator"`
}

func main() {
	homedir, _ := os.UserHomeDir()
	flag.StringVarP(&outDir, "downloadDir", "d", filepath.Join(homedir, "Downloads"), "Directory Circeload will download maps into")
	flag.StringVarP(&mirrorName, "mirror", "m", "chimu", "Mirror to download from (kitsu or chimu)")
	flag.BoolVarP(&mirrorFallback, "fallback", "f", true, "Fallback to other mirrors if main mirror fails")
	flag.BoolVarP(&noVideo, "noVideo", "V", false, "Download mapset without video")
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
		//not implemented
		//fmt.Println("mapset can be the url or just the id")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if maxResults < 2 && maxResults > 100 {
		logerror("Search amount must be between 2 and 100.")
		os.Exit(1)
	}

	mirrorOpts := mirrorOptions{
		noVideo: noVideo,
		maxResults: maxResults,
	}
	mirror := getMirror(mirrorName, mirrorOpts)

	if mirror == nil {
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
		var set osuMapset
		idInt, err := strconv.Atoi(v)
		if err != nil {
			// will assume its a search query
			escapedSearch := url.PathEscape(v)
			info(fmt.Sprintf("Searching for query \"%s\"", v))
			sets, _ := mirror.Search(escapedSearch)
			if len(sets) == 0 {
				logerror("No results found.")
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

		set, err = mirror.GetMapset(idInt)
		if err != nil {
			fmt.Println("Error getting mapset:", err)
			if mirrorFallback {
				// if no other mirrors, exit
				if len(mirrors) == 0 {
					fmt.Println("All mirrors tried, exiting..")
					os.Exit(1)
				}
				mirror = getMirror(mirrors[0], mirrorOpts)
				info("Falling back to " + mirrors[0])
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
		name := strings.Replace(fmt.Sprintf("%d %s - %s", set.SetID, set.Artist, set.Title), "/", "", -1)
		err = downloadMapset(set.SetID, name, mirror)
		if err != nil {
			// i dont really like the repeating code here but i dont know how to do it better
			logerror(fmt.Sprint("Error downloading mapset:", err))
			if mirrorFallback {
				// if no other mirrors, exit
				if len(mirrors) == 0 {
					logerror("All mirrors tried, exiting..")
					os.Exit(1)
				}
				info("Falling back to " + mirrors[0])
				mirror = getMirror(mirrors[0], mirrorOpts)
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

func downloadMapset(mapsetID int, name string, mirror mapsetMirror) error {
	mapsetResp, err := mirror.GetMapsetData(mapsetID)
	if err != nil {
		return err
	}

	info("Downloading " + name)

	// write body to file
	file, err := os.Create(fmt.Sprintf("%s/%s.osz", outDir, name))
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
	io.Copy(file, barWriter)
	mapsetResp.Body.Close()
	bar.Finish()

	return nil
}

func getMirror(name string, opts mirrorOptions) mapsetMirror {
	switch name {
	case "kitsu":
		return kitsuMirror{opts: opts}
	case "chimu":
		return chimuMirror{opts: opts}
	// perhaps in the future, copilot
	/*
	case "osu":
		return osuMirror{}
	*/
	}

	return nil
}

