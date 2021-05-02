package data

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/pokekrishna/chitchat/internal/config"
	"github.com/pokekrishna/chitchat/pkg/log"
	_ "github.com/lib/pq"
	"crypto/rand"
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

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Error("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}