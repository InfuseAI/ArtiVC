---
weight: 2
---

### Q: Why don't we use object storage directly
Yes. we can use object storage to version the artifact using different prefixes. However, it does not support atomic operation for multi-files upload. That is, we cannot know if there is one file missing when uploading, or if there are any files added to the folder. In ArtiV, a commit is an atomic operation to advance one version.

### Q: Why don't we use git
Git is not suitable for [versioning big files](https://stackoverflow.com/questions/17888604/git-with-large-files). Whenever we clone a git repository, it will download ALL files from every branch and commit to the local. 

There are good solutions like [Git LFS](https://git-lfs.github.com/) or [DVC](https://dvc.org/) out there to version large files in the git repository. However, we consider that if we require a dependency of some storage like s3 or NFS, why don't we do the versioning just only on the storage?

### Q: Can we do the versioning with multiple datasets or models?
Yes. Just create different repositories with different prefixes in the same s3 bucket or different folders in NFS

### Q: What storage does ArtiV support?
Currently, we can support local, NFS(by local), and AWS S3. We will support GCS(Google Cloud Storage), Azure Blob Storage in the future.

### Q: Do we support Windows?
Not yet, we will implement and test it soon.