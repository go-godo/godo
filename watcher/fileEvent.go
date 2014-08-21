package watcher

import (
	"gopkg.in/fsnotify.v1"
	//"log"
)

// FileEvent is a wrapper around github.com/howeyc/fsnotify.FileEvent
type FileEvent struct {
	*fsnotify.Event
	Name     string
	UnixNano int64
}

func newFileEvent(originEvent fsnotify.Event, unixNano int64) *FileEvent {
	//log.Printf("to channel %+v\n", originEvent)
	return &FileEvent{Event: &originEvent, Name: originEvent.Name, UnixNano: unixNano}
}
