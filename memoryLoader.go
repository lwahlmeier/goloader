package main // import "github.com/lwahlmeier/goloader"

import (
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
		log.Debug("Memload:  max:{} to:{}, rate:{} to:{}", memMax, tmpMemMax, memRate, tmpMemRate)
		if tmpMemMax == memMax && tmpMemRate == memRate {
			continue
		}
		log.Info("MemLoad Changed: max:{} to:{}, rate:{} to:{}", memMax, tmpMemMax, memRate, tmpMemRate)
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
	log.Info("MemLoader: perCycle:{}", memPerCycle)
	currentCycle := uint64(0)
	for {
		select {
		case <-stop:
			buffer = nil
			log.Info("MemLoader: Stopping Memory processing")
			return
		case <-time.After(sleepTime):
			if currentCycle == cycles {
				buffer = nil
				buffer = make([][]byte, 0)
				log.Info("MemLoader: Ending Memory grow Cycle, Flushing Memory Buffers")
				runtime.GC()
				runtime.GC()
				runtime.GC()
				currentCycle = 0
				continue
			}
			log.Debug("MemLoader: Adding {} bytes", memPerCycle)
			b := make([]byte, memPerCycle)
			buffer = append(buffer, b)
			currentCycle++
		}
	}
}
