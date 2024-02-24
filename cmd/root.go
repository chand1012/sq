/*
Copyright © 2024 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chand1012/sq/pkg/constants"
	"github.com/chand1012/sq/pkg/db"
	"github.com/chand1012/sq/pkg/file_types"
	"github.com/chand1012/sq/pkg/logger"
	"github.com/chand1012/sq/pkg/utils"
)

var log = logger.DefaultLogger

var inputFilePath string
var tableName string
var quiet bool
var outputFormat string // csv, json, jsonl, sqlite
var outputFilePath string
var columnNames bool
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:    "sq",
	Short:  "Convert and query JSON, JSONL, CSV, and SQLite with ease!",
	Long:   `Like jq, but for SQL! Simply pipe in your data or specify a file and run your SQL queries!`,
	Run:    run,
	Args:   cobra.MaximumNArgs(1),
	PreRun: prerun,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputFilePath, "read", "r", "", "input file path")
	rootCmd.Flags().StringVarP(&tableName, "table", "t", "", "table name")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "output format. csv, json, jsonl")
	rootCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "output file path")
	rootCmd.Flags().BoolVarP(&columnNames, "columns", "c", false, "print the columns names and exit")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print verbose output. Prints full stack trace for debugging.")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "execute the query and exit without printing anything")
}

func run(cmd *cobra.Command, args []string) {
	var err error
	var d *sql.DB
	// var filePath string
	var content string
	// always has one argument, which is a sql query
	query := args[0]
	// tread the input as bytes
	// no matter where it came from
	// or what it is
	var input []byte

	// if the input file path is not empty
	// load the bytes from the file
	if inputFilePath != "" {
		d, _, err = db.LoadFile(inputFilePath)
		if err != nil {
			file, err := os.ReadFile(inputFilePath)
			if err != nil {
				logger.HandlePanic(log, err, verbose)
			}
			input = file
		}
	} else {
		// check stdin
		input, err = utils.ReadStdin()
		if err != nil {
			logger.HandlePanic(log, err, verbose)
		}
	}

	if tableName == "" {
		// if an error occurs it stays empty
		tableName, err = utils.GetTableName(query)
		if err != nil {
			log.Warn(err.Error())
			tableName = constants.TableName
		}
	}

	// if the database hasn't been loaded
	if d == nil {
		// resolve the type

		// if input is empty, panic
		if len(input) == 0 {
			logger.HandlePanic(log, errors.New("input is empty"), verbose)
		}

		switch file_types.Resolve(input) {
		case file_types.SQLite:
			d, _, err = db.LoadStdin(input)
		case file_types.JSONL:
			d, _, err = db.FromJSONL(input, tableName)
		case file_types.JSON:
			d, _, err = db.FromJSON(input, tableName)
		case file_types.CSV:
			d, _, err = db.FromCSV(input, tableName)
		default:
			logger.HandlePanic(log, errors.New("unsupported file type"), verbose)
		}
	}

	if err != nil {
		logger.HandlePanic(log, err, verbose)
	}
	defer d.Close()

	if columnNames {
		columns, err := db.GetColumnNames(d, tableName)
		if err != nil {
			logger.HandlePanic(log, err, verbose)
		}
		os.Stdout.WriteString(strings.Join(columns, ",") + "\n")
		os.Exit(0)
	}

	// run the query
	rows, err := d.Query(query)
	if err != nil {
		logger.HandlePanic(log, err, verbose)
	}

	if quiet {
		os.Exit(0)
	}

	switch outputFormat {
	case "json":
		content, err = db.RowsToJSON(rows)
	case "jsonl":
		content, err = db.RowsToJSONL(rows)
	case "csv":
		content, err = db.RowsToCSV(rows)
	case "sqlite":
		if outputFilePath == "" {
			logger.HandlePanic(log, errors.New("output file path required for sqlite output"), verbose)
		}
		err = db.RowsToSQLite(rows, tableName, outputFilePath)
		if err != nil {
			logger.HandlePanic(log, err, verbose)
		}
		os.Exit(0)
	default:
		logger.HandlePanic(log, errors.New("unsupported output format"), verbose)
	}

	if err != nil {
		logger.HandlePanic(log, err, verbose)
	}

	if outputFilePath != "" {
		err = os.WriteFile(outputFilePath, []byte(content), 0644)
		if err != nil {
			logger.HandlePanic(log, err, verbose)
		}
	} else {
		os.Stdout.WriteString(content)
	}
}

// sets up the logger and output format
func prerun(cmd *cobra.Command, args []string) {
	if verbose {
		log = logger.VerboseLogger
	}
	if outputFilePath != "" && outputFormat == "" {
		outputFormat = file_types.ResolveByPath(outputFilePath).String()
	} else if outputFormat == "" {
		outputFormat = "csv"
	}
	// if outputFormat is not one of json, jsonl, or csv, panic with an error
	if outputFormat != "json" && outputFormat != "jsonl" && outputFormat != "csv" && outputFormat != "sqlite" {
		logger.HandlePanic(log, errors.New("unsupported output format"), verbose)
	}
}
