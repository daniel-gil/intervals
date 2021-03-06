package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/daniel-gil/intervals"
)

const (
	MinX = 0
	MaxX = 40
)

func main() {
	filename := "data.txt"
	xys, err := readData(filename)
	if err != nil {
		log.Fatalf("could not read %s: %v", filename, err)
	}
	intvls := initIntervals(xys)
	ip := intervals.NewPlot(intvls.IsLowInclusive(), intvls.IsHighInclusive())
	err = ip.PlotData("out.png", intvls, true, true, true, true)
	if err != nil {
		log.Fatalf("could not plot data: %v", err)
	}
}

type xy struct{ x, y int }

func readData(path string) ([]xy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var xys []xy

	// read line by line using a scanner (because we don't know if the file will be huge)
	s := bufio.NewScanner(f)
	for s.Scan() {
		var low, high int
		_, err := fmt.Sscanf(s.Text(), "%d,%d", &low, &high)
		if err != nil {
			log.Printf("discarding bad data point %v: %v", s.Text(), err)
			continue
		}
		if low > high {
			log.Printf("discarding bad data point (low, high)=(%v): low can not be greater than high", s.Text())
			continue
		}
		xys = append(xys, xy{low, high})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("could not scan: %v", err)
	}
	return xys, nil
}

func initIntervals(xys []xy) intervals.Intervals {
	// initialize Intervals
	minLow := MinX
	maxHigh := MaxX
	lowInclusive := true
	highInclusive := true
	selfAdjustMinLow := false
	selfAdjustMaxHigh := true
	intvls := intervals.New(minLow, maxHigh, lowInclusive, highInclusive, selfAdjustMinLow, selfAdjustMaxHigh)

	for _, xy := range xys {
		err := intvls.AddInterval(&intervals.Interval{Low: xy.x, High: xy.y})
		if err != nil {
			fmt.Printf("invalid interval discarded: %v\n", err)
		}
	}
	return intvls
}
