package main

import (
	"errors"
	"fmt"
	"math"
	"os"
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
	ErrDurationTooSmall = errors.New("provided duration is too small: less than 1s")
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
)

func calcDowntime(u float64, d time.Duration) time.Duration {
	avail := float64(d.Nanoseconds()) * (1 - u/100)
	rounded := math.Round(avail)
	return time.Duration(rounded).Truncate(time.Second)
}

func formatDuration(d time.Duration) (string, error) {
	if d < time.Second {
		return "", ErrDurationTooSmall
	}
	// put the total in seconds
	total := uint64(d / time.Second)
	var buf [24]byte
	i := len(buf) - 1
	buf[i] = 's'
	// format seconds
	i = fmtInt(buf[:i], total%60)
	total /= 60

	// handle minutes
	if total > 0 {
		i--
		buf[i] = ' '
		i--
		buf[i] = 'm'
		i = fmtInt(buf[:i], total%60)
		total /= 60
	}

	// handle hours
	if total > 0 {
		i--
		buf[i] = ' '
		i--
		buf[i] = 'h'
		i = fmtInt(buf[:i], total%24)
		total /= 24
	}

	// handle days and stop
	if total > 0 {
		i--
		buf[i] = ' '
		i--
		buf[i] = 'd'
		i = fmtInt(buf[:i], total)
	}
	return string(buf[i:]), nil
}

func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
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
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', tabwriter.AlignRight)
	_, _ = w.Write([]byte("Period\tDowntime\t\n"))
	for idx, period := range order {
		downtime := calcDowntime(uptime, durations[idx])
		formattedDuration, err := formatDuration(downtime)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Fprintf(w, "%s\t%s\t\n", period, formattedDuration)
	}
	w.Flush()
}
