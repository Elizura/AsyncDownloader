package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter some text (press Ctrl+D or Ctrl+Z to end):")
	scanner.Scan()
	// accept url from user
	givenUrl := scanner.Text()
	fmt.Println("You entered:", givenUrl)
	// parse the url
	url, err := url.Parse(givenUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	fmt.Println("Scheme:", url.Scheme, url.Host, url.Path, url.RawQuery, url.Fragment)
}
