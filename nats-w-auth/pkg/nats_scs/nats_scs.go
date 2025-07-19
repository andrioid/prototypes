// Adapted from https://gist.github.com/mrguamos/2640c2bbbb4bb4d5b73ba7d816734759
package nats_scs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

const DEFAULT_TIMEOUT = 10 * time.Second
const BUCKET_NAME = "user_session"

// Implements scs.Store interface
// TODO: Use StoreCtx instead: https://github.com/alexedwards/scs?tab=readme-ov-file#using-custom-session-stores-with-contextcontext
type NatsStore struct {
	js jetstream.JetStream
	kv jetstream.KeyValue
}

func New(js jetstream.JetStream) *NatsStore {
	parentCtx := context.Background()

	kv, err := initializeKeyValue(js, parentCtx)
	if err != nil {
		log.Fatalf("Failed to initialize key-value store: %v", err)
	}

	return &NatsStore{
		js: js,
		kv: kv,
	}
}

func initializeKeyValue(js jetstream.JetStream, parentCtx context.Context) (jetstream.KeyValue, error) {
	ctx, cancel := context.WithTimeout(parentCtx, DEFAULT_TIMEOUT)
	defer cancel()

	kv, err := js.KeyValue(ctx, BUCKET_NAME)
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			kv, err = js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
				Bucket:  BUCKET_NAME,
				Storage: jetstream.FileStorage,
			})
			if err != nil {
				return nil, err
			}
			return kv, nil
		}
		return nil, err
	}
	return kv, nil
}

func (s *NatsStore) Delete(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()
	err := s.kv.Delete(ctx, token)
	if err != nil {
		return fmt.Errorf("natsstore error deleting key: %w", err)
	}
	return nil
}

func (s *NatsStore) Find(token string) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()
	entry, err := s.kv.Get(ctx, token)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			fmt.Println("key not found", token)
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("natsstore error getting key: %w", err)
	}
	return entry.Value(), true, nil
}

func (s *NatsStore) Commit(token string, data []byte, expiry time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()
	_, err := s.kv.Put(ctx, token, data)
	if err != nil {
		return err
	}
	fmt.Println("wrote", token, expiry)
	return nil
}

func (s *NatsStore) All() (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	kl, err := s.kv.ListKeys(ctx, jetstream.IgnoreDeletes())
	if err != nil {
		return nil, fmt.Errorf("natsstore error listing keys: %w", err)
	}
	data := make(map[string][]byte)
	for k := range kl.Keys() {
		data[k], _, err = s.Find(k)
		log.Println(fmt.Errorf("natsstore error getting key: %w", err))
	}
	return data, nil
}
