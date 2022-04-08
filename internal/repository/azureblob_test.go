package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseAzureBlobUrl(t *testing.T) {
	testCases := []struct {
		repo           string
		storageAccount string
		container      string
		prefix         string
	}{
		{
			repo:           "https://artivc.blob.core.windows.net/avc",
			storageAccount: "artivc",
			container:      "avc",
			prefix:         "",
		},
		{
			repo:           "https://artivc.blob.core.windows.net/avc/",
			storageAccount: "artivc",
			container:      "avc",
			prefix:         "",
		},
		{
			repo:           "https://artivc.blob.core.windows.net/avc/abc",
			storageAccount: "artivc",
			container:      "avc",
			prefix:         "abc",
		},
		{
			repo:           "https://artivc.blob.core.windows.net/avc/abc/",
			storageAccount: "artivc",
			container:      "avc",
			prefix:         "abc/",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.repo, func(t *testing.T) {
			storageAccount, container, prefix, err := ParseAzureBlobUrl(tC.repo)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, tC.storageAccount, storageAccount)
			assert.Equal(t, tC.container, container)
			assert.Equal(t, tC.prefix, prefix)

		})
	}
}
