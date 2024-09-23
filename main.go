package main

import (
	"fmt"
	"net/http"
    "flow/blocks"
)

func main() {
	handler := blocks.NewBlockHandler()

	http.HandleFunc("/process", handler.Process)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
