// https://gist.github.com/napicella/777e83c0ef5b77bf72c0a5d5da9a4b4e

// Companion code for the Linux terminals blog series: https://dev.to/napicella/linux-terminals-tty-pty-and-shell-192e
// I have simplified the code to highlight the interesting bits for the purpose of the blog post:
// - windows resizing is not addressed
// - client does not catch signals (CTRL + C, etc.) to gracefully close the tcp connection
//
// Build: go build -o remote main.go
// In one terminal run: ./remote -server
// In another terminal run: ./remote
//
// Run on multiple machines:
// In the client function, replace the loopback address with IP of the machine, then rebuild
// Beware the unecrypted TCP connection!
package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// Create command
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, e := pty.Start(c)
	if e != nil {
		log.Println(e)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.
	sc := bufio.NewScanner(os.Stdin)
	go func() {
		//		io.Copy(ptmx, os.Stdout)
		for {
			if sc.Scan() {
				log.Println(sc.Text())
				ptmx.Write(append(sc.Bytes(), '\n'))
			}
		}
	}()
	// output to console
	buf := make([]byte, 1024)
	for {
		n, err := ptmx.Read(buf)
		if err != nil {
			log.Println(err)
		} else {
			os.Stdin.Write(buf[:n])
		}

	}
	//	io.Copy(os.Stdout, ptmx)
}
