package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Create clientBinaries.json taking into account .tar.gz and .zip archives in current directory.

type gethArchiveMeta struct {
	Download struct {
		URL  string `json:"url"`
		Type string `json:"type"`
		MD5  string `json:"md5"`
		Bin  string `json:"bin"`
	} `json:"download"`
	Bin      string `json:"bin"`
	Commands struct {
		Sanity struct {
			Args   []string `json:"args"`
			Output []string `json:"output"`
		} `json:"sanity"`
	} `json:"commands"`
}

// Keep swarm config out.
// To think: should platforms be hardcoded?
type mistConfig struct {
	Clients struct {
		Geth struct {
			Version   string `json:"version"`
			Platforms struct {
				Linux struct {
					X64  gethArchiveMeta `json:"x64"`
					IA32 gethArchiveMeta `json:"ia32"`
				} `json:"linux"`
				Mac struct {
					X64 gethArchiveMeta `json:"x64"`
				} `json:"mac"`
				Win struct {
					X64  gethArchiveMeta `json:"x64"`
					IA32 gethArchiveMeta `json:"ia32"`
				} `json:"win"`
			} `json:"platforms"`
		} `json:"Geth"`
	} `json:"clients"`
}

func doMist(cmdline []string) {
	var (
		URL        = flag.String("url", "", `Destination from where to download binaries`)
		path       = flag.String("path", "./", `Path to archives with binaries`)
		version    = flag.String("version", "1.0", `Version of binaries`)
		binaryName = flag.String("binary", "geth", `Base name for binary, e.g. geth`)
	)
	flag.CommandLine.Parse(cmdline)

	files, err := ioutil.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	var config mistConfig
	config.Clients.Geth.Version = *version

	for _, file := range files {
		fn := file.Name()
		if strings.HasSuffix(fn, ".zip") || strings.HasSuffix(fn, ".tar.gz") {
			archiveMeta, err := extractArchiveMeta(*path, fn, *URL, *binaryName)
			if err != nil {
				log.Printf("Skipping file %s, reason: %s\n", fn, err.Error())
				continue
			}
			archiveMeta.fillInCommands(*version, *binaryName)
			if strings.Contains(fn, "linux") {
				if strings.Contains(fn, "amd64") {
					config.Clients.Geth.Platforms.Linux.X64 = archiveMeta
				} else if strings.Contains(fn, "386") {
					config.Clients.Geth.Platforms.Linux.IA32 = archiveMeta
				}
			}
			if strings.Contains(fn, "windows") {
				if strings.Contains(fn, "amd64") {
					config.Clients.Geth.Platforms.Win.X64 = archiveMeta
				} else if strings.Contains(fn, "386") {
					config.Clients.Geth.Platforms.Win.IA32 = archiveMeta
				}
			}
			if strings.Contains(fn, "darwin") && strings.Contains(fn, "amd64") {
				config.Clients.Geth.Platforms.Mac.X64 = archiveMeta
			}
		}
	}

	encoded, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("clientBinaries.json", encoded, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)
}

func extractArchiveMeta(path, fn, URL, binaryName string) (meta gethArchiveMeta, err error) {
	binaryNames, archiveType, md5, err := InvestigateArchive(path+fn, binaryName)
	if err != nil {
		return
	}
	meta.Download.URL = URL + "/" + fn
	meta.Download.Type = archiveType
	meta.Download.MD5 = md5
	meta.Download.Bin = binaryNames[1]
	meta.Bin = binaryNames[0]
	return
}

func (meta *gethArchiveMeta) fillInCommands(version, binaryName string) {
	meta.Commands.Sanity.Args = []string{"version"}
	meta.Commands.Sanity.Output = []string{binaryName, version}
}
