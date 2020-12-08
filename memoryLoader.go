package main // import "github.com/lwahlmeier/goloader"

import (
	"crypto/rand"
	"runtime"
	"time"
)

func MemoryWatcher() {
	memMax := uint64(0)
	memRate := time.Duration(0)
	var stopChannel chan bool = nil
	for {
		time.Sleep(time.Second * 5)
		tmpMemMax := config.GetUint64("mem.max")
		tmpMemRate := config.GetDuration("mem.rate")
		log.Debug("memload:  max:{} to:{}, rate:{} to:{}", memMax, tmpMemMax, memRate, tmpMemRate)
		if tmpMemMax == memMax && tmpMemRate == memRate {
			continue
		}
		log.Info("memload Changed: max:{} to:{}, rate:{} to:{}", memMax, tmpMemMax, memRate, tmpMemRate)
		if stopChannel != nil {
			stopChannel <- true
			close(stopChannel)
		}
		memMax = tmpMemMax
		memRate = tmpMemRate
		if memMax == 0 || memRate == time.Duration(0) {
			stopChannel = nil
			continue
		}
		stopChannel = make(chan bool)
		go MemLoader(stopChannel, memMax, memRate)
	}
}

func MemLoader(stop chan bool, maxMem uint64, rate time.Duration) {
	sleepTime := time.Second
	cycles := uint64(rate.Milliseconds() / sleepTime.Milliseconds())
	memPerCycle := maxMem / cycles
	buffer := make([][]byte, 0)
	currentCycle := uint64(0)
	for {
		select {
		case <-stop:
			buffer = nil
			return
		case <-time.After(sleepTime):
			if currentCycle == cycles {
				buffer = make([][]byte, 0)
				log.Debug("Flushing Memory Buffers")
				runtime.GC()
				runtime.GC()
				runtime.GC()
				currentCycle = 0
				continue
			}
			log.Debug("Adding {} bytes", memPerCycle)
			b := make([]byte, memPerCycle)
			rand.Read(b)
			buffer = append(buffer, b)
			currentCycle++
		}
	}
}
