package connect

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net"
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
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Fatal(err)
	}

	return session
}

func ConnectSecure() *mgo.Session {
	cas_file := "/etc/pki/certs/private_cas_chain.crt"

	roots := x509.NewCertPool()

	if ca, err := ioutil.ReadFile(cas_file); err == nil {
		roots.AppendCertsFromPEM(ca)
	}

	tlsConfig := &tls.Config{}
	tlsConfig.RootCAs = roots
	//////////////////////////////////////////////////////////

	servers := strings.Split(viper.GetString("servers"), ",")
	db := viper.GetString("db")
	rs := viper.GetString("rs")
	dbuser := viper.GetString("dbuser")
	password := viper.GetString("password")

	//////////////////////////////////////////////////////////

	dialInfo := &mgo.DialInfo{
		Addrs:          servers,
		Timeout:        30 * time.Second,
		Database:       db,
		ReplicaSetName: rs,
		Username:       dbuser,
		Password:       password,

		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), tlsConfig)
		},
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Fatal(err)
	}

	return session
}
