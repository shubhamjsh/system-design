// Merkle Tree-based tamper-evident audit ledger in Go

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type MTEvent struct {
	UserID    string
	Action    string
	Sequence  uint64
	Timestamp string
}

type MTLedger struct {
	Events      []MTEvent
	TrustedRoot string
}

func hashEvent(e MTEvent) []byte {
	data, _ := json.Marshal(e)
	h := sha256.Sum256(append([]byte{0}, data...))
	return h[:]
}

func hashPair(left, right []byte) []byte {
	data := make([]byte, 0, 1+len(left)+len(right))
	data = append(data, 1)
	data = append(data, left...)
	data = append(data, right...)
	h := sha256.Sum256(data)
	return h[:]
}

func buildMT(events []MTEvent) string {
	if len(events) == 0 {
		return ""
	}
	level := make([][]byte, 0, len(events))

	for _, e := range events {
		level = append(level, hashEvent(e))
	}
	/*
	      level
	      ├── level[0] = h1  // event 1 hash
	      ├── level[1] = h2  // event 2 hash
	      ├── level[2] = h3  // event 3 hash
	      └── level[3] = h4  // event 4 hash

	   	next
	      ├── next[0] = p12  // parent of h1 and h2
	      └── next[1] = p34  // parent of h3 and h4

	   	level = next
	      Now level has two hashes: [p12, p34].

	      Second pass
	      next
	      └── next[0] = root
	*/
	for len(level) > 1 {
		next := make([][]byte, 0, (len(level)+1)/2)
		for i := 0; i < len(level); i += 2 {
			if i+1 < len(level) {
				next = append(next, hashPair(level[i], level[i+1]))
			} else {
				next = append(next, level[i])
			}
		}
		level = next
	}
	return hex.EncodeToString(level[0])
}

func (l *MTLedger) append(userID, action string) MTEvent {
	e := MTEvent{
		UserID:    userID,
		Action:    action,
		Sequence:  uint64(len(l.Events) + 1),
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}
	l.Events = append(l.Events, e)
	return e
}

func (l *MTLedger) checkpoint() {
	l.TrustedRoot = buildMT(l.Events)
}

func (l *MTLedger) verify() error {
	actualRoot := buildMT(l.Events)

	if actualRoot != l.TrustedRoot {
		return fmt.Errorf("tampering detected: expected %s, got %s", l.TrustedRoot, actualRoot)
	}
	return nil
}

func main() {
	l := &MTLedger{}

	l.append("user123", "LOGIN")
	l.append("user123", "TRANSFER 10000")
	l.append("user121", "LOGIN")
	l.append("user123", "LOGOUT")

	l.checkpoint()

	fmt.Println("Before tampering:", l.verify())

	l.Events[1].Action = "TRANSFER 100"
	fmt.Println("After tampering:", l.verify())
}
