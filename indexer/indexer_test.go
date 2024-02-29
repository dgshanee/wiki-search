package indexer

import (
	"testing"
)

func Test_index(t *testing.T) {
	idxr := NewIndexer()
	test := []struct {
		name      string
		data      string
		url       string
		expectErr bool
	}{
		{"good_insert", "word", "https://www.google.com", false},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			err := idxr.Index(tc.data, tc.url)
			if err != nil {
				if !tc.expectErr {
					t.Fatalf("Unexpected error at %s, %s", tc.data, tc.url)
				}

				return
			}

			if err == nil && tc.expectErr {
				t.Errorf("Expected error at %s, %s", tc.data, tc.url)
				return
			}
		})
	}
}
