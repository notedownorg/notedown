package documents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocuments_Client(t *testing.T) {
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(dir, "testclient")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, client.documents, 1)
}

func loadtestDocuments_Client(count int, t *testing.T) {
	dir, err := generateTestData("client", count)
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(dir, "testclient")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, client.documents, count)
}

// Load test rather than benchmark as we are testing the ability to handle a large number of files not the speed (for now)
func TestDocuments_Client_Loadtest_10(t *testing.T)    { loadtestDocuments_Client(10, t) }
func TestDocuments_Client_Loadtest_100(t *testing.T)   { loadtestDocuments_Client(100, t) }
func TestDocuments_Client_Loadtest_1000(t *testing.T)  { loadtestDocuments_Client(1000, t) }
func TestDocuments_Client_Loadtest_10000(t *testing.T) { loadtestDocuments_Client(10000, t) }

// func TestDocuments_Client_Loadtest_50000(t *testing.T) { loadtestDocuments_Client(50000, t) }
