package watcher

import (
	"github.com/go-fsnotify/fsnotify"
	"time"
)

// FileEvent is a wrapper around github.com/howeyc/fsnotify.FileEvent
type FileEvent struct {
	*fsnotify.Event
	Name string
	Time time.Time
}

func newFileEvent(originEvent fsnotify.Event) *FileEvent {
	return &FileEvent{Event: &originEvent, Name: originEvent.Name, Time: time.Now()}
}
