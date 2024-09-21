package writer_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liamawhite/nl/pkg/workspace/documents/writer"
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
		if err := os.MkdirAll(dir, 0777); err != nil {
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

type Document struct {
	writer.Document
	Contents []byte
}

func loadDocument(t *testing.T, root string, path string) Document {
	contents, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()
	hash.Write(contents)

	return Document{Document: writer.Document{Path: strings.TrimPrefix(path, root), Hash: fmt.Sprintf("%x", hash.Sum(nil))}, Contents: contents}
}

type stringer struct {
	text string
}

func (s stringer) String() string {
	return s.text
}

func Text(text string) fmt.Stringer {
	return stringer{text}
}
