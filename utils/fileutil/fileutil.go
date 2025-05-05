package fileutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// CleanUp 删除指定路径的文件或目录，并根据 deleteDirectoryContents 参数删除该文件所在目录下的所有文件和文件夹
func CleanUp(path string, deleteDirectoryContents ...bool) error {
	// 如果没有传入 deleteDirectoryContents 参数，默认值为 false
	deleteDirContents := false
	if len(deleteDirectoryContents) > 0 && deleteDirectoryContents[0] {
		deleteDirContents = true
	}

	// 确保路径存在
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("error checking path: %v", err)
	}

	// 如果是文件，直接删除
	if !isDirectory(path) {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("error removing file: %v", err)
		}
		fmt.Printf("File %s removed successfully\n", path)

		// 如果需要删除文件所在目录下的所有文件和文件夹
		if deleteDirContents {
			dir := filepath.Dir(path)
			err := deleteDirectoryContentsInDir(dir)
			if err != nil {
				return fmt.Errorf("error deleting files and directories in directory: %v", err)
			}
			fmt.Printf("All files and subdirectories in directory %s have been removed.\n", dir)
		}

		return nil
	}

	// 如果是目录，删除目录下所有文件和文件夹
	err = deleteDirectoryContentsInDir(path)
	if err != nil {
		return fmt.Errorf("error deleting files and directories in directory: %v", err)
	}

	// 删除空目录
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("error removing directory: %v", err)
	}
	fmt.Printf("Directory %s removed successfully\n", path)

	return nil
}

// isDirectory 检查路径是否是一个目录
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// deleteDirectoryContentsInDir 删除目录下的所有文件和子目录
func deleteDirectoryContentsInDir(dir string) error {
	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 排除根目录，删除所有文件和子目录
		if filePath != dir {
			if info.IsDir() {
				// 删除子目录及其中的内容
				err := os.RemoveAll(filePath)
				if err != nil {
					return fmt.Errorf("error removing directory %s: %v", filePath, err)
				}
				fmt.Printf("Directory %s removed successfully\n", filePath)
			} else {
				// 删除文件
				err := os.Remove(filePath)
				if err != nil {
					return fmt.Errorf("error removing file %s: %v", filePath, err)
				}
				fmt.Printf("File %s removed successfully\n", filePath)
			}
		}
		return nil
	})

	return err
}

// DownloadFile 根据给定的 URL 下载文件并保存到本地，返回本地文件路径
func DownloadFile(url, downloadDir string) (string, error) {
	// 解析文件名，默认从 URL 获取文件名
	fileName := filepath.Base(url)
	if fileName == "" {
		return "", fmt.Errorf("unable to extract file name from URL: %s", url)
	}

	// 如果没有指定下载目录，则默认使用当前目录
	if downloadDir == "" {
		downloadDir = "."
	}

	// 确保下载目录存在
	err := os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating directory %s: %v", downloadDir, err)
	}

	// 构建本地文件路径
	filePath := filepath.Join(downloadDir, fileName)

	// 发送 GET 请求下载文件
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading file from URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file, HTTP status code: %d", resp.StatusCode)
	}

	// 创建本地文件
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer outFile.Close()

	// 将下载的内容写入文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error saving file %s: %v", filePath, err)
	}

	fmt.Printf("File downloaded successfully: %s\n", filePath)

	// 返回下载后的本地文件路径
	return filePath, nil
}
