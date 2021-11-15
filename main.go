package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	outDir string
)

type osuMapSet struct {
	SetID   int
	Artist string
	Title string
}

func main() {
	flag.StringVarP(&outDir, "outputDir", "o", "~/Downloads", "Directory Circeload will download maps into")
	flag.Parse()

	for _, v := range flag.Args() {
		idInt, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("Invalid mapset ID:", v)
			continue
		}

		resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/api/s/%d", idInt))
		var set osuMapSet
		json.NewDecoder(resp.Body).Decode(&set)

		name := strings.Replace(fmt.Sprintf("%d %s - %s", idInt, set.Artist, set.Title), "/", "", -1)
		fmt.Printf("Downloading %s\n", name)
		downloadMapset(idInt, name)
	}
}

func downloadMapset(mapsetID int, name string) {
	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/d/%d", mapsetID))
	if err != nil {
		panic(err)
	}

	// write body to file
	file, err := os.Create(fmt.Sprintf("%s/%s.osz", outDir, name))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, resp.Body)
}

