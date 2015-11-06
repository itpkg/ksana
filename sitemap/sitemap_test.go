package sitemap_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	ks "github.com/itpkg/ksana/sitemap"
)

func TestSitemap(t *testing.T) {
	sm := ks.New()
	for i := 0; i < 10; i++ {
		sm.Add(fmt.Sprintf("http://www.aaa.com/%d", i), time.Now(), ks.Dialy, 0.7)
	}

	if err := sm.Write(os.Stdout, true); err != nil {
		t.Errorf("error on to xml: %v", err)
	}
}
