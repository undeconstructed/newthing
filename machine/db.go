package machine

import (
	"encoding/binary"
	"encoding/json"

	"github.com/coreos/bbolt"
)

var eventsBytes = []byte("events")

type db struct {
	file string
	bolt *bolt.DB
}

func newDB(file string) *db {
	return &db{file: file}
}

func (db *db) start() error {
	bolt0, err := bolt.Open(db.file, 0600, nil)
	if err != nil {
		return err
	}
	bolt0.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(eventsBytes)
		return err
	})
	db.bolt = bolt0
	return nil
}

func (db *db) putEvent(event Event) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.Bucket(eventsBytes).CreateBucketIfNotExists([]byte(event.Bucket))
		if err != nil {
			return err
		}
		// XXX timestamp better
		id, _ := bucket.NextSequence()
		value, _ := json.Marshal(event)
		return bucket.Put(itob(id), value)
	})
}

func (db *db) getEvents(bucket string) ([]Event, error) {
	out := []Event{}
	err := db.bolt.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket(eventsBytes).Bucket([]byte(bucket))
		if buck == nil {
			return nil
		}
		buck.ForEach(func(k []byte, v []byte) error {
			event := Event{}
			json.Unmarshal(v, &event)
			out = append(out, event)
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (db *db) stop() {
	db.bolt.Close()
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
