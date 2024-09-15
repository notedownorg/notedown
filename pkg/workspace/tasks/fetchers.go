package tasks

type TaskFetcher func(c *Client) ([]Task, error)

func FetchAllTasks() TaskFetcher {
	return func(c *Client) ([]Task, error) {
		var tasks []Task
		c.mutex.RLock()
		for _, document := range c.cache {
			for _, task := range document {
				tasks = append(tasks, *task)
			}
		}
		c.mutex.RUnlock()
		return tasks, nil
	}
}

func FetchTasksForDocument(document string) TaskFetcher {
	return func(c *Client) ([]Task, error) {
		var tasks []Task
		c.mutex.RLock()
		for _, task := range c.cache[document] {
			tasks = append(tasks, *task)
		}
		c.mutex.RUnlock()
		return tasks, nil
	}
}
