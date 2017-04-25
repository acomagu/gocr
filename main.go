package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"flag"

	"cloud.google.com/go/vision"
)

type result struct {
	Filename string `json:"filename"`
	Text     string `json:"text"`
}

func main() {
	vc, err := createVisionClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	resultChan := make(chan result)

	var wg sync.WaitGroup
	concurrency, filenames := parseFlags()
	filename := make(chan string)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(vc, filename, &wg, resultChan)
	}

	wait := make(chan bool)
	go output(resultChan, wait)
	for _, fn := range filenames {
		filename <- fn
	}

	close(filename)
	wg.Wait()

	close(resultChan)
	<-wait
}

func parseFlags() (int, []string) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var concurrency int
	f.IntVar(&concurrency, "concurrency", 3, "Number of concurrent worker to run")
	filenames := parseFlagPart(f, os.Args[1:])
	return concurrency, filenames
}

func parseFlagPart(f *flag.FlagSet, args []string) []string {
	f.Parse(args)
	if left := f.Args(); len(left) > 0 {
		return append(parseFlagPart(f, left[1:]), left[0])
	}
	return []string{}
}

func output(resultChan <-chan result, end chan<- bool) {
	results := []result{}
	for res := range resultChan {
		results = append(results, res)
	}
	jsonStr, err := json.Marshal(results)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(string(jsonStr))
	end <- true
}

func worker(vc *vision.Client, filename chan string, wg *sync.WaitGroup, resultChan chan<- result) {
	for fn := range filename {
		f, err := os.Open(fn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		str, err := detectText(vc, f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			filename <- fn
			continue
		}
		resultChan <- result{
			Filename: fn,
			Text:     str,
		}
		fmt.Fprintf(os.Stderr, "%s: Done.\n", fn)
	}
	wg.Done()
}

func createVisionClient() (*vision.Client, error) {
	ctx := context.Background()
	visionClient, err := vision.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return visionClient, nil
}

func detectText(vc *vision.Client, imgContent io.ReadCloser) (string, error) {
	ctx := context.Background()
	img, err := vision.NewImageFromReader(imgContent)
	if err != nil {
		return "", err
	}

	resultSlice, err := vc.DetectTexts(ctx, img, 10)
	if err != nil {
		return "", err
	}
	if len(resultSlice) == 0 {
		return "", fmt.Errorf("any characters are not detected")
	}
	return resultSlice[0].Description, nil
}
