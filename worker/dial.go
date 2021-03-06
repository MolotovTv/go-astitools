package astiworker

import (
	"net/http"
	"time"

	"github.com/molotovtv/go-astilog"
	"github.com/molotovtv/go-astiws"
	"github.com/pkg/errors"
)

// DialOptions represents dial options
type DialOptions struct {
	Addr        string
	Client      *astiws.Client
	Header      http.Header
	OnDial      func() error
	OnReadError func(err error)
}

// Dial dials with options
func (w *Worker) Dial(o DialOptions) {
	// Create task
	t := w.NewTask()

	// Execute the rest in a goroutine
	go func() {
		// Dial
		go func() {
			const sleepError = 5 * time.Second
			for {
				// Check context error
				if w.ctx.Err() != nil {
					break
				}

				// Dial
				astilog.Infof("astiworker: dialing %s", o.Addr)
				if err := o.Client.DialWithHeaders(o.Addr, o.Header); err != nil {
					astilog.Error(errors.Wrapf(err, "astiworker: dialing %s failed", o.Addr))
					time.Sleep(sleepError)
					continue
				}

				// Custom callback
				if o.OnDial != nil {
					if err := o.OnDial(); err != nil {
						astilog.Error(errors.Wrapf(err, "astiworker: custom on dial callback on %s failed", o.Addr))
						time.Sleep(sleepError)
						continue
					}
				}

				// Read
				if err := o.Client.Read(); err != nil {
					if o.OnReadError != nil {
						o.OnReadError(err)
					} else {
						astilog.Error(errors.Wrapf(err, "astiworker: reading on %s failed", o.Addr))
					}
					time.Sleep(sleepError)
					continue
				}
			}
		}()

		// Wait for context to be done
		<-w.ctx.Done()

		// Close client
		if err := o.Client.Close(); err != nil {
			astilog.Error(errors.Wrapf(err, "astiworker: closing dial on %s failed", o.Addr))
		}

		// Task is done
		t.Done()
	}()

}
