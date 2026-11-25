package output

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestVaillantCacheSeedRoundTrip(t *testing.T) {
	const seed = "testdata/semantic_cache_seed.json"
	snapshot, ok := NewSemanticCacheStore(seed, t.Logf).Load()
	if !ok {
		t.Fatal("seed did not load")
	}
	store := NewSemanticCacheStore(filepath.Join(t.TempDir(), "semantic-cache.json"), t.Logf)
	store.now = func() time.Time { return time.Date(2030, time.April, 10, 8, 0, 0, 0, time.UTC) }
	if err := store.Save(snapshot); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(store.path)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/semantic_cache_seed.saved.json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("saved bytes:\n%s\nwant:\n%s", got, want)
	}
}
