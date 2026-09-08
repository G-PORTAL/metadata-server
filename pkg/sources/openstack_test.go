package sources_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/g-portal/metadata-server/pkg/sources"
)

func renderOpenStackMetaData(t *testing.T, m *sources.Metadata) map[string]any {
	t.Helper()

	writer := httptest.NewRecorder()
	r := m.OpenStackMetaData()
	r.WriteContentType(writer)

	if err := r.Render(writer); err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(writer.Body)
	if err != nil {
		t.Fatal(err)
	}

	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("meta_data.json is not valid JSON: %v (%s)", err, string(body))
	}

	return out
}

func TestOpenStackMetaDataRendersTagsAsMeta(t *testing.T) {
	metadata := metadataMockup()
	metadata.Tags = map[string]string{
		"role": "web",
		"env":  "production",
	}

	out := renderOpenStackMetaData(t, metadata)

	meta, ok := out["meta"].(map[string]any)
	if !ok {
		t.Fatalf("expected \"meta\" object in meta_data.json, got %v", out["meta"])
	}

	if meta["role"] != "web" || meta["env"] != "production" {
		t.Fatalf("unexpected meta content: %v", meta)
	}
}

func TestOpenStackMetaDataOmitsMetaWithoutTags(t *testing.T) {
	out := renderOpenStackMetaData(t, metadataMockup())

	if _, present := out["meta"]; present {
		t.Fatalf("expected no \"meta\" key without tags, got %v", out["meta"])
	}
}
