package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("使用方法: cleanup <压缩包文件名>")
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

type ArchiveFormat int

const (
	FormatUnknown ArchiveFormat = iota
	FormatZip
	FormatTar
	FormatTarGz
	FormatTarBz2
	FormatTarXz
)

func detectArchiveFormat(filePath string) ArchiveFormat {
	ext := strings.ToLower(filepath.Ext(filePath))
	baseName := strings.ToLower(filepath.Base(filePath))
	
	switch {
	case ext == ".zip":
		return FormatZip
	case ext == ".tar":
		return FormatTar
	case ext == ".gz" && strings.HasSuffix(baseName, ".tar.gz"):
		return FormatTarGz
	case ext == ".tgz":
		return FormatTarGz
	case ext == ".bz2" && strings.HasSuffix(baseName, ".tar.bz2"):
		return FormatTarBz2
	case ext == ".xz" && strings.HasSuffix(baseName, ".tar.xz"):
		return FormatTarXz
	default:
		return FormatUnknown
	}
}

func getTopLevelItems(archiveFilePath string) ([]string, error) {
	format := detectArchiveFormat(archiveFilePath)
	
	switch format {
	case FormatZip:
		return getTopLevelItemsFromZip(archiveFilePath)
	case FormatTar:
		return getTopLevelItemsFromTar(archiveFilePath)
	case FormatTarGz:
		return getTopLevelItemsFromTarGz(archiveFilePath)
	case FormatTarBz2:
		return getTopLevelItemsFromTarBz2(archiveFilePath)
	case FormatTarXz:
		return getTopLevelItemsFromTarXz(archiveFilePath)
	default:
		return nil, fmt.Errorf("不支持的压缩格式")
	}
}

func getTopLevelItemsFromZip(archiveFilePath string) ([]string, error) {
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

func getTopLevelItemsFromTar(archiveFilePath string) ([]string, error) {
	file, err := os.Open(archiveFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return extractTopLevelItemsFromTarReader(file)
}

func getTopLevelItemsFromTarGz(archiveFilePath string) ([]string, error) {
	file, err := os.Open(archiveFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	return extractTopLevelItemsFromTarReader(gzReader)
}

func getTopLevelItemsFromTarBz2(archiveFilePath string) ([]string, error) {
	file, err := os.Open(archiveFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bz2Reader := bzip2.NewReader(file)
	return extractTopLevelItemsFromTarReader(bz2Reader)
}

func getTopLevelItemsFromTarXz(archiveFilePath string) ([]string, error) {
	cmd := exec.Command("xz", "-dc", archiveFilePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("无法创建xz解压管道: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("无法启动xz解压: %v", err)
	}
	defer func() {
		stdout.Close()
		cmd.Wait()
	}()

	return extractTopLevelItemsFromTarReader(stdout)
}

func extractTopLevelItemsFromTarReader(reader io.Reader) ([]string, error) {
	tarReader := tar.NewReader(reader)
	topLevelItemsMap := make(map[string]bool)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if header.Name == "" {
			continue
		}

		pathParts := strings.Split(strings.Trim(header.Name, "/"), "/")
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
