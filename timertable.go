package timertable

import (
	"fmt"
	"sync"
	"time"
)

type TimerTable struct {
	table map[string]Timer
	mu    sync.RWMutex
}

type Timer struct {
	interval uint
	timer    *time.Timer
	begin    time.Time
	end      time.Time
}

func New() *TimerTable {
	return &TimerTable{table: make(map[string]Timer)}
}

func (tt *TimerTable) Add(key string, after uint, interval uint, f func()) bool {
	tt.mu.Lock()
	_, ok := tt.table[key]
	tt.mu.Unlock()
	if ok {
		fmt.Printf("TimerTable::Add %s fail\n", key)
		return false
	}

	fmt.Printf("TimerTable::Add %s\n", key)
	tt.Set(key, after, interval, f)
	return true
}

func (tt *TimerTable) Set(key string, after uint, interval uint, f func()) {
	if after < 1 {
		after = 1
	}

	tmr := Timer{
		interval: interval,
		begin:    time.Now(),
		end:      time.Now().Add(time.Duration(after) * time.Second),
	}

	tmr.timer = time.AfterFunc(time.Duration(after)*time.Second, func() {
		f()
		if interval > 0 {
			tmr.timer.Reset(time.Duration(tmr.interval) * time.Second)
		} else {
			tmr.timer.Stop()
			tt.mu.Lock()
			delete(tt.table, key)
			tt.mu.Unlock()
		}
	})

	tt.mu.Lock()
	defer tt.mu.Unlock()
	tt.table[key] = tmr
}

func (tt *TimerTable) IsExist(key string) bool {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	_, ok := tt.table[key]

	return ok
}

func (tt *TimerTable) TimeStamp(key string) []int64 {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	tmr, ok := tt.table[key]
	if ok {
		return []int64{tmr.begin.Unix(), tmr.end.Unix()}
	}

	return nil
}

func (tt *TimerTable) RemainTime(key string) int {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	tmr, ok := tt.table[key]
	if ok {
		rt := int(tmr.end.Sub(time.Now()).Seconds() + 1) // floor
		if rt < 0 {
			rt = 0
		}

		return rt
	}

	return 0
}

func (tt *TimerTable) Clear(key string) bool {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	fmt.Printf("TimerTable::Clear %s\n", key)

	tmr, ok := tt.table[key]

	if ok {
		tmr.timer.Stop()
		delete(tt.table, key)

		return true
	}

	return false
}

func (tt *TimerTable) ClearAll() {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	fmt.Println("TimerTable::ClearAll")
	for key, tmr := range tt.table {
		tmr.timer.Stop()
		delete(tt.table, key)
	}
}
