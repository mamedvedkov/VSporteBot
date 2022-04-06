package google_sheet

import "sync"

type idToNamecache struct {
	mu       *sync.RWMutex
	idToName map[string]string
	nameToId map[string]string
}

func (i *idToNamecache) update(data map[string]string) {
	nameToId := make(map[string]string)

	for key, val := range data {
		nameToId[val] = key
	}

	i.mu.Lock()

	i.idToName = data
	i.nameToId = nameToId

	i.mu.Unlock()
}

func (i *idToNamecache) getName(id string) (string, bool) {
	i.mu.RLock()

	name, ok := i.idToName[id]

	i.mu.RUnlock()

	return name, ok
}

func (i *idToNamecache) getId(name string) (string, bool) {
	i.mu.RLock()

	id, ok := i.nameToId[name]

	i.mu.RUnlock()

	return id, ok
}

type paymentsData struct {
	value     string
	detailUrl string
}

type paymentsCache struct {
	mu   *sync.RWMutex
	data map[string]paymentsData
}

func (p *paymentsCache) update(data map[string]paymentsData) {
	p.mu.Lock()

	p.data = data

	p.mu.Unlock()
}

func (p *paymentsCache) getData(name string) (paymentsData, bool) {
	p.mu.RLock()

	data, ok := p.data[name]

	p.mu.RUnlock()

	return data, ok
}

type requisitesData struct {
	phone      string
	cardNumber string
	inn        string
	rs         string
}

type requisitesCache struct {
	data map[string]requisitesData
	mu   *sync.RWMutex
}

func (r *requisitesCache) update(data map[string]requisitesData) {
	r.mu.Lock()

	r.data = data

	r.mu.Unlock()
}

func (r *requisitesCache) getData(name string) (requisitesData, bool) {
	r.mu.RLock()

	data, ok := r.data[name]

	r.mu.RUnlock()

	return data, ok
}
