package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/abiosoft/ishell"
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

// RawnameList produces an array of `GitBranch#Rawname` values
func (gbc *GitBranchCollection) RawnameList() []string {
	var rl []string

	for i := range gbc.items {
		rl = append(rl, gbc.items[i].Rawname)
	}

	return rl
}

var (
	collection GitBranchCollection
	selection  []string
)

func main() {
	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "selectBranches",
		Help: "",
		Func: func(c *ishell.Context) {
			branches := collection.RawnameList()
			choices := c.Checklist(branches,
				"Which branches would you like to remove ?",
				nil)

			for _, v := range choices {
				selection = append(selection, branches[v])
			}

			shouldDelete := verifyRemoval(c, selection)

			if !shouldDelete {
				os.Exit(0)
			}

			processBranchRemovals(selection)
		},
	})

	// run shell
	shell.Process("selectBranches")
}

func verifyRemoval(c *ishell.Context, selection []string) bool {
	c.Println("Do you really want to remove: \n\n", strings.Join(selection, ", "))

	c.Print("\nyes/no (y/n): ")
	response := c.ReadLine()

	return response == "yes" || response == "y" || response == ""
}

func processBranchRemovals(selection []string) {
	for i := range selection {
		branch := selection[i]

		fmt.Printf("Removing branch: %v\n", branch)
		_, err := exec.Command("git", "branch", "-D", branch).Output()

		if err != nil {
			log.Fatalf("cmd.Output() failed with %s\n", err)
		}
	}
}

func init() {
	run()
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
