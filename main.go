package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/allofthesepeople/github-org-prs/pullrequests"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	orgName      string
	apiKey       string
	returnFormat string
	orderby      string

	RootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if orgName == "" || apiKey == "" {
				fmt.Println("'org' and 'key' flags must be set")
				return
			}

			orders := strings.Split(orderby, ",")
			if len(orders) != 2 {
				fmt.Println("'orderby' should be formatted: 'columnName,direction'")
				return
			}
			orderCol := orders[0]
			orderDirection := orders[1]

			_, prs, err := pullrequests.GetPRs(orgName, apiKey)
			if err != nil {
				fmt.Println(err)
				return
			}

			prs = prs.Sort(orderCol, orderDirection)

			printToScreen(prs)
		},
	}
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&orgName, "org", "", "Github organization shortname")
	RootCmd.PersistentFlags().StringVar(&apiKey, "key", "", "Github API key")
	RootCmd.PersistentFlags().StringVarP(&returnFormat, "format", "f", "table", "The format to print to screen: table|json")
	RootCmd.PersistentFlags().StringVarP(&orderby, "orderby", "o", "UpdatedAt,desc", "Order the results: columnName,asc|desc ")

}

func printToScreen(prs pullrequests.PullRequestContainer) {
	switch returnFormat {
	case "table":
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(prs.Headers())
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.SetHeaderAlignment(3)
		table.SetColumnSeparator("")
		for _, pr := range prs {
			table.Append(pr.ToStrings())
		}
		table.Render()
		return
	case "json":
		prsBytes, err := json.Marshal(prs)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("%s", prsBytes))
	}
}
