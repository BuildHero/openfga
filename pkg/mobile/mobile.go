package mobile

import (
	"context"
	"database/sql"
	"log"
	"net/url"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"

	"github.com/openfga/openfga/assets"
	"github.com/openfga/openfga/pkg/server"
	"github.com/openfga/openfga/pkg/storage"
	"github.com/openfga/openfga/pkg/storage/sqlcommon"
	"github.com/openfga/openfga/pkg/storage/sqlite"
	"github.com/pressly/goose/v3"
)

var serverInstance *server.Server

var ctx = context.Background()

func InitServer(dbPath string) {
	var datastore storage.OpenFGADatastore
	var err error

	datastoreOptions := []sqlcommon.DatastoreOption{}

	dsCfg := sqlcommon.NewConfig(datastoreOptions...)

	datastore, err = sqlite.New(dbPath, dsCfg)

	if err != nil {
		println("error: ", err)
		panic(err)
	}

	serverInstance = server.MustNewServerWithOpts(
		server.WithDatastore(datastore),
	)

	println("serverInstance: ", serverInstance)
}

func CreateStore(storeName string) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.CreateStoreRequest{
		Name: storeName,
	}

	resp, err := serverInstance.CreateStore(ctx, &req)

	if err != nil {
		println("CreateStoreRequest error: ", err.Error())
	}

	println("CreateStoreResponse: ", resp)
}

func ListObjects(storeId string,
	authorizationModelId string,
	objectType string,
	relation string,
	user string,
) {
	if serverInstance == nil {
		log.Fatalf("server instance is nil")
	}

	req := openfgav1.ListObjectsRequest{
		StoreId:              storeId,
		AuthorizationModelId: authorizationModelId,
		Type:                 objectType,
		Relation:             relation,
		User:                 user,
	}

	println("ListObjectsRequest: ", req.User)

	resp, err := serverInstance.ListObjects(ctx, &req)

	if err != nil {
		println("ListObjectsRequest error: ", err.Error())
	}

	println("ListObjectsResponse: ", resp)
}

func MigrateDatabase(dbPath string) {
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
		log.Fatalf("invalid database uri: %v\n", err)
	}

	uri = dbURI.String()

	db, err := sql.Open(driver, uri)
	if err != nil {
		log.Fatalf("failed to open a connection to the datastore: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close the datastore: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}

	// TODO use goose.OpenDBWithDriver which already sets the dialect
	if err := goose.SetDialect(dialect); err != nil {
		log.Fatalf("failed to initialize the migrate command: %v", err)
	}

	goose.SetBaseFS(assets.EmbedMigrations)

	currentVersion, err := goose.GetDBVersion(db)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("current version %d", currentVersion)

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatal(err)
	}
}
