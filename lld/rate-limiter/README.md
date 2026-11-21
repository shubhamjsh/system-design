# Token Bucket Rate Limiter

This package provides a small, thread-safe, in-process token bucket. The bucket
starts full, accepts requests while tokens are available, and replenishes one
token per configured interval up to its capacity.

```go
limiter := tokenbucket.New(5, time.Second) // capacity 5; one token per second

if limiter.Allow() {
	// Handle the request.
} else {
	// Reject or defer the request.
}
```

Import `system-design/lld/rate-limiter/token-bucket` as `tokenbucket`. A
non-positive capacity or refill interval creates a limiter that rejects all
requests.

## Sliding Window

The `sliding-window` package allows up to a configured number of requests in
the trailing time window. It records each allowed request, so the limit is
exact rather than approximate.

```go
limiter := slidingwindow.New(50, time.Second) // at most 50 requests per second

if limiter.Allow() {
	// Handle the request.
} else {
	// Reject or defer the request.
}
```

Import `system-design/lld/rate-limiter/sliding-window` as `slidingwindow`. A
request exactly one window old no longer counts toward the limit. A
non-positive limit or window creates a limiter that rejects all requests.
