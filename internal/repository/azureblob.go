package repository

import (
	"context"
	"errors"
	"fmt"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/infuseai/artivc/internal/log"
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

	r := &AzureBlobRepository{
		Client:   &containerClient,
		BasePath: repo,
		Prefix:   prefix,
	}

	// check if the client has enough permission
	dir, err := os.MkdirTemp("", "artivc-azblob-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // clean up

	err = r.Download("refs/latest", filepath.Join(dir, "latest"), nil)
	if err != nil {
		var internalError *azblob.InternalError
		if !errors.As(err, &internalError) {
			return nil, err
		}

		var errStorage *azblob.StorageError
		if !internalError.As(&errStorage) {
			return nil, internalError
		}

		if errStorage.ErrorCode == azblob.StorageErrorCodeBlobNotFound {
			// blob not found. but authentication/authorization check is ok. PASS
		} else if errStorage.ErrorCode == azblob.StorageErrorCodeAuthorizationPermissionMismatch {
			// authorization permission mismatch
			log.Debugln(errStorage.Error())
			fmt.Fprintf(os.Stderr, "Authorization permission mismatch. Please assign 'Storage Blob Data Contributor' role to the logged-in account in the storage account '%s'\n", accountName)
			fmt.Fprintln(os.Stderr, "Please see https://docs.microsoft.com/azure/storage/blobs/assign-azure-role-data-access")
			fmt.Fprintln(os.Stderr, "")
			return nil, fmt.Errorf("authorization permission mismatch")
		} else {
			// other error
			return nil, errStorage
		}
	}

	return r, nil
}

func (repo *AzureBlobRepository) Upload(localPath, repoPath string, m *Meter) error {
	ctx := context.Background()

	// file
	src, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// upload
	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobClient := repo.Client.NewBlockBlobClient(blobPath)

	blobClient.UploadFileToBlockBlob(
		ctx,
		src,
		azblob.HighLevelUploadToBlockBlobOption{
			Progress: func(bytesTransferred int64) {
				if m != nil {
					m.SetBytes(bytesTransferred)
				}
			},
			Parallelism: 10,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *AzureBlobRepository) Download(repoPath, localPath string, m *Meter) error {
	ctx := context.Background()

	// file
	dest, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// download
	blobPath := filepath.Join(repo.Prefix, repoPath)
	blobClient := repo.Client.NewBlockBlobClient(blobPath)
	err = blobClient.DownloadBlobToFile(ctx, 0, 0, dest, azblob.HighLevelDownloadFromBlobOptions{
		Progress: func(bytesTransferred int64) {
			if m != nil {
				m.SetBytes(bytesTransferred)
			}
		},
		Parallelism: 10,
	})
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
