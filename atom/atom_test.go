package atom_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	ka "github.com/itpkg/ksana/atom"
)

func TestAtom(t *testing.T) {
	f := ka.NewFeed("id-111", "main title", "sub title", "http://www.aaa.com")

	for i := 0; i < 10; i++ {
		e := ka.NewEntry(fmt.Sprintf("article-%d", i), fmt.Sprintf("title %d", i), fmt.Sprintf("http://www.aaa.com/articles/%d", i), fmt.Sprintf("summary-%d", i), time.Now())
		e.SetAuthor("My Name", "name@aaa.com")
		e.SetContent("<h1>Hello</h1>")

		f.Add(e)
	}
	if err := f.Write(os.Stdout, true); err != nil {
		fmt.Sprintf("error on write: %v", err)
	}
}
