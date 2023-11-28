package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	err := fmt.Errorf("some error")
	errJSON, _ := json.Marshal(err)
	fmt.Println(string(errJSON))
}
