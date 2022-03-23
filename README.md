# ArtiVC

[ArtiVC](https://artivc.netlify.app/) (**Arti**fact **V**ersion **C*ontrl) is a version control system for large files. 

To store and share large files, we may use NFS or object storage (e.g. s3, MinIO). However, if we would like to do versioning on top of them, it is not a trivial thing. ArtiVC is a CLI tool to enable you to version files on your storage without pain. You don't need to install any additional server or gateway and we turn your storage into the versioned repository.

[![asciicast](https://asciinema.org/a/XP7AWyrCszebm1KuF3PcYfnRK.svg)](https://asciinema.org/a/XP7AWyrCszebm1KuF3PcYfnRK?autoplay=1)


Try it out from the [Getting Started](https://artivc.netlify.app/usage/getting-started/) guide

# Features

- **Use your own storage**: If you store data in NFS or S3, just use the storage you already use.
- **No additional server required**: ArtiVC is a CLI tool. No server or gateway is required to install or operate.
- **Multiple backend support**: Currently, we support local, NFS (by local repo), and s3. And more in the future

- **Reproducible**: A commit is stored in a single file and cannot be changed. There is no way to add/remove/modify a single file in a commit.
- **Expose your data publicly**: Expose your repository with a public HTTP endpoint, then you can download your data in this way
  ```
  avc get -o /tmp/dataset https://mybucket.s3.ap-northeast-1.amazonaws.com/path/to/my/data@v0.1.0
  ```
- **Smart storage and transfer**: For the same content of files, there is only one instance stored in the artifact repository. If a file has been uploaded by other commits, no upload is required because we know the file is already there in the repository. Under the hood, we use [content-addressable storage](https://en.wikipedia.org/wiki/Content-addressable_storage) to put the objects.


# Documentation

For more detail, please read the [ArtiVC documentation](https://artivc.netlify.app/usage/getting-started/)