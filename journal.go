// This program reads the latest journal entry and moves incomplete tasks to a new
// journal file. It also creates standard (blank) sections for filling out during the day.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	journalDir := os.Getenv("JOURNAL_DIR")

	files, err := os.ReadDir(journalDir)
	if err != nil {
		panic("Directory JOURNAL_DIR (" + journalDir + ") does not exist")
	}

	latestFile, found := findLatestFilename(files)
	createNewJournalFile(journalDir, latestFile, found)
}

// Create a new journal file for today copying tasks from the base file denoted by `baseFilename`
func createNewJournalFile(dir string, baseFilename string, shouldCopy bool) {
	var file io.Reader
	if shouldCopy {
		baseFile, err := os.Open(dir + baseFilename)
		if err != nil {
			panic("Unable to open most recent file: " + dir + baseFilename)
		}
		defer baseFile.Close()
		file = baseFile
	} else {
		file = bytes.NewBufferString("")
	}

	year, month, day := time.Now().Date()
	newFileName := fmt.Sprintf("Journal-%04d-%02d-%02d.md", year, month, day)

	if _, err := os.Stat(dir + newFileName); !os.IsNotExist(err) {
		panic("File already exists: " + dir + newFileName)
	}

	newFile, err := os.Create(dir + newFileName)
	if err != nil {
		panic("Unable to open file for writing: " + dir + newFileName)
	}
	defer newFile.Close()

	writer := bufio.NewWriter(newFile)
	writer.WriteString(fmt.Sprintf("# Daily Work Journal %04d-%02d-%02d\n\n", year, month, day))
	writer.WriteString("## Tasks\n\n")
	copyIncompleteTasks(file, writer)
	writer.WriteString("\n\n")
	writer.WriteString("## Interactions \n\n")
	writer.WriteString("## Notes \n\n")

	writer.Flush()
}

// Copies lines from the tasks section that have a incomplete notation. (e.g. "[]")
// Each of these lines is appended with an asterisk to easily see how many days it has been
// deferred for
func copyIncompleteTasks(file io.Reader, writer *bufio.Writer) {
	scanner := bufio.NewScanner(file)

	inTasks := false
SCANNING:
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		switch {
		case strings.HasPrefix(line, "## Tasks"):
			inTasks = true
		case strings.HasPrefix(line, "## Interactions"):
			break SCANNING
		case inTasks:
			if strings.Contains(line, "[]") {
				writer.WriteString(strings.Replace(line, "[]", "*[]", 1))
			}
		}
	}
}

var regex = regexp.MustCompile(".*-(\\d{4}-\\d{2}-\\d{2}).*")

func findLatestFilename(files []os.DirEntry) (string, bool) {
	var largestYet int64
	var largestFilename string
	var found = false
	for _, file := range files {
		if !file.IsDir() {
			matches := regex.FindStringSubmatch(file.Name())
			if len(matches) < 2  {
				continue
			}

			dateString := matches[1]
			time, err := time.Parse("2006-01-02", dateString)
			if err != nil {
				continue
			}

			timeUnix := time.Unix()
			if largestYet < timeUnix {
				found = true
				largestYet = timeUnix
				largestFilename = file.Name()
			}
		}
	}

	return largestFilename, found
}
