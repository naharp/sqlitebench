package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"
)

func main() {
	var (
		dirs      = []string{"crawshaw", "mattn", "modernc", "tailscale", "zombiezen"}
		poolsizes = []int{1, 4, 8, 50, 100}
	)

	for _, dir := range dirs {
		for _, poolsize := range poolsizes {
			func() {
				server := servercmd(dir, poolsize)
				server.Start()
				defer syscall.Kill(-server.Process.Pid, syscall.SIGKILL)

				// Wait for server to start
				time.Sleep(1 * time.Second)

				f := resultfile(dir, poolsize)
				defer f.Close()

				// Fire 100.000 requests to server
				hey := heycmd("http://localhost:8080", 100000)
				stdout, _ := hey.StdoutPipe()
				go io.Copy(f, stdout)
				hey.Run()
			}()
		}
	}

	tabw := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprint(tabw, "package")
	for _, poolsize := range poolsizes {
		fmt.Fprint(tabw, "\trps @ C", poolsize)
	}
	fmt.Fprint(tabw, "\n")

	reqsec := regexp.MustCompile(`Requests/sec:\s+(\d+\.\d+)`)

	for _, dir := range dirs {
		fmt.Fprint(tabw, dir)
		for _, poolsize := range poolsizes {
			filename := fmt.Sprintf("result_%s_poolsize_%d.txt", dir, poolsize)
			data, err := os.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}
			segs := reqsec.FindStringSubmatch(string(data))
			if len(segs) == 0 {
				log.Fatal(fmt.Sprintf("couldn't find req/sec in %s", filename))
			}

			fmt.Fprintf(tabw, "\t%s", segs[1])
		}
		fmt.Fprint(tabw, "\n")
	}
	tabw.Flush()
}

func servercmd(dir string, poolsize int) *exec.Cmd {
	cmd := exec.Command("go", "run", filepath.Join(dir, "main.go"), "-poolsize", strconv.Itoa(poolsize))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}

func heycmd(url string, requests int) *exec.Cmd {
	return exec.Command("hey", "-n", strconv.Itoa(requests), url)
}

func resultfile(dir string, poolsize int) *os.File {
	f, err := os.Create(fmt.Sprintf("result_%s_poolsize_%d.txt", dir, poolsize))
	if err != nil {
		log.Fatal(err)
	}
	return f
}
