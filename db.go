package main

import (
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
)

var dbConf = &mgo.DialInfo{
	Addrs:    []string{os.Getenv("MONGO_HOST")},
	Timeout:  60 * time.Second,
	Database: os.Getenv("MONGO_NAME"),
	Username: os.Getenv("MONGO_USER"),
	Password: os.Getenv("MONGO_PASS"),
}

type dbSession struct {
	session *mgo.Session
}

func mgoSession() (*dbSession, error) {
	session, err := mgo.DialWithInfo(dbConf)
	if err != nil {
		return nil, err
	}

	return &dbSession{session: session}, nil
}

func (s *dbSession) Close() {
	s.session.Close()
}
