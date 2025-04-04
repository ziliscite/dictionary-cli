package main

import "github.com/ziliscite/dictionary-cli/internal"

func main() {
	result, err := internal.Search("rain")
	if err != nil {
		return
	}

	internal.DisplayJisho(result)
}
