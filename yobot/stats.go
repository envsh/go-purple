package main

import (
	"fmt"
	"sync/atomic"
	"time"

	deadlock "github.com/sasha-s/go-deadlock"
	// "sync"
)

type RunStats struct {
	masterReconnectTimes int32
	ircConnTimeoutCount  int32
	ircReconnCount       int32
	ircConnMaxTime       time.Duration
	ircConnMinTime       time.Duration
	ircConnCount         int32
	ircConnAvgTime       time.Duration
	sendbusTimeoutCount  int32
	handleEventMaxTime   time.Duration
	handleEventMinTime   time.Duration
	handleEventCount     int32
	handleEventAvgTime   time.Duration
	mu                   deadlock.RWMutex
	// mu sync.RWMutex
}

func (this *RunStats) collect() string {
	return fmt.Sprintf("runstats:\n"+
		"ircConnTOCount: %v\n"+
		"ircReconnCount: %v\n"+
		"ircConnMaxTime: %v\n"+
		"ircConnMinTime: %v\n"+
		"ircConnAvgTime: %v\n"+
		"sendbusTOCount: %v\n"+
		"handleEventMaxTime: %v\n"+
		"handleEventMinTime: %v\n"+
		"handleEventAvgTime: %v",
		this.ircConnTimeoutCount,
		this.ircReconnCount,
		this.ircConnMaxTime,
		this.ircConnMinTime,
		this.ircConnAvgTime,
		this.sendbusTimeoutCount,
		this.handleEventMaxTime,
		this.handleEventMinTime,
		this.handleEventAvgTime)
}

func NewRunStats() *RunStats {
	this := &RunStats{}
	this.ircConnMinTime = 2 << 31 * time.Second
	this.handleEventMinTime = 2 << 31 * time.Second

	return this
}
func (this *RunStats) masterIrcReconnect() {
	atomic.AddInt32(&this.masterReconnectTimes, 1)
}
func (this *RunStats) ircReconnect() {
	atomic.AddInt32(&this.ircReconnCount, 1)
}
func (this *RunStats) ircConnTimeout() {
	atomic.AddInt32(&this.ircConnTimeoutCount, 1)
}
func (this *RunStats) ircConnTime(btime time.Time) {
	atomic.AddInt32(&this.ircConnCount, 1)

	this.mu.Lock()
	defer this.mu.Unlock()
	etime := time.Now()
	dtime := etime.Sub(btime)
	if dtime > this.ircConnMaxTime {
		this.ircConnMaxTime = dtime
	}
	if dtime < this.ircConnMinTime {
		this.ircConnMinTime = dtime
	}

	if this.ircConnCount == 1 {
		this.ircConnAvgTime = dtime
	} else {
		this.ircConnAvgTime = time.Duration((this.ircConnAvgTime.Nanoseconds()*int64(this.ircConnCount-1) + dtime.Nanoseconds()) / int64(this.ircConnCount))
	}
}

func (this *RunStats) sendbusTimeout() {
	this.sendbusTimeoutCount += 1
	atomic.AddInt32(&this.sendbusTimeoutCount, 1)
}

func (this *RunStats) handleEventTime(btime time.Time) {
	atomic.AddInt32(&this.handleEventCount, 1)

	this.mu.Lock()
	defer this.mu.Unlock()
	etime := time.Now()
	dtime := etime.Sub(btime)
	if dtime > this.handleEventMaxTime {
		this.handleEventMaxTime = dtime
	}
	if dtime < this.handleEventMinTime {
		this.handleEventMinTime = dtime
	}

	if this.handleEventCount == 1 {
		this.handleEventAvgTime = dtime
	} else {
		this.handleEventAvgTime = time.Duration((this.handleEventAvgTime.Nanoseconds()*int64(this.handleEventCount-1) + dtime.Nanoseconds()) / int64(this.handleEventCount))
	}
}

// 消息，事件，用户统计
type EventStats struct {
	eventCount uint64
}

func NewEventStats() *EventStats {
	this := &EventStats{}
	return this
}
