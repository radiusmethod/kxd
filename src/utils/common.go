package utils

import (
	"fmt"
	"log"
	"os"
)

func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func WriteFile(config, loc string) {
	s := []byte("")
	if config != "unset" {
		s = []byte(config)
	}
	err := os.WriteFile(fmt.Sprintf("%s/.kxd", loc), s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func CheckError(err error) {
	if err.Error() == "^D" {
		// https://github.com/manifoldco/promptui/issues/179
		log.Fatalf("<Del> not supported")
	} else if err.Error() == "^C" {
		os.Exit(1)
	} else {
		log.Fatal(err)
	}
}

func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v\n", err)
	}
	return homeDir
}
