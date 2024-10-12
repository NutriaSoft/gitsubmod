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
	log.Println(newSub)
	collector.AddSubmodule(newSub, submodules)
	err = collector.SaveSubmodulesToFile(submodules)
	if err != nil {
		log.Println(err)
	}
	updateSub := &pb.Submodule{
		Url:    "https://github.com/new/sub",
		Branch: "main",
		Name:   "newsub",
	}
	updated := collector.UpdateSubmodule(submodules, "newsub", updateSub)
	log.Println("Updated: ", updated)
	submodule, found := collector.FindSubmodule(submodules, "newsub")
	if found {
		log.Println(submodule)
	}
	deleted := collector.DeleteSubmodule(submodules, "newsub")
	log.Println("DELETED: ", deleted)
}
