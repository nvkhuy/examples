package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// Build info
var (
	Env              = ""
	Name             = ""
	BuildDate        = ""
	BuildNumber      = ""
	Version          = ""
	GitCommit        = ""
	GitCommitFull    = ""
	GitCommitMsg     = ""
	GitCommitMsgFull = ""
	GitBranch        = ""
	GitState         = ""
	GitSummary       = ""
)

type buildInfo struct {
	Env              string `json:"env"`
	Name             string `json:"name"`
	GitCommit        string `json:"git_commit"`
	GitCommitFull    string `json:"git_commit_full"`
	GitCommitMsg     string `json:"git_commit_msg"`
	GitCommitMsgFull string `json:"git_commit_msg_full"`
	GitBranch        string `json:"git_branch"`
	GitState         string `json:"git_state"`
	GitSummary       string `json:"git_summary"`
	BuildDate        string `json:"build_date"`
	Version          string `json:"version"`
}

func (i *buildInfo) print() {
	data, err := json.MarshalIndent(&i, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func showInfo() {
	var i = buildInfo{
		Env:              Env,
		Name:             Name,
		GitCommit:        GitCommit,
		GitCommitFull:    GitCommitFull,
		GitCommitMsg:     GitCommitMsg,
		GitCommitMsgFull: GitCommitMsgFull,
		GitBranch:        GitBranch,
		GitState:         GitState,
		GitSummary:       GitSummary,
		BuildDate:        BuildDate,
		Version:          Version,
	}

	i.print()
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "info",
	Short: "App Infomartion",
	Run: func(cmd *cobra.Command, args []string) {
		showInfo()
	},
}
