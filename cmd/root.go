/*
Copyright © 2024 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"database/sql"
	"os"

	"github.com/spf13/cobra"

	"github.com/chand1012/sq/pkg/db"
	"github.com/chand1012/sq/pkg/file_types"
	"github.com/chand1012/sq/pkg/utils"
)

var inputFilePath string
var tableName string
var quiet bool
var outputFormat string // csv, json, jsonl
var outputFilePath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sq",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run:    run,
	Args:   cobra.ExactArgs(1),
	PreRun: prerun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sq.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "input file path")
	rootCmd.Flags().StringVarP(&tableName, "table", "t", "", "table name")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "execute the query and exit without printing anything")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "csv", "output format. csv, json, jsonl")
	rootCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "output file path")
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
				panic(err)
			}
			input = file
		}
	} else {
		// check stdin
		input, err = utils.ReadStdin()
		if err != nil {
			panic(err)
		}
	}

	if tableName == "" {
		// if an error occurs it stays empty
		tableName, _ = utils.GetTableName(query)
	}

	// if the database hasn't been loaded
	if d == nil {
		// resolve the type

		// if input is empty, panic
		if len(input) == 0 {
			panic("no input")
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
			panic("unsupported file type")
		}
	}

	if err != nil {
		panic(err)
	}

	// run the query
	rows, err := d.Query(query)
	if err != nil {
		panic(err)
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
	default:
		panic("unsupported output format")
	}

	if err != nil {
		panic(err)
	}

	if outputFilePath != "" {
		err = os.WriteFile(outputFilePath, []byte(content), 0644)
		if err != nil {
			panic(err)
		}
	} else {
		os.Stdout.WriteString(content)
	}
}

func prerun(cmd *cobra.Command, args []string) {
	// if outputFormat is not one of json, jsonl, or csv, panic with an error
	if outputFormat != "json" && outputFormat != "jsonl" && outputFormat != "csv" {
		panic("unsupported output format")
	}
}
