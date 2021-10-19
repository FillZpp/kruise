package manager

import (
	"context"
	"errors"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	errRunnableGroupStopped = errors.New("can't accept new runnable as stop procedure is already engaged")
)

// readyRunnable encapsulates a runnable with
// a ready check.
type readyRunnable struct {
	Runnable
	Check runnableCheck
	Ready bool
}

// runnableCheck can be passed to Add() to let the runnable group determine that a
// runnable is ready.
type runnableCheck func(ctx context.Context) bool

// runnables handles all the runnable for a manager by grouping them accordingly to their
// type (webhooks, caches etc.).
type runnables struct {
	Webhooks       *runnableGroup
	Caches         *runnableGroup
	LeaderElection *runnableGroup
	Others         *runnableGroup
}

// newRunnables creates a new runnables object.
func newRunnables(errChan chan error) *runnables {
	return &runnables{
		Webhooks:       newRunnableGroup(errChan),
		Caches:         newRunnableGroup(errChan),
		LeaderElection: newRunnableGroup(errChan),
		Others:         newRunnableGroup(errChan),
	}
}

// Add adds a runnable and its ready check to the closest
// group of runnable that they belong to.
func (r *runnables) Add(fn Runnable, ready runnableCheck) error {
	switch runnable := fn.(type) {
	case hasCache:
		return r.Caches.Add(fn, func(ctx context.Context) bool {
			// Run the ready check for the cache a fixed number of times
			// backing off a bit; this is to give time to the runnables
			// to start up before their health check returns true.
			for i := 0; i < 10; i++ {
				<-time.After(time.Duration(i) * 10 * time.Millisecond)
				if !runnable.GetCache().WaitForCacheSync(ctx) {
					continue
				}
				return true
			}
			return false
		})
	case *webhook.Server:
		return r.Webhooks.Add(fn, ready)
	case LeaderElectionRunnable:
		if !runnable.NeedLeaderElection() {
			return r.Others.Add(fn, ready)
		}
		return r.LeaderElection.Add(fn, ready)
	default:
		return r.LeaderElection.Add(fn, ready)
	}
}

// runnableGroup manages a group of runnable that are
// meant to be running together until Stop is called.
//
// Runnables can be added to a group after the group has started
// but not after it's stopped or while shutting down.
type runnableGroup struct {
	ctx    context.Context
	cancel context.CancelFunc

	start     sync.Mutex
	startOnce sync.Once
	started   bool

	stop     sync.RWMutex
	stopOnce sync.Once
	stopped  bool

	// errChan is the error channel passed by the caller
	// when the group is created.
	// All errors are forwarded to this channel once they occur.
	errChan chan error

	// ch is the internal channel where the runnables are read off from.
	ch chan *readyRunnable

	// wg is an internal sync.WaitGroup that allows us to properly stop
	// and wait for all the runnables to finish before returning.
	wg *sync.WaitGroup

	// group is a sync.Map that contains every runnable ever.
	// The key of the map is the runnable itself (key'd by pointer),
	// while the value is its ready state.
	//
	// The group of runnable is append-only, runnables scheduled
	// through this group are going to be stored in this internal map
	// until the application exits. The limit is the available memory.
	group *sync.Map
}

func newRunnableGroup(errChan chan error) *runnableGroup {
	r := &runnableGroup{
		errChan: errChan,
		ch:      make(chan *readyRunnable),
		wg:      new(sync.WaitGroup),
		group:   new(sync.Map),
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	return r
}

// Started returns true if the group has started.
func (r *runnableGroup) Started() bool {
	r.start.Lock()
	defer r.start.Unlock()
	return r.started
}

// StartAndWaitReady starts all the runnables previously
// added to the group and waits for all to report ready.
func (r *runnableGroup) StartAndWaitReady(ctx context.Context) error {
	r.Start()
	return r.WaitReady(ctx)
}

// Start starts the group, it can only be called once.
func (r *runnableGroup) Start() {
	r.startOnce.Do(func() {
		go r.reconcile()
		r.start.Lock()
		r.started = true
		r.group.Range(func(key, _ interface{}) bool {
			r.ch <- key.(*readyRunnable)
			return true
		})
		r.start.Unlock()
	})
}

// reconcile is our main entrypoint for every runnable added
// to this group. Its primary job is to read off the internal channel
// and schedule runnables while tracking their state.
func (r *runnableGroup) reconcile() {
	for runnable := range r.ch {
		// Handle stop.
		// If the shutdown has been called we want to avoid
		// adding new goroutines to the WaitGroup because Wait()
		// panics if Add() is called after it.
		{
			r.stop.RLock()
			if r.stopped {
				// Drop any runnables if we're stopped.
				r.errChan <- errRunnableGroupStopped
				r.stop.RUnlock()
				continue
			}

			// Why is this here?
			// When StopAndWait is called, if a runnable is in the process
			// of being added, we could end up in a situation where
			// the WaitGroup is incremented while StopAndWait has called Wait(),
			// which would result in a panic.
			r.wg.Add(1)
			r.stop.RUnlock()
		}

		// Start the runnable.
		go func(rn *readyRunnable) {
			go func() {
				if rn.Check(r.ctx) {
					r.group.Store(rn, true)
				}
			}()

			// If we return the runnable ended cleanly
			// or returned an error to the channel.
			//
			// We should always decrement the waigroup and
			// mark the runnable as ready.
			defer r.wg.Done()
			defer r.group.Store(rn, true)

			// Start the runnable.
			if err := rn.Start(r.ctx); err != nil {
				r.errChan <- err
			}
		}(runnable)
	}
}

// WaitReady polls until the group is ready until the context is cancelled.
func (r *runnableGroup) WaitReady(ctx context.Context) error {
	return wait.PollImmediateInfiniteWithContext(ctx,
		100*time.Millisecond,
		func(_ context.Context) (bool, error) {
			if !r.Started() {
				return false, nil
			}
			ready, total := 0, 0
			r.group.Range(func(_, value interface{}) bool {
				total++
				if rd, ok := value.(bool); ok && rd {
					ready++
				}
				return true
			})
			return ready == total, nil
		},
	)
}

// Add should be able to be called before and after Start, but not after shutdown.
// Add should return an error when called during shutdown.
func (r *runnableGroup) Add(rn Runnable, ready runnableCheck) error {
	r.stop.RLock()
	if r.stopped {
		r.stop.RUnlock()
		return errRunnableGroupStopped
	}
	r.stop.RUnlock()

	// If we don't have a readiness check, always return true.
	if ready == nil {
		ready = func(_ context.Context) bool { return true }
	}

	readyRunnable := &readyRunnable{
		Runnable: rn,
		Check:    ready,
	}

	// Handle start.
	// If the overall runnable group isn't started yet
	// we want to buffer the runnables and let Start()
	// queue them up again later.
	{
		r.start.Lock()

		// Store the runnable in the internal buffer.
		r.group.Store(readyRunnable, false)

		// Check if we're already started.
		if !r.started {
			r.start.Unlock()
			return nil
		}
		r.start.Unlock()
	}

	// Enqueue the runnable.
	r.ch <- readyRunnable
	return nil
}

// StopAndWait waits for all the runnables to finish before returning.
func (r *runnableGroup) StopAndWait(ctx context.Context) {
	r.stopOnce.Do(func() {
		r.Start()
		r.stop.Lock()
		// Store the stopped variable so we don't accept any new
		// runnables for the time being.
		r.stopped = true
		r.stop.Unlock()

		// Cancel the internal channel.
		r.cancel()
		// Close the reconciler channel.
		close(r.ch)

		done := make(chan struct{})
		go func() {
			defer close(done)
			// Wait for all the runnables to finish.
			r.wg.Wait()
		}()

		select {
		case <-done:
			// We're done, exit.
		case <-ctx.Done():
			// Calling channel has expired, exit.
		}
	})
}
