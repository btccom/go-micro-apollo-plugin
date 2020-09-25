package apollo

import (
	"errors"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/zouyx/agollo/v3/storage"
	"strings"
	"time"
)

type watcher struct {
	namespaceName string
	name          string
	exit          chan bool
	eventChan     chan *storage.ChangeEvent
}

func (w *watcher) OnChange(changeEvent *storage.ChangeEvent) {
	log.Info("change listener.")
	log.Info(changeEvent.Changes)
	log.Info(changeEvent.Namespace)
	w.eventChan <- changeEvent
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case event := <-w.eventChan:
		var format string
		var content []byte

		namespaceParts := strings.Split(w.namespaceName, ".")
		if len(namespaceParts) > 1 {
			format = namespaceParts[len(namespaceParts)-1]
			content = []byte(event.Changes["content"].NewValue)
		} else {
			// TODO 遍历 event.Changes
		}

		cs := &source.ChangeSet{
			Timestamp: time.Now(),
			Format:    format,
			Source:    w.name,
			Data:      content,
		}
		cs.Checksum = cs.Sum()
		return cs, nil
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	default:
	}
	return nil
}

func newWatcher(name string, namespaceName string) (*watcher, error) {
	return &watcher{
		namespaceName: namespaceName,
		name:          name,
		exit:          make(chan bool),
		eventChan:     make(chan *storage.ChangeEvent),
	}, nil
}
