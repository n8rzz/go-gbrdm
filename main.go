package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// GitBranch is a representation of a single git branch
type GitBranch struct {
	Rawname string
}

// GitBranchCollection is a Collection object of `GitBranch` instances
type GitBranchCollection struct {
	items []GitBranch
}

// AddItem `GitBranch` instance to collection
func (gbc *GitBranchCollection) AddItem(gitBranch GitBranch) {
	gbc.items = append(gbc.items, gitBranch)
}

var collection GitBranchCollection

func init() {
	run()
}

func main() {
	fmt.Printf("%v \n\n\n", collection.items)
	fmt.Print("Now were ready...")
}

func run() {
	out, err := exec.Command("git", "branch").Output()

	if err != nil {
		log.Fatalf("cmd.Output() failed with %s\n", err)
	}

	branchList := readLocalGitBranchList(out)

	hydrateBranchCollection(branchList)
}

func readLocalGitBranchList(cmdOut []byte) []string {
	str := string(cmdOut)
	trimmedStr := strings.TrimRight(str, "\n")
	localGitBranchList := strings.Split(trimmedStr, "\n")

	return localGitBranchList
}

func hydrateBranchCollection(branchList []string) {
	for i := range branchList {
		s := strings.TrimLeft(branchList[i], " ")
		gb := GitBranch{s}

		collection.AddItem(gb)
	}
}
