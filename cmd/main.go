package main

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/santosh/gingo/docs"
	"github.com/sirupsen/logrus"

	route "app/api/route"
	"app/config"
	"app/db"
	"app/db/migration"
)

// @title           Auth API
// @version         1.0
// @description     Test of Aprel
// @host            localhost:8080
// @BasePath        /v1
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name X-Token
func main() {
	setupLogrus()

	env := config.NewEnv()

	connection := connectToDB(env)

	DBName := env.DBDriver
	runDBMigration(connection, DBName, env.MigrationURL)

	store := db.NewStore(connection)

	if env.AppEnv == "development" {
		// This is necessary to fill in the users table with hashed passwords
		migration.FillUsers(store)
	}

	runGinServer(env, store)
}

func setupLogrus() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s:%d", filename, f.Line)
			return "", filename
		},
	})
}

func connectToDB(env *config.Env) *sql.DB {
	var counts int

	for {
		connection, err := openDB(env.DBDriver, env.DBSource)
		if err != nil {
			logrus.Print("postgres not yet ready ...")
			counts++
		} else {
			logrus.Print("connected to Postgres!")
			return connection
		}

		if counts > env.MaxDBConnectionAttempts {
			logrus.Fatalf("failed to connect to Postgresql: %v", err)
		}

		logrus.Print("backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dbDriver, dbSource string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runDBMigration(db *sql.DB, DBname, migrationURL string) {
	driver, _ := postgres.WithInstance(db, &postgres.Config{})

	migration, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		DBname, // "postgres"
		driver)
	if err != nil {
		logrus.Fatalf("cannot create new migrate instance: %v", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatalf("failed to run migrate up: %v", err)
	}

	logrus.Printf("db migrated successfully")
}

func runGinServer(env *config.Env, store db.Store) {
	ginEngine := gin.Default()

	// CORS middleware.
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	ginEngine.Use(cors.New(config))

	connectSwaggerToGin(ginEngine)

	router := ginEngine.Group("/v1/")
	route.Setup(env, store, router)

	logrus.Printf("server running on address: %s", env.ServerAddress)
	ginEngine.Run(env.ServerAddress)
}

func connectSwaggerToGin(ginEngine *gin.Engine) {
	// Serve the Swagger UI files.
	ginEngine.Static("/swagger/", "./doc")
}
