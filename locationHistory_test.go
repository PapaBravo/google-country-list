package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_readLocationHistoryFile(t *testing.T) {
	content, _ := ioutil.ReadFile("./data/Semantic Location History/2022/2022_FEBRUARY.json")
	readLocationHistoryFile(content)
}

func Test_getLocationHistoryFiles(t *testing.T) {
	files := getLocationHistoryFiles("./data/Semantic Location History/")

	fmt.Println(files[0:5])
}

func Test_getDateFromName(t *testing.T) {
	name := "MARCH MARCH APRIL"

	actual := getDateFromName(name)

	expected := "03 03 04"
	if actual != expected {
		t.Fatalf("Expected %s to be %s", actual, expected)
	}
}
