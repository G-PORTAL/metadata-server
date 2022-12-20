package sources

import (
	"fmt"
	"github.com/g-portal/metadata-server/pkg/config"
	"log"
	"net/http"
	"sync"
)

type registrationType map[config.SourceType]Source

var registration = make(registrationType)

var mtx sync.RWMutex

func GetMetadata(r *http.Request) (*Metadata, error) {
	for _, source := range registration {
		if result, err := source.GetMetadata(r); err == nil && result != nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("no matching metadata found")
}

func Load() ([]Source, error) {
	mtx.RLock()
	defer mtx.RUnlock()

	list := make([]Source, 0)
	cfg := config.GetConfig()
	for id := range registration {
		if ds, ok := registration[id]; ok {
			if err := ds.Initialize(cfg.Sources.GetConfig(id)); err != nil {
				log.Printf("Failed to initialize datasource %s: %v", id, err)
				continue
			}

			list = append(list, ds)
		}
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("no datasources found")
	}

	return list, nil
}

func Register(t config.SourceType, source Source) {
	mtx.Lock()
	defer mtx.Unlock()

	log.Printf("Registering datasource type %s.", t)
	registration[t] = source
}
