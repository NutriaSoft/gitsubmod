package collector

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	pb "submoduleop/protos"

	"google.golang.org/protobuf/proto"
)

const (
	FILENAME = "submodules"
)

func GetHomeLocation() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	subDir := fmt.Sprintf("%s/.submoduleop/", dirname)
	if _, err := os.Stat(subDir); err != nil {
		if os.IsNotExist(err) {
			log.Println(filepath.Dir(subDir))
			err := os.MkdirAll(filepath.Dir(subDir), 0777)
			if err != nil {
				return "", fmt.Errorf("Cannot create: %s", subDir)
			}
		}
	}
	location := fmt.Sprintf("%s%s", subDir, FILENAME)
	return location, nil
}

func AddSubmodule(newSubModule *pb.Submodule, submodules *pb.SubmoduleList) {
	if submodules == nil {
		submodules = &pb.SubmoduleList{}
	}

	if submodules.Submodules != nil {
		_, found := FindSubmodule(submodules, newSubModule.Name)
		if found {
			return
		}
	}

	if submodules.Submodules == nil {
		submodules.Submodules = make([]*pb.Submodule, 0)
	}

	submodules.Submodules = append(submodules.Submodules, newSubModule)
}

func UpdateSubmodule(submodules *pb.SubmoduleList, name string, updatedSubmodule *pb.Submodule) bool {
	for i, submodule := range submodules.Submodules {
		if submodule.Name == name {
			submodules.Submodules[i] = updatedSubmodule
			return true
		}
	}
	return false
}

func DeleteSubmodule(submodules *pb.SubmoduleList, name string) bool {
	for i, submodule := range submodules.Submodules {
		if submodule.Name == name {
			submodules.Submodules = append(submodules.Submodules[:i], submodules.Submodules[i+1:]...)
			return true
		}
	}
	return false
}

func FindSubmodule(submodules *pb.SubmoduleList, name string) (*pb.Submodule, bool) {
	for _, submodule := range submodules.Submodules {
		if submodule.Name == name {
			return submodule, true
		}
	}
	return nil, false
}

func SaveSubmodulesToFile(submodules *pb.SubmoduleList) error {
	location, err := GetHomeLocation()
	if err != nil {
		return err
	}
	data, err := proto.Marshal(submodules)
	if err != nil {
		return err
	}
	return os.WriteFile(location, data, 0644)
}

func LoadSubmodulesFromFile() (*pb.SubmoduleList, error) {
	location, err := GetHomeLocation()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(location)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}
	submodules := &pb.SubmoduleList{}
	err = proto.Unmarshal(data, submodules)
	if err != nil {
		return nil, err
	}
	return submodules, nil
}
