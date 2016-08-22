package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/aclements/go-gg/gg"
	"github.com/aclements/go-gg/table"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

type event struct {
	diff int // +1 for open, -1 for close
	time time.Time
}

type eventsByTime []event

func (e eventsByTime) Len() int           { return len(e) }
func (e eventsByTime) Less(i, j int) bool { return e[i].time.Before(e[j].time) }
func (e eventsByTime) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func loadEvents(fn string) ([]event, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	var events []event
	for fields, err := r.Read(); err != io.EOF; fields, err = r.Read() {
		t, err := time.Parse(time.RFC3339, fields[2])
		if err != nil {
			return nil, fmt.Errorf("%v: %s", fields, err)
		}
		events = append(events, event{1, t})

		if fields[3] != "" {
			t, err := time.Parse(time.RFC3339, fields[3])
			if err != nil {
				return nil, fmt.Errorf("%v: %s", fields, err)
			}
			events = append(events, event{-1, t})
		}
	}
	sort.Sort(eventsByTime(events))
	return events, nil
}

func plotOpenClose(fn string, outf string, start time.Time) error {
	events, err := loadEvents(fn)
	if err != nil {
		return err
	}
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
	p.Add(NewVLine(float64(time.Date(2016, 2, 17, 0, 0, 0, 0, time.UTC).Unix()))) // 1.6
	// p.Add(NewVerticalLine(float64(time.Date(2016, 4, 12, 0, 0, 0, 0, time.UTC).Unix()))) // 1.6.1
	// p.Add(NewVerticalLine(float64(time.Date(2016, 4, 20, 0, 0, 0, 0, time.UTC).Unix()))) // 1.6.2

	p.Add(NewVLine(float64(time.Date(2015, 8, 19, 0, 0, 0, 0, time.UTC).Unix()))) // 1.5
	// p.Add(NewVerticalLine(float64(time.Date(2015, 9, 8, 0, 0, 0, 0, time.UTC).Unix())))  // 1.5.1
	// p.Add(NewVerticalLine(float64(time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC).Unix()))) // 1.5.2
	// p.Add(NewVerticalLine(float64(time.Date(2015, 1, 13, 0, 0, 0, 0, time.UTC).Unix()))) // 1.5.3
	// p.Add(NewVerticalLine(float64(time.Date(2015, 4, 12, 0, 0, 0, 0, time.UTC).Unix()))) // 1.5.4

	p.Add(NewVLine(float64(time.Date(2014, 12, 10, 0, 0, 0, 0, time.UTC).Unix()))) // 1.4

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

func plotOpenCloseGG(fn string, outf string, start time.Time) error {
	events, err := loadEvents(fn)
	if err != nil {
		return err
	}
	var time []time.Time
	var issues []int
	var types []string
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

		// time = append(time, e.time)
		// issues = append(issues, opened)
		// types = append(types, "opened")
		//
		// time = append(time, e.time)
		// issues = append(issues, closed)
		// types = append(types, "closed")

		time = append(time, e.time)
		issues = append(issues, open)
		types = append(types, "open")
	}

	data := table.NewBuilder(nil).
		Add("time", time).
		Add("issues", issues).
		Add("type", types).
		Done()

	//fmt.Println(data)

	p := gg.NewPlot(data)
	p.SetScale("x", gg.NewTimeScaler())
	p.Add(gg.LayerLines{Color: "type"})

	f, err := os.Create(outf)
	if err != nil {
		return err
	}
	defer f.Close()
	return p.WriteSVG(f, 1000, 800)
}

func plotByWeek(fn string, outf string, start time.Time) error {
	events, err := loadEvents(fn)
	if err != nil {
		return err
	}
	var openned plotter.XYs
	var closed plotter.XYs
	for _, e := range events {
		t := e.time
		if !start.IsZero() && t.Before(start) {
			continue
		}
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) // round to day
		t = t.AddDate(0, 0, -int(t.Weekday()))                                // round to week
		d := float64(t.Unix())
		if e.diff > 0 {
			if len(openned) == 0 || d != openned[len(openned)-1].X {
				openned = append(openned, struct{ X, Y float64 }{d, 1.0})
			} else {
				openned[len(openned)-1].Y++
			}
		} else {
			if len(closed) == 0 || d != closed[len(closed)-1].X {
				closed = append(closed, struct{ X, Y float64 }{d, 1.0})
			} else {
				closed[len(closed)-1].Y++
			}
		}
	}

	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Issues opened by week for " + fn
	p.X.Tick.Marker = plot.UnixTimeTicks{Format: "2006-01-02"}
	p.Y.Label.Text = "Issues"
	p.Add(plotter.NewGrid())

	err = plotutil.AddLines(p,
		"Issued closed", closed,
		"Issued opened", openned)

	if err != nil {
		return err
	}

	return p.Save(30*vg.Centimeter, 20*vg.Centimeter, outf)
}
