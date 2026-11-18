package webmoney

import (
	"regexp"
	"strconv"
	"sync"
	"testing"
)

// Регрессия: Reqn строился только из времени без синхронизации — вызовы в
// одной миллисекунде (или при откате часов) давали одинаковые номера, а WM
// требует строго возрастающий reqn на WMID.

func TestReqnSequentialStrictlyIncreasing(t *testing.T) {
	var prev int64
	for i := 0; i < 2000; i++ {
		s := Reqn()
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			t.Fatalf("Reqn %q is not numeric: %v", s, err)
		}
		if i > 0 && v <= prev {
			t.Fatalf("Reqn not strictly increasing: %d after %d (call %d)", v, prev, i)
		}
		prev = v
	}
}

func TestReqnConcurrentUnique(t *testing.T) {
	t.Run("modern", func(t *testing.T) {
		t.Setenv("USE_OLD_REQN", "")
		checkConcurrentReqn(t, 17)
	})
	t.Run("legacy", func(t *testing.T) {
		t.Setenv("USE_OLD_REQN", "1")
		checkConcurrentReqn(t, 14)
	})
}

func checkConcurrentReqn(t *testing.T, digits int) {
	t.Helper()
	const workers, perWorker = 50, 200
	results := make(chan string, workers*perWorker)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				results <- Reqn()
			}
		}()
	}
	wg.Wait()
	close(results)

	pattern := regexp.MustCompile(`^\d{` + strconv.Itoa(digits) + `}$`)
	seen := make(map[string]bool, workers*perWorker)
	for s := range results {
		if !pattern.MatchString(s) {
			t.Fatalf("Reqn %q does not match %d-digit format", s, digits)
		}
		if seen[s] {
			t.Fatalf("duplicate Reqn under concurrency: %s", s)
		}
		seen[s] = true
	}
	if len(seen) != workers*perWorker {
		t.Fatalf("got %d unique reqn, want %d", len(seen), workers*perWorker)
	}
}
