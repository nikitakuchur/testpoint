package reporter

import (
	"sync"
	"testpoint/internal/comparator"
)

type Reporter interface {
	report(input <-chan comparator.RespDiff)
}

func GenerateReport(input <-chan comparator.RespDiff, reporters ...Reporter) {
	var chans []chan comparator.RespDiff
	for i := 0; i < len(reporters); i++ {
		chans = append(chans, make(chan comparator.RespDiff))
	}

	wg := sync.WaitGroup{}
	for i, reporter := range reporters {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reporter.report(chans[i])
		}()
	}

	for diff := range input {
		for _, c := range chans {
			c <- diff
		}
	}
	for _, c := range chans {
		close(c)
	}

	wg.Wait()
}
