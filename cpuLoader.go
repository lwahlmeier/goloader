package main // import "github.com/lwahlmeier/goloader"

import (
	"time"
)

func CpuWatcher() {
	cpuPct := float64(0)
	cpuDelay := time.Duration(0)
	cpuThreads := uint64(0)
	stopChannels := make([]chan bool, 0)
	for {
		time.Sleep(time.Second * 5)
		tmpCpuPct := config.GetFloat64("cpu.pct")
		tmpCpuDelay := config.GetDuration("cpu.delay")
		tmpCpuThreads := config.GetUint64("cpu.threads")
		log.Debug("cpuLoad:  pct:{} to:{}, delay:{} to:{}, threads:{} to:{}", cpuPct, tmpCpuPct, cpuDelay, tmpCpuDelay, cpuThreads, tmpCpuThreads)
		if tmpCpuPct == cpuPct && tmpCpuDelay == cpuDelay && tmpCpuThreads == cpuThreads {
			continue
		}
		log.Info("CpuLoad Change:  pct:{} to:{}, delay:{} to:{}, threads:{} to:{}", cpuPct, tmpCpuPct, cpuDelay, tmpCpuDelay, cpuThreads, tmpCpuThreads)
		for _, v := range stopChannels {
			v <- true
		}
		cpuDelay = tmpCpuDelay
		cpuPct = tmpCpuPct
		cpuThreads = tmpCpuThreads
		stopChannels = make([]chan bool, cpuThreads)
		for i := uint64(0); i < cpuThreads; i++ {
			tmpChan := make(chan bool)
			stopChannels[i] = tmpChan
			go CpuLoader(tmpChan, cpuPct, cpuDelay)
		}
	}
}

func CpuLoader(stop chan bool, pctLoad float64, delay time.Duration) {
	loopSize := int64(1000)
	lastStartTime := time.Now()
	targetTime := float64(delay.Nanoseconds()) * pctLoad
	sleepTime := time.Duration(delay.Nanoseconds() - int64(targetTime))
	for {
		select {
		case <-stop:
			return
		case <-time.After(sleepTime):
			q := float64(100)
			lastStartTime = time.Now()
			for i := int64(0); i < loopSize; i++ {
				q++
			}
			pt := time.Since(lastStartTime)
			processTime := float64(pt.Nanoseconds())
			// log.Info("Current loops to:{}", loopSize)
			// log.Info("ProcessTime was:{}, target:{}", pt, time.Duration(int64(targetTime)))i
			pct := (targetTime / processTime) / 10
			if processTime > targetTime {
				loopSize = loopSize - int64(float64(loopSize)*pct)
				// log.Info("Decreasing loops to:{}, decreasing by:{}%", loopSize, fmt.Sprintf("%.2f", pct))
			} else if processTime < targetTime {
				loopSize = loopSize + int64(float64(loopSize)*pct)
				// log.Info("Increasing loops to:{}, pct Increase:{}%", loopSize, fmt.Sprintf("%.2f", pct))
			}
		}
	}
}
