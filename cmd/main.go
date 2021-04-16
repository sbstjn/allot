package main

import (
	"fmt"
	"os"

	allot "github.com/sdslabs/allot/pkg"
)

func main() {
	cmd := allot.New("revert <commits:integer> commits on <project:string> at (stage|prod)")
	match, err := cmd.Match("revert 12 commits on example at prod")

	if err == nil {
		fmt.Println("Request did not match command.")
		os.Exit(1)
	}
	commits, _ := match.Integer("commits")
	project, _ := match.String("project")
	env, _ := match.Match(2)

	fmt.Printf("Revert \"%d\" on \"%s\" at \"%s\"", commits, project, env)
}
