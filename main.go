package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	dirPtr := flag.String("dir", ":/required",
		"parent directory containing numbered subdirectories containing keys")
	nPtr := flag.Uint64("n", 5, "number subdirectories (containing keys) to create")
	hostPtr := flag.String("hostname", "localhost", "folder to generate peers.json")
	portStartPtr := flag.Uint64("port-start", 12000, "port to start counting at")

	ensureCliArgs([]string{"dir", "n"})

	abspath, err := filepath.Abs(*dirPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Abs(%q) error: %v", *dirPtr, err)
		os.Exit(1)
	}

	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		os.MkdirAll(abspath, 0700)
	}

	wg := new(sync.WaitGroup)
	nStr := strconv.FormatUint(*nPtr, 10)
	padding := fmt.Sprintf("%%0%dd", len(nStr))
	for i := uint64(0); i < *nPtr; i++ {
		wg.Add(1)
		go GeneratePemKey(path.Join(abspath, fmt.Sprintf(padding, i)), wg)
	}
	wg.Wait()

	fmt.Println("[")
	err = filepath.Walk(*dirPtr, visitF(*nPtr, *portStartPtr, *hostPtr))
	fmt.Println("]")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func ensureCliArgs(required []string) {
	// From: https://stackoverflow.com/a/31795922
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			// or possibly use `log.Fatalf` instead of:
			fmt.Fprintf(os.Stderr, "missing required -%s argument/flag\n", req)
			os.Exit(2) // the same exit code flag.Parse uses
		}
	}
}
