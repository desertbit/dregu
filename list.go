/*
 * DREGU - Docker Registry Utility
 * Copyright (c) 2020 DesertBit
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/desertbit/grumble"
)

const (
	flagListPrefix = "prefix"
)

var cmdList = &grumble.Command{
	Name: "list",
	Help: "list all repositories with their versions matching the given name; usage: <repository-name>",
	Flags: func(f *grumble.Flags) {
		f.BoolL(flagListPrefix, false, "all repositories are listed whose name contains the given name as a prefix")
	},
	AllowArgs: true,
	Run:       runList,
}

func init() {
	App.AddCommand(cmdList)
}

func runList(ctx *grumble.Context) (err error) {
	// Check args.
	if len(ctx.Args) != 1 {
		return errors.New("invalid arguments, see help")
	}
	repoName := ctx.Args[0]
	prefix := ctx.Flags.Bool(flagListPrefix)

	// Query the repositories.
	reps, err := reg.Repositories()
	if err != nil {
		log.Fatalln(err)
	}

	// Collect all repos that match the name.
	matchingRepos := make([]string, 0, len(reps))
	for _, rep := range reps {
		if prefix {
			if strings.HasPrefix(rep, repoName) {
				matchingRepos = append(matchingRepos, rep)
			}
		} else if rep == repoName {
			matchingRepos = append(matchingRepos, rep)
		}
	}
	if len(matchingRepos) == 0 {
		msg := "name"
		if prefix {
			msg = "prefix"
		}

		fmt.Printf("no repository found with %s '%s'\n", msg, repoName)
		return
	}

	var tags []string

	// List all tags of each matching repo.
	for _, mr := range matchingRepos {
		tags, err = reg.Tags(mr)
		if err != nil {
			return
		}

		ts := "[ "
		for i, t := range tags {
			if i > 0 {
				ts += ", "
			}
			ts += t
		}
		ts += " ]"

		fmt.Printf("%s: %s\n", mr, ts)
	}
	return
}
