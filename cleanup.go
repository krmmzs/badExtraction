package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("使用方法: go run cleanup.go <压缩包文件名>")
		os.Exit(1)
	}

	archiveFilePath := os.Args[1]

	if !fileExists(archiveFilePath) {
		fmt.Printf("错误: 压缩包文件 '%s' 不存在\n", archiveFilePath)
		os.Exit(1)
	}

	topLevelItems, err := getTopLevelItems(archiveFilePath)
	if err != nil {
		fmt.Printf("错误: 无法读取压缩包内容: %v\n", err)
		os.Exit(1)
	}

	if len(topLevelItems) == 0 {
		fmt.Println("压缩包为空")
		return
	}

	existingItems := findExistingItems(topLevelItems)
	if len(existingItems) == 0 {
		fmt.Println("没有找到需要清理的文件或文件夹")
		return
	}

	fmt.Println("发现以下文件/文件夹可能来自压缩包解压:")
	for _, item := range existingItems {
		fmt.Printf("  - %s\n", item)
	}

	if !confirmDeletion() {
		fmt.Println("取消删除操作")
		return
	}

	deleteItems(existingItems)
	fmt.Println("清理完成")
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func getTopLevelItems(archiveFilePath string) ([]string, error) {
	zipReader, err := zip.OpenReader(archiveFilePath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	topLevelItemsMap := make(map[string]bool)

	for _, file := range zipReader.File {
		if file.Name == "" {
			continue
		}

		pathParts := strings.Split(strings.Trim(file.Name, "/"), "/")
		if len(pathParts) > 0 && pathParts[0] != "" {
			topLevelItemsMap[pathParts[0]] = true
		}
	}

	var topLevelItems []string
	for item := range topLevelItemsMap {
		topLevelItems = append(topLevelItems, item)
	}

	return topLevelItems, nil
}

func findExistingItems(topLevelItems []string) []string {
	var existingItems []string

	for _, item := range topLevelItems {
		if fileExists(item) {
			existingItems = append(existingItems, item)
		}
	}

	return existingItems
}

func confirmDeletion() bool {
	fmt.Print("确认删除这些文件和文件夹吗? (y/N): ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		response := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return response == "y" || response == "yes"
	}

	return false
}

func deleteItems(items []string) {
	for _, item := range items {
		err := os.RemoveAll(item)
		if err != nil {
			fmt.Printf("删除 '%s' 时出错: %v\n", item, err)
		} else {
			fmt.Printf("已删除: %s\n", item)
		}
	}
}
