package main

import (
	"bufio"
	"fmt"
	"in-memory-kv/compute"
	"in-memory-kv/storage"
	"os"
)

func main() {
	engine := storage.NewMemoryEngine()
	computeLayer := compute.NewCompute(engine)

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
