package main

import (
	"flag"
	"fmt"
	"github.com/jbarham/primegen.go"
	"os"
	"strconv"
)

var low, high uint64 = 2, 1000000000

func usage() {
	fmt.Fprintf(os.Stderr, "usage: primes [[low=%d] high=%d]\n", low, high)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	var err error
	if len(args) == 1 {
		high, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			usage()
		}
	} else if len(args) == 2 {
		low, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			usage()
		}
		high, err = strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			usage()
		}
	} else if len(args) > 2 {
		usage()
	}

	primegen.CountClasses(low, high)
}
