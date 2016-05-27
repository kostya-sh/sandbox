package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

type event struct {
	diff int
	time time.Time
}

type eventsByTime []event

func (e eventsByTime) Len() int           { return len(e) }
func (e eventsByTime) Less(i, j int) bool { return e[i].time.Before(e[j].time) }
func (e eventsByTime) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func plotOpenClose(fn string, outf string, start time.Time) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)
	var events []event
	for fields, err := r.Read(); err != io.EOF; fields, err = r.Read() {
		t, err := time.Parse(time.RFC3339, fields[2])
		if err != nil {
			return fmt.Errorf("%v: %s", fields, err)
		}
		events = append(events, event{1, t})

		if fields[3] != "" {
			t, err := time.Parse(time.RFC3339, fields[3])
			if err != nil {
				return fmt.Errorf("%v: %s", fields, err)
			}
			events = append(events, event{-1, t})
		}
	}
	sort.Sort(eventsByTime(events))

	var issuesOpened plotter.XYs
	var issuesClosed plotter.XYs
	var openIssues plotter.XYs
	opened := 0
	closed := 0
	open := 0
	for _, e := range events {
		if e.diff > 0 {
			opened++
		} else {
			closed++
		}
		open += e.diff

		if !start.IsZero() && e.time.Before(start) {
			continue
		}

		issuesOpened = append(issuesOpened, struct{ X, Y float64 }{
			X: float64(e.time.Unix()),
			Y: float64(opened),
		})
		issuesClosed = append(issuesClosed, struct{ X, Y float64 }{
			X: float64(e.time.Unix()),
			Y: float64(closed),
		})
		openIssues = append(openIssues, struct{ X, Y float64 }{
			X: float64(e.time.Unix()),
			Y: float64(open),
		})
	}

	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Issues for " + fn
	p.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	p.Y.Label.Text = "Issues"
	p.Add(plotter.NewGrid())

	err = plotutil.AddLines(p, "Open issues", openIssues)
	// err = plotutil.AddLines(p,
	// 	"Issues opened", issuesOpened,
	// 	"Issues closed", issuesClosed,
	// )

	if err != nil {
		return err
	}

	return p.Save(30*vg.Centimeter, 20*vg.Centimeter, outf)
}
