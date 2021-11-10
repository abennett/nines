package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"text/tabwriter"
	"time"
)

const (
	DAY     = time.Hour * 24
	WEEK    = DAY * 7
	MONTH   = WEEK * 4
	QUARTER = MONTH * 3
	YEAR    = QUARTER * 4
)

var (
	durations = []time.Duration{
		DAY,
		WEEK,
		MONTH,
		QUARTER,
		YEAR,
	}
	order = []string{"Daily", "Weekly", "Monthly", "Quarterly", "Yearly"}

    durRX = regexp.MustCompile(`(\D)(\d)`)
    zeroRX = regexp.MustCompile(` 0\D`)
)

func calcDowntime(u float64, d time.Duration) time.Duration {
	avail := float64(d.Nanoseconds()) * (1 - u/100)
	rounded := math.Round(avail)
	return time.Duration(rounded).Truncate(time.Second)
}

func formatDuration(d time.Duration) string {
    s := durRX.ReplaceAllString(d.String(), `$1 $2`)
    return zeroRX.ReplaceAllLiteralString(s, "")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("invalid number of arguments")
		os.Exit(1)
	}
	uptime, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("unable to convert %s to float\n", os.Args[1])
		os.Exit(1)
	}
    w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', tabwriter.AlignRight)
    w.Write([]byte("Period\tDowntime\n"))
	for idx, period := range order {
        downtime := calcDowntime(uptime, durations[idx])
		fmt.Fprintf(w, "%s:\t%s\n", period, formatDuration(downtime))
	}
    w.Flush()
}
