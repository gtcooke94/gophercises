/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"errors"
	"github.com/boltdb/bolt"
	"github.com/gtcooke94/gophercises/task/task_helpers"
	"github.com/spf13/cobra"
	"strings"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the task list",
	Long: `Adds a task to the list list. For example:

	$ task add Do laundry

Would add the "Do laundry" to the task list`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Add requires a task to add")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		addTask(task)
		fmt.Println(task)
	},
}

func addTask(task string) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TASK_BUCKET))
		seqId, _ := b.NextSequence()
		id := task_helpers.Itob(int(seqId))
		err := b.Put(id, []byte(task))
		fmt.Printf("Adding %d: %s\n", id, task)
		return err
	})
	defer db.Close()
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
