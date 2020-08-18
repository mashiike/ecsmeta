package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mashiike/ecsmeta"
)

func main() {
	var (
		endpoint string
		silent   bool
	)
	flag.StringVar(&endpoint, "endpoint", "", "custom ECS metadata endpoint")
	flag.BoolVar(&silent, "silent", false, "no log flag")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	opts := []ecsmeta.Option{}
	if !silent {
		opts = append(opts, ecsmeta.WithLogger(log.New(os.Stderr, "", log.LstdFlags)))
	}
	if endpoint != "" {
		opts = append(opts, ecsmeta.WithEndpoint(endpoint))
	}
	obj := ecsmeta.New(opts...)
	val, err := obj.Query(strings.Join(flag.Args(), " "))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, val)
}
