package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type PlaceVisit struct {
	Location struct {
		Lat  int    `json:"latitudeE7"`
		Lon  int    `json:"longitudeE7"`
		Name string `json:"name"`
	} `json:"location"`
	Duration struct {
		Start time.Time `json:"startTimestamp"`
		End   time.Time `json:"endTimestamp"`
	} `json:"duration"`
}

type TimeLine struct {
	TimelineObjects []struct {
		PlaceVisit PlaceVisit `json:"placeVisit"`
	} `json:"timelineObjects"`
}

func readLocationHistoryFile(content []byte) []PlaceVisit {
	tl := TimeLine{}
	err := json.Unmarshal(content, &tl)

	if err != nil {
		log.Fatalln("could not unmarshal file", err)
	}
	var placeVisits []PlaceVisit
	for _, val := range tl.TimelineObjects {
		if val.PlaceVisit.Location.Lon != 0 {
			placeVisits = append(placeVisits, val.PlaceVisit)
		}
	}
	return placeVisits
}

func WriteLocationHistory(dirname string, c chan PlaceVisit) {
	files := getLocationHistoryFiles(dirname)

	for _, f := range files {
		content, err := ioutil.ReadFile(f.Path)
		if err != nil {
			log.Printf("Could not read %s: %v\n", f.Path, err)
		}
		log.Println("Reading " + f.Path)
		visits := readLocationHistoryFile(content)

		for _, v := range visits {
			c <- v
		}
	}
	close(c)
}

type TimeLineFileInfo struct {
	Path string
	Date string
}

func getLocationHistoryFiles(dirname string) []TimeLineFileInfo {
	var files []TimeLineFileInfo

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err == nil && strings.HasSuffix(info.Name(), ".json") {
			files = append(files, TimeLineFileInfo{path, getDateFromName(info.Name())})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(files, func(p, q int) bool { return files[p].Date < files[q].Date })
	return files
}

var nameMap = [][]string{
	{"JANUARY", "01"},
	{"FEBRUARY", "02"},
	{"MARCH", "03"},
	{"APRIL", "04"},
	{"MAY", "05"},
	{"JUNE", "06"},
	{"JULY", "07"},
	{"AUGUST", "08"},
	{"SEPTEMBER", "09"},
	{"OCTOBER", "10"},
	{"NOVEMBER", "11"},
	{"DECEMBER", "12"},
}

func getDateFromName(name string) string {
	result := name
	for _, mapping := range nameMap {
		result = strings.ReplaceAll(result, mapping[0], mapping[1])
	}
	return result
}
