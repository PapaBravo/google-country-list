package main

import (
	"encoding/csv"
	"fmt"
	"github.com/twpayne/go-geom/encoding/geojson"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	content, err := ioutil.ReadFile("./data/countries.geojson")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var fc geojson.FeatureCollection
	err = fc.UnmarshalJSON(content)

	if err != nil {
		log.Fatal("Error when unmarshalling country file: ", err)
	}

	history := AnalyzeHistory(fc)
	writeToCSV(history)
}

func writeToCSV(history []CountryVisit) {
	f, err := os.Create("./data/out/history.csv")
	if err != nil {
		log.Fatal("Error when creating output file: ", err)
	}

	writer := csv.NewWriter(f)

	err = writer.Write([]string{"Country", "Start", "End"})
	if err != nil {
		log.Fatal("Error when writing header: ", err)
	}

	for _, v := range history {
		record := []string{
			v.Country,
			fmt.Sprint(v.Start.Format("2006-01-02")),
			fmt.Sprint(v.End.Format("2006-01-02")),
		}
		err = writer.Write(record)
		if err != nil {
			log.Println("Could not write record: ", err)
		}
	}
	writer.Flush()
}
