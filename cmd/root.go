package benchmark

import (
	"os"

	"github.com/spf13/cobra"
)

// TODO: args, run
var benchCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "benchmark SELECT query performance for given queries",
	Long:  "given a csv file or standard input and number of workers, find various benchmarking values",
	Args: cobra.  
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// TODO
func Execute() {
	err := benchCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// TODO
func init() {
	benchCmd.AddCommand(benchCmd)

	// adding Flag and making required
	benchCmd.PersistentFlags().Int32P("workers", "w", 1, "number of workers")
	benchCmd.MarkPersistentFlagRequired("workers")
}	
