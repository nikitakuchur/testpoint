package reporter

import (
	"sync"
	"testpoint/internal/comparator"
)

// Reporter is an interface that can report mismatches in different ways depending on the implementation.
// For example, we might want to log mismatches, or write them to a CSV file.
type Reporter interface {
	Report(input <-chan comparator.RespDiff)
}

// GenerateReport fans out the data from the input channel to the given reporters.
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
			reporter.Report(chans[i])
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
