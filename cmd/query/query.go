package query

import (
	"github.com/c-bata/go-prompt"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/table"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "query",
		Short:            "Query infrastructure data from pgstorage",
		Long:             "Query infrastructure data from pgstorage",
		PersistentPreRun: global.DefaultWrappedInit(),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ui.Warningln("Please select table.")

			queryClient, _ := NewQueryClient(ctx)
			p := prompt.New(func(in string) {
				strArr := strings.Split(in, "/")
				s := strArr[0]

				res, err := queryClient.Storage.Query(ctx, s)
				if err != nil {
					ui.Errorln(err)
				} else {
					tables, e := res.ReadRows(-1)
					if e != nil && e.HasError() {
						return
					}
					header := tables.GetColumnNames()
					body := tables.GetMatrix()
					var tableBody [][]string
					for i := range body {
						var row []string
						for j := range body[i] {
							row = append(row, utils.Strava(body[i][j]))
						}
						tableBody = append(tableBody, row)
					}

					if len(strArr) > 1 && strArr[1] == "g" {
						table.ShowRows(header, tableBody, []string{}, true)
					} else {
						table.ShowTable(header, tableBody, []string{}, true)
					}

				}
				if s == "exit;" || s == ".exit" {
					os.Exit(0)
				}
			}, queryClient.completer,
				prompt.OptionTitle("Table"),
				prompt.OptionPrefix("> "),
				prompt.OptionAddKeyBind(prompt.KeyBind{
					Key: prompt.ControlC,
					Fn: func(buffer *prompt.Buffer) {
						os.Exit(0)
					},
				}),
			)
			p.Run()
		},
	}
	return cmd
}
