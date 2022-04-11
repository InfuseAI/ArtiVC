package repository

import (
	"context"
	"fmt"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureBlobRepository struct {
	Client   *azblob.ContainerClient
	Prefix   string
	BasePath string
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

func makeAzureServiceAccountUrl(accountName string) string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
}

func NewAzureBlobRepository(repo string) (*AzureBlobRepository, error) {
	ctx := context.Background()

	accountName, container, prefix, err := ParseAzureBlobUrl(repo)
	if err != nil {
		return nil, err
	}

	serviceUrl := makeAzureServiceAccountUrl(accountName)
	var serviceClient azblob.ServiceClient

	var accountKey string
	if value := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"); value != "" {
		accountKey = value
	}

	if accountKey != "" {
		credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			return nil, err
		}

		serviceClient, err = azblob.NewServiceClientWithSharedKey(serviceUrl, credential, nil)
		if err != nil {
			return nil, err
		}
	} else {
		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}

		serviceClient, err = azblob.NewServiceClient(serviceUrl, credential, nil)
		if err != nil {
			return nil, err
		}
	}

	containerClient := serviceClient.NewContainerClient(container)
	_, err = containerClient.GetProperties(ctx, nil)
	if err != nil {
		// log.Debugln(err)
		// return nil, fmt.Errorf("cannot access the container '%s' under the storage account '%s'. Pleasee use the option '--debug' to see the detail", container, accountName)
		return nil, err
	}

	return &AzureBlobRepository{
		Client:   &containerClient,
		BasePath: repo,
		Prefix:   prefix,
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
	blobClient := repo.Client.NewBlockBlobClient(blobPath)
	blobClient.Upload(
		ctx,
		src,
		nil)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AzureBlobRepository) Download(repoPath, localPath string, m *Meter) error {
	ctx := context.Background()

	// src
	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobClient := repo.Client.NewBlockBlobClient(blobPath)
	get, err := blobClient.Download(ctx, nil)
	if err != nil {
		return err
	}

	src := get.Body(nil)
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
	blobClient := repo.Client.NewBlockBlobClient(blobPath)
	_, err := blobClient.Delete(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AzureBlobRepository) Stat(repoPath string) (FileInfo, error) {
	ctx := context.Background()

	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobClient := repo.Client.NewBlockBlobClient(blobPath)
	_, err := blobClient.GetProperties(ctx, nil)
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
	pager := repo.Client.ListBlobsHierarchy("/", &azblob.ContainerListBlobHierarchySegmentOptions{Prefix: &prefix})
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blobInfo := range resp.Segment.BlobItems {
			n := *blobInfo.Name
			name := n[len(prefix):]
			entries = append(entries, &SimpleFileInfo{
				name:  name,
				isDir: false,
			})
		}

		for _, blobPrefix := range resp.Segment.BlobPrefixes {
			p := *blobPrefix.Name
			name := p[len(prefix) : len((p))-1]
			entries = append(entries, &SimpleFileInfo{
				name:  name,
				isDir: true,
			})
		}
	}

	return entries, nil
}
