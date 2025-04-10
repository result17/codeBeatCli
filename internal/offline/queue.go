package offline

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/result17/codeBeatCli/internal/heartbeat"
	bolt "go.etcd.io/bbolt"
)

const (
	dbBucket = "heartbeats"
)

type Queue struct {
	Bucket string
	tx     *bolt.Tx
}

// NewQueue creates a new instance of Queue.
func NewQueue(tx *bolt.Tx) *Queue {
	return &Queue{
		Bucket: dbBucket,
		tx:     tx,
	}
}

func (q *Queue) checkBucketExist() (*bolt.Bucket, error) {
	bucket := q.tx.Bucket([]byte(q.Bucket))
	if bucket == nil {
		return nil, errors.New("failed to load bucket")
	}
	return bucket, nil
}

func (q *Queue) checkBucketExistIfNotCreate() (*bolt.Bucket, error) {
	bucket, _ := q.checkBucketExist()
	if bucket != nil {
		bucket, err := q.tx.CreateBucket([]byte(q.Bucket))
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %s", err)
		}
		return bucket, nil
	}
	return bucket, nil
}

func (q *Queue) Count() (int, error) {
	bucket, error := q.checkBucketExistIfNotCreate()
	if bucket == nil {
		return 0, error
	}
	return bucket.Inspect().KeyN, nil
}

func (q *Queue) PopMany(limit int) ([]heartbeat.Heartbeat, error) {
	bucket, error := q.checkBucketExistIfNotCreate()
	if bucket == nil {
		return nil, error
	}

	var (
		heartbeats []heartbeat.Heartbeat
		ids        []string
	)

	c := bucket.Cursor()

	for key, value := c.First(); key != nil; key, value = c.Next() {
		if len(heartbeats) >= limit {
			break
		}
		var h heartbeat.Heartbeat
		err := json.Unmarshal(value, &h)

		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal heartbeat data: %s", err)
		}
		heartbeats = append(heartbeats, h)
		ids = append(ids, string(key))
	}
	for _, id := range ids {
		if err := bucket.Delete([]byte(id)); err != nil {
			return nil, fmt.Errorf("failed to delete key %q: %s", id, err)
		}
	}
	return heartbeats, nil
}

func (q *Queue) PushMany(hs []heartbeat.Heartbeat) error {
	bucket, error := q.checkBucketExistIfNotCreate()
	if bucket == nil {
		return error
	}

	for _, h := range hs {
		data, err := json.Marshal(h)

		if err != nil {
			return fmt.Errorf("failed to json marshal heartbeat: %s", err)
		}

		err = bucket.Put([]byte(h.ID()), data)
		if err != nil {
			return fmt.Errorf("failed to store heartbeat with id %q: %s", h.ID(), err)
		}
	}

	return nil
}

func (q *Queue) ReadMany(limit int) ([]heartbeat.Heartbeat, error) {
	bucket, error := q.checkBucketExistIfNotCreate()
	if bucket == nil {
		return nil, error
	}

	var heartbeats = make([]heartbeat.Heartbeat, 0)

	// load values
	c := bucket.Cursor()
	for key, value := c.First(); key != nil; key, value = c.Next() {
		if len(heartbeats) >= limit {
			break
		}
		var h heartbeat.Heartbeat
		err := json.Unmarshal(value, &h)

		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal heartbeat data: %s", err)
		}
		heartbeats = append(heartbeats, h)
	}
	return heartbeats, nil
}
