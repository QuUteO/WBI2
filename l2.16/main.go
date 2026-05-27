package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	crawler2 "wget/crawler"
)

func main() {
	var (
		startURL  string
		maxDepth  int
		outputDir string
	)

	flag.StringVar(&startURL, "url", "", "start URL for mirroring")
	flag.IntVar(&maxDepth, "depth", 1, "recursion depth for HTML links")
	flag.StringVar(&outputDir, "out", "downloaded_site", "directory where mirrored files will be saved")
	flag.Parse()

	if startURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("не удалось определить выходной каталог: %v", err)
	}

	crawler, err := crawler2.NewCrawler(startURL, maxDepth, absOutputDir)
	if err != nil {
		log.Fatalf("ошибка инициализации краулера: %v", err)
	}

	if err := crawler.Start(startURL); err != nil {
		log.Fatalf("ошибка зеркалирования: %v", err)
	}

	fmt.Printf("Mirror saved to %s\n", absOutputDir)
}
