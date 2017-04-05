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
		mRoles := make([]mgo.Role, len(roles))

		for i, r := range roles {
			mRoles[i] = mgo.Role(Roles[r])
		}

		thisUser := mgo.User{
			Username: viper.GetString("user"),
			Password: viper.GetString("userPassword"),
			Roles:    mRoles,
		}

		session := connect.Connect()
		defer session.Close()

		if err := session.DB(viper.GetString("targetDb")).UpsertUser(&thisUser); err != nil {
			log.Fatal("Failed to create user", thisUser)
		}
	},
}

var user, userPassword, userRoles string

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
	createCmd.Flags().StringVarP(&userRoles, "targetDb", "t", "", "The DB the user to be added in")

	viper.BindPFlag("user", createCmd.Flags().Lookup("user"))
	viper.BindPFlag("userPassword", createCmd.Flags().Lookup("userPassword"))
	viper.BindPFlag("userRoles", createCmd.Flags().Lookup("userRoles"))
	viper.BindPFlag("targetDb", createCmd.Flags().Lookup("targetDb"))
}
