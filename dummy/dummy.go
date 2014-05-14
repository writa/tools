package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	url string = "http://aaa.com/"

	hashkeys = []string{}

	location string = "http://$.test.com"

	wg sync.WaitGroup

	path []string

	num *int = flag.Int("n", 0, "number of request")
	con *int = flag.Int("c", 0, "concurrency")

	client *http.Client
)

func init() {
	flag.Parse()
	if *num == 0 || *con == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, v := range hashkeys {
		path = append(path, url+"cv?hashkey="+v)
	}

	for i := 0; i < 100; i++ {
		path = append(path, url+"?url="+strings.Replace(location, "$", strconv.Itoa(i), -1))
	}

	rand.Seed(time.Now().UTC().UnixNano())

	client = &http.Client{CheckRedirect: noRedirect}
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("no redirect")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("start\n.\n.\n.")

	var count int = *num / *con

	for i := 0; i < *con; i++ {
		wg.Add(1)
		go start(count)
	}

	wg.Wait()
	fmt.Println(".\n.\n.\nend")
}

func start(count int) {
	defer wg.Done()

	for i := 0; i < count; i++ {
		request()
	}
}

func request() {
	req := path[rand.Intn(len(path))]
	fmt.Println("request", req)
	_, err := client.Get(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}
