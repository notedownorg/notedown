package workspace

// import (
// "fmt"
// "testing"
// "time"
//
// "github.com/stretchr/testify/assert"
// )
//
// func TestWorkspace_DailyNotePath(t *testing.T) {
// 	tmp := copyTestData(t, "daily-note")
// 	fmt.Println("Created temp dir: ", tmp)
//
// 	ws, err := New(tmp)
// 	assert.NoError(t, err)
// 	time.Sleep(1 * time.Second) // remove once we have a way to wait for the initial state to be built
//
// 	// Create a new daily note
// 	date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
// 	path, err := ws.DailyNotePath(date)
// 	assert.NoError(t, err)
//
// 	// Check the daily note was Created
// 	assert.FileExists(t, path)
// }
