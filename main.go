package main

import (
	"bufio"
	"fmt"
	manager "main/manager"
	"net/url"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter some text (press Ctrl+D or Ctrl+Z to end):")
	scanner.Scan()
	// accept url from user
	givenUrl := scanner.Text()
	// parse the url
	url, err := url.Parse(givenUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	manager.Download(*url)
}
