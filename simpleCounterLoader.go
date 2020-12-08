package main // import "github.com/lwahlmeier/goloader"

import (
	"math"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var simpleCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "goloader_simple_total",
	Help: "The number of total",
})

func SimpleCounterWatcher() {
	increaseMin := float64(0)
	increaseMax := float64(0)
	rate := time.Duration(0)
	var stopChannel chan bool = nil
	for {
		time.Sleep(time.Second * 5)
		tmpIncreaseMin := config.GetFloat64("simpleCounter.increaseMin")
		tmpIncreaseMax := config.GetFloat64("simpleCounter.increaseMax")
		tmpRate := config.GetDuration("simpleCounter.rate")
		log.Debug("SimpleCounter:  increaseMin:{} to:{}, increaseMax:{} to:{} , rate:{} to:{}", increaseMin, tmpIncreaseMin, increaseMax, tmpIncreaseMax, rate, tmpRate)
		if tmpIncreaseMin == increaseMin && tmpIncreaseMax == increaseMax && tmpRate == rate {
			continue
		}
		log.Info("SimpleCounter Changed: increaseMin:{} to:{}, increaseMax:{} to:{} , rate:{} to:{}", increaseMin, tmpIncreaseMin, increaseMax, tmpIncreaseMax, rate, tmpRate)
		if stopChannel != nil {
			stopChannel <- true
			close(stopChannel)
		}
		increaseMin = tmpIncreaseMin
		increaseMax = tmpIncreaseMax
		rate = tmpRate
		if increaseMax == 0 || rate == time.Duration(0) {
			stopChannel = nil
			continue
		}
		if increaseMin < 0 {
			log.Warn("increaseMin must be 0 or greater, was:{} setting to 0", increaseMin)
			increaseMin = 0
		}
		stopChannel = make(chan bool)
		go SimpleCounterLoader(stopChannel, increaseMin, increaseMax, rate)
	}
}

func SimpleCounterLoader(stopChannel chan bool, increaseMin float64, increaseMax float64, rate time.Duration) {
	rand.Seed(time.Now().Unix())

	for {
		select {
		case <-stopChannel:
			return
		case <-time.After(rate):
			incr := math.Floor(increaseMin + rand.Float64()*(increaseMax-increaseMin))
			log.Debug("Increasing counter by:{}", incr)
			simpleCounter.Add(incr)
		}
	}

}
