package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("git", "branch").Output()

	if err != nil {
		log.Fatalf("cmd.Output() failed with %s\n", err)
	}

	str := string(out)
	trimmedStr := strings.TrimRight(str, "\n")
	result := strings.Split(trimmedStr, "\n")

	for i := range result {
		fmt.Printf("\n%v - %v\n", i, result[i])
	}

	fmt.Println(len(result))
}
