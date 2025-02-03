package main

import (
	"fmt"
	"os"
	"path/filepath"
	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookBody struct {
  Hostname string `json:"hostname"`
  TargetFolder string `json:"target_folder"`
  FileCount int `json:"file_count"`
  FolderCount int `json:"folder_count"`
}

func getDownloadFolderPath() string {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return ""
  }
  return filepath.Join(homeDir, "Downloads")
}

func main() {
  downloadsFolderPath := getDownloadFolderPath()
  if downloadsFolderPath == "" {
    panic("Failed to get downloads folder path")
  }

  files, err := os.ReadDir(downloadsFolderPath)
  if err != nil {
    panic(err)
  }

  folderCount := 0
  fileCount := 0

  for _, file := range files {
    if file.IsDir() {
      folderCount++
    } else {
      fileCount++
    }
    // fmt.Printf(
    //   "Name: %-20s IsDir: %-6v, Type: %v\n",
    //   file.Name(),
    //   file.IsDir(),
    //   file.Type(),
    // )
  }

  hostname, err := os.Hostname()
  if err != nil {
    fmt.Println("Failed to get hostname")
    return
  }
  url := os.Getenv("WEBHOOK_URL")

  data := WebhookBody{
    Hostname: hostname,
    TargetFolder: downloadsFolderPath,
    FileCount: fileCount,
    FolderCount: folderCount,
  }

  jsonData, err := json.Marshal(data)
  if err != nil {
    panic(err)
  }

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
  if err != nil {
    panic(err)
  }
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
}
