package writer

import "path/filepath"

var _ LineWriter = &Client{}

type Client struct {
	root string
}

func NewClient(root string) *Client {
	return &Client{root: root}
}

func (c Client) abs(doc string) string {
	return filepath.Join(c.root, doc)
}
