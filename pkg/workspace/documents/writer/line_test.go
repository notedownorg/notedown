package writer_test

import (
	"testing"

	"github.com/liamawhite/nl/pkg/workspace/documents/writer"
	"github.com/stretchr/testify/assert"
)

func TestLines_Basic(t *testing.T) {
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client := writer.NewClient(dir)
	basic := func() Document { return loadDocument(t, dir, "basic.md") }

	// Updates (do these first as they rely most on the original line numbers)
	assert.NoError(t, client.UpdateLine(basic().Document, 7, Text("This line was updated")))
	assert.Error(t, client.UpdateLine(basic().Document, 0, Text("No such line as they're 1-indexed")))
	assert.Error(t, client.UpdateLine(basic().Document, 999, Text("No such line, oob")))
	assert.Error(t, client.UpdateLine(basic().Document, writer.AtBeginning, Text("Must provide an absolute line number")))
	assert.Error(t, client.UpdateLine(basic().Document, writer.AtEnd, Text("Must provide an absolute line number")))

	// Deletes
	assert.NoError(t, client.RemoveLine(basic().Document, 5))
	assert.Error(t, client.RemoveLine(basic().Document, 0))
	assert.Error(t, client.RemoveLine(basic().Document, 999))
	assert.Error(t, client.RemoveLine(basic().Document, writer.AtBeginning))
	assert.Error(t, client.RemoveLine(basic().Document, writer.AtEnd))

	// Adds (do these last so we don't change the delete/update line numbers)
	assert.NoError(t, client.AddLine(basic().Document, 999, Text("This line was added at line 999")))
	assert.NoError(t, client.AddLine(basic().Document, writer.AtBeginning, Text("This line was added at the beginning")))
	assert.NoError(t, client.AddLine(basic().Document, 3, Text("This line was added at line 3")))
	assert.NoError(t, client.AddLine(basic().Document, writer.AtEnd, Text("This line was added at the end")))

	// Verify the files are all correct
	basicWant := loadDocument(t, "testdata/golden", "basic.md")
	assert.Equal(t, string(basicWant.Contents), string(basic().Contents))
}

func TestLines_Empty(t *testing.T) {
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client := writer.NewClient(dir)
	empty := func() Document { return loadDocument(t, dir, "empty.md") }

	// Beginning
	assert.NoError(t, client.AddLine(empty().Document, writer.AtBeginning, Text("This line was added at the beginning")))
	assert.NoError(t, client.UpdateLine(empty().Document, 1, Text("This line was updated")))
	assert.NoError(t, client.RemoveLine(empty().Document, 1))

	// End
	assert.NoError(t, client.AddLine(empty().Document, writer.AtEnd, Text("This line was added at the end")))
	assert.NoError(t, client.UpdateLine(empty().Document, 1, Text("This line was updated")))
	assert.NoError(t, client.RemoveLine(empty().Document, 1))

	// 1
	assert.NoError(t, client.AddLine(empty().Document, 1, Text("This line was added at line 1")))
	assert.NoError(t, client.UpdateLine(empty().Document, 1, Text("This line was updated")))
	assert.NoError(t, client.RemoveLine(empty().Document, 1))

	emptyWant := loadDocument(t, "testdata/golden", "empty.md")
	assert.Equal(t, string(emptyWant.Contents), string(empty().Contents))
}

func TestLines_Frontmatter(t *testing.T) {
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client := writer.NewClient(dir)
	frontmatter := func() Document { return loadDocument(t, dir, "frontmatter.md") }

	assert.NoError(t, client.AddLine(frontmatter().Document, writer.AtBeginning, Text("This line was added at the beginning but should be after frontmatter")))

	// 0 == AtBeginning which is fine, its inserted after the frontmatter
	assert.Error(t, client.AddLine(frontmatter().Document, 1, Text("Can't add frontmatter by line")))
	assert.Error(t, client.AddLine(frontmatter().Document, 2, Text("Can't add frontmatter by line")))

	assert.Error(t, client.RemoveLine(frontmatter().Document, 0))
	assert.Error(t, client.RemoveLine(frontmatter().Document, 1))
	assert.Error(t, client.RemoveLine(frontmatter().Document, 2))

	assert.Error(t, client.UpdateLine(frontmatter().Document, 0, Text("Can't update frontmatter by line")))
	assert.Error(t, client.UpdateLine(frontmatter().Document, 1, Text("Can't update frontmatter by line")))
	assert.Error(t, client.UpdateLine(frontmatter().Document, 2, Text("Can't update frontmatter by line")))

	frontmatterWant := loadDocument(t, "testdata/golden", "frontmatter.md")
	assert.Equal(t, string(frontmatterWant.Contents), string(frontmatter().Contents))
}

func TestLines_StaleWrites(t *testing.T) {
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client := writer.NewClient(dir)
	basic := loadDocument(t, dir, "basic.md")

	// Make changes to the file
	assert.NoError(t, client.AddLine(basic.Document, writer.AtEnd, Text("This line was added at the end")))

	// Now when we go to write using the original document/hash, we should get an error
	assert.Error(t, client.AddLine(basic.Document, writer.AtEnd, Text("This line is being written to a stale document")))
}
