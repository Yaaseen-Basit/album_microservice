package main

import (
	"encoding/json"
	"fmt"
)

type Album struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func main() {
	jsonData := `[{"name": "DS1", "id": "1"}, {"name": "DS2", "id": "2"}]`

	var albums []Album
	err := json.Unmarshal([]byte(jsonData), &albums)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Result:", albums)
}
