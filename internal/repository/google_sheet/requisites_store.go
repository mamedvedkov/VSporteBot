package google_sheet

import (
	"context"
	"fmt"
)

func (g *Google) RequisitesInfo(ctx context.Context, name string) (string, error) {
	data, ok := g.requisitesCache.getData(name)
	if !ok {
		return "", fmt.Errorf("%s not found in cache", name)
	}

	res := fmt.Sprintf("Номер телефона:\t%s\nНомер карты:\t%s\n", data.phone,
		"**** "+data.cardNumber[len(data.cardNumber)-4:])

	if data.inn != "" {
		res += fmt.Sprintf("ИНН ФЛ:\t%s\n", data.inn)
	}

	if data.rs != "" {
		res += fmt.Sprintf("Р/С ФЛ:\t%s\n", data.rs)
	}

	return res, nil
}
