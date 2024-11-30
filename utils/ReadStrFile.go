package utils

import (
	"bufio"
	"encoding/csv"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ReadSrtFilesInInputDirAndWriteWordCountsToCSV reads all .srt files in the input directory, extracts all English words, counts their occurrences, orders them by highest occurrences, and writes the results to a CSV file in the output directory, keeping the original file name
func ReadSrtFilesInInputDirAndWriteWordCountsToCSV() error {
	inputDir := "input"
	outputDir := "output"
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".srt" {
			wordCounts := make(map[string]int)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				matches := wordRegex.FindAllString(line, -1)
				for _, word := range matches {
					wordCounts[strings.ToLower(word)]++
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				return err
			}

			csvFileName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) + ".csv"
			csvFilePath := filepath.Join(outputDir, csvFileName)
			csvFile, err := os.Create(csvFilePath)
			if err != nil {
				return err
			}
			defer csvFile.Close()

			writer := csv.NewWriter(csvFile)
			defer writer.Flush()

			type wordCount struct {
				Word  string
				Count int
			}

			var wordCountList []wordCount
			for word, count := range wordCounts {
				wordCountList = append(wordCountList, wordCount{Word: word, Count: count})
			}

			sort.Slice(wordCountList, func(i, j int) bool {
				return wordCountList[i].Count > wordCountList[j].Count
			})

			for _, wc := range wordCountList {
				err := writer.Write([]string{wc.Word, strconv.Itoa(wc.Count)})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return err
}
