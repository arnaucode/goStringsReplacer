package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const rawFolderPath = "./originalFiles"
const newFolderPath = "./newFiles"

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func replaceFromTable(table map[string]string, originalContent string) string {
	var newContent string
	newContent = originalContent
	for i := 0; i < len(table); i++ {
		//first, get the map keys
		var keys []string
		for key, _ := range table {
			keys = append(keys, key)
		}
		//now, replace the keys with the values
		for i := 0; i < len(keys); i++ {
			newContent = strings.Replace(newContent, keys[i], table[keys[i]], -1)
		}
	}
	return newContent

}
func readConfigTable(path string) map[string]string {
	table := make(map[string]string)
	f, err := os.Open(path)
	check(err)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if currentLine != "" {
			val1 := strings.Split(currentLine, " ")[0]
			val2 := strings.Split(currentLine, " ")[1]
			table[val1] = val2
		}
	}
	return table
}

func readFile(folderPath string, filename string) string {
	dat, err := ioutil.ReadFile(folderPath + "/" + filename)
	check(err)
	return string(dat)
}

func writeFile(path string, newContent string) {
	err := ioutil.WriteFile(path, []byte(newContent), 0644)
	check(err)
}

func parseDir(folderPath string, newDir string, table map[string]string) {
	files, _ := ioutil.ReadDir(folderPath)
	for _, f := range files {
		fileNameSplitted := strings.Split(f.Name(), ".")
		if len(fileNameSplitted) == 1 {
			newDir := newDir + "/" + f.Name()
			oldDir := rawFolderPath + "/" + f.Name()
			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				_ = os.Mkdir(newDir, 0700)
			}
			parseDir(oldDir, newDir, table)
		} else {
			fileContent := readFile(folderPath, f.Name())
			newContent := replaceFromTable(table, fileContent)
			writeFile(newDir+"/"+f.Name(), newContent)
			fmt.Println(newDir + "/" + f.Name())
		}
	}
}

func main() {
	table := readConfigTable("replacerConfig.txt")
	fmt.Println(table)
	parseDir(rawFolderPath, newFolderPath, table)
}
