package formatters

import (
	"os"

	"github.com/georgboe/cover/models"
	"github.com/olekukonko/tablewriter"
)

type TableFormatter struct{}

func (t TableFormatter) Render(pairs []models.FileAndTest) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Source", "Test"})

	for _, v := range pairs {

		var sourceFile string
		var testFile string

		if v.Sourcefile != nil {
			sourceFile = v.Sourcefile.Filename
		}

		if v.Testfile != nil {
			testFile = v.Testfile.Filename
		}

		table.Append([]string{sourceFile, testFile})
	}
	table.Render()
}
