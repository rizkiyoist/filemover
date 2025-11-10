package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: filemover <contents_dir> <folders_dir>")
		return
	}

	contentsDir := os.Args[1]
	foldersDir := os.Args[2]

	folders, err := ioutil.ReadDir(foldersDir)
	if err != nil {
		fmt.Printf("Error reading folders directory: %v\n", err)
		return
	}

	files, err := ioutil.ReadDir(contentsDir)
	if err != nil {
		fmt.Printf("Error reading contents directory: %v\n", err)
		return
	}

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}
		folderName := folder.Name()
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if strings.HasPrefix(fileName, folderName) {
				srcPath := filepath.Join(contentsDir, fileName)
				dstPath := filepath.Join(foldersDir, folderName, fileName)
				highlightsDir := filepath.Join(foldersDir, folderName, "highlights")
				if err := os.MkdirAll(highlightsDir, 0755); err != nil {
					fmt.Printf("Failed to create highlights folder: %v\n", err)
					continue
				}
				dstPath = filepath.Join(highlightsDir, fileName)
				err := os.Rename(srcPath, dstPath)
				if err != nil {
					// Try copy and delete for cross-drive moves
					copyErr := copyFile(srcPath, dstPath)
					if copyErr != nil {
						fmt.Printf("Failed to move %s: %v\n", fileName, err)
						continue
					}
					delErr := os.Remove(srcPath)
					if delErr != nil {
						fmt.Printf("Copied but failed to delete %s: %v\n", fileName, delErr)
						continue
					}
					fmt.Printf("Moved %s to %s (copy/delete)\n", fileName, dstPath)
				} else {
					fmt.Printf("Moved %s to %s\n", fileName, dstPath)
				}
			}
		}
	}

	fmt.Println("Done.")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Ensure destination folder exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}
