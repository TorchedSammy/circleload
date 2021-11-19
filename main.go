package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/manifoldco/promptui"
)

var (
	outDir string
	mirrorName string
	mirrorFallback bool
	noVideo bool
)

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
	flag.BoolVarP(&mirrorFallback, "fallback", "f", false, "Fallback to other mirrors if main mirror fails")
	flag.BoolVarP(&noVideo, "noVideo", "V", false, "Don't download map with video")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: Circeload [flags] <mapset> [mapset] ...")
		//not implemented
		//fmt.Println("mapset can be the url or just the id")
		flag.PrintDefaults()
		os.Exit(1)
	}

	mirror := getMirror(mirrorName)
	mirrorOpts := mirrorOptions{
		noVideo: noVideo,
	}

	if mirror == nil {
		fmt.Println("Invalid mirror", mirrorName)
		fmt.Println("Valid mirrors are: kitsu, chimu")
		os.Exit(1)
	}

	mirrors := []string{"kitsu", "chimu"}
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
			info(fmt.Sprintf("Searching for query \"%s\"", v))
			sets, _ := mirror.Search(v)
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
				fmt.Println("Falling back to other mirror")
				mirror = getMirror(mirrors[0])
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
		err = downloadMapset(set.SetID, name, mirror, mirrorOpts)
		if err != nil {
			// i dont really like the repeating code here but i dont know how to do it better
			fmt.Println("Error downloading mapset:", err)
			if mirrorFallback {
				// if no other mirrors, exit
				if len(mirrors) == 0 {
					fmt.Println("All mirrors tried, exiting..")
					os.Exit(1)
				}
				fmt.Println("Falling back to chimu")
				mirror = getMirror(mirrors[0])
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

func downloadMapset(mapsetID int, name string, mirror mapsetMirror, opts mirrorOptions) error {
	mapset, err := mirror.GetMapsetData(mapsetID, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %s\n", name)

	// write body to file
	file, err := os.Create(fmt.Sprintf("%s/%s.osz", outDir, name))
	if err != nil {
		return err
	}
	defer file.Close()

	// mapset is a ReadCloser
	io.Copy(file, mapset)
	mapset.Close()
	return nil
}

func getMirror(name string) mapsetMirror {
	switch name {
	case "kitsu":
		return kitsuMirror{}
	case "chimu":
		return chimuMirror{}
	// perhaps in the future, copilot
	/*
	case "osu":
		return osuMirror{}
	*/
	}

	return nil
}

