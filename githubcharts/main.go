package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var err error
	switch os.Args[1] {
	case "get":
		err = get(os.Args[2], os.Args[3], fmt.Sprintf("%s-%s.csv", os.Args[2], os.Args[3]))
	case "oc":
		csv := fmt.Sprintf("%s-%s.csv", os.Args[2], os.Args[3])
		png := fmt.Sprintf("%s-%s.png", os.Args[2], os.Args[3])
		err = plotOpenClose(csv, png)
	}
	if err != nil {
		log.Fatal(err)
	}
}
