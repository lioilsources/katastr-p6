package cuzk

import (
	"context"
	"fmt"
)

// SearchUnits searches for property units by cadastral area, building number, and unit number.
func (c *Client) SearchUnits(ctx context.Context, areaCode int, buildingNo, unitNo string) (*UnitSearchResponse, error) {
	path := fmt.Sprintf("/Jednotky/Vyhledani?katastralniUzemi=%d&cisloStavby=%s&cisloJednotky=%s", areaCode, buildingNo, unitNo)
	var resp UnitSearchResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("search units: %w", err)
	}
	return &resp, nil
}

// GetUnit returns unit detail by ISKN ID.
func (c *Client) GetUnit(ctx context.Context, id int64) (*Unit, error) {
	path := fmt.Sprintf("/Jednotky/%d", id)
	var u Unit
	if err := c.get(ctx, path, &u); err != nil {
		return nil, fmt.Errorf("get unit: %w", err)
	}
	return &u, nil
}
