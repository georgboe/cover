package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func getFileContent(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	a, _ := ioutil.ReadAll(file)

	return string(a)
}

func BenchmarkGetFiles(b *testing.B) {
	content := getFileContent("java.txt")
	reader := strings.NewReader(content)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getFiles(reader)
	}
}

// func BenchmarkGetPairs(b *testing.B) {
// 	content := getFileContent("java.txt")
// 	reader := strings.NewReader(content)
// 	filesOriginal := getFiles(reader)
// 	sort.Sort(bySortname(filesOriginal))

// 	b.N = 10000

// 	files := make([]file, len(filesOriginal))
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		b.StopTimer()
// 		copy(files, filesOriginal)
// 		b.StartTimer()
// 		createPairs(files)
// 	}
// }

func BenchmarkGetData(b *testing.B) {

	content := getFileContent("java.txt")
	reader := strings.NewReader(content)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getData(reader)
	}
}
