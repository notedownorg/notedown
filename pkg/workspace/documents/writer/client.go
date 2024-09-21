package writer

import "path/filepath"

type Client struct {
	root string
}

func NewClient(root string) *Client {
	return &Client{root: root}
}

func (c Client) abs(doc string) string {
	return filepath.Join(c.root, doc)
}
