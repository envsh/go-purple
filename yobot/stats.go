package main

import (
	"fmt"
	"time"

	deadlock "github.com/sasha-s/go-deadlock"
)

type RunStats struct {
	ircConnTimeoutCount int
	ircReconnCount      int
	ircConnMaxTime      time.Duration
	ircConnMinTime      time.Duration
	ircConnCount        int
	ircConnAvgTime      time.Duration
	sendbusTimeoutCount int
	handleEventMaxTime  time.Duration
	handleEventMinTime  time.Duration
	handleEventCount    int
	handleEventAvgTime  time.Duration
	mu                  deadlock.RWMutex
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

func (this *RunStats) ircReconnect() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.ircReconnCount += 1
}
func (this *RunStats) ircConnTimeout() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.ircConnTimeoutCount += 1
}
func (this *RunStats) ircConnTime(btime time.Time) {
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
	this.ircConnCount += 1
	if this.ircConnCount == 1 {
		this.ircConnAvgTime = dtime
	} else {
		this.ircConnAvgTime = time.Duration((this.ircConnAvgTime.Nanoseconds()*int64(this.ircConnCount-1) + dtime.Nanoseconds()) / int64(this.ircConnCount))
	}
}

func (this *RunStats) sendbusTimeout() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.sendbusTimeoutCount += 1
}

func (this *RunStats) handleEventTime(btime time.Time) {
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
	this.handleEventCount += 1
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
