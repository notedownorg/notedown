package documents

import (
	"fmt"
	"os"
	"path"

	cp "github.com/otiai10/copy"
)

func setupTestDir(name string) (string, error) {
	// If we're running in a CI environment, we dont want to create temp directories
	// This ensures we can store the artifacts for debugging
	dir := os.Getenv("GITHUB_WORKSPACE")
	if dir == "" {
		var err error
		dir, err = os.MkdirTemp("", fmt.Sprintf("nl-%v-", name))
		if err != nil {
			return "", err
		}
	} else {
		dir = fmt.Sprintf("%v/testdata/%v", dir, name)
		if err := os.Mkdir(dir, 0644); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func copyTestData(name string) (string, error) {
	dir, err := setupTestDir(name)
	if err != nil {
		return "", err
	}
	if err := cp.Copy("testdata/workspace", dir); err != nil {
		return "", err
	}
	return dir, nil
}

func generateTestData(name string, fileCount int) (string, error) {
	dir, err := setupTestDir(name)
	if err != nil {
		return "", err
	}
	for i := 0; i < fileCount; i++ {
		content := fmt.Sprintf("# Test Document %v", i) // maybe put more meaningful content here
		if err := writeFile(dir, fmt.Sprintf("%v.md", i), content); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func writeFile(dir string, name string, content string) error {
	return os.WriteFile(path.Join(dir, name), []byte(content), 0644)
}
