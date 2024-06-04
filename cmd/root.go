package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "benchmark SELECT query performance for given queries",
	Long:  "given a csv file or standard input and number of workers, find various benchmarking values",
	Args: func(cmd *cobra.Command, args []string) error {
		// checks at least 1 arg
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		// checks if a valid csv file
		// will leave checking presence of csv file to run
		matchCSV, err := regexp.MatchString("^[\\w,\\s-]+\\.csv$", args[0])
		if err != nil {
			return err
		}

		// check if contains at least csv header line of
		// hostname,start_time,end_time
		matchARG, err := regexp.MatchString("^hostname,start_time,end_time(.|\n)*$", args[0])
		if err != nil {
			return err
		}

		if !matchCSV || !matchARG {
			return fmt.Errorf("args does not match either csv file or required csv-formatted input")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// create TimescaleDB connection pool
		connStr := os.Getenv("DATABASE_CONNECTION_STRING")
		ctx := context.Background()
		dbpool, err := pgxpool.New(ctx, connStr)
		if err != nil {
			return fmt.Errorf("unable to connect to database: %v", err)
		}
		defer dbpool.Close()

		// get number of workers and initialize respective workers and channels
		// putting buffer of 3 for each channel in case some backup
		numWork, _ := cmd.Flags().GetInt("workers")
		var ch = make([]chan []string, numWork)
		out := make(chan float64, numWork)
		var wg sync.WaitGroup
		for i := range ch {
			ch[i] = make(chan []string, 3)
			wg.Add(1)
			go worker(ch[i], out, &wg, dbpool)
		}

		// initalize tools for maintaining each worker takes on queries for the same hostname
		m := make(map[string]int)
		next := 1 % numWork

		// begin collecting values
		var results_wg sync.WaitGroup
		results_wg.Add(1)
		go collectResults(out, &results_wg)

		// first determine if reading csv or csv data stream
		isFile, _ := regexp.MatchString("^[\\w,\\s-]+\\.csv$", args[0])
		var reader *csv.Reader

		// based on how data is passed in, prepare and run as needed
		// if csv file was passed in, attempt to open and access
		if isFile {
			file, err := os.Open(args[0])

			if err != nil {
				return fmt.Errorf("error while opening the file: %v", err)
			}

			defer file.Close()

			reader = csv.NewReader(file)

			// check csv file is proper format
			line, err := reader.Read()
			if err != nil {
				return fmt.Errorf("error while reading file: %v", err)
			}

			if line[0] != "hostname" || line[1] != "start_time" || line[2] != "end_time" {
				return fmt.Errorf("csv headers do not match required: %q, %q, %q", line[0], line[1], line[2])
			}
		} else {
			// otherwise, read csv formatted data stream from args TODO
			reader = csv.NewReader(strings.NewReader(args[0]))
			reader.Read()
			//discard first line of headers as already checked
		}

		// begin processing
		// keep reading from file until hitting EOF or error
		for {

			line, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error while reading file: %v", err)
			}

			// check if host has been queried before
			// if not, send to next worker
			check := m[line[0]] - 1
			if check == -1 {
				m[line[0]] = next
				ch[next] <- line
				next = (next + 1) % numWork
			} else {
				ch[check] <- line
			}

		}

		for i := range ch {
			close(ch[i])
		}
		wg.Wait()
		close(out)
		results_wg.Wait()

		return nil
	},
}

func worker(job chan []string, out chan float64, wg *sync.WaitGroup, dbpool *pgxpool.Pool) {
	defer wg.Done()
	for work := range job {
		// run each query and record time taken
		var min float64
		var max float64
		start := time.Now()
		err := dbpool.QueryRow(context.Background(), "SELECT MIN(usage) min, MAX(usage) max FROM cpu_usage WHERE host = $1 AND ts BETWEEN $2 AND $3", work[0], work[1], work[2]).Scan(&min, &max)
		stop := time.Now()
		if err != nil {
			log.Fatal("error while iterating dataset: ", err)
		}
		//note for now actual result of query not important
		elapsed := stop.Sub(start)
		out <- elapsed.Seconds()
	}
}

func collectResults(out chan float64, wg *sync.WaitGroup) {
	defer wg.Done()
	// begin collecting values
	count := 0
	totalTime := float64(0)
	minVal := math.Inf(1)
	maxVal := math.Inf(-1)
	vals := []float64{}

	for result := range out {
		count++
		totalTime += result
		minVal = math.Min(minVal, float64(result))
		maxVal = math.Max(maxVal, float64(result))
		vals = append(vals, float64(result))
	}

	// perform necessary metrics calculations

	sort.Float64s(vals)
	var median float64
	if count%2 == 0 {
		median = (vals[count/2-1] + vals[count/2]) / 2
	} else {
		median = vals[count/2]
	}
	average := float64(totalTime) / float64(count)

	// output metrics
	fmt.Print("Returning follow metrics for queries performed:\n")
	fmt.Printf("total queries: %d\n", count)
	fmt.Printf("total processing time: %f seconds\n", totalTime)
	fmt.Printf("minimum query time: %f seconds\n", minVal)
	fmt.Printf("median query time: %f seconds\n", median)
	fmt.Printf("maximum query time: %f seconds\n", maxVal)
	fmt.Printf("average query time: %f seconds\n", average)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// adding Flag and making required
	rootCmd.Flags().IntP("workers", "w", 1, "number of workers")
	rootCmd.MarkFlagRequired("workers")
}
