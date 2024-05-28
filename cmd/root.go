package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "benchmark SELECT query performance for given queries",
	Long:  "given a csv file or standard input and number of workers, find various benchmarking values",
	Args: func(cmd *cobra.Command, args []string) error {
		// checks at least 1 arg
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		// checks if 1 arg that it is a valid csv file
		// will leave checking presence of csv file to run
		if len(args) == 1 {
			if _, err := regexp.MatchString("^[\\w,\\s-]+\\.csv$", args[0]); err != nil {
				return err
			}
			return nil
		}
		// if multiple args, will check that first line follows proper csv format of
		// hostname,start_time,end_time
		if _, err := regexp.MatchString("hostname,start_time,end_time", args[0]); err != nil {
			return err
		}
		return fmt.Errorf("invalid input")
	},
	//TODO STILL WORKING ON
	Run: func(cmd *cobra.Command, args []string) {
		// get number of workers and initialize respective channels
		// putting buffer of 3 for each channel in case some backup
		numWork, _ := cmd.Flags().GetInt("workers")
		var ch = make([]chan string, numWork)
		for i := range ch {
			ch[i] = make(chan string, 3)
		}
		outCh := make(chan int64)

		// initialze wait group
		// while data will be split for workers to take on the same hostname, 1 waitgroup used to synchronize
		wg := sync.WaitGroup{}

		// first determine if reading csv or csv data stream
		isFile, _ := regexp.MatchString("^[\\w,\\s-]+\\.csv$", args[0])

		// based on how data is passed in, prepare and run as needed
		// if csv file was passed in, attempt to open and access
		if isFile {
			file, err := os.Open(args[0])

			if err != nil {
				log.Fatal("Error while reading the file", err)
			}

			defer file.Close()

			reader := csv.NewReader(file)

			// begin processing
		} else {
			// otherwise, read csv formatted data stream from args

		}

		// perform necessary metrics calculations

	},
}

// TODO
func query(vals []string, out chan int64) {

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(rootCmd)

	// adding Flag and making required
	rootCmd.Flags().IntP("workers", "w", 1, "number of workers")
	rootCmd.MarkFlagRequired("workers")
}
