package google_sheet

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/mamedvedkov/tools/processes"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	baseURL string = "https://sheets.googleapis.com/v4"
	scope   string = "https://spreadsheets.google.com/feeds"
)

type SpreadsheetsIds struct {
	Payment             string
	PaymentConsolidated string
	Schedule            string
	Financial           string
}

type Google struct {
	logger logr.Logger

	client *http.Client

	sheets *sheets.Service

	spreadsheetIds SpreadsheetsIds
	mu             *sync.RWMutex

	idToNamecache     *idToNamecache
	pastMonthCache    *paymentsCache
	currentMonthCache *paymentsCache
	requisitesCache   *requisitesCache

	entityToInn map[string]string
}

func MustGoogle(logger logr.Logger, jsonKey []byte, ids SpreadsheetsIds, entityToInn map[string]string) *Google {
	g, err := NewGoogle(logger, jsonKey, ids, entityToInn)
	if err != nil {
		panic(err)
	}

	return g
}

func NewGoogle(logger logr.Logger, jsonKey []byte, ids SpreadsheetsIds, entityToInn map[string]string) (*Google, error) {
	sheetsService, err := sheets.NewService(context.Background(), option.WithCredentialsJSON(jsonKey))
	if err != nil {
		return nil, err
	}

	g := &Google{
		logger:         logger.WithName("GoogleAdapter"),
		sheets:         sheetsService,
		spreadsheetIds: ids,
		mu:             &sync.RWMutex{},
		idToNamecache: &idToNamecache{
			mu:       &sync.RWMutex{},
			idToName: make(map[string]string),
		},
		pastMonthCache: &paymentsCache{
			mu:   &sync.RWMutex{},
			data: make(map[string]paymentsData),
		},
		currentMonthCache: &paymentsCache{
			mu:   &sync.RWMutex{},
			data: make(map[string]paymentsData),
		},
		requisitesCache: &requisitesCache{
			mu:   &sync.RWMutex{},
			data: make(map[string]requisitesData),
		},
		entityToInn: entityToInn,
	}

	err = g.reloadCache(context.Background())
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Google) ReloadSpreadSheets(interval time.Duration) processes.Process {
	ticker := time.NewTicker(interval)

	return func(ctx context.Context) (err error) {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				err := g.reloadCache(ctx)
				if err != nil {
					g.logger.Error(err, "error reloading cache")
				}
			}
		}
	}
}

func (g *Google) reloadCache(ctx context.Context) error {
	err := processes.RunParallelAndWait(ctx,
		g.reloadIdToNameCache,
		g.reloadPaymentsCache("отчет", g.pastMonthCache),
		g.reloadPaymentsCache("отчетТМ", g.currentMonthCache),
		g.reloadRequisites,
	)
	if err != nil {
		return fmt.Errorf("error in one of reloads: %w", err)
	}

	return nil
}

func (g *Google) reloadIdToNameCache(ctx context.Context) error {
	names, err := g.getColumn(ctx, g.spreadsheetIds.Payment, "отчет рассылка!A:A")
	if err != nil {
		return err
	}

	ids, err := g.getColumn(ctx, g.spreadsheetIds.Payment, "отчет рассылка!H:H")
	if err != nil {
		return err
	}

	data := make(map[string]string)

	for i := 1; i < len(ids); i++ {
		if ids[i] == "" {
			continue
		}

		_, err := strconv.Atoi(ids[i])
		if err != nil {
			g.logger.Error(err, "wrong id", "value", ids[i])
			continue
		}

		data[ids[i]] = names[i]
	}

	g.idToNamecache.update(data)

	return nil
}

func (g *Google) reloadPaymentsCache(title string, dest *paymentsCache) processes.Process {
	return func(ctx context.Context) (err error) {
		columnWithName, err := g.getColumn(ctx, g.spreadsheetIds.Payment,
			fmt.Sprintf("%s!%s", title, "A:A"))
		if err != nil {
			return err
		}

		columnWithSum, err := g.getColumn(ctx, g.spreadsheetIds.Payment,
			fmt.Sprintf("%s!%s", title, "B:B"))
		if err != nil {
			return err
		}

		columnWithDetail, err := g.getColumn(ctx, g.spreadsheetIds.Payment,
			fmt.Sprintf("%s!%s", title, "C:C"))
		if err != nil {
			return err
		}

		idxLimit := findMin(
			len(columnWithName),
			len(columnWithSum),
			len(columnWithDetail),
		)

		data := make(map[string]paymentsData)

		for i := 1; i < idxLimit; i++ {
			data[columnWithName[i]] = paymentsData{
				value:     columnWithSum[i],
				detailUrl: columnWithDetail[i],
			}
		}

		dest.update(data)

		return nil
	}
}

func (g *Google) reloadRequisites(ctx context.Context) error {
	columnWithName, err := g.getColumn(ctx, g.spreadsheetIds.PaymentConsolidated, "Реквизиты!A:A")
	if err != nil {
		return err
	}

	columnWithInn, err := g.getColumn(ctx, g.spreadsheetIds.PaymentConsolidated, "Реквизиты!C:C")
	if err != nil {
		return err
	}

	columnWithRS, err := g.getColumn(ctx, g.spreadsheetIds.PaymentConsolidated, "Реквизиты!D:D")
	if err != nil {
		return err
	}

	columnWithPhone, err := g.getColumn(ctx, g.spreadsheetIds.PaymentConsolidated, "Реквизиты!G:G")
	if err != nil {
		return err
	}

	columnWithCard, err := g.getColumn(ctx, g.spreadsheetIds.PaymentConsolidated, "Реквизиты!I:I")
	if err != nil {
		return err
	}

	idxLimit := findMin(
		len(columnWithCard),
		len(columnWithInn),
		len(columnWithRS),
		len(columnWithName),
		len(columnWithPhone),
	)

	data := make(map[string]requisitesData)

	for i := 1; i < idxLimit; i++ {
		data[columnWithName[i]] = requisitesData{
			phone:      columnWithPhone[i],
			cardNumber: columnWithCard[i],
			inn:        columnWithInn[i],
			rs:         columnWithRS[i],
		}
	}

	g.requisitesCache.update(data)

	return nil
}

func findMin(lens ...int) int {
	if len(lens) == 0 {
		return 0
	}

	if len(lens) == 1 {
		return lens[1]
	}

	min := float64(lens[0])

	for idx := range lens {
		min = math.Min(min, float64(lens[idx]))
	}

	return int(min)
}
