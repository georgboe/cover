package formatters

import (
	"fmt"

	"github.com/georgboe/cover/models"
)

type ExcelFormatter struct{}

func (t ExcelFormatter) Render(pairs []models.FileAndTest) {
	for _, v := range pairs {

		var sourceFile string
		var testFile string

		if v.Sourcefile != nil {
			sourceFile = v.Sourcefile.Filename
		}

		if v.Testfile != nil {
			testFile = v.Testfile.Filename
		}

		fmt.Printf("%v\t%v\n", sourceFile, testFile)
	}
}
