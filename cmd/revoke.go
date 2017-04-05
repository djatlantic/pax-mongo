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
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

// revokeCmd represents the revoke command
var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke roles from user associated with a DB",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		roles := strings.Split(viper.GetString("revokeUserRoles"), ",")
		role_db := make([]interface{}, len(roles))
		db := viper.GetString("revokeTargetDb")

		for i, v := range roles {
			role_db[i] = bson.M{"role": v, "db": db}
		}

		session := connect.Connect()
		defer session.Close()

		result := bson.M{}

		if err := session.DB(db).Run(bson.D{{"revokeRolesFromUser", viper.GetString("revokeUser")},
			{"roles", role_db},
			{"writeConcern", bson.M{"w": "majority"}}}, &result); err != nil {
			log.Fatal("Faield to revoke roles: ", err)
		}
	},
}

func init() {
	userCmd.AddCommand(revokeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// revokeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// revokeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	revokeCmd.Flags().StringVarP(&user, "revokeUser", "u", "", "User name to be revoked with role(s)")
	revokeCmd.Flags().StringVarP(&userRoles, "revokeUserRoles", "r", "", "List of mongodb roles for the user to be revoked")
	revokeCmd.Flags().StringVarP(&userRoles, "revokeTargetDb", "t", "", "The DB the user to be revoked from")

	viper.BindPFlag("revokeUser", revokeCmd.Flags().Lookup("revokeUser"))
	viper.BindPFlag("revokeUserRoles", revokeCmd.Flags().Lookup("revokeUserRoles"))
	viper.BindPFlag("revokeTargetDb", revokeCmd.Flags().Lookup("revokeTargetDb"))
}
