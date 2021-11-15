package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"
)

var (
	outDir string
)

func main() {
	flag.StringVarP(&outDir, "outputDir", "o", "~/Downloads", "Directory Circeload will download maps into")
	flag.Parse()

	for _, v := range flag.Args() {
		idInt, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("Invalid mapset ID:", v)
		}
		downloadMapset(idInt)
	}
}

func downloadMapset(mapsetID int) {
	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/d/%d", mapsetID))
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)
	// write body to file
	file, err := os.Create(fmt.Sprintf("%s/%d.osz", outDir, mapsetID))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, resp.Body)
}
