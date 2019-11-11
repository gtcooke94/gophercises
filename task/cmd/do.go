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

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"strconv"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do <task>",
	Short: "Complete a task",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("do called")
		taskNum, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		doTask(taskNum)
	},
}

func doTask(task int) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TASK_BUCKET))
		completeBucket := tx.Bucket([]byte(COMPLETE_BUCKET))

		c := b.Cursor()

		taskCounter := 0
		var k, v []byte
		for k, v = c.First(); k != nil; k, v = c.Next() {
			if taskCounter == task {
				break
			}
			taskCounter++
		}
		b.Delete(k)
		completeBucket.Put(k, v)
		return nil
	})
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
