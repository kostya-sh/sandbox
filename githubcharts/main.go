package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var err error
	switch os.Args[1] {
	case "get":
		err = get(os.Args[2], os.Args[3], fmt.Sprintf("%s-%s.csv", os.Args[2], os.Args[3]))
	case "open", "open-gg", "byweek":
		csv := fmt.Sprintf("%s-%s.csv", os.Args[2], os.Args[3])
		png := fmt.Sprintf("%s-%s_%s.png", os.Args[2], os.Args[3], os.Args[1])
		var start time.Time
		if len(os.Args) >= 4 {
			start, err = time.Parse("2006-01-02", os.Args[4])
			if err != nil {
				log.Fatal(err)
			}
		}
		if os.Args[1] == "open" {
			err = plotOpenClose(csv, png, start)
		} else if os.Args[1] == "open-gg" {
			err = plotOpenCloseGG(csv, strings.Replace(png, ".png", ".svg", -1), start)
		} else {
			err = plotByWeek(csv, png, start)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}
