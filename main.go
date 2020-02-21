package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	s "strings"
	"sync"
	"sync/atomic"
	"time"
)

//CGoRoutinDefault default gouroutin count
const CGoRoutinDefault int = 5

//CCrowlerTimeOut two second timeout
const CCrowlerTimeOut int = 1

//total number of go word
var total int = 0

//GetParamGoroutine from cmd with flag -g
func GetParamGoroutine() int {

	var cmdVar int
	flag.IntVar(&cmdVar, "g", CGoRoutinDefault, "Goroutin number")
	flag.Parse()
	//	spew.Dump(cmdVar)
	return cmdVar
}

// Crowler in url, out number of match GO word
func Crowler(url string, URLTimeOut int) int {
	res := 0
	client := http.Client{
		Timeout: time.Duration(URLTimeOut) * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return -1
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			//fmt.Printf("%s", err)
			return -1
		}
		//body := s.ToUpper(string(contents))
		res = s.Count(string(contents), "Go")
	}
	return res
}

// Run crowler goroutin index, *totalcount, channel, waitgroup
func WorkCrowler(i int, count *int64, cmdChan chan string, wg *sync.WaitGroup) {
	//func WorkCrowler(i int, count *int64, cmdChan chan string, resChan chan int, wg *sync.WaitGroup) {
	for url := range cmdChan {
		res := Crowler(url, CCrowlerTimeOut)
		total := atomic.LoadInt64(count)
		if res > 0 {
			total = atomic.AddInt64(count, int64(res))
		}
		fmt.Println("total:", total, " N: ", i, " res ", res, " Url:", url)
		//resChan <- res
	}
	wg.Done()
}

func main() {

	var TotalGo int64 = 0
	var startGoNum int = 0

	goRoutinCount := GetParamGoroutine() //get total goroutin default or default count
	scanner := bufio.NewScanner(os.Stdin)

	//mutex := new(sync.Mutex)
	cmdChan := make(chan string)
	resChan := make(chan int)
	wg := new(sync.WaitGroup)
	for scanner.Scan() {
		if startGoNum < goRoutinCount { //Run max goRoutinCount
			wg.Add(1)
			//go WorkCrowler(startGoNum, &TotalGo, cmdChan, resChan, wg)
			go WorkCrowler(startGoNum, &TotalGo, cmdChan, wg)
			startGoNum++
		}
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		cmdChan <- line //send url to channel
		// select {
		// case cmdChan <- line:
		// case res := <-resChan:
		// 	if res > 0 {
		// 		TotalGo = TotalGo + int64(res)
		// 		fmt.Println(TotalGo)
		// 	}

		// }

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Println("end")
	close(cmdChan)
	close(resChan)
	wg.Wait()

	fmt.Println("End, total word 'Go' count:", TotalGo)
}
