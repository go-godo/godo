// clinotify provides an example file system watching command line app. It
// scans the file system, and every 15 seconds prints out the files being
// watched and their current state.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/davars/godo/watcher/fswatch"
)

var dur time.Duration

func init() {
	if len(os.Args) == 1 {
		fmt.Println("[+] not watching anything, exiting.")
		os.Exit(1)
	}
	dur, _ = time.ParseDuration("15s")
}

func main() {
	var w *fswatch.Watcher

	auto_watch := flag.Bool("a", false, "auto add new files in directories")
	flag.Parse()
	paths := flag.Args()
	if *auto_watch {
		w = fswatch.NewAutoWatcher(paths...)
	} else {
		w = fswatch.NewWatcher(paths...)
	}
	fmt.Println("[+] listening...")

	l := w.Start()
	go func() {
		for {
			n, ok := <-l
			if !ok {
				return
			}
			var status_text string
			switch n.Event {
			case fswatch.CREATED:
				status_text = "was created"
			case fswatch.DELETED:
				status_text = "was deleted"
			case fswatch.MODIFIED:
				status_text = "was modified"
			case fswatch.PERM:
				status_text = "permissions changed"
			case fswatch.NOEXIST:
				status_text = "doesn't exist"
			case fswatch.NOPERM:
				status_text = "has invalid permissions"
			case fswatch.INVALID:
				status_text = "is invalid"
			}
			fmt.Printf("[+] %s %s\n", n.Path, status_text)
		}
	}()
	go func() {
		for {
			<-time.After(dur)
			if !w.Active() {
				fmt.Println("[!] not watching anything")
				os.Exit(1)
			}
			fmt.Printf("[-] watching: %+v\n", w.State())
		}
	}()
	time.Sleep(60 * time.Second)
	fmt.Println("[+] stopping...")
	w.Stop()
	time.Sleep(5 * time.Second)
}
