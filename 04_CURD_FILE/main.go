package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	const filepath string = "./Logs.txt"

	// 	Common flags:
	// os.O_RDONLY → read only
	// os.O_WRONLY → write only
	// os.O_RDWR → read & write
	// os.O_APPEND → append to file
	// os.O_CREATE → create file if not exists
	// os.O_TRUNC → truncate (clear) file

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("1. Read Logs file")
		fmt.Println("2. Write to Logs file")
		fmt.Println("3. Delete Logs file")
		fmt.Println("4. List all files")
		fmt.Println("5. Create a new file")
		fmt.Println("6. Create a new folder")
		fmt.Println("7. Create new nested folder")
		fmt.Println("Enter 'exit' to exit.")
		fmt.Println("Choose a option")
		input, err := reader.ReadString('\n') // reads until Enter
		if err != nil {
			return
		}
		input = strings.TrimSpace(input) // remove \n

		switch input {
		case "1":
			ReadAndPrintFile(filepath)
		case "2":
			fmt.Printf("what do you want to write?")
			cmd, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			CreateAndAppendFile(filepath, cmd)
		case "3":
			DeleteFile(filepath)
		case "4":
			ListAllFiles()
		case "5":
			CreateNewFile()
		case "6":
			CreateNewFolder()
		case "7":
			CreateNewNestedFolders()
		case "exit":
			defer fmt.Printf("exited.")
			return
		default:
			fmt.Println("invalid choice, choose agagin >")
		}
	}

}

func ReadAndPrintFile(filepath string) {
	// file, err := os.OpenFile(filepath, os.O_RDONLY, 0644|0755)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("file not exists.")
		// panic(err)
		return
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func CreateNewFile() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter file name: ")
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("something wrong happened")
		return
	}
	filename = strings.TrimSpace(filename)

	_, err = os.OpenFile(filename, os.O_CREATE, 0644)
	if err != nil {
		return
	}
}

func CreateAndAppendFile(filepath, content string) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// file.WriteString(content)
	writer := bufio.NewWriter(file)
	writer.WriteString(content)
	defer writer.Flush()
}

func DeleteFile(filepath string) {
	// Check if file exists
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist:", filepath)
		return
	}

	err = os.Remove(filepath)
	if err != nil {
		fmt.Printf("Error Occured: %s\n", err)
		os.Create(filepath)
		fmt.Println("File Overwritten.")
		return
	}
	fmt.Println("file successfully deleted")
}

func ListAllFiles() {
	dir := "."

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, entry := range entries {

		info, err := entry.Info() // get FileInfo
		if err != nil {
			return
		}

		if entry.IsDir() {
			fmt.Printf("[DIR ] %-10s Size: %d bytes, Modified: %v\n", info.Name(), info.Size(), info.ModTime())
		} else {
			fmt.Printf("[FILE] %-10s Size: %d bytes, Modified: %v\n", info.Name(), info.Size(), info.ModTime())
		}
	}
}

func CreateNewFolder() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter folder name: ")
	folderName, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	folderName = strings.TrimSpace(folderName)

	err = os.Mkdir(folderName, 0755)
	if err != nil {
		return
	}
	fmt.Println("Folder created successfully.")
}

func CreateNewNestedFolders() {
	fmt.Print("Enter folder names (seperate folders with '/'):")
	reader := bufio.NewReader(os.Stdin)
	folders, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	err = os.MkdirAll(folders, 0755)
	if err != nil {
		return
	}

	fmt.Print("Folders created successfully.")
}
