package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Command string

type Person struct {
	Name string
}

type Commit struct {
	Hash    string
	Author  Person
	Message string
	Files   []string
}

const (
	cmdConfig   Command = "config"
	cmdAdd      Command = "add"
	cmdLog      Command = "log"
	cmdCommit   Command = "commit"
	cmdCheckout Command = "checkout"
)

const (
	vcsDir     = "./vcs/"
	commitsDir = vcsDir + "commits/"
	configFile = vcsDir + "config.txt"
	indexFile  = vcsDir + "index.txt"
	logFile    = vcsDir + "log.txt"
)

var allowedCommands = []Command{
	cmdConfig,
	cmdAdd,
	cmdLog,
	cmdCommit,
	cmdCheckout,
}

var helpMsg = map[Command]string{
	cmdConfig:   "Get and set a username.",
	cmdAdd:      "Add a file to the index.",
	cmdLog:      "Show commit logs.",
	cmdCommit:   "Save changes.",
	cmdCheckout: "Restore a file.",
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		displayHelp()
		return
	}

	var cmd = Command(os.Args[1])
	var args = os.Args[2:]

	switch cmd {
	case cmdConfig:
		doConfig(args)
	case cmdAdd:
		doAdd(args)
	case cmdLog:
		doLog()
	case cmdCommit:
		doCommit(args)
	case cmdCheckout:
		doCheckout(args)
	default:
		fmt.Printf("'%s' is not a SVCS command.\n", cmd)
	}
}

func displayHelp() {
	fmt.Println("These are SVCS commands:")
	for _, command := range allowedCommands {
		fmt.Printf("%-8s  %s\n", command, helpMsg[command])
	}
}

func doConfig(args []string) {
	if len(args) == 1 {
		setUsername(args[0])
	}

	if username, err := getUsername(); err == nil {
		fmt.Printf("The username is %s.\n", username)
	} else {
		fmt.Println("Please, tell me who you are.")
	}
}

func setUsername(username string) {
	os.MkdirAll(vcsDir, os.ModePerm)

	file, err := os.Create(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(username)
}

func getUsername() (string, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	return scanner.Text(), scanner.Err()
}

func doAdd(args []string) {
	if len(args) != 1 {
		showTrackedFiles()
		return
	}

	os.MkdirAll(vcsDir, os.ModePerm)

	filename := args[0]

	if _, err := os.Stat(filename); err != nil {
		fmt.Printf("Can't find '%s'.\n", filename)
		return
	}

	file, err := os.OpenFile(indexFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(filename + "\n")

	fmt.Printf("The file '%s' is tracked.\n", filename)
}

func showTrackedFiles() {
	files := getTrackedFiles()
	if len(files) == 0 {
		fmt.Println("Add a file to the index.")
		return
	}

	fmt.Println("Tracked files:")
	for _, filename := range files {
		fmt.Println(filename)
	}
}

func getTrackedFiles() []string {
	file, err := os.Open(indexFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	files := make([]string, 0, 5)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	return files
}

func doLog() {
	commits := getCommits()
	if len(commits) == 0 {
		fmt.Println("No commits yet.")
		return
	}

	commits[0].show()

	for i := 1; i < len(commits); i++ {
		fmt.Println()
		commits[i].show()
	}
}

func (c *Commit) show() {
	fmt.Println("commit", c.Hash)
	fmt.Println("Author:", c.Author.Name)
	fmt.Println(c.Message)
}

func getCommits() []Commit {
	file, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		log.Fatal(err)
	}
	defer file.Close()

	var commits []Commit

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.SplitN(line, " ", 3)
		commits = append(commits, Commit{
			Hash:    data[0],
			Author:  Person{Name: data[1]},
			Message: data[2],
		})
	}

	slices.Reverse(commits)

	return commits
}

func doCommit(args []string) {
	if len(args) != 1 {
		fmt.Println("Message was not passed.")
		return
	}

	commit := Commit{Message: args[0]}

	commit.Files = getTrackedFiles()
	if len(commit.Files) == 0 {
		fmt.Println("Nothing to commit.")
		return
	}

	username, err := getUsername()
	if err != nil {
		fmt.Println("Please, tell me who you are.")
		return
	}

	commit.Author.Name = username
	if commit.Save() {
		fmt.Println("Changes are committed.")
	} else {
		fmt.Println("Nothing to commit.")
	}
}

func doCheckout(args []string) {
	if len(args) != 1 {
		fmt.Println("Commit id was not passed.")
		return
	}

	commitID := args[0]
	commitPath := filepath.Join(commitsDir, commitID)

	if _, err := os.Stat(commitPath); os.IsNotExist(err) {
		fmt.Println("Commit does not exist.")
		return
	}

	files, err := os.ReadDir(commitPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		srcPath := filepath.Join(commitPath, file.Name())
		destPath := file.Name()
		if err := restoreFile(destPath, srcPath); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Switched to commit %s.\n", commitID)
}

func restoreFile(dest, src string) error {
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, source); err != nil {
		return err
	}

	return nil
}

func (c *Commit) Save() (saved bool) {
	computeHash(c)
	if isNewChanges(c) {
		commitFiles(c)
		addLog(c)
		saved = true
	}
	return
}

func computeHash(c *Commit) {
	sha256Hash := sha256.New()
	for _, filename := range c.Files {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.Copy(sha256Hash, file); err != nil {
			log.Fatal(err)
		}

		file.Close()
	}

	c.Hash = hex.EncodeToString(sha256Hash.Sum(nil))
}

func isNewChanges(c *Commit) bool {
	path := filepath.Join(commitsDir, c.Hash)

	if _, err := os.Stat(path); err == nil {
		return false
	}

	return true
}

func commitFiles(c *Commit) {
	commitPath := filepath.Join(commitsDir, c.Hash)
	os.MkdirAll(commitPath, os.ModePerm) 

	for _, filename := range c.Files {
		srcPath := filename
		destPath := filepath.Join(commitPath, filename)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			fmt.Printf("File '%s' does not exist in commit '%s'.\n", filename, c.Hash)
			continue
		}
		copyFile(destPath, srcPath)
	}
}

func copyFile(dest, src string) {
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	source, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := io.Copy(file, source); err != nil {
		log.Fatal(err)
	}
}

func addLog(c *Commit) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("%s %s %s\n", c.Hash, c.Author.Name, c.Message))
}
