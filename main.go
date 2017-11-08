package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/allofthesepeople/github-org-prs/pullrequests"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	orgNameFlag string
	apiKeyFlag  string
	outputFlag  string
	sortByFlag  string
	columnsFlag string
	filtersFlag []string

	RootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if orgNameFlag == "" || apiKeyFlag == "" {
				fmt.Println("'org' and 'key' flags must be set")
				return
			}

			// Get all PRs
			_, prs, err := pullrequests.GetPRs(orgNameFlag, apiKeyFlag)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Filter
			filterArgs, err := setFilterArgs()
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, filterArg := range filterArgs {
				prs, err = prs.Filter(filterArg[0], filterArg[1], filterArg[2])
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			// Sort
			sortArgs, err := setSortArgs()
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, sortArg := range sortArgs {
				prs = prs.Sort(sortArg[0], sortArg[1])
			}

			// Set Columns
			var cols []string
			if columnsFlag == "all" {
				cols = pullrequests.Columns
			} else {
				cols = strings.Split(columnsFlag, ",")
			}

			// Print to screen
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
	RootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "table", "The format to print to screen: table|json")
	RootCmd.PersistentFlags().StringVarP(&sortByFlag, "sortby", "s", "UpdatedAt__desc", "Order the results: columnName__asc|desc")
	RootCmd.PersistentFlags().StringVarP(&columnsFlag, "columns", "c", "URL,Approved", "List of columns to return")
	RootCmd.PersistentFlags().StringArrayVarP(&filtersFlag, "filters", "f", []string{}, "List of filters to apply")
}

func printToScreen(prs pullrequests.PullRequestContainer, columns []string) {
	switch outputFlag {
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

func setSortArgs() ([][]string, error) {
	orderingOpts := strings.Split(sortByFlag, ",")
	var orderingArgs [][]string
	for _, o := range orderingOpts {
		opts := strings.Split(o, "__")

		if len(opts) != 2 {
			return nil, errors.New("'orderby' should be formatted: 'columnName__direction,'")
		}

		inList := false
		for _, col := range pullrequests.Columns {
			if col != opts[0] {
				inList = true
				break
			}
		}
		if inList == false {
			return nil, errors.New("'orderby' column name not recognised")
		}

		if !(opts[1] == "asc" || opts[1] == "desc") {
			return nil, errors.New("'orderby' direction should be 'asc' or 'desc'")
		}

		orderingArgs = append(orderingArgs, opts)
	}
	return orderingArgs, nil
}

func setFilterArgs() ([][]string, error) {

	var filterArgs [][]string

	for _, f := range filtersFlag {
		opts := strings.Split(f, ",")

		if len(opts) != 3 {
			return nil, errors.New("'filters' should be formatted: '{operator},{columnName},{value},'")
		}

		// Operator is available?
		if !(opts[0] == "eq" || opts[0] == "neq") {
			return nil, errors.New("'filters' operator not recognised: should be 'eq' or 'neq'")
		}

		// ColumnName available?
		inList := false
		for _, col := range pullrequests.Columns {
			if col != opts[1] {
				inList = true
				break
			}
		}
		if inList == false {
			return nil, errors.New("'filters' column name not recognised")
		}

		filterArgs = append(filterArgs, opts)
	}
	return filterArgs, nil
}
