package utils

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func createFolder(path string) (string, error) {
	path = fmt.Sprintf("%s/%s", path, time.Now().Format("01022006"))
	err := os.MkdirAll(path, 0700)

	if os.IsExist(err) == false && err != nil {
		return path, err
	}

	return path, nil
}

func write(path string, domain *string, d *string, t string, chErr *chan error) {
	path, err := createFolder(path)
	if err != nil {
		*chErr <- err
		log.Fatalln(err)
	}

	file := fmt.Sprintf("%s/%s_%s.log", path, *domain, t)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		*chErr <- err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)
	r := strings.NewReader(*d)
	if _, err := io.Copy(f, r); err != nil {
		*chErr <- err
		log.Println(err)
	}
}

func WriteBody(out *string, domain *string, body *string, chErr *chan error) {
	write(*out, domain, body, "body", chErr)
}

func WriteHeader(out *string, domain *string, header *string, chErr *chan error) {
	write(*out, domain, header, "headers", chErr)
}

func WriteCsv(path string, domain *string, ips *[]net.IP, chErr *chan error) {
	_, err := createFolder(path)
	if err != nil {
		*chErr <- err
		log.Fatalln(err)
	}

	file := fmt.Sprintf("%s/lookup.csv", path)

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		*chErr <- err
	}
	defer f.Close()
	for _, ip := range *ips {
		r := strings.NewReader(fmt.Sprintf("%s,%s\n", ip.String(), *domain))
		if _, err := io.Copy(f, r); err != nil {
			log.Println(err)
		}
	}
}
