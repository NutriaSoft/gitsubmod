package main

import (
	"log"
	"submoduleop/collector"
	pb "submoduleop/protos"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	submodules, err := collector.LoadSubmodulesFromFile()
	if err != nil {
		log.Println(err)
	}
	newSub := &pb.Submodule{
		Url:    "https://github.com/new/sub",
		Branch: "master",
		Name:   "newsub",
	}
	collector.AddSubmodule(newSub, submodules)
	collector.SaveSubmodulesToFile(submodules)
}
