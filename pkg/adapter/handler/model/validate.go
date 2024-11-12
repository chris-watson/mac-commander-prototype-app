package model

import "errors"

func (c *CommandRequest) Validate() error {
	if c.Type == "" {
		return errors.New("type is required")
	}

	if c.Type == "ping" && c.Payload == "" {
		return errors.New("host is required")
	}
	return nil
}
