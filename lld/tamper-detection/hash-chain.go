package tamper_detection

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	UserID       string
	Action       string
	Sequence     uint64
	PreviousHash string
	Timestamp    string
	Hash         string
}

type ledger struct {
	event []Event
}

func calculateHash(e Event) string {
	payload := struct {
		UserID       string
		Action       string
		Sequence     uint64
		PreviousHash string
		Timestamp    string
	}{UserID: e.UserID, Action: e.Action, Sequence: e.Sequence, PreviousHash: e.PreviousHash, Timestamp: e.Timestamp}

	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func (l *ledger) append(userID string, action string) Event {
	var sequence uint64 = 1
	var previous string
	for i := len(l.event) - 1; i >= 0; i-- {
		e := l.event[i]
		if e.UserID == userID {
			previous = e.Hash
			sequence = e.Sequence + 1
			break
		}
	}
	e := Event{UserID: userID, Action: action, Sequence: sequence, PreviousHash: previous, Timestamp: time.Now().String()}
	e.Hash = calculateHash(e)
	l.event = append(l.event, e)
	return e
}

func (l *ledger) verify(userID, hash string, sequence uint64) error {
	for _, e := range l.event {
		if e.UserID == userID && e.Hash == hash && e.Sequence == sequence {
			return nil
		}
	}
	return fmt.Errorf("Hash mismatch tampering detected")
}

func main() {
	l := &ledger{}
	l.append("user123", "LOGIN")
	last := l.append("user123", "TRANSFER 10,000")
	checkpointSequense, checkpointHash := last.Sequence, last.Hash

	fmt.Println("before tampering", l.verify("user123", checkpointHash, checkpointSequense))
	l.event[1].Sequence = 200
	fmt.Println("after tampering", l.verify("user123", checkpointHash, checkpointSequense))

	l.append("user123", "Logout")

}
