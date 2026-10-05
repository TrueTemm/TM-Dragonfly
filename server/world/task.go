/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package world

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrWorldClosed = errors.New("world: world closed")

	ErrEntityClosed = errors.New("world: entity closed")

	ErrEntityNotInWorld = errors.New("world: entity not in this world")

	ErrTaskCancelled = errors.New("world: scheduled task cancelled")

	ErrTaskPanicked = errors.New("world: scheduled task panicked")

	ErrEntityType = errors.New("world: unexpected entity type")
)

type PanicError struct {
	Value any

	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("world: scheduled task panicked: %v", e.Value)
}

func (e *PanicError) Unwrap() error { return ErrTaskPanicked }

func rethrowPanic(err error) {
	if pe, ok := errors.AsType[*PanicError](err); ok {
		panic(pe.Value)
	}
}

func callContext(ctx context.Context) (context.Context, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx, ctx.Err()
}

func executeWithRecovery(w *World, f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr := &PanicError{Value: r, Stack: debug.Stack()}

			w.conf.Log.Error("scheduled task panicked", "panic", panicErr.Value)
			w.conf.Log.Debug("scheduled task panic stack", "stack", string(panicErr.Stack))
			err = panicErr
		}
	}()
	return f()
}

func awaitTask[T any](ctx context.Context, task *Task, result *T) (T, error) {
	var zero T
	completed := func() (T, error) {
		if err := task.Err(); err != nil {
			rethrowPanic(err)
			return zero, err
		}
		return *result, nil
	}
	select {
	case <-task.Done():
		return completed()
	case <-ctx.Done():
		if task.Cancel() {
			return zero, ctx.Err()
		}
		<-task.Done()
		return completed()
	}
}

const (
	taskPending int32 = iota
	taskRunning
	taskDone
	taskCancelled
)

type Task struct {
	done  chan struct{}
	state atomic.Int32

	errMu sync.Mutex
	err   error

	cancelMu sync.Mutex
	onCancel func()
}

func newTask() *Task {
	return &Task{done: make(chan struct{})}
}

func NewFinishedTask(err error) *Task {
	t := newTask()
	t.failIfPending(err)
	return t
}

var closedDone = func() <-chan struct{} {
	c := make(chan struct{})
	close(c)
	return c
}()

func (t *Task) Done() <-chan struct{} {
	if t == nil || t.done == nil {
		return closedDone
	}
	return t.done
}

func (t *Task) Err() error {
	if t == nil || t.done == nil {
		return ErrTaskCancelled
	}
	select {
	case <-t.done:
		t.errMu.Lock()
		defer t.errMu.Unlock()
		return t.err
	default:
		return nil
	}
}

func (t *Task) Wait(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if t == nil || t.done == nil {
		return ErrTaskCancelled
	}
	select {
	case <-t.done:
		return t.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Task) OnDone(f func(err error)) {
	if t == nil || t.done == nil {
		go f(ErrTaskCancelled)
		return
	}
	go func() {
		<-t.done
		f(t.Err())
	}()
}

func (t *Task) Cancel() bool {
	if t == nil || t.done == nil || !t.state.CompareAndSwap(taskPending, taskCancelled) {
		return false
	}
	t.setErr(ErrTaskCancelled)
	close(t.done)
	t.runCancel()
	return true
}

func (t *Task) begin() bool {
	return t != nil && t.state.CompareAndSwap(taskPending, taskRunning)
}

func (t *Task) failIfPending(err error) bool {
	if t == nil || !t.state.CompareAndSwap(taskPending, taskRunning) {
		return false
	}
	t.finish(err)
	return true
}

func (t *Task) finish(err error) {
	t.setErr(err)
	t.state.Store(taskDone)
	close(t.done)
}

func (t *Task) setErr(err error) {
	t.errMu.Lock()
	t.err = err
	t.errMu.Unlock()
}

func (t *Task) pending() bool {
	return t != nil && t.state.Load() == taskPending
}

func (t *Task) setCancel(f func()) {
	if t == nil || f == nil {
		return
	}
	t.cancelMu.Lock()
	cancelled := t.state.Load() == taskCancelled
	if !cancelled {
		t.onCancel = f
	}
	t.cancelMu.Unlock()
	if cancelled {
		f()
	}
}

func (t *Task) runCancel() {
	t.cancelMu.Lock()
	f := t.onCancel
	t.cancelMu.Unlock()
	if f != nil {
		f()
	}
}

func (w *World) Do(f func(tx *Tx)) *Task {
	return w.scheduleTask(newTask(), func(tx *Tx) error {
		f(tx)
		return nil
	})
}

func (w *World) DoAfter(delay time.Duration, f func(tx *Tx)) *Task {
	t := newTask()
	run := func(tx *Tx) error {
		f(tx)
		return nil
	}
	if delay <= 0 {
		return w.scheduleTask(t, run)
	}
	if w == nil || w.queue == nil || w.closed.Load() {
		t.failIfPending(ErrWorldClosed)
		return t
	}
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-timer.C:
			w.scheduleTask(t, run)
		case <-t.Done():
		case <-w.closeStarted:
			t.failIfPending(ErrWorldClosed)
		case <-w.closing:
			t.failIfPending(ErrWorldClosed)
		case <-w.queueClosing:
			t.failIfPending(ErrWorldClosed)
		}
	}()
	return t
}

func Call[T any](ctx context.Context, w *World, f func(tx *Tx) (T, error)) (T, error) {
	var zero T
	ctx, err := callContext(ctx)
	if err != nil {
		return zero, err
	}
	var result T
	task := w.scheduleTask(newTask(), func(tx *Tx) error {
		var err error
		result, err = f(tx)
		return err
	})
	return awaitTask(ctx, task, &result)
}

func CallEntity[T any](ctx context.Context, h *EntityHandle, f func(tx *Tx, e Entity) (T, error)) (T, error) {
	return CallRef(ctx, NewEntityRef[Entity](h), f)
}

func (w *World) scheduleTask(task *Task, f func(tx *Tx) error) *Task {
	if task == nil {
		task = newTask()
	}
	if w == nil || w.queue == nil || w.closed.Load() {
		task.failIfPending(ErrWorldClosed)
		return task
	}
	if !task.pending() {
		return task
	}
	st := scheduledTransaction{task: task, f: f}
	w.scheduleMu.Lock()
	if w.closed.Load() {
		w.scheduleMu.Unlock()
		task.failIfPending(ErrWorldClosed)
		return task
	}
	if w.conf.Synchronous {
		w.scheduleMu.Unlock()
		st.Run(w)
		return task
	}
	select {
	case <-w.closing:
		task.failIfPending(ErrWorldClosed)
	case <-w.queueClosing:
		task.failIfPending(ErrWorldClosed)
	case w.queue <- st:
	default:
		w.scheduling.Add(1)
		go w.queueScheduled(st)
	}
	w.scheduleMu.Unlock()
	return task
}

func (w *World) queueScheduled(st scheduledTransaction) {
	defer w.scheduling.Done()
	if w.closed.Load() {
		st.task.failIfPending(ErrWorldClosed)
		return
	}
	select {
	case <-w.closing:
		st.task.failIfPending(ErrWorldClosed)
	case <-w.queueClosing:
		st.task.failIfPending(ErrWorldClosed)
	case <-st.task.Done():
	case w.queue <- st:
	}
}

type scheduledTransaction struct {
	task *Task
	f    func(tx *Tx) error
}

func (st scheduledTransaction) Run(w *World) {
	if !st.task.begin() {
		return
	}
	tx := newTx(w)
	err := executeWithRecovery(w, func() error { return st.f(tx) })
	tx.close()
	tx.runDeferred()
	st.task.finish(err)
}
