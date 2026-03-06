package debouncer

import (
	"log/slog"
	"time"

	"github.com/fsnotify/fsnotify"
)

func New(events <-chan fsnotify.Event, waitTime time.Duration) <-chan struct{} {
	buildSignal := make(chan struct{})

	go func() {
		var timer *time.Timer
		var timerC <-chan time.Time

		for {
			select {
			case event, ok := <-events:
				if !ok {
					return
				}

				slog.Debug("Event caught, resetting debounce timer", "file", event.Name)

				if timer != nil {
					timer.Stop()
				}

				timer = time.NewTimer(waitTime)
				timerC = timer.C

			case <-timerC:
				slog.Info("Changes settled, firing build signal!")

				buildSignal <- struct{}{}

				timerC = nil
			}
		}
	}()

	return buildSignal
}
