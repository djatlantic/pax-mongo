package connect

import (
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"strings"
	"time"
)

func Connect() *mgo.Session {
	servers := strings.Split(viper.GetString("servers"), ",")
	db := viper.GetString("db")
	rs := viper.GetString("rs")
	dbuser := viper.GetString("dbuser")
	password := viper.GetString("password")

	dialInfo := &mgo.DialInfo{
		Addrs:          servers,
		Timeout:        30 * time.Second,
		Database:       db,
		ReplicaSetName: rs,
		Username:       dbuser,
		Password:       password,

		//DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
		//      return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Fatal(err)
	}

	return session
}
