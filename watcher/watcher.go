// Package watcher implements filesystem notification,.
package watcher

import (
	//"fmt"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gokyle/fswatch"
	"github.com/mgutz/str"
)

const (
	// IgnoreThresholdRange is the amount of time in ns to ignore when
	// receiving watch events for the same file
	IgnoreThresholdRange = 50 * 1000000 // convert to ms
	// WatchDelay is the time to poll the file system
	WatchDelay = 1100 * time.Millisecond
)

// Watcher is a wrapper around which adds some additional features:
//
// - recursive directory watch
// - buffer to even chan
// - even time
//
// Original work from https://github.com/bronze1man/kmg
type Watcher struct {
	*fswatch.Watcher
	Event chan *FileEvent
	Error chan error
	//default ignore all file start with "."
	IsIgnorePath func(path string) bool
	//default is nil,if is nil ,error send through Error chan,if is not nil,error handle by this func
	ErrorHandler func(err error)
	isClosed     bool
	quit         chan bool
}

// NewWatcher creates an instance of watcher.
func NewWatcher(bufferSize int) (watcher *Watcher, err error) {
	fswatch.WatchDelay = WatchDelay
	fswatcher := fswatch.NewAutoWatcher()

	if err != nil {
		return nil, err
	}
	watcher = &Watcher{
		Watcher:      fswatcher,
		Error:        make(chan error, 10),
		Event:        make(chan *FileEvent, bufferSize),
		IsIgnorePath: DefaultIsIgnorePath,
	}
	go watcher.eventLoop()
	return
}

// Close closes the watcher channels.
func (w *Watcher) Close() error {
	if w.isClosed {
		return nil
	}
	w.Watcher.Stop()
	w.quit <- true
	w.isClosed = true
	return nil
}

var eventText = map[int]string{
	fswatch.CREATED:  "was created",
	fswatch.DELETED:  "was deleted",
	fswatch.MODIFIED: "was modified",
	fswatch.PERM:     "permissions changed",
	fswatch.NOEXIST:  "doesn't exist",
	fswatch.NOPERM:   "has invalid permissions",
	fswatch.INVALID:  "is invalid",
}

func debugEvent(event int, path string) {
	fmt.Printf("%s %s\n", path, eventText[event])
}

func (w *Watcher) eventLoop() {
	cache := map[string]*os.FileInfo{}
	mu := &sync.Mutex{}

	coutput := w.Watcher.Start()
	for {
		event, ok := <-coutput
		if !ok {
			return
		}

		//fmt.Printf("event %+v\n", event)
		if w.IsIgnorePath(event.Path) {
			continue
		}

		//debugEvent(event.Event, event.Path)

		// you can not stat a delete file...
		if event.Event == fswatch.DELETED || event.Event == fswatch.NOEXIST {
			// adjust with arbitrary value because it was deleted
			// before it got here
			w.Event <- newFileEvent(event.Event, event.Path, time.Now().UnixNano()-10)
			continue
		}

		fi, err := os.Stat(event.Path)
		if os.IsNotExist(err) {
			//fmt.Println("not exists", event)
			continue
		}

		// fsnotify is sending multiple MODIFY events for the same
		// file which is likely OS related. The solution here is to
		// compare the current stats of a file against its last stats
		// (if any) and if it falls within a nanoseconds threshold,
		// ignore it.
		mu.Lock()
		oldFI := cache[event.Path]
		cache[event.Path] = &fi
		mu.Unlock()

		if oldFI != nil && fi.ModTime().UnixNano() < (*oldFI).ModTime().UnixNano()+IgnoreThresholdRange {
			continue
		}

		//fmt.Println("sending fi", fi.ModTime().UnixNano()/1000000, event.Name)
		w.Event <- newFileEvent(event.Event, event.Path, fi.ModTime().UnixNano())

		if err != nil {
			//rename send two events,one old file,one new file,here ignore old one
			if os.IsNotExist(err) {
				continue
			}
			w.errorHandle(err)
			continue
		}
		// if fi.IsDir() {
		// 	w.WatchRecursive(event.Name)
		// }
		// case err := <-w.Watcher.Errors:
		// 	w.errorHandle(err)
		// case _ = <-w.quit:
		// 	break
		// }
	}
}
func (w *Watcher) errorHandle(err error) {
	if w.ErrorHandler == nil {
		w.Error <- err
		return
	}
	w.ErrorHandler(err)
}

// GetErrorChan gets error chan.
func (w *Watcher) GetErrorChan() chan error {
	return w.Error
}

// GetEventChan gets event chan.
func (w *Watcher) GetEventChan() chan *FileEvent {
	return w.Event
}

// WatchRecursive watches a directory recursively. If a dir is created
// within directory it is also watched.
func (w *Watcher) WatchRecursive(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	folders, err := w.getSubFolders(path)
	for _, v := range folders {
		w.Watcher.Add(v)
	}
	return nil
}

func (w *Watcher) getSubFolders(path string) (paths []string, err error) {
	err = filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}
		if w.IsIgnorePath(newPath) {
			return filepath.SkipDir
		}
		paths = append(paths, newPath)
		return nil
	})
	return paths, err
}

// DefaultIsIgnorePath checks whether a path is ignored. Currently defaults
// to hidden files on *nix systems, ie they start with a ".".
func DefaultIsIgnorePath(path string) bool {
	if strings.HasPrefix(path, ".") || strings.Contains(path, "/.") {
		return true
	}

	// ignore node
	if strings.HasPrefix(path, "node_modules") || strings.Contains(path, "/node_modules") {
		return true
	}

	return isDotFile(path) || isVimFile(path)
}

func isDotFile(path string) bool {
	if strings.HasPrefix(path, ".") || strings.Contains(path, "/.") {
		return true
	}
	return false
}

func isVimFile(path string) bool {
	base := filepath.Base(path)
	return str.IsNumeric(base)
}
