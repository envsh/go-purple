package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/gops/agent"
	"github.com/spf13/cobra"
)

func main() {
	fmt.Println(getBuildInfo(true))
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalln(err)
	}

	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func main_grp2mob() {
	tvm := newToxVM()
	tvm.run()
}

func main_tyxdo() {
	ti := tyxdoToxNew()
	ti.run()
}

///
var RootCmd = &cobra.Command{
	Use: "xbts", Long: "tox bot collects, like grp2mob, tyxdo",
}

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	RootCmd.AddCommand(grp2mobCmd, tyxdoCmd)
}

var grp2mobArgs = &struct{}{}
var grp2mobCmd = &cobra.Command{
	Use: "grp2mob", Short: "relay group message to friend message for mobile usage.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("heheh", args)
		// main_grp2mob()
	},
}

var tyxdoArgs = &struct{}{}
var tyxdoCmd = &cobra.Command{
	Use: "tyxdo", Short: "make effort to keep toxync's group stable online.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("heheh", args)
		main_tyxdo()
	},
}

///
// build info
var GitCommit, GitBranch, GitState, GitSummary, BuildDate, Version string

func getBuildInfo(full bool) string {
	trim := func(s string) string {
		if strings.HasPrefix(s, "GOVVV-") {
			return s[6:]
		}
		return s
	}
	commit := trim(GitCommit)
	branch := trim(GitBranch)
	// state := trim(GitState)
	summary := trim(GitSummary)
	date := trim(BuildDate)
	version := trim(Version)

	if full {
		return fmt.Sprintf("govvv: v%s branch:%s git:%s build:%s summary:%s, ",
			version, branch, commit, date, summary)
	}
	return fmt.Sprintf("govvv: v%s git:%s build:%s", version, commit, date)
}
