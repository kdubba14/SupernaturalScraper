package parsers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// GetScript getting the page with script
func GetScript(r string) string {
	f, err := http.Get(r)
	if err != nil {
		fmt.Println("========ERROR FETCHING", err)
		os.Exit(1)
	}
	defer f.Body.Close()
	fmt.Println("Reading...")
	transcriptPage, err := ioutil.ReadAll(f.Body)
	if err != nil {
		fmt.Printf("========ERROR READING: \n%s\n", err)
		os.Exit(1)
	}

	return string(transcriptPage)
}

// func main() {
// 	fmt.Println(GetScript("/1.01_Pilot_(transcript)"))
// }
