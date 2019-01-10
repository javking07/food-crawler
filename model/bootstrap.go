package model

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/javking07/food-crawler/conf"
	"github.com/rs/zerolog/log"
)

func BootstrapCassandra(config *conf.DatabaseConfig) (*gocql.Session, error) {
	log.Info().Msg("using cassandra")
	cluster := gocql.NewCluster(config.Host)
	cluster.Port = config.Port
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.Keyspace = config.DatabaseName
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.User,
		Password: config.Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	session.SetTrace(gocql.NewTraceWriter(session, os.Stdout))
	return session, nil
}

func BootstrapPostgress(config *conf.DatabaseConfig) (*sql.DB, error) {
	var connectString string
	if password := os.Getenv("APP_DB_PASSWORD"); password != "" {
		connectString = fmt.Sprintf(
			"host=%s port=%v user=%s "+
				"password=%s dbname=%s sslmode=%s",
			config.Host,
			config.Port,
			config.User,
			password,
			config.DatabaseName,
			config.SslMode)
	} else {
		connectString = fmt.Sprintf(
			"host=%s port=%v user=%s "+
				"password=%s dbname=%s sslmode=%s",
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.DatabaseName,
			config.SslMode)
	}

	var err error
	db, err := sql.Open("postgres", connectString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
