package command_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/justasbieliauskas/drivemv/command"
	"github.com/justasbieliauskas/drivemv/fs"
)

func TestUploadFresh(t *testing.T) {
	filename := "drivemv-upload-new-test.txt"
	content := "foo"
	file, err := createFile(filename, content)
	if err != nil {
		t.Fatal("Error while creating a file with content", err)
	}
	file.Close()
	defer os.Remove(filename)
	command := command.New()
	err = command.Run([]string{filename, "/"})
	if err != nil {
		t.Fatal("Error while executing `drivemv`", err)
	}
	root, err := fs.NewRoot(os.Environ())
	if err != nil {
		t.Fatal("Error while creating a drive root", err)
	}
	gFile, err := root.GetFileByName(filename)
	if err != nil {
		t.Fatal("Error while getting file in drive", err)
	}
	defer gFile.Delete()
	gContent, err := gFile.Content()
	if err != nil {
		t.Fatal("Error while getting uploaded file's content", err)
	}
	if gContent != content {
		t.Errorf("Files do not match!\nExpected:\n%s\nGot:\n%s\n", content, gContent)
	}
}

func TestOverwriteExisting(t *testing.T) {
	filename := "drivemv-update-existing-test.txt"
	firstStr := "foo"
	secondStr := "bar"
	content := firstStr + secondStr
	file, err := createFile(filename, firstStr)
	if err != nil {
		t.Fatal("Error while creating a file with first string", err)
	}
	defer os.Remove(filename)
	defer file.Close()
	command := command.New()
	err = command.Run([]string{filename, "/"})
	if err != nil {
		t.Fatal("Error while executing drivemv", err)
	}
	root, err := fs.NewRoot(os.Environ())
	if err != nil {
		t.Fatal("Error while creating a drive root", err)
	}
	gFirstFile, err := root.GetFileByName(filename)
	if err != nil {
		t.Fatal("Error while getting file in drive the first time", err)
	}
	defer gFirstFile.Delete()
	gContent, err := gFirstFile.Content()
	if err != nil {
		t.Fatal("Error while getting file's content the first time", err)
	}
	if gContent != firstStr {
		t.Errorf("Files do not match!\nExpected:\n%s\nGot:\n%s\n", firstStr, gContent)
	}
	_, err = file.WriteString(secondStr)
	if err != nil {
		t.Fatal("Error while writing second string to file", err)
	}
	err = command.Run([]string{filename, "/"})
	if err != nil {
		t.Fatal("Error while executing drivemv the second time", err)
	}
	gSecondFile, err := root.GetFileByName(filename)
	if err != nil {
		t.Fatal("Error while getting file in drive the second time", err)
	}
	defer gSecondFile.Delete()
	gContent, err = gSecondFile.Content()
	if err != nil {
		t.Fatal("Error while getting file's content the second time", err)
	}
	if gContent != content {
		t.Errorf("Second time files do not match!\nExpected:\n%s\nGot:\n%s\n", firstStr, gContent)
	}
}

func TestVarsMissing(t *testing.T) {
	command := command.New()
	command.Env = []string{
		"DRIVE_CLIENT_ID=1h23j5js4jd3rsd6sj57",
		"DRIVE_PROJECT_ID=projectid",
		// missing DRIVE_CLIENT_SECRET
		"DRIVE_ACCESS_TOKEN=d79sg78s789789",
		"DRIVE_REFRESH_TOKEN=d8b8fgr8fb8rb7",
		"DRIVE_TOKEN_EXPIRY=2019-03-23T21:56:46.085692+02:00",
	}
	err := command.Run([]string{"source", "target"})
	if err == nil {
		t.Error("Command should have failed, but did not")
	}
}

func createFile(name, content string) (*os.File, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Error while creating a file", err))
	}
	_, err = file.WriteString(content)
	if err != nil {
		file.Close()
		os.Remove(name)
		return nil, errors.New(fmt.Sprintln("Error while writing to file", err))
	}
	return file, nil
}
