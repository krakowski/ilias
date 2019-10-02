package submissions

import (
	"encoding/csv"
	"github.com/krakowski/ilias/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	shouldPrintCsv bool
	includeEmpty bool
)

var submissionsListCommand = &cobra.Command{
	Use:   "list [exerciseId] [assignmentId]",
	Short: "Lists all submissions within an submissions",
	SilenceErrors: true,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new API client
		client, err := api.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		submissions, err := client.GetSubmissions(args[0], args[1], includeEmpty)
		if err != nil {
			log.Fatal(err)
		}

		if shouldPrintCsv {
			printCsv(submissions)
		} else {
			printTable(submissions)
		}
	},
}

func printCsv(submissions []api.SubmissionInfo)  {
	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"Kennung", "Nachname", "Vorname", "Abgabe"})

	for _, submission := range submissions {
		writer.Write(submission.ToRow())
	}

	writer.Flush()
}

func printTable(submissions []api.SubmissionInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kennung", "Nachname", "Vorname", "Abgabe"})

	for _, submission := range submissions {
		table.Append(submission.ToRow())
	}

	table.Render()
}

func init() {
	submissionsListCommand.Flags().BoolVar(&shouldPrintCsv, "csv", false, "Prints the table in csv format")
	submissionsListCommand.Flags().BoolVar(&includeEmpty, "empty", false, "Includes empty submissions")
}
