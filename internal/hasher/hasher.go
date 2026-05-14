package hasher

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/idleberg/go-hashman/internal/algo"
)

// Result holds the output of a single hash computation.
type Result struct {
	Algorithm algo.Algorithm
	Hash      string
	Duration  time.Duration
	Err       error
}

// HashFile computes all selected algorithms for the given file concurrently.
// Concurrency is bounded by maxWorkers. Each goroutine opens its own file
// descriptor and streams in 1 MB chunks, matching the TypeScript implementation.
func HashFile(filePath string, algorithms []algo.Algorithm, maxWorkers int) []Result {
	info, err := os.Lstat(filePath)
	if err != nil {
		return errorResults(algorithms, err)
	}
	if !info.Mode().IsRegular() {
		return errorResults(algorithms, fmt.Errorf("%s is not a regular file", filePath))
	}

	results := make([]Result, len(algorithms))
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for i, a := range algorithms {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, alg algo.Algorithm) {
			defer wg.Done()
			defer func() { <-sem }()

			start := time.Now()
			h := alg.NewHash()

			f, err := os.Open(filePath)
			if err != nil {
				results[idx] = Result{Algorithm: alg, Err: err}
				return
			}
			defer f.Close()

			buf := make([]byte, 1024*1024)
			if _, err := io.CopyBuffer(h, f, buf); err != nil {
				results[idx] = Result{Algorithm: alg, Err: err}
				return
			}

			results[idx] = Result{
				Algorithm: alg,
				Hash:      hex.EncodeToString(h.Sum(nil)),
				Duration:  time.Since(start),
			}
		}(i, a)
	}

	wg.Wait()
	return results
}

func errorResults(algorithms []algo.Algorithm, err error) []Result {
	results := make([]Result, len(algorithms))
	for i, a := range algorithms {
		results[i] = Result{Algorithm: a, Err: err}
	}
	return results
}
