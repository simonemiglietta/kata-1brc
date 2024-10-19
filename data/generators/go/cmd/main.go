package main

import (
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"lvciot/generate/internal/engine"
	"os"
	"time"
)

func main() {
	p := message.NewPrinter(language.English)

	inputStationsFilename := flag.String("input", "../../weather_stations.csv", "input stations filename")
	outputMeasurementsFilename := flag.String("output", "measurements.txt", "output measurements filename")
	totalStations := flag.Uint("stations", 10_000, "total used stations")
	totalMeasurements := flag.Uint("measurements", 0, "total measurements produced")
	flag.Parse()

	if *totalMeasurements == 0 {
		fmt.Println("Number of measurement is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Starting with following configuration:")
	fmt.Printf("- input       : %s\n", *inputStationsFilename)
	fmt.Printf("- output      : %s\n", *outputMeasurementsFilename)
	_, _ = p.Printf("- stations    : %d\n", *totalStations)
	_, _ = p.Printf("- measurements: %d\n", *totalMeasurements)

	stationsFile, err := os.Open(*inputStationsFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := stationsFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	measurementsFile, err := os.Create(*outputMeasurementsFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := measurementsFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("Stations retrieving started")
	stations, err := engine.RetrieveStations(stationsFile, int(*totalStations))
	if err != nil {
		panic(err)
	}
	fmt.Println("Stations retrieving ended")

	fmt.Println("Measurements generation started")
	bar := progressbar.Default(int64(*totalMeasurements))
	ticker := time.NewTicker(time.Second)
	var c int
	go func() {
		for {
			select {
			case <-ticker.C:
				_ = bar.Set(c)
			}
		}
	}()

	err = engine.GenerateMeasurements(measurementsFile, stations, int(*totalMeasurements), &c)

	fmt.Println()
	fmt.Println("Measurements generation ended")
}
