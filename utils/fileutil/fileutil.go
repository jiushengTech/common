package fileutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FormatFileSize 格式化文件大小，使用二进制单位
func FormatFileSize(fileSize int64) string {
	const (
		B  = 1
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
		PB = TB * 1024
	)

	switch {
	case fileSize < KB:
		return fmt.Sprintf("%dB", fileSize)
	case fileSize < MB:
		return fmt.Sprintf("%.2fKB", float64(fileSize)/KB)
	case fileSize < GB:
		return fmt.Sprintf("%.2fMB", float64(fileSize)/MB)
	case fileSize < TB:
		return fmt.Sprintf("%.2fGB", float64(fileSize)/GB)
	case fileSize < PB:
		return fmt.Sprintf("%.2fTB", float64(fileSize)/TB)
	default:
		return fmt.Sprintf("%.2fPB", float64(fileSize)/PB)
	}
}

// CleanUp 删除指定路径的文件或目录，并根据 deleteDirectoryContents 参数删除该文件所在目录下的所有文件和文件夹
func CleanUp(path string, deleteDirectoryContents bool) error {

	// 检查路径是否存在
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("error checking path: %v", err)
	}

	// 如果是文件
	if !info.IsDir() {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("error removing file: %v", err)
		}
		fmt.Printf("File %s removed successfully\n", path)

		// 如果需要删除文件所在目录下的所有文件和文件夹
		if deleteDirectoryContents {
			dir := filepath.Dir(path)
			if err := removeDirectoryContents(dir); err != nil {
				return fmt.Errorf("error deleting directory contents: %v", err)
			}
			fmt.Printf("All files and subdirectories in directory %s have been removed.\n", dir)
		}
		return nil
	}

	// 如果是目录，删除目录下所有内容
	if err := removeDirectoryContents(path); err != nil {
		return fmt.Errorf("error deleting directory contents: %v", err)
	}

	// 删除空目录
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error removing directory: %v", err)
	}
	fmt.Printf("Directory %s removed successfully\n", path)

	return nil
}

// removeDirectoryContents 删除目录下的所有文件和子目录，但保留目录本身
func removeDirectoryContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("error removing %s: %v", path, err)
		}
		if entry.IsDir() {
			fmt.Printf("Directory %s removed successfully\n", path)
		} else {
			fmt.Printf("File %s removed successfully\n", path)
		}
	}

	return nil
}

// DownloadFile 根据给定的 URL 下载文件并保存到本地，返回本地文件路径
func DownloadFile(url, downloadDir string) (string, error) {
	// 解析文件名
	fileName := filepath.Base(url)
	if fileName == "" || fileName == "." || fileName == "/" {
		return "", fmt.Errorf("unable to extract file name from URL: %s", url)
	}

	// 设置默认下载目录
	if downloadDir == "" {
		downloadDir = "."
	}

	// 确保下载目录存在
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", fmt.Errorf("error creating directory %s: %v", downloadDir, err)
	}

	// 构建本地文件路径
	filePath := filepath.Join(downloadDir, fileName)

	// 下载文件内容
	data, err := DownloadFileToBytes(url)
	if err != nil {
		return "", err
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("error writing file %s: %v", filePath, err)
	}

	fmt.Printf("File downloaded successfully: %s\n", filePath)
	return filePath, nil
}

// DownloadFileToBytes 从指定URL下载文件，返回二进制数据
func DownloadFileToBytes(url string) ([]byte, error) {
	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// 设置User-Agent头
	req.Header.Set("User-Agent", "Go-FileUtil/1.0")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading file from URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file, HTTP status code: %d", resp.StatusCode)
	}

	// 读取响应体
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Printf("File downloaded to memory: %d bytes from %s\n", len(data), url)
	return data, nil
}

// DownloadFileWithProgressToBytes 从指定URL下载文件，返回二进制数据，支持进度回调
func DownloadFileWithProgressToBytes(url string, progressCallback func(downloaded, total int64)) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Go-FileUtil/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading file from URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file, HTTP status code: %d", resp.StatusCode)
	}

	// 获取文件大小
	contentLength := resp.ContentLength

	// 创建进度读取器
	var reader io.Reader = resp.Body
	if progressCallback != nil && contentLength > 0 {
		reader = &progressReader{
			reader:   resp.Body,
			total:    contentLength,
			callback: progressCallback,
		}
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Printf("File downloaded to memory: %d bytes from %s\n", len(data), url)
	return data, nil
}

// progressReader 实现带进度回调的读取器
type progressReader struct {
	reader     io.Reader
	total      int64
	downloaded int64
	callback   func(downloaded, total int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.downloaded += int64(n)
	if pr.callback != nil {
		pr.callback(pr.downloaded, pr.total)
	}
	return n, err
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory 检查路径是否是目录
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
