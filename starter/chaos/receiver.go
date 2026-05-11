// Chaos receiver — the unreliable downstream your service must forward events to.
//
// Listens on localhost:9000 and exposes POST /events. Behavior is controlled by flags:
//
//	--fail-rate=0.10     fraction of requests that return HTTP 500
//	--slow-rate=0.10     fraction of requests that sleep 2–5s before responding 200
//	--timeout-rate=0.05  fraction of requests that never respond (connection held open)
//	--seed=0             RNG seed; non-zero gives deterministic, reproducible chaos
//	--addr=:9000         listen address
//
// Defaults add up to ~25% disrupted requests, which is enough that "ignore errors and
// move on" will lose events. Use --seed to reproduce a specific failure pattern while
// debugging.
//
// Run:  go run ./chaos/receiver.go
//       go run ./chaos/receiver.go --fail-rate=0.3 --seed=42
//
// On success, returns 200 with body: {"received": "<event-id>"}
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type event struct {
	ID string `json:"id"`
}

func main() {
	failRate := flag.Float64("fail-rate", 0.10, "fraction of requests that return HTTP 500")
	slowRate := flag.Float64("slow-rate", 0.10, "fraction of requests that respond after a 2-5s delay")
	timeoutRate := flag.Float64("timeout-rate", 0.05, "fraction of requests that never respond")
	seed := flag.Int64("seed", 0, "RNG seed; 0 = time-based")
	addr := flag.String("addr", ":9000", "listen address")
	flag.Parse()

	if *failRate+*slowRate+*timeoutRate > 1.0 {
		log.Fatalf("fail-rate + slow-rate + timeout-rate must be <= 1.0")
	}

	src := *seed
	if src == 0 {
		src = time.Now().UnixNano()
	}
	var mu sync.Mutex
	rng := rand.New(rand.NewSource(src))
	roll := func() float64 {
		mu.Lock()
		defer mu.Unlock()
		return rng.Float64()
	}
	slowSleep := func() time.Duration {
		mu.Lock()
		defer mu.Unlock()
		return 2*time.Second + time.Duration(rng.Int63n(int64(3*time.Second)))
	}

	log.Printf("chaos receiver listening on %s (seed=%d, fail=%.2f, slow=%.2f, timeout=%.2f)",
		*addr, src, *failRate, *slowRate, *timeoutRate)

	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var e event
		if err := json.Unmarshal(body, &e); err != nil || e.ID == "" {
			http.Error(w, "invalid event", http.StatusBadRequest)
			return
		}

		x := roll()
		switch {
		case x < *timeoutRate:
			log.Printf("[timeout]  id=%s — holding connection", e.ID)
			<-r.Context().Done()
			return
		case x < *timeoutRate+*failRate:
			log.Printf("[fail-500] id=%s", e.ID)
			http.Error(w, "downstream is having a bad day", http.StatusInternalServerError)
			return
		case x < *timeoutRate+*failRate+*slowRate:
			d := slowSleep()
			log.Printf("[slow]     id=%s sleep=%s", e.ID, d)
			select {
			case <-time.After(d):
			case <-r.Context().Done():
				return
			}
		default:
			log.Printf("[ok]       id=%s", e.ID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"received": e.ID})
	})

	srv := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
