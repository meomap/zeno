// zeno is a command line tool to analyse ansible playbook dependencies
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/meomap/zeno/loader"
	"github.com/meomap/zeno/search"
)

func main() {
	var (
		filesIn = flag.String("files", "", "names of changed files from command 'git diff $BEFORE $AFTER --name-only'")
		debug   = flag.Bool("debug", false, "enable for verbose logging")
		pbsIn   = flag.String("playbooks", "", "comma separated list of playbooks to examined")
	)
	flag.Parse()

	// required args
	if *pbsIn == "" || *filesIn == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *debug == false {
		log.SetOutput(ioutil.Discard)
	}
	diffFiles := strings.Split(*filesIn, "\n")
	repoDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	lenDiffs := len(diffFiles)
	// construct absolute path for input files
	for i := 0; i < lenDiffs; i++ {
		diffFiles[i] = path.Join(repoDir, diffFiles[i])
	}
	pbFiles := strings.Split(*pbsIn, ",")
	lenPbs := len(pbFiles)
	log.Printf("Match against [%d] files", lenDiffs)
	log.Printf("Examine [%d] playbooks: %s\n", lenPbs, *pbsIn)

	ds := new(loader.FileLoader)
	var (
		matched bool
		out     []string
	)
	for i := 0; i < lenPbs; i++ {
		name := pbFiles[i]
		if matched, err = search.MatchPlaybook(name, diffFiles, repoDir, ds); err != nil {
			log.Fatal(err)
		} else if matched {
			out = append(out, name)
		}
	}
	fmt.Println(strings.Join(out, ","))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}
