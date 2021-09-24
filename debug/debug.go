package debug

import (
	"fmt"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func TraceMemStats() {
	go func() {
		for {
			var ms runtime.MemStats
			tick := time.After(time.Second)
			select {
			case <-tick:
				runtime.ReadMemStats(&ms)
				mb := 1024 * 1024.0
				logstr := fmt.Sprintf("Alloc=%.2fMB  TotalAlloc=%.2fMB  Sys=%.2fMB  NumGC=%v",
					float64(ms.Alloc)/mb, float64(ms.TotalAlloc)/mb, float64(ms.Sys)/mb, ms.NumGC)
				log.Println(logstr)
			}
		}
	}()
}
