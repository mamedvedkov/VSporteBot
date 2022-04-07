package google_sheet

import (
	"context"
)

func (g *Google) IsRegistred(ctx context.Context, id string) (name string, ok bool) {
	return g.idToNamecache.getName(id)
}
