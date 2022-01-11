package main

import (
	// System
	"flag"
	"log"
	_ "net/http/pprof"

	// Homebrew
	"github.com/je09/siteStalker/siteScan"
)

func main() {
	filePath := flag.String("file", "sites.txt", "Input file")
	outPath := flag.String("out", "reports", "Output dir")
	milliseconds := flag.Int("ms", 10, "Milliseconds per requests")
	showError := flag.Bool("error", false, "Show errors")
	userAgent := flag.String(
		"agent",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"User Agent")
	flag.Parse()

	err := siteScan.Start(filePath, outPath, userAgent, milliseconds, showError)
	if err != nil {
		log.Fatalln(err)
	}
}
