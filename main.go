package main

import (
	"encoding/json"
	"io"
	"opusEDB/economics/learnx"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	LearnXOrder := learnx.Order{}
	err = json.Unmarshal(data, &LearnXOrder)
	if err != nil {
		panic(err)
	}
	err = LearnXOrder.HandleOrder()
	if err != nil {
		panic(err)
	}
	io.WriteString(os.Stdout, "Order handled successfully")
}
