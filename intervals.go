package interval

import (
	"fmt"
	"math"
	"sort"
)

const (
	defaultMinLow  = 0
	defaultMaxHigh = math.MaxInt64
)

type Intervals interface {
	// Add appends a new interval
	Add(itvl *Interval)

	// Sort sorts the intervals list by the Low property (ascending)
	Sort()

	// Gaps first sorts (if necessary) and then returns the interval gaps
	Gaps() []*Interval

	// Overlapped first sorts (if necessary) and then returns the overlapped intervals
	Overlapped() []*Interval

	// FindIntervalsForValue returns all the intervals which contains the passed value
	FindIntervalsForValue(value int) []*Interval

	// Print first sorts (if necessary) and then displays graphically the interval sequence
	Print() string
}

type intervals struct {
	Intervals []*Interval
	MinLow    int
	MaxHigh   int
	Sorted    bool
}

func NewIntervalsDefault() Intervals {
	return NewIntervals(defaultMinLow, defaultMaxHigh)
}

func NewIntervals(minLow int, maxHigh int) Intervals {
	return &intervals{
		MinLow:    minLow,
		MaxHigh:   maxHigh,
		Intervals: []*Interval{},
		Sorted:    false,
	}
}

func (intvls *intervals) Add(itvl *Interval) {
	intvls.Intervals = append(intvls.Intervals, itvl)
	intvls.Sorted = false
}

func (intvls *intervals) FindIntervalsForValue(value int) []*Interval {
	var matches []*Interval
	for _, intvl := range intvls.Intervals {
		if inBetweenInclusive(value, intvl.Low, intvl.High) {
			matches = append(matches, intvl)
		}
	}
	return matches
}

func (intvls *intervals) Sort() {
	if !intvls.Sorted {
		sort.Sort(ByLow(intvls.Intervals))
	}
	intvls.Sorted = true
}

func (intvls *intervals) Gaps() []*Interval {
	intvls.Sort()
	gaps := []*Interval{}
	lastHigh := intvls.MinLow
	for _, intvl := range intvls.Intervals {
		if intvl.Low > lastHigh {
			gaps = append(gaps, &Interval{Low: lastHigh, High: intvl.Low - 1})
		}
		lastHigh = intvl.High + 1
	}
	if lastHigh < intvls.MaxHigh {
		gaps = append(gaps, &Interval{Low: lastHigh, High: intvls.MaxHigh})
	}
	return gaps
}

func (intvls *intervals) Overlapped() []*Interval {
	intvls.Sort()
	list := []*Interval{}
	lastMinLow := math.MaxInt64
	lastMaxHigh := math.MinInt64
	for i, intvl := range intvls.Intervals {
		if i > 0 {
			lowInBetween := inBetweenInclusive(lastMinLow, intvl.Low, intvl.High) || inBetweenInclusive(intvl.Low, lastMinLow, lastMaxHigh)
			highInBetween := inBetweenInclusive(lastMaxHigh, intvl.Low, intvl.High) || inBetweenInclusive(intvl.High, lastMinLow, lastMaxHigh)
			if lowInBetween || highInBetween {
				greaterLow := max(intvl.Low, lastMinLow)
				lowerHigh := min(intvl.High, lastMaxHigh)
				list = append(list, &Interval{Low: greaterLow, High: lowerHigh})
			}
		}
		if intvl.Low < lastMinLow {
			lastMinLow = intvl.Low
		}
		if intvl.High > lastMaxHigh {
			lastMaxHigh = intvl.High
		}
	}
	return list
}

func (intvls *intervals) isAnOverlap(value int, overlapped []*Interval) bool {
	for _, ovrlp := range overlapped {
		if inBetweenInclusive(value, ovrlp.Low, ovrlp.High) {
			return true
		}
	}
	return false
}

func (intvls *intervals) Print() string {
	intvls.Sort()

	// Available Symbols:  ( ◯ ◌ ◍ ◎ ● ◉ ) , ( □ ■ ), ( ░ ▒ ▓ █ )
	emptySymbol := "◌"
	fullSymbol := "◎"
	overlapSymbol := "●"
	separator := "║"

	introText := fmt.Sprintf("\n==================================\n SUMMARY (minLow=%d, maxHigh=%d)\n==================================", intvls.MinLow, intvls.MaxHigh)
	legend := fmt.Sprintf("\n • Legend: %v (empty), %v (full), %v (overlap)", emptySymbol, fullSymbol, overlapSymbol)
	intervalText := "\n • Intervals: "
	gapsText := "\n • Gaps: "
	overlapText := "\n • Overlapped: "
	graph := ""
	index := intvls.MinLow
	blockSize := 10
	numSeparators := 0

	overlapped := intvls.Overlapped()
	for i, ovrlp := range overlapped {
		if i != 0 {
			overlapText += ", "
		}
		overlapText += fmt.Sprintf("[%d,%d]", ovrlp.Low, ovrlp.High)
	}

	for i, intvl := range intvls.Intervals {
		if i != 0 {
			intervalText += ", "
		}
		intervalText += fmt.Sprintf("[%d,%d]", intvl.Low, intvl.High)
		for i := index; i < intvl.Low; i++ {
			index++
			graph += emptySymbol
			if index%10 == 0 {
				graph += separator
				numSeparators++
			}
		}

		for i := index; i <= intvl.High; i++ {
			if intvls.isAnOverlap(index, overlapped) {
				graph += overlapSymbol
			} else {
				graph += fullSymbol
			}
			index++
			if index%blockSize == 0 {
				graph += separator
				numSeparators++
			}
		}
	}
	gaps := intvls.Gaps()
	for i, gap := range gaps {
		if i != 0 {
			gapsText += ", "
		}
		gapsText += fmt.Sprintf("[%d,%d]", gap.Low, gap.High)
	}

	for i := index; i < intvls.MaxHigh; i++ {
		graph += emptySymbol
	}
	axisLegend := fmt.Sprintf(" %v", intvls.MinLow)
	for i := intvls.MinLow; i < intvls.MaxHigh+numSeparators-2; i++ {
		axisLegend += " "
	}
	axisLegend += fmt.Sprintf("%v", intvls.MaxHigh)
	graphText := fmt.Sprintf("\n\n%s\n╠%s╣", axisLegend, graph)
	return "\n" + introText + legend + intervalText + gapsText + overlapText + graphText + "\n"
}
