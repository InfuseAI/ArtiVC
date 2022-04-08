package repository

import (
	"context"
	"fmt"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/infuseai/artivc/internal/log"
)

type AzureBlobRepository struct {
	ContainerURL azblob.ContainerURL
	Prefix       string
	Pipeline     pipeline.Pipeline
	BasePath     string
}

func IsAzureStorageUrl(repoUrl string) bool {
	url, err := neturl.Parse(repoUrl)
	if err != nil {
		return false
	}

	return strings.HasSuffix(url.Host, ".blob.core.windows.net")
}

func ParseAzureBlobUrl(urlString string) (storageAccount, container, prefix string, err error) {
	url, err := neturl.Parse(urlString)
	if err != nil {
		return
	}

	storageAccount = url.Host[:len(url.Host)-len(".blob.core.windows.net")]
	comps := strings.Split(url.Path, "/")
	if len(comps) < 2 {
		err = fmt.Errorf("invalid azure blob url: " + urlString)
		return
	}

	container = comps[1]
	prefix = strings.Join(comps[2:], "/")

	return
}

func NewAzureBlobRepository(repo string) (*AzureBlobRepository, error) {
	url, err := neturl.Parse(repo)
	if err != nil {
		return nil, err
	}

	accountName, container, prefix, err := ParseAzureBlobUrl(repo)
	if err != nil {
		return nil, err
	}

	var storageAccountKeys string
	if value := os.Getenv("AZBLOB_STORAGE_ACCOUNT_KEY"); value != "" {
		storageAccountKeys = value
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, storageAccountKeys)
	if err != nil {
		log.Debugln(err)
		return nil, fmt.Errorf("cannot login to storage account " + accountName)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	if err != nil {
		return nil, err
	}

	serviceURL := azblob.NewServiceURL(*u, pipeline)
	containerURL := serviceURL.NewContainerURL(container)

	_, err = containerURL.GetAccountInfo(context.Background())
	if err != nil {
		log.Debugln(err)
		return nil, fmt.Errorf("cannot access the container '%s' under the storage account '%s'. Pleasee use the option '--debug' to see the detail", container, accountName)
	}

	return &AzureBlobRepository{
		ContainerURL: containerURL,
		Pipeline:     pipeline,
		BasePath:     repo,
		Prefix:       prefix,
	}, nil
}

func (repo *AzureBlobRepository) Upload(localPath, repoPath string, m *Meter) error {
	ctx := context.Background()

	// src
	src, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// dest
	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobUrl := repo.ContainerURL.NewBlockBlobURL(blobPath)
	blobUrl.Upload(
		ctx,
		src,
		azblob.BlobHTTPHeaders{},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
		azblob.DefaultAccessTier,
		nil,
		azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (repo *AzureBlobRepository) Download(repoPath, localPath string, m *Meter) error {
	ctx := context.Background()

	// src
	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobURL := repo.ContainerURL.NewBlockBlobURL(blobPath)
	get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return err
	}

	src := get.Body(azblob.RetryReaderOptions{})
	defer src.Close()

	// dest
	dest, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// copy
	_, err = CopyWithMeter(dest, src, m)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AzureBlobRepository) Delete(repoPath string) error {
	ctx := context.Background()

	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobURL := repo.ContainerURL.NewBlockBlobURL(blobPath)
	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}

	return nil
}

func (repo *AzureBlobRepository) Stat(repoPath string) (FileInfo, error) {
	ctx := context.Background()

	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobURL := repo.ContainerURL.NewBlockBlobURL(blobPath)
	_, err := blobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	return &SimpleFileInfo{
		name: filepath.Base(repoPath),
	}, nil
}

func (repo *AzureBlobRepository) List(repoPath string) ([]FileInfo, error) {
	ctx := context.Background()
	entries := make([]FileInfo, 0)
	prefix := filepath.Join(repo.Prefix, repoPath) + "/"
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlobs, err := repo.ContainerURL.ListBlobsHierarchySegment(ctx, marker, "/", azblob.ListBlobsSegmentOptions{Prefix: prefix})
		if err != nil {
			return nil, err
		}

		marker = listBlobs.NextMarker
		for _, blobInfo := range listBlobs.Segment.BlobItems {
			name := blobInfo.Name[len(prefix):]
			entries = append(entries, &SimpleFileInfo{
				name:  name,
				isDir: false,
			})
		}

		for _, blobPrefix := range listBlobs.Segment.BlobPrefixes {
			name := blobPrefix.Name[len(prefix) : len(blobPrefix.Name)-1]
			entries = append(entries, &SimpleFileInfo{
				name:  name,
				isDir: true,
			})
		}
	}

	return entries, nil
}
