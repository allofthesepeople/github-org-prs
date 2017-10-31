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
	orgNameFlag      string
	apiKeyFlag       string
	returnFormatFlag string
	orderbyFlag      string
	columnsFlag      string

	RootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if orgNameFlag == "" || apiKeyFlag == "" {
				fmt.Println("'org' and 'key' flags must be set")
				return
			}

			orderingOpts := strings.Split(orderbyFlag, ",")
			var orderings [][]string
			for _, o := range orderingOpts {
				opts := strings.Split(o, "__")

				if len(opts) != 2 {
					fmt.Println("'orderby' should be formatted: 'columnName__direction,'")
					return
				}

				inList := false
				for _, col := range pullrequests.Columns {
					if col != opts[0] {
						inList = true
						break
					}
				}
				if inList == false {
					fmt.Println("'orderby' column name not recognised")
					return
				}

				if !(opts[1] == "asc" || opts[1] == "desc") {
					fmt.Println("'orderby' direction should be 'asc' or 'desc'")
					return
				}

				orderings = append(orderings, opts)
			}

			_, prs, err := pullrequests.GetPRs(orgNameFlag, apiKeyFlag)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, ord := range orderings {
				prs = prs.Sort(ord[0], ord[1])
			}

			var cols []string
			if columnsFlag == "all" {
				cols = pullrequests.Columns
			} else {
				cols = strings.Split(columnsFlag, ",")
			}

			printToScreen(prs, cols)
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
	RootCmd.PersistentFlags().StringVar(&orgNameFlag, "org", "", "Github organization shortname")
	RootCmd.PersistentFlags().StringVar(&apiKeyFlag, "key", "", "Github API key")
	RootCmd.PersistentFlags().StringVarP(&returnFormatFlag, "format", "f", "table", "The format to print to screen: table|json")
	RootCmd.PersistentFlags().StringVarP(&orderbyFlag, "orderby", "o", "UpdatedAt__desc", "Order the results: columnName__asc|desc")
	RootCmd.PersistentFlags().StringVarP(&columnsFlag, "columns", "c", "URL,Approved", "List of columns to return")
}

func printToScreen(prs pullrequests.PullRequestContainer, columns []string) {
	switch returnFormatFlag {
	case "table":
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(columns)
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.SetHeaderAlignment(3)
		table.SetColumnSeparator("")
		for _, pr := range prs {
			table.Append(pr.ToStrings(columns))
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
