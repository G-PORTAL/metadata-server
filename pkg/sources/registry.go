package sources

import (
	"github.com/g-portal/metadata-server/pkg/config"
	"log"
	"sync"
)

type registrationType map[config.SourceType]Source

var registration = make(registrationType)

var mtx sync.RWMutex

func Load() ([]Source, error) {
	mtx.RLock()
	defer mtx.RUnlock()

	list := make([]Source, 0)
	cfg := config.GetConfig()
	for id := range registration {
		if source, ok := registration[id]; ok {
			if err := source.Initialize(cfg.Sources.GetConfig(id)); err != nil {
				log.Printf("Failed to initialize datasource %s: %v", id, err)

				continue
			}

			list = append(list, source)
		}
	}

	if len(list) == 0 {
		return nil, ErrNoDatasourceFound
	}

	return list, nil
}

func Register(t config.SourceType, source Source) {
	mtx.Lock()
	defer mtx.Unlock()

	log.Printf("Registering datasource type %s.", t)
	registration[t] = source
}
