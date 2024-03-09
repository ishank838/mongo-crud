package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMongoErrCollectionAlreadyInitialised = errors.New("collection already  initialised")
	ErrMongoErrCollectionNotInitialised     = errors.New("collection not  initialised")
)

type MongoStore interface {
	InitCollection(col string) error
	Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, collection string, document []interface{}) (*mongo.InsertManyResult, error)
	UpdateMany(ctx context.Context, collection string, update interface{}, filter interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)
	GetMany(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

var mongoInstance MongoStore

type mongoStore struct {
	db         *mongo.Database
	colections map[string]*mongo.Collection
}

func NewMongoStore(configURL string, database string) (MongoStore, error) {
	if mongoInstance != nil {
		return mongoInstance, nil
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://ishank:ishank@cluster0.fe2aokt.mongodb.net/?retryWrites=true&w=majority&compressors=snappy,zlib,zstd"))
	if err != nil {
		panic(err)
	}

	db := client.Database(database)

	mongoInstance := mongoStore{
		db:         db,
		colections: make(map[string]*mongo.Collection),
	}

	return mongoInstance, nil
}

func (m mongoStore) InitCollection(col string) error {
	if ok := m.isCollectionExists(col); ok {
		return ErrMongoErrCollectionNotInitialised
	}
	collection := m.db.Collection(col)
	m.colections[col] = collection
	return nil
}

func (m mongoStore) Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	if ok := m.isCollectionExists(collection); !ok {
		return nil, ErrMongoErrCollectionNotInitialised
	}
	return m.colections[collection].InsertOne(ctx, document)
}

func (m mongoStore) InsertMany(ctx context.Context, collection string, document []interface{}) (*mongo.InsertManyResult, error) {
	if ok := m.isCollectionExists(collection); !ok {
		return nil, ErrMongoErrCollectionNotInitialised
	}
	return m.colections[collection].InsertMany(ctx, document)
}

func (m mongoStore) UpdateMany(ctx context.Context, collection string, update interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	if ok := m.isCollectionExists(collection); !ok {
		return nil, ErrMongoErrCollectionNotInitialised
	}
	return m.colections[collection].UpdateMany(ctx, filter, update)
}

func (m mongoStore) Delete(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	if ok := m.isCollectionExists(collection); !ok {
		return nil, ErrMongoErrCollectionNotInitialised
	}
	return m.colections[collection].DeleteMany(ctx, filter)
}

func (m mongoStore) GetMany(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	if ok := m.isCollectionExists(collection); !ok {
		return nil, ErrMongoErrCollectionNotInitialised
	}
	return m.colections[collection].Find(ctx, filter, opts...)
}

func (m mongoStore) isCollectionExists(collection string) bool {
	if _, ok := m.colections[collection]; ok {
		return true
	}
	return false
}
