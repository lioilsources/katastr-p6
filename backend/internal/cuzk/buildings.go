package cuzk

import (
	"context"
	"fmt"
)

// SearchBuildings searches for buildings by cadastral area code and number.
func (c *Client) SearchBuildings(ctx context.Context, areaCode int, number string) (*BuildingSearchResponse, error) {
	path := fmt.Sprintf("/Stavby/Vyhledani?katastralniUzemi=%d&cislo=%s", areaCode, number)
	var resp BuildingSearchResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("search buildings: %w", err)
	}
	return &resp, nil
}

// GetBuilding returns building detail by ISKN ID.
func (c *Client) GetBuilding(ctx context.Context, id int64) (*Building, error) {
	path := fmt.Sprintf("/Stavby/%d", id)
	var b Building
	if err := c.get(ctx, path, &b); err != nil {
		return nil, fmt.Errorf("get building: %w", err)
	}
	return &b, nil
}
