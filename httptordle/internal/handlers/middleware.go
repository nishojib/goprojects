package handlers

import (
	"goprojects/httptordle/internal/problem"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

const ratePerSecond = 2
const limiterBurst = 4

// rateLimit middleware that applies rate limiting on HTTP requests before passing them to the next handler.
func rateLimit(
	w http.ResponseWriter,
	r *http.Request,
	next func(http.ResponseWriter, *http.Request),
) {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	ip := realip.FromRequest(r)

	mu.Lock()

	if _, found := clients[ip]; !found {
		clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(ratePerSecond), limiterBurst)}
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			problem.Of(http.StatusTooManyRequests).
				Append(problem.Detail("rate limit exceeded")).
				WriteTo(w)
			return
		}

		mu.Unlock()
	}

	next(w, r)
}
