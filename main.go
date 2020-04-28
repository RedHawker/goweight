package main

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"

	"github.com/RedHawker/goweight/pkg"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	// jsonOutput overrides other output options
	jsonOutput = kingpin.Flag("json", "Output json (overrides other output options)").Short('j').Bool()

	byteOutput = kingpin.Flag("bytes", "Output size in bytes, instead of human-readable").Short('b').Default("false").Bool()
	excludeInternal = kingpin.Flag("exclude-internal", "Exclude Go internal packages").Short('x').Default("false").Bool()
	includeSum = kingpin.Flag("sum", "Include sum line in output").Short('s').Default("false").Bool()

	buildTags  = kingpin.Flag("tags", "Build tags").String()
	packages   = kingpin.Arg("packages", "Packages to build").String()
)

// Grabbed from vcs.go, don't know if it's perfect for what I'm trying to achieve
func isInternal(importPath string) (bool) {
	slash := strings.Index(importPath, "/")
	if slash < 0 {
		slash = len(importPath)
	}
	host := importPath[:slash]
	if !strings.Contains(host, ".") {
		return true
	}
	return false
}

func main() {
	kingpin.Version(fmt.Sprintf("%s (%s)", version, commit))
	kingpin.Parse()
	weight := pkg.NewGoWeight()
	if *buildTags != "" {
		weight.BuildCmd = append(weight.BuildCmd, "-tags", *buildTags)
	}
	if *packages != "" {
		weight.BuildCmd = append(weight.BuildCmd, *packages)
	}

	work := weight.BuildCurrent()
	modules := weight.Process(work)

	if *jsonOutput {
		m, _ := json.Marshal(modules)
		fmt.Print(string(m))
	} else {
		var total uint64 = 0
		for _, module := range modules {
			if module.Name != "runtime" && *excludeInternal {
				if isInternal(module.Name) {
					continue
				}
			}
			total += module.Size
			if !*byteOutput {
				fmt.Printf("%8s %s\n", module.SizeHuman, module.Name)
			} else {
				fmt.Printf("%10d %s\n", module.Size, module.Name)
			}
		}
		if *includeSum {
			if !*byteOutput {
				fmt.Printf("%8s Sum\n", humanize.Bytes(total))
			} else {
				fmt.Printf("%10d Sum\n", total)
			}
		}
	}
}
