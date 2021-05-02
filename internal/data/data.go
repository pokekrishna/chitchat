package data

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/pokekrishna/chitchat/internal/config"
	"github.com/pokekrishna/chitchat/pkg/log"
	_ "github.com/lib/pq"
	"time"
)

var db *sql.DB
func Initialize() error {
	var err error
	log.Info("Initializing DB ... ")
	postgresConnectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=5",
		config.GetDbUser(),config.GetDbPassword(),config.GetDbHost(), config.GetDbPort(), config.GetDbName())
	db, err = sql.Open("postgres", postgresConnectionURL)
	if err != nil{
		return err
	}

	// Wait until 5 seconds for ping
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return err
	}

	log.Info("Done")
	return nil
}

func Encrypt(plainText string) (encryptedText string) {
	encryptedText = fmt.Sprintf("%x", sha1.Sum([]byte(plainText)))
	return
}