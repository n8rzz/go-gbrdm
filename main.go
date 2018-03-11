package main

import (
	"log"
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

var collection GitBranchCollection

func init() {
	run()
}

func main() {
	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "select-branches",
		Help: "",
		Func: func(c *ishell.Context) {
			branches := collection.RawnameList()
			choices := c.Checklist(branches,
				"Which branches would you like to remove ?",
				nil)
			out := func() (c []string) {
				for _, v := range choices {
					c = append(c, branches[v])
				}
				return
			}
			c.Println("Your choices are", strings.Join(out(), ", "))
		},
	})

	// run shell
	shell.Process("select-branches")
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
