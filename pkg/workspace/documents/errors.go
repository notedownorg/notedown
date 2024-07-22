package documents

func (c *Client) Errors() <-chan error {
	return c.errors
}
