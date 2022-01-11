package siteScan

import (
	// System
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"

	// Homebrew
	"github.com/je09/siteStalker/types"
)

func getHost(site types.Site, ch chan types.Site, chErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	r, _ := regexp.Compile("^[a-z][a-z0-9+\\-.]*://([a-z0-9\\-._~%!$&'()*+,;=]+@)?([a-z0-9\\-._~%]+|↵\n\\[[a-z0-9\\-._~%!$&'()*+,;=:]+\\])")
	p := r.FindStringSubmatch(site.Address)

	if len(p) >= 2 {
		site.Domain = p[2]
		ch <- site
	}

	if len(p) < 2 {
		chErr <- errors.New(fmt.Sprintf("address: %s can't get domain name", site.Address))
	}
}

func fixScheme(link string) string {
	r, _ := regexp.Compile("^(?:(https?):)?(\\/\\/.*\\.[^\\/]+)$")
	p := r.MatchString(link)
	if p != true {
		return "https://" + link
	}

	return link
}

func readBody(d io.ReadCloser) string {
	body, err := ioutil.ReadAll(d)
	if err != nil {
		log.Println(err)
	}

	return string(body)
}

func readHeader(d http.Header) string {
	var r string
	for k, v := range d {
		for _, h := range v {
			r += fmt.Sprintf("%v: %v\n", k, h)
		}
	}

	return r
}

func getSite(link string, userAgent string, ch chan types.Site, chErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	link = fixScheme(link)
	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		chErr <- err
		return
	}
	if req == nil {
		chErr <- errors.New(fmt.Sprintf("request to %s can't be processed", link))
		return
	}
	req.Header.Set("User-Agent", userAgent)
	r, err := client.Do(req)
	if r != nil {
		ch <- types.Site{Status: r.StatusCode, Header: readHeader(r.Header), Body: readBody(r.Body), Address: link}
	} else {
		chErr <- errors.New(fmt.Sprintf("no response from %s\n", link))
	}
}

func lookupSite(site types.Site, ch chan types.Site, chErr chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	ips, err := net.LookupIP(site.Domain)
	if err != nil {
		chErr <- err
	}
	site.Ips = ips
	ch <- site
}
