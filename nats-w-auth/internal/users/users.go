// Manages users within the system
package users

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const USER_BUCKET = "user"
const OAUTH_BUCKET = "user_oauth"
const PROVIDER_SEPARATOR = "."

type UserManager struct {
	js jetstream.JetStream
	nc *nats.Conn
}

func New(clientURL string) *UserManager {
	nc, err := nats.Connect(clientURL)
	if err != nil {
		log.Fatal("Failed to create nats client", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Failed to create jetstream", err)
	}

	um := &UserManager{
		js: js,
		nc: nc,
	}
	um.setupBuckets()

	return um
}

type UserModel struct {
	ID   string
	Name string
}

type OAuthModel struct {
	Provider string
	// UserID specific to the provider
	UserID string
	// Normalized email address
	Email string
	// Latest access token
	Token string
	// Refresh token if available. Used to refresh access token
	RefreshToken string
}

// Returns a user associated with the provider uid or nil
func (um *UserManager) GetOAuthUser(ctx context.Context, provider, puid string) (*UserModel, error) {
	// TODO: Validate provider and puid
	kv, err := um.js.KeyValue(ctx, OAUTH_BUCKET)
	if err != nil {
		return nil, err
	}

	// Grab user_id mapping from provider, if any
	val, err := kv.Get(ctx, provider+PROVIDER_SEPARATOR+puid)
	if err != nil {
		return nil, err
	}
	userID := string(val.Value())

	user, err := um.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (um *UserManager) GetUser(ctx context.Context, userID string) (*UserModel, error) {
	kv, err := um.js.KeyValue(ctx, USER_BUCKET)
	if err != nil {
		return nil, err
	}

	ve, err := kv.Get(ctx, userID)
	// convert value into usermodel
	user := &UserModel{}
	err = json.Unmarshal(ve.Value(), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (um *UserManager) setupBuckets() {
	ctx := context.Background()

	_, err := um.js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  USER_BUCKET,
		Storage: jetstream.FileStorage,
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = um.js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  OAUTH_BUCKET,
		Storage: jetstream.FileStorage,
	})

	if err != nil {
		log.Fatal(err)
	}
}
