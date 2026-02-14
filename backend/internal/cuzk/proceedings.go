package cuzk

import (
	"context"
	"fmt"
)

// GetProceeding returns proceeding detail by ISKN ID.
func (c *Client) GetProceeding(ctx context.Context, id int64) (*Proceeding, error) {
	path := fmt.Sprintf("/Rizeni/%d", id)
	var p Proceeding
	if err := c.get(ctx, path, &p); err != nil {
		return nil, fmt.Errorf("get proceeding: %w", err)
	}
	return &p, nil
}
