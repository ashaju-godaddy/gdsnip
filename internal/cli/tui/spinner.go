package tui

import (
	"context"
	"fmt"
	"time"
)

// Spinner represents a simple CLI spinner
type Spinner struct {
	message string
	frames  []string
	delay   time.Duration
	active  bool
	done    chan bool
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:   100 * time.Millisecond,
		done:    make(chan bool),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.active = true
	go s.run()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	if s.active {
		s.active = false
		s.done <- true
		fmt.Print("\r\033[K") // Clear the line
	}
}

// run is the spinner animation loop
func (s *Spinner) run() {
	i := 0
	for {
		select {
		case <-s.done:
			return
		default:
			frame := s.frames[i%len(s.frames)]
			fmt.Printf("\r%s %s", InfoStyle.Render(frame), s.message)
			time.Sleep(s.delay)
			i++
		}
	}
}

// RunWithSpinner runs a function with a spinner and returns the result
func RunWithSpinner(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	defer spinner.Stop()

	return fn()
}

// RunWithSpinnerAndResult runs a function with a spinner and returns a result
func RunWithSpinnerAndResult[T any](message string, fn func() (T, error)) (T, error) {
	spinner := NewSpinner(message)
	spinner.Start()
	defer spinner.Stop()

	return fn()
}

// RunWithContext runs a function with a spinner and context
func RunWithContext(ctx context.Context, message string, fn func(context.Context) error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	defer spinner.Stop()

	errChan := make(chan error, 1)
	go func() {
		errChan <- fn(ctx)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
