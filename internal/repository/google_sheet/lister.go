package google_sheet

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mamedvedkov/VSporteBot/internal"
	"github.com/mamedvedkov/VSporteBot/internal/notifier"
)

func (g *Google) ListIds() []notifier.NotifyData {
	ctx := context.Background()

	t := time.Now()
	month := t.Month()
	year := t.Year()

	var pastMonth string
	if month == 1 {
		pastMonth = internal.RussianMonths[12]
	} else {
		pastMonth = internal.RussianMonths[int(month)-1]
	}

	return g.fetchMonth(ctx, pastMonth, year)
}

func (g *Google) fetchMonth(ctx context.Context, month string, year int) []notifier.NotifyData {
	columnWithName, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!A:A", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch name for notify")
	}

	columnWithState, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!C:C", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch state for notify")
	}

	columnWithPrePaimentState, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!J:J", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch name for notify")
	}

	columnWithPrePaimentSum, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!K:K", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch name for notify")
	}

	columnWithPrePaimentEntity, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!M:M", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch name for notify")
	}

	columnWithPaimentState, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!V:V", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch name for notify")
	}

	columnWithPaimentSum, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!U:U", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch name for notify")
	}

	columnWithPaimentEntity, err := g.getColumn(ctx, g.spreadsheetIds.Financial, fmt.Sprintf("%s %v!W:W", month, year))
	if err != nil {
		g.logger.Error(err, "cant fetch name for notify")
	}

	minLen := findMin(
		len(columnWithName),
		len(columnWithState),
		len(columnWithPrePaimentState),
		len(columnWithPrePaimentSum),
		len(columnWithPrePaimentEntity),
		len(columnWithPaimentState),
		len(columnWithPaimentSum),
		len(columnWithPaimentEntity),
	)

	res := make([]notifier.NotifyData, 0)
	for idx := 2; idx < minLen; idx++ {
		if columnWithState[idx] != "Оплачен" && columnWithState[idx] != "Частично" {
			continue
		}

		if columnWithPrePaimentState[idx] == "СЗ" {
			idString, ok := g.idToNamecache.getId(columnWithName[idx])
			if !ok {
				continue
			}

			id, err := strconv.Atoi(idString)
			if err != nil {
				g.logger.Error(err, "cant parse id", "id", idString)
				continue
			}

			res = append(res, notifier.NotifyData{
				Id:     int64(id),
				Entity: columnWithPrePaimentEntity[idx],
				Inn:    "",
				Sum:    columnWithPrePaimentSum[idx],
				Month:  month,
			})
		}

		if columnWithPaimentState[idx] == "СЗ" {
			idString, ok := g.idToNamecache.getId(columnWithName[idx])
			if !ok {
				continue
			}

			id, err := strconv.Atoi(idString)
			if err != nil {
				g.logger.Error(err, "cant parse id", "id", idString)
				continue
			}

			res = append(res, notifier.NotifyData{
				Id:     int64(id),
				Entity: columnWithPaimentEntity[idx],
				Inn:    g.entityToInn[columnWithPaimentEntity[idx]],
				Sum:    columnWithPaimentSum[idx],
				Month:  month,
			})
		}
	}

	return res
}
