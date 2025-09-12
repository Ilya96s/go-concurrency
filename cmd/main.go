package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"in-memory-kv/internal/compute"
	"in-memory-kv/internal/storage"
	"os"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	engine := storage.NewMemoryEngine()
	computeLayer := compute.NewCompute(engine, logger)

	fmt.Println("In-Memory KV database (Type EXIT to quit)")
	fmt.Println("query = set_command | get_command | del_command")
	fmt.Println("set_command = \"SET\" argument argument")
	fmt.Println("get_command = \"GET\" argument")
	fmt.Println("get_command = \"GET\" argument")
	fmt.Println("del_command = \"DEL\" argument")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		if line == "EXIT" {
			break
		}

		result := computeLayer.Handle(line)
		fmt.Println(result)
	}
}
