package main

import (
	"bufio"
	"flag"
	"net/http"
	"os"
	"time"

	"idstam/gahttp"

	"github.com/fatih/color"
)

var concurrency int
var times int
var status int
var actualHost string

func Banner() {
	color.HiCyan(`
 ____ _____________.____   __________              ___.
|    |   \______   \    |  \______   \_______  ____\_ |__   ____
|    |   /|       _/    |   |     ___/\_  __ \/  _ \| __ \_/ __ \
|    |  / |    |   \    |___|    |     |  | \(  <_> ) \_\ \  ___/
|______/  |____|_  /_______ \____|     |__|   \____/|___  /\___  >
                 \/        \/                           \/     \/   V 2.1`)
	color.HiYellow("           URLProbe:- Urls Status Code & ContentLength Checker")
	color.HiRed("              https://github.com/1ndianl33t")

	color.HiCyan("-------------------------------------------------------------------------")
}
func printStatus(req *http.Request, resp *http.Response, err error) {
	if err != nil {
		return
	}
	StatusCheck(req, resp)

}

func ParseArguments() {
	flag.IntVar(&concurrency, "c", 500, "Number of workers to use..default 500")
	flag.IntVar(&status, "s", 1, "If enabled..then check for specific status")
	flag.IntVar(&times, "t", 05, "Set rate limit")
	flag.StringVar(&actualHost, "h", "", "Host adress to use instead of host in url. Will put the url host into the host-header. Can include :port if needed")
	flag.Parse()
}

func StatusCheck(req *http.Request, resp *http.Response) {
	if status != 1 {
		if status == resp.StatusCode {
			if status == 404 {
				color.HiRed("[%d] L %d : %s %s\n", resp.StatusCode, resp.ContentLength, req.Host, req.URL)
			}
			if status != 404 {
				color.HiCyan("[%d] L %d : %s %s\n", resp.StatusCode, resp.ContentLength, req.Host, req.URL)
			}

		}
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			color.HiGreen("[%d] L %d : %s %s\n", resp.StatusCode, resp.ContentLength, req.Host, req.URL)

		}
		if resp.StatusCode >= 300 && resp.StatusCode <= 308 {
			color.HiBlue("[%d] L %d : %s %s\n", resp.StatusCode, resp.ContentLength, req.Host, req.URL)

		}
		if resp.StatusCode >= 400 && resp.StatusCode <= 451 {
			if resp.StatusCode == 404 {
				color.HiRed("[%d] L %d : %s %s\n", resp.StatusCode, resp.ContentLength, req.Host, req.URL)
			} else {
				color.HiCyan("[%d] L %d : %s %s\n", resp.StatusCode, resp.ContentLength, req.Host, req.URL)
			}

		}

		if resp.StatusCode >= 500 && resp.StatusCode <= 511 {
			color.HiCyan("[%d] L %d : %s %s\n", resp.StatusCode, resp.ContentLength, req.Host, req.URL)

		}
	}
}

func main() {
	Banner()
	ParseArguments()
	p := gahttp.NewPipeline()
	p.SetConcurrency(concurrency)
	p.SetRateLimit(time.Duration(times) * time.Second)
	urls := gahttp.Wrap(printStatus, gahttp.CloseBody)
	sc := bufio.NewScanner(os.Stdin)

	if actualHost == "" {
		for sc.Scan() {
			p.Get(sc.Text(), urls)
		}
	} else {
		for sc.Scan() {
			p.GetFromHost(sc.Text(), actualHost, urls)
		}
	}
	p.Done()
	p.Wait()
}
