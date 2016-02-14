package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type param struct {
	domainName string
	count      int
	symbols    []byte
	length     int
}

var config param

func init() {
	const (
		domainName = "Domain name"
		count = "Count of email address"
		length = "Length of email address"
	)
	flag.StringVar(&config.domainName, "domain", "example.com", domainName)
	flag.IntVar(&config.count, "count", 100, count)
	flag.IntVar(&config.length, "length", 6, length)
	flag.StringVar(&config.domainName, "d", "example.com", domainName)
	flag.IntVar(&config.count, "c", 100, count)
	flag.IntVar(&config.length, "l", 6, length)
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomString(symbols []byte, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func main() {
	flag.Parse()
	config.symbols = []byte("abcdefghijklmnopqrstuvwxyz0123456789")

	emails := make(map[string]bool)
	for {
		if(len(emails) >= config.count) {
			break;
		}

		name := GenerateRandomString(config.symbols, config.length)
		if emails[name] {
			continue
		}

		emails[name] = true
		fmt.Println(name + `@` + config.domainName + ` ` + name)
	}
}
