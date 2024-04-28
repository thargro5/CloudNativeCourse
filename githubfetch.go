package githubfetch

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// GrabGame retrieves specific lines from a text file hosted on GitHub.
func GrabGame(username, repo, filePath string) error {
	// Append the desired file extension (.cpp in this case)
	filePath = fmt.Sprintf("%s.cpp", filePath)

	// GitHub raw file URL
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", username, repo, filePath)

	// Send HTTP GET request to download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %s", err)
	}
	defer resp.Body.Close()

	// Create a new file to write the downloaded content
	out, err := os.Create("downloaded-file.cpp")
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer out.Close()

	// Copy the HTTP response body (file content) to the new file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %s", err)
	}

	fmt.Println("File downloaded successfully!")
	return nil
}
