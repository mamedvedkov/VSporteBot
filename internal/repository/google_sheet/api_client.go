package google_sheet

import (
	"context"
	"fmt"
)

const columnsDimension = "COLUMNS"

func (g *Google) getColumn(ctx context.Context, spreadsheetId, _range string) ([]string, error) {
	resp, err := g.sheets.Spreadsheets.
		Values.Get(spreadsheetId, _range).
		MajorDimension(columnsDimension).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("cant get column: %w", err)
	}

	if len(resp.Values) != 1 {
		return nil, fmt.Errorf("wrong resp len: %v", len(resp.Values))
	}

	column := make([]string, 0)

	for idx := range resp.Values[0] {
		switch val := resp.Values[0][idx].(type) {
		case string:
			column = append(column, val)
		default:
			return nil, fmt.Errorf("non string value in collumn")
		}
	}

	return column, nil
}
