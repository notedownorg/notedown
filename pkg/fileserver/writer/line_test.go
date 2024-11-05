// Copyright 2024 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package writer_test

import (
	"fmt"
	"testing"

	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/stretchr/testify/assert"
)

func TestLine_AddLine(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		line     fmt.Stringer
		lines    []string
		checksum string
		want     []string
		wantErr  bool
	}{
		{
			name:     "Add line at beginning",
			number:   writer.AT_BEGINNING,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"new line", "line 1", "line 2", "line 3"},
		},
		{
			name:     "Add line at beginning with empty checksum",
			number:   writer.AT_BEGINNING,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "",
			want:     []string{"new line", "line 1", "line 2", "line 3"},
		},
		{
			name:     "Add line at end",
			number:   writer.AT_END,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 1", "line 2", "line 3", "new line"},
		},
		{
			name:     "Add line at end with empty checksum",
			number:   writer.AT_END,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "",
			want:     []string{"line 1", "line 2", "line 3", "new line"},
		},
		{
			name:     "Add line at > number of lines",
			number:   999,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 1", "line 2", "line 3", "new line"},
		},
		{
			name:     "Add line at 0",
			number:   0,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"new line", "line 1", "line 2", "line 3"},
		},
		{
			name:     "Add line at 1",
			number:   1,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"new line", "line 1", "line 2", "line 3"},
		},
		{
			name:     "Add line at 2",
			number:   2,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 1", "new line", "line 2", "line 3"},
		},
		{
			name:     "Add line at 2 with empty checksum",
			number:   2,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "",
			wantErr:  true,
		},
		{
			name:     "Add line at negative",
			number:   -1,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"new line", "line 1", "line 2", "line 3"},
		},
		{
			name:     "Add line with newline character",
			number:   1,
			line:     Text("new\nline"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := writer.AddLine(tt.number, tt.line)(tt.checksum, tt.lines)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestLine_RemoveLine(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		lines    []string
		checksum string
		want     []string
		wantErr  bool
	}{
		{
			name:     "Remove line at 1",
			number:   1,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 2", "line 3"},
		},
		{
			name:     "Remove line at 2",
			number:   2,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 1", "line 3"},
		},
		{
			name:     "Remove line at 3",
			number:   3,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 1", "line 2"},
		},
		{
			name:     "Remove line at end",
			number:   writer.AT_END,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Remove line at beginning",
			number:   writer.AT_BEGINNING,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Remove line at 0",
			number:   0,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Remove line at -1",
			number:   -1,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Remove line at > number of lines",
			number:   999,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Remove line with empty checksum",
			number:   1,
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := writer.RemoveLine(tt.number)(tt.checksum, tt.lines)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestLine_UpdateLine(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		line     fmt.Stringer
		lines    []string
		checksum string
		want     []string
		wantErr  bool
	}{
		{
			name:     "Update line at 1",
			number:   1,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"new line", "line 2", "line 3"},
		},
		{
			name:     "Update line at 2",
			number:   2,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 1", "new line", "line 3"},
		},
		{
			name:     "Update line at 3",
			number:   3,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			want:     []string{"line 1", "line 2", "new line"},
		},
		{
			name:     "Update line at end",
			number:   writer.AT_END,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Update line at beginning",
			number:   writer.AT_BEGINNING,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Update line at 0",
			number:   0,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Update line at -1",
			number:   -1,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Update line at > number of lines",
			number:   999,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
		{
			name:     "Update line with empty checksum",
			number:   1,
			line:     Text("new line"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "",
			wantErr:  true,
		},
		{
			name:     "Update line with newline character",
			number:   1,
			line:     Text("new\nline"),
			lines:    []string{"line 1", "line 2", "line 3"},
			checksum: "hash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := writer.UpdateLine(tt.number, tt.line)(tt.checksum, tt.lines)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
