package mobile_openfga

import (
	"context"
	"database/sql"
	"log"
	"net/url"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/openfga/openfga/assets"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/storage/sqlcommon"
	"github.com/openfga/openfga/pkg/storage/sqlite"
	"github.com/pressly/goose/v3"
)

var serverInstance *server.Server

var ctx = context.Background()

type ctxKey string

func InitServer(dbPath string) error {
	var datastore storage.OpenFGADatastore
	var err error

	ctx = context.WithValue(ctx, ctxKey("request-validated"), true)

	datastoreOptions := []sqlcommon.DatastoreOption{}

	dsCfg := sqlcommon.NewConfig(datastoreOptions...)

	datastore, err = sqlite.New(dbPath, dsCfg)

	if err != nil {
		println("error: ", err)
		return err
	}

	serverInstance = server.MustNewServerWithOpts(
		server.WithDatastore(datastore),
		server.WithExperimentals(server.ExperimentalFeatureFlag("enable-list-users")),
	)

	return nil
}

func MigrateDatabase(dbPath string) error {
	var uri, driver, dialect, migrationsPath string

	driver = "sqlite3"
	dialect = "sqlite3"
	migrationsPath = assets.SQLiteMigrationDir

	if uri == "" {
		uri = dbPath
	}

	// Parse the database uri with the sqlite drivers function for it and update username/password, if set via flags
	dbURI, err := url.Parse(uri)

	if err != nil {
		return err
	}

	uri = dbURI.String()

	db, err := sql.Open(driver, uri)

	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close the datastore: %v", err)
		}
	}()

	if err != nil {
		return err
	}

	// TODO use goose.OpenDBWithDriver which already sets the dialect
	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	goose.SetBaseFS(assets.EmbedMigrations)

	currentVersion, err := goose.GetDBVersion(db)

	if err != nil {
		return err
	}

	log.Printf("current version %d", currentVersion)

	if err := goose.Up(db, migrationsPath); err != nil {
		return err
	}

	return nil
}

func ReadAuthorizationModels(
	encodedReadAuthorizationModelsRequest []byte,
) ([]byte, error) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.ReadAuthorizationModelsRequest{}

	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	unmarshalErr := unmarshaler.Unmarshal(encodedReadAuthorizationModelsRequest, &req)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	resp, err := serverInstance.ReadAuthorizationModels(ctx, &req)

	if err != nil {
		return nil, err
	}

	jsonResult, err := protojson.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

func WriteAuthorizationModel(
	encodedWriteAuthorizationModelRequest []byte,
) ([]byte, error) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.WriteAuthorizationModelRequest{}

	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	unmarshalErr := unmarshaler.Unmarshal(encodedWriteAuthorizationModelRequest, &req)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	resp, err := serverInstance.WriteAuthorizationModel(ctx, &req)

	if err != nil {
		return nil, err
	}

	jsonResult, err := protojson.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

func ListStores(
	encodedListStoreRequest []byte,
) ([]byte, error) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.ListStoresRequest{}

	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	unmarshalErr := unmarshaler.Unmarshal(encodedListStoreRequest, &req)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	resp, err := serverInstance.ListStores(ctx, &req)

	if err != nil {
		return nil, err
	}

	jsonResult, err := protojson.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

func CreateStore(storeName string) ([]byte, error) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.CreateStoreRequest{
		Name: storeName,
	}

	resp, err := serverInstance.CreateStore(ctx, &req)

	if err != nil {
		return nil, err
	}

	jsonResult, err := protojson.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

func Write(
	encodedWriteRequest []byte,
) error {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.WriteRequest{}

	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	unmarshalErr := unmarshaler.Unmarshal(encodedWriteRequest, &req)

	if unmarshalErr != nil {
		return unmarshalErr
	}

	_, err := serverInstance.Write(ctx, &req)

	if err != nil {
		return err
	}

	return nil
}

func Check(
	encodedCheckRequest []byte,
) ([]byte, error) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.CheckRequest{}

	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	unmarshalErr := unmarshaler.Unmarshal(encodedCheckRequest, &req)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	resp, err := serverInstance.Check(ctx, &req)

	if err != nil {
		return nil, err
	}

	jsonResult, err := protojson.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

func ListObjects(
	encodedListObjectsRequest []byte,
) ([]byte, error) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.ListObjectsRequest{}

	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	unmarshalErr := unmarshaler.Unmarshal(encodedListObjectsRequest, &req)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	resp, err := serverInstance.ListObjects(ctx, &req)

	if err != nil {
		return nil, err
	}

	jsonResult, err := protojson.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

func ListUsers(
	encodedListUsersRequest []byte,
) ([]byte, error) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.ListUsersRequest{}

	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	unmarshalErr := unmarshaler.Unmarshal(encodedListUsersRequest, &req)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	resp, err := serverInstance.ListUsers(ctx, &req)

	if err != nil {
		return nil, err
	}

	jsonResult, err := protojson.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}
