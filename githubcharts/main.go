package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var err error
	switch os.Args[1] {
	case "get":
		err = get(os.Args[2], os.Args[3], fmt.Sprintf("%s-%s.csv", os.Args[2], os.Args[3]))
	case "oc":
		csv := fmt.Sprintf("%s-%s.csv", os.Args[2], os.Args[3])
		png := fmt.Sprintf("%s-%s.png", os.Args[2], os.Args[3])
		var start time.Time
		if len(os.Args) >= 4 {
			start, err = time.Parse("2006-01-02", os.Args[4])
			if err != nil {
				log.Fatal(err)
			}
		}
		err = plotOpenClose(csv, png, start)
	}
	if err != nil {
		log.Fatal(err)
	}
}
