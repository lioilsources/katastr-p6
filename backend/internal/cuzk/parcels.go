package cuzk

import (
	"context"
	"fmt"
)

// SearchParcels searches for parcels by cadastral area code and number.
func (c *Client) SearchParcels(ctx context.Context, areaCode int, number string) (*ParcelSearchResponse, error) {
	path := fmt.Sprintf("/Parcely/Vyhledani?katastralniUzemi=%d&kmenoveCislo=%s", areaCode, number)
	var resp ParcelSearchResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("search parcels: %w", err)
	}
	return &resp, nil
}

// GetParcel returns parcel detail by ISKN ID.
func (c *Client) GetParcel(ctx context.Context, id int64) (*Parcel, error) {
	path := fmt.Sprintf("/Parcely/%d", id)
	var p Parcel
	if err := c.get(ctx, path, &p); err != nil {
		return nil, fmt.Errorf("get parcel: %w", err)
	}
	return &p, nil
}

// PolygonParcels finds parcels within a polygon area defined by S-JTSK coordinates.
// x, y are positive S-JTSK values; radius is in meters.
func (c *Client) PolygonParcels(ctx context.Context, x, y float64, radius int) (*ParcelSearchResponse, error) {
	// Build a small bounding box around the point.
	r := float64(radius)
	path := fmt.Sprintf("/Parcely/Polygon?souradniceX=%.0f&souradniceX=%.0f&souradniceX=%.0f&souradniceX=%.0f&souradniceY=%.0f&souradniceY=%.0f&souradniceY=%.0f&souradniceY=%.0f",
		x-r, x-r, x+r, x+r,
		y-r, y+r, y+r, y-r,
	)
	var resp ParcelSearchResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("polygon parcels: %w", err)
	}
	return &resp, nil
}

// NeighborParcels returns neighboring parcels for a given parcel ID.
func (c *Client) NeighborParcels(ctx context.Context, id int64) (*NeighborParcelsResponse, error) {
	path := fmt.Sprintf("/Parcely/SousedniParcely/%d", id)
	var resp NeighborParcelsResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("neighbor parcels: %w", err)
	}
	return &resp, nil
}
