package astiworker

import (
	"context"
	"net/http"

	"github.com/molotovtv/go-astilog"
	"github.com/pkg/errors"
)

// Serve spawns a server
func (w *Worker) Serve(addr string, h http.Handler) {
	// Create server
	s := &http.Server{Addr: addr, Handler: h}

	// Create task
	t := w.NewTask()

	// Execute the rest in a goroutine
	astilog.Infof("astiworker: serving on %s", addr)
	go func() {
		// Serve
		var chanDone = make(chan error)
		go func() {
			if err := s.ListenAndServe(); err != nil {
				chanDone <- err
			}
		}()

		// Wait for context or chanDone to be done
		select {
		case <-w.ctx.Done():
			if w.ctx.Err() != context.Canceled {
				astilog.Error(errors.Wrap(w.ctx.Err(), "astiworker: context error"))
			}
		case err := <-chanDone:
			if err != nil {
				astilog.Error(errors.Wrap(err, "astiworker: serving failed"))
			}
		}

		// Shutdown
		astilog.Infof("astiworker: shutting down server on %s", addr)
		if err := s.Shutdown(context.Background()); err != nil {
			astilog.Error(errors.Wrapf(err, "astiworker: shutting down server on %s failed", addr))
		}

		// Task is done
		t.Done()
	}()
}
