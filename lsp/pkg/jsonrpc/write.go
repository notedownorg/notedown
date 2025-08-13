package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
)

func Write(w *bufio.Writer, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}
	headers := fmt.Sprintf("%s: %d\r\n\r\n", ContentLengthHeader, len(body))
	if _, err := w.WriteString(headers); err != nil {
		return fmt.Errorf("error writing headers: %w", err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("error writing body: %w", err)
	}
	return w.Flush()
}

