package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	// parsers "./parsers"
)

const SN_WIKI string = "http://www.supernaturalwiki.com"

func main() {
	fmt.Println("Ready to start scraping...")

	fmt.Println("Fetching...")
	resp, err := http.Get(fmt.Sprintf("%v/Category:Transcripts", SN_WIKI))
	if err != nil {
		fmt.Println("========ERROR FETCHING", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Println("Reading...")
	dirPage, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("========ERROR READING: \n%s\n", err)
		os.Exit(1)
	}
	// fmt.Printf("---------------PAGE---------->\n%s\n", string(dirPage))

	var directory string = string(dirPage)
	var links []string
	tracking := false
	trackingStart := 0

	for i, s := range directory {
		if tracking == false {
			var prev4 string
			// making sure we can go back 4 spaces
			if i > 3 {
				prev4 = string(directory[i-4]) + string(directory[i-3]) + string(directory[i-2]) + string(directory[i-1])
			}
			if string(s) == "=" && prev4 == "href" {
				tracking = true
				trackingStart = i + 2
				fmt.Println("Tracking new href...")
			}
		} else {
			if (string(s) == "\"" || string(s) == "'") && i > trackingStart {
				links = append(links, directory[trackingStart:i])
				tracking = false
				fmt.Println("Found new href!", directory[trackingStart:i])
			}
		}
	}

	fmt.Println("Filtering hrefs for tanscript pages...")
	var filteredLinks []string
	for _, s := range links {
		if strings.HasSuffix(s, "(transcript)") {
			filteredLinks = append(filteredLinks, s)
			fmt.Println(s)
		}
	}

	fmt.Println("=======================")
	var lines [][]string

	for _, s := range filteredLinks {
		// if i == 0 {
		// transcriptPage := parsers.GetScript(fmt.Sprintf("%v%v", SN_WIKI, s))
		transcriptPage := GetScript(fmt.Sprintf("%v%v", SN_WIKI, s))

		scriptTags := GetLines(transcriptPage)

		lineObj := ToObjs(scriptTags)

		lines = append(lines, lineObj...)
		// }
	}

	fmt.Println("=================")
	for _, sl := range lines {
		fmt.Println(sl)
	}
	fmt.Println("=================")

	// writing to csv file
	file, err := os.Create("supernatural_lines.csv")
	if err != nil {
		fmt.Println("Unable to create file: supernatural_lines.csv ==>", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write headers first
	headers := []string{"Character", "Line"}
	err = writer.Write(headers)
	if err != nil {
		fmt.Println("Cannot write to file: supernatural_lines.csv ==>", err)
		os.Exit(1)
	}

	// writing each lines with character
	for _, value := range lines {
		err := writer.Write(value)
		if err != nil {
			fmt.Println("Cannot write to file: supernatural_lines.csv ==>", err)
			os.Exit(1)
		}
	}

}
