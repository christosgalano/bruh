/*
TODO: add description
*/
package main

import (
	"fmt"
	"log"

	"github.com/christosgalano/bruh/internal/apiversions"
	"github.com/christosgalano/bruh/internal/bicep"
)

func main() {
	dir := "/Users/galano/Developer/Christos/Development/Go/bruh/parse"
	bicepDirectory, err := bicep.ParseDirectory(dir)
	if err != nil {
		log.Fatalf("parse error: %s", err)
	}
	fmt.Printf("Initial bicepDirectory:\n%s\n\n", bicepDirectory)

	err = apiversions.UpdateBicepDirectory(bicepDirectory)
	if err != nil {
		log.Fatalf("failed to update API versions: %s", err)
	}
	fmt.Printf("bicepDirectory after api versions:\n%s\n\n", bicepDirectory)

	err = bicep.UpdateDirectory(bicepDirectory, true, false)
	if err != nil {
		log.Fatalf("failed to update API versions: %s", err)
	}
	fmt.Printf("bicepDirectory after api versions:\n%s\n\n", bicepDirectory)
}
