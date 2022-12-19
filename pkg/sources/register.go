package sources

import (
	"fmt"
	"github.com/g-portal/metadata-server/pkg/config"
	"log"
	"sync"
)

var registration = make(map[config.SourceType]Source)

var mtx sync.RWMutex

func Load() (Source, error) {
	mtx.RLock()
	defer mtx.RUnlock()

	cfg := config.GetConfig()
	for id := range registration {

		if cfg.Source.Type == id {
			if ds, ok := registration[id]; ok {
				if err := ds.Initialize(cfg.Source.Config); err != nil {
					return nil, err
				}

				return ds, nil
			}
		}
	}

	return nil, fmt.Errorf("no datasource found for type %s", cfg.Source.Type)

}

func Register(t config.SourceType, source Source) {
	mtx.Lock()
	defer mtx.Unlock()

	log.Printf("Registering datasource type %s.", t)
	registration[t] = source
}
