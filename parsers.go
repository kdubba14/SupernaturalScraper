package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ScriptLine struct {
	Character string
	Line      string
}

// GetScript getting the page with script
func GetScript(r string) string {
	f, err := http.Get(r)
	if err != nil {
		fmt.Println("========ERROR FETCHING", err)
		os.Exit(1)
	}
	defer f.Body.Close()
	fmt.Printf("Reading from %v...\n", r)
	transcriptPage, err := ioutil.ReadAll(f.Body)
	if err != nil {
		fmt.Printf("========ERROR READING: \n%s\n", err)
		os.Exit(1)
	}

	return string(transcriptPage)
}

// GetLines will pull all script lines from a script page
func GetLines(s string) []string {
	var lines []string
	tracking := false
	trackingStart := 0

	for i, l := range s {
		if tracking == false {
			var prev2 string
			// making sure we can go back 2 spaces
			if i > 2 {
				prev2 = string(s[i-2]) + string(s[i-1])
			}
			if prev2+string(l) == "<p>" {
				tracking = true
				trackingStart = i + 1
				// fmt.Println("Tracking new line...")
			}
		} else {
			if (string(l) == "<" && strings.ToLower(string(s[i+1])) != "b") && i > trackingStart {
				lines = append(lines, s[trackingStart:i])
				tracking = false
				// fmt.Println("Found new line!", s[trackingStart:i])
			}
		}
	}

	var filteredLines []string
	for _, l := range lines {
		if strings.Contains(l, "<br>") == true || strings.Contains(l, "<br />") == true || strings.Contains(l, "<br/>") == true {
			filteredLines = append(filteredLines, l)
		}
	}

	return filteredLines
}

func ToObjs(lines []string) [][]string {
	var r [][]string

	for _, l := range lines {
		var arr []string

		spl := strings.Split(l, "<br />")
		line := strings.ReplaceAll(spl[1], "&#8212;", " - ")

		if spl[0] != "" && line != "" {
			arr = append(arr, strings.TrimSpace(spl[0]))
			arr = append(arr, strings.TrimSpace(line))
			r = append(r, arr)
		}
	}

	// removing empty ScriptLines
	// var ret []ScriptLine
	// for _, l := range r {
	// 	if l.Character != "" && l.Line != "" {
	// 		ret = append(ret, l)
	// 	}
	// }

	return r
}

// func main() {
// 	fmt.Println(GetScript("/1.01_Pilot_(transcript)"))
// }
