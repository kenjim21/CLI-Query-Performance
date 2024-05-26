package cmd

import (
	"os"

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
			if err:= regexp.MatchString("^[\w,\s-]+\.csv$", args[0]); err != nil {
				return err
			}
			else {
				return nil
			}
		}
		// if multiple args, will check that first line follows proper csv format of 
		// hostname,start_time,end_time
		if err:= regexp.MatchString("hostname,start_time,end_time", args[0]); err != nil {
		  return err
		}
		return fmt.Error("invalid input")
	  },
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(benchCmd)

	// adding Flag and making required
	rootCmd.Flags().Int32P("workers", "w", 1, "number of workers")
	rootCmd.MarkFlagRequired("workers")
}
