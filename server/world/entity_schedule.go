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

import "time"

func (e *EntityHandle) Do(f func(tx *Tx, e Entity)) *Task {
	return e.schedule(func(tx *Tx, e Entity) error {
		f(tx, e)
		return nil
	})
}

func (e *EntityHandle) DoAfter(delay time.Duration, f func(tx *Tx, e Entity)) *Task {
	return e.scheduleAfter(delay, func(tx *Tx, e Entity) error {
		f(tx, e)
		return nil
	})
}

func (e *EntityHandle) schedule(f func(tx *Tx, e Entity) error) *Task {
	task := newTask()
	if e == nil {
		task.failIfPending(ErrEntityClosed)
		return task
	}
	task.setCancel(e.wakeScheduled)
	w := e.trackCloseSchedule(task)
	if !task.pending() {
		return task
	}
	run := func() {
		if w != nil {
			defer w.scheduling.Done()
		}
		e.runScheduled(task, f, w)
	}
	if e.currentWorldSynchronous() {
		run()
	} else {
		go run()
	}
	return task
}

func (e *EntityHandle) scheduleSync(f func(tx *Tx, e Entity) error) *Task {
	task := newTask()
	if e == nil {
		task.failIfPending(ErrEntityClosed)
		return task
	}
	task.setCancel(e.wakeScheduled)
	w := e.trackCloseSchedule(task)
	if !task.pending() {
		return task
	}
	if w != nil {
		defer w.scheduling.Done()
	}
	e.runScheduled(task, f, w)
	return task
}

func (e *EntityHandle) scheduleAfter(delay time.Duration, f func(tx *Tx, e Entity) error) *Task {
	if delay <= 0 {
		return e.schedule(f)
	}
	task := newTask()
	if e == nil {
		task.failIfPending(ErrEntityClosed)
		return task
	}
	task.setCancel(e.wakeScheduled)
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		for {
			closeStarted, worldChanged := e.currentWorldSignals()
			select {
			case <-timer.C:
				if e.currentWorldClosing() {
					task.failIfPending(ErrWorldClosed)
					return
				}
				e.runScheduled(task, f, nil)
				return
			case <-task.Done():
				return
			case <-closeStarted:
				if cs, _ := e.currentWorldSignals(); cs == closeStarted {
					task.failIfPending(ErrWorldClosed)
					return
				}
			case <-worldChanged:
			case <-e.closed:
				task.failIfPending(ErrEntityClosed)
				return
			}
		}
	}()
	return task
}

func (e *EntityHandle) trackCloseSchedule(task *Task) *World {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()
	w := e.w
	if w == nil || w == closeWorld || !w.closed.Load() {
		return nil
	}
	w.scheduleMu.Lock()
	defer w.scheduleMu.Unlock()
	if !w.closeAcceptingEntityTasks.Load() {
		task.failIfPending(ErrWorldClosed)
		return nil
	}
	select {
	case <-w.queueClosing:
		task.failIfPending(ErrWorldClosed)
		return nil
	default:
		w.scheduling.Add(1)
		return w
	}
}

func (e *EntityHandle) wakeScheduled() {
	e.cond.L.Lock()
	e.cond.Broadcast()
	e.cond.L.Unlock()
}

func (e *EntityHandle) currentWorldSignals() (<-chan struct{}, <-chan struct{}) {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()
	if e.worldChanged == nil {
		e.worldChanged = make(chan struct{})
	}
	if e.w == nil || e.w == closeWorld {
		return nil, e.worldChanged
	}
	return e.w.closeStarted, e.worldChanged
}

func (e *EntityHandle) currentWorldSynchronous() bool {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()
	return e.w != nil && e.w != closeWorld && e.worldReady && e.w.conf.Synchronous
}

func (e *EntityHandle) currentWorldClosing() bool {
	closeStarted, _ := e.currentWorldSignals()
	return cancelled(closeStarted)
}

func (e *EntityHandle) runScheduled(task *Task, f func(tx *Tx, e Entity) error, allowedCloseWorld *World) {
	run := e.execWorld(func(tx *Tx, ent Entity) {
		if !task.begin() {
			return
		}
		err := executeWithRecovery(tx.w, func() error { return f(tx, ent) })
		tx.runDeferred()
		task.finish(err)
	}, false, task.Done(), allowedCloseWorld)
	if !run || task.pending() {
		err := ErrEntityClosed
		if e.currentWorldClosing() {
			err = ErrWorldClosed
		}
		task.failIfPending(err)
	}
}
