package siteScan

import (
	// System
	"bufio"
	"log"
	"os"
	"sync"
	"time"

	// Homebrew
	"github.com/je09/siteStalker/types"
	"github.com/je09/siteStalker/utils"
)

func Start(path *string, out *string, userAgent *string, msec *int, showErr *bool) error {
	chGet := make(chan types.Site)
	chDomain := make(chan types.Site)
	chLook := make(chan types.Site)
	chError := make(chan error)

	wg := &sync.WaitGroup{}
	f, err := os.Open(*path)
	if err != nil {
		log.Fatalln(err)
	}
	scn := bufio.NewScanner(f)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for scn.Scan() {
			link := scn.Text()
			wg.Add(1)
			go getSite(link, *userAgent, chGet, chError, wg)
			time.Sleep(time.Duration(*msec)*time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case err = <- chError:
				if *showErr {
					log.Print(err)
				}
				break
			case msg1 := <-chGet:
				wg.Add(1)
				go getHost(msg1, chDomain, chError, wg)
			case msg2 := <- chDomain:
				wg.Add(1)
				go lookupSite(msg2, chLook, chError, wg)
			case msg3 := <- chLook:
				log.Printf("address: %s status: %d, ips: %s",
					msg3.Domain, msg3.Status, msg3.Ips)
				wg.Add(1)
				go writeSite(msg3, *out, chError, wg)
			}
		}
	}()

	wg.Wait()
	err = f.Close()
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func writeSite(site types.Site, out string, chErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	if site.Body != "" {
		utils.WriteBody(&out, &site.Domain, &site.Body, &chErr)
	}
	if site.Header != "" {
		utils.WriteHeader(&out, &site.Domain, &site.Header, &chErr)
	}
	if site.Ips != nil {
		utils.WriteCsv(out, &site.Domain, &site.Ips, &chErr)
	}
}
