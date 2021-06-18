package database

import (
	"context"
	"fmt"
)

// Queue is abstracting the logic of a FIFO queue
type Queue interface {
	Push(listName string, data []byte) error
	BPop(listName string) ([]string, error)
}

func (r *redisClient) Push(listName string, data []byte) error {
	result := r.rdb.RPush(context.Background(), listName, data)
	if result.Err() != nil {
		return fmt.Errorf("error adding %v to list %s: %w", string(data), listName, result.Err())
	}

	return nil
}

func (r *redisClient) BPop(listName string) ([]string, error) {
	result := r.rdb.BLPop(context.Background(), 0, listName) // timeout of zero to block indefinitely
	if result.Err() != nil {
		return nil, fmt.Errorf("error getting value from list %s: %w", listName, result.Err())
	}

	value := result.Val()
	if len(value) <= 1 {
		return nil, fmt.Errorf("error reading value for list %s: value empty", listName)
	}

	// returns values[1:] to do not contain the listName
	return value[1:], nil
}
