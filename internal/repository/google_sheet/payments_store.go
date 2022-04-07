package google_sheet

import (
	"context"
	"fmt"
	"strings"
)

func (g *Google) CurrentMonthPayment(ctx context.Context, name string) (string, error) {
	data, ok := g.currentMonthCache.getData(name)
	if !ok {
		return "", fmt.Errorf("no name = %s in cache", name)
	}

	return strings.TrimRight(data.value, "*расчет не полный"), nil
}

func (g *Google) PastMonthPayment(ctx context.Context, name string) (string, error) {
	data, ok := g.pastMonthCache.getData(name)
	if !ok {
		return "", fmt.Errorf("no name in cache")
	}

	return data.value, nil
}

func (g *Google) PaymentDetail(ctx context.Context, name string) (string, error) {
	data, ok := g.pastMonthCache.getData(name)
	if !ok {
		return "", fmt.Errorf("no name in cache")
	}

	return data.detailUrl, nil
}
