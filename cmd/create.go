// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	//"fmt"

	"github.com/djatlantic/pax-mongo/connect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	//"reflect"
	"strings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create or update a user",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		roles := strings.Split(viper.GetString("userRoles"), ",")
		roleDb := make([]interface{}, len(roles))
		db := viper.GetString("targetDb")

		log.Println("targetDb: ", db)

		for i, v := range roles {
			roleDb[i] = bson.M{"role": v, "db": db}
		}

		if viper.IsSet("userExtraRoles") && viper.IsSet("targetExtraDb") {
			exRoles := strings.Split(viper.GetString("userExtraRoles"), ",")
			exDb := viper.GetString("targetExtraDb")

			for _, r := range exRoles {
				roleDb = append(roleDb, bson.M{"role": r, "db": exDb})
			}
		}

		var session *mgo.Session

		secure := viper.GetBool("tls")
		if secure {
			session = connect.ConnectSecure()
		} else {
			session = connect.Connect()
		}
		defer session.Close()

		log.Println("Creating user: ", viper.GetString("user"))

		result := bson.M{}
		//if err := session.DB(db).Run(bson.M{"createUser": viper.GetString("user"),
		//	"pwd":          viper.GetString("userPassword"),
		//	"roles":        roleDb,
		//	"writeConcern": bson.M{"w": "majority"}}, &result); err != nil {
		//	log.Fatal("Failed to create user: ", err)
		//}

		if err := session.DB(db).Run(bson.D{{"createUser", viper.GetString("user")},
			{"pwd", viper.GetString("userPassword")},
			{"roles", roleDb},
			{"writeConcern", bson.M{"w": "majority"}}}, &result); err != nil {
			log.Fatal("Failed to create user: ", err)
		}

	},
}

var user, userPassword, userRoles, targetDb, userExtraRoles, targetExtraDb string

func init() {
	userCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.Flags().StringVarP(&user, "user", "u", "", "User name to be created")
	createCmd.Flags().StringVarP(&userPassword, "userPassword", "p", "", "User password to be created")
	createCmd.Flags().StringVarP(&userRoles, "userRoles", "r", "", "List of mongodb roles for the user to be created")
	createCmd.Flags().StringVarP(&targetDb, "targetDb", "t", "", "The DB the user to be added in")
	createCmd.Flags().StringVarP(&userExtraRoles, "userExtraRoles", "a", "", "List of additional mongodb roles for 2nd db")
	createCmd.Flags().StringVarP(&targetExtraDb, "targetExtraDb", "x", "", "The 2nd DB the for the extra roles")

	viper.BindPFlag("user", createCmd.Flags().Lookup("user"))
	viper.BindPFlag("userPassword", createCmd.Flags().Lookup("userPassword"))
	viper.BindPFlag("userRoles", createCmd.Flags().Lookup("userRoles"))
	viper.BindPFlag("targetDb", createCmd.Flags().Lookup("targetDb"))
	viper.BindPFlag("userExtraRoles", createCmd.Flags().Lookup("userExtraRoles"))
	viper.BindPFlag("targetExtraDb", createCmd.Flags().Lookup("targetExtraDb"))
}
