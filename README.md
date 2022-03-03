# ArtiV

ArtiV (**Arti**fact **V**ersions) is a version control system for large files. 

To store and share large files, we may use NFS or object storage (e.g. s3, MinIO). However, if we would like to do versioning on top of them, it is not a trivial thing. ArtiV is a CLI tool to enable you to version files on your storage without pain. You don't need to install any additional server or gateway and we turn your storage into the versioned repository.

# Install
Download the latest ArtiV in the [release](https://github.com/InfuseAI/ArtiV/releases) page

# Quick Start

## Push data
1. Select the bucket and path to put your versioned data. Here we use `s3://mybucket/path/to/repo`. You can use local repository instead (e.g. `/path/to/repo`)
1. Prepare the `~/.aws/credentials` if you use the S3 repository
1. Push your data to the repository
   ```
   cd /path/to/my/data
   art init s3://mybucket/path/to/repo
   art push
   ```

1. See the commit log
   ```
   art log
   ```

1. Add more data to your repository

   ```
   echo "hello" > hello
   art push --dry-run
   art push
   ```

1. Tag a version
   ```
   art tag v0.1.0
   art log
   ```


## Pull data
1. Pull data from the repository
   ```
   cd /path/to/download
   art init s3://mybucket/path/to/repo
   art pull
   ```

1. or download data with single `get` command
    ```
    art get -o /path/to/download s3://mybucket/path/to/repo    
    ```
    or download the specific version
    ```
    art get -o /path/to/download s3://mybucket/path/to/repo@v0.1.0
    ```

# Features

- **Use your own storage**: If you store data in NFS or S3, just use the storage you already use.
- **No additional server required**: ArtiV is a CLI tool. No server or gateway is required to install or operate.
- **Multiple backend support**: Currently, we support local, NFS (by local repo), and s3. And more in the future

- **Reproducible**: A commit is stored in a single file and cannot be changed. There is no way to add/remove/modify a single file in a commit.
- **Expose your data publicly**: Expose your repository with a public HTTP endpoint, then you can download your data in this way
  ```
  art get -o /tmp/dataset https://mybucket.s3.ap-northeast-1.amazonaws.com/path/to/my/data@v0.1.0
  ```
- **Smart storage and transfer**: For the same content of files, there is only one instance stored in the artifact repository. If a file has been uploaded by other commits, no upload is required because we know the file is already there in the repository. Under the hood, we use [content-addressable storage](https://en.wikipedia.org/wiki/Content-addressable_storage) to put the objects.

# Documentation
Please see [documentation](docs/index.md)



