//go:build darwin || netbsd || freebsd || openbsd || dragonfly

package internal

import (
	"math/rand"
	"syscall"
	"time"
)

var _ ITimer = &Timer{}

type Timer struct {
	fd     int
	poller *Poller
	pd     PollData
}

func NewTimer(poller *Poller) (*Timer, error) {
	t := &Timer{
		fd:     rand.Int(), // TODO figure out something better
		poller: poller,
	}
	t.pd.Fd = t.fd
	return t, nil
}

func (t *Timer) Set(dur time.Duration, cb func()) error {
	// Make sure there's not another timer setup on the same fd.
	if err := t.Unset(); err != nil {
		return err
	}
	t.pd.Set(ReadEvent, func(_ error) { cb() })

	err := t.poller.set(t.fd, createEvent(
		syscall.EV_ADD|syscall.EV_ENABLE|syscall.EV_ONESHOT,
		syscall.EVFILT_TIMER,
		&t.pd,
		dur))
	if err == nil {
		t.poller.pending++
		t.pd.Flags |= ReadFlags
	}
	return nil
}

func (t *Timer) Unset() error {
	if t.pd.Flags&ReadFlags != ReadFlags {
		return nil
	}
	err := t.poller.set(t.fd, createEvent(
		syscall.EV_DELETE|syscall.EV_DISABLE,
		syscall.EVFILT_TIMER, &t.pd, 0))
	if err == nil {
		t.poller.pending--
	}
	return err
}

func (t *Timer) Close() error {
	return t.Unset()
}
