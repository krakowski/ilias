package members

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/krakowski/ilias/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	shouldPrintIdsOnly bool
	shouldPrintCsv bool
	shouldPrintJson bool
	includeEmpty bool
)

var membersListCommand = &cobra.Command{
	Use:   "list [courseId]",
	Short: "Lists all members within a course",
	SilenceErrors: true,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new API client
		client, err := api.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		members, err := client.GetMembers(args[0], includeEmpty)
		if err != nil {
			log.Fatal(err)
		}

		if shouldPrintCsv {
			printCsv(members)
		} else if shouldPrintIdsOnly {
			printIds(members)
		} else if shouldPrintJson {
			printJson(members)
		} else {
			printTable(members)
		}
	},
}

func printJson(members []api.MemberInfo) {
	type member struct {
		Usernames []string `json:"usernames"`
	}

	var usernames []string
	for _, member := range members {
		usernames = append(usernames, member.UserId)
	}

	buffer, err := json.Marshal(member{Usernames:usernames})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buffer))
}

func printIds(members []api.MemberInfo) {
	for _, member := range members {
		fmt.Println(member.UserId)
	}
}

func printCsv(members []api.MemberInfo)  {
	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"Kennung", "Vorname", "Nachname", "Rolle"})

	for _, member := range members {
		writer.Write(member.ToRow())
	}

	writer.Flush()
}

func printTable(members []api.MemberInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kennung", "Vorname", "Nachname", "Rolle"})

	for _, member := range members {
		table.Append(member.ToRow())
	}

	table.Render()
}

func init() {
	membersListCommand.Flags().BoolVar(&shouldPrintJson, "json", false, "Prints all user ids within an json object")
	membersListCommand.Flags().BoolVar(&shouldPrintIdsOnly, "ids", false, "Prints only user ids")
	membersListCommand.Flags().BoolVar(&shouldPrintCsv, "csv", false, "Prints the table in csv format")
	membersListCommand.Flags().BoolVar(&includeEmpty, "empty", false, "Includes empty submissions")
}
