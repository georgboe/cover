package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/georgboe/cover/formatters"
	"github.com/georgboe/cover/models"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	excel = kingpin.Flag("excel", "Use excel formatter.").Short('x').Bool()
)

func getFiles(reader io.Reader) []models.File {
	scanner := bufio.NewScanner(reader)
	files := make([]models.File, 0, 10)
	for scanner.Scan() {
		path := scanner.Text()
		if len(path) == 0 {
			continue
		}
		f := models.File{
			Filename: filepath.Base(path),
		}
		f.Basename = strings.ToLower(
			strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename)))
		f.Testname = f.Basename + "test"
		f.IsTest = strings.HasSuffix(f.Basename, "test")
		files = append(files, f)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return files
}

func getPairs(files []models.File) []models.FileAndTest {
	fileCount := len(files)
	pairs := make([]models.FileAndTest, 0, fileCount)
	fileList := make(map[string]*models.File, fileCount)
	for i := range files {
		f := &files[i]
		fileList[f.Basename] = f
	}

	skippedFiles := make(map[*models.File]bool)
	for i := range files {
		f := &files[i]

		_, ok := skippedFiles[f]
		if ok {
			continue
		}

		testFile, ok := fileList[f.Testname]

		if ok {
			pairs = append(pairs, models.FileAndTest{Sourcefile: f, Testfile: testFile})
			skippedFiles[testFile] = false
		} else if f.IsTest {
			pairs = append(pairs, models.FileAndTest{Testfile: f})
		} else {
			pairs = append(pairs, models.FileAndTest{Sourcefile: f})
		}
	}
	return pairs
}

func getData(reader io.Reader) []models.FileAndTest {
	files := getFiles(reader)
	sort.Slice(files, func(i, j int) bool { return files[i].Basename < files[j].Basename })
	return getPairs(files)
}

func main() {
	kingpin.Parse()

	pairs := getData(os.Stdin)

	var formatter formatters.Formatter

	formatter = formatters.TableFormatter{}

	if *excel {
		formatter = formatters.ExcelFormatter{}
	}

	formatter.Render(pairs)
}
