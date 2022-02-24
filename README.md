# ArtiV

ArtiV is a version control system for large files. Artifact store is a crucial part of the machine learning project. We use it to store different versions of datasets or models. ArtiV provides a better way to manage your artifact stored in local, NFS, or cloud object storage.

*WARNING: This project is in a very early stage of development. There is still no workable release yet.*

# Usage
Download the latest version in the [release](https://github.com/InfuseAI/ArtiV/releases) page


- Put and get the latest version in one line
  ```
  art put /path/to/my/data s3://bucket/path/to/mydata
  art get -o /tmp/mydata s3://bucket/path/to/mydata
  ```

- Put and get the artifact with a specific tag
  ```
  art put /path/to/my/data s3://bucket/path/to/mydata@v0.1.0
  art get -o /tmp/mydata-v0.1.0 s3://bucket/path/to/mydata@v0.1.0
  ```

- Manage data in a workspace

  ```
  cd /path/to/mydata

  # init a workspace
  art init s3://bucket/path/to/mydata

  # pull latest data from the repository
  art pull

  # See the commit log
  art log
  
  # Push the workspace data to remote
  art push

  # Tag the latest commit as v0.2.0
  art tag v0.2.0
  ```

- For more information please see the command help
  ```
  art -h
  art <subcommand> -h
  ```  

# Features

- **No additional server required**: You only need one command, we turn your storage as versioned storage.
- **Multiple backend support**: Currently, we support local and NFS. And we will support more storage in the future
- **Http compatible**: If you make your s3 bucket public, you can expose your dataset by http
  ```
  art get -o /tmp/dataset https://mybucket.s3.ap-northeast-1.amazonaws.com/path/to/my/data@v0.1.0
  ```
- **Reproducible**: A commit is stored in a single file and cannot be changed. There is no way to add/remove/modify a single file in a commit.
- **Smart storage and transfer**: For the same content of files, there is only one instance stored in the artifact repository. If a file has been uploaded by other commits, no upload is required because we know the file is already there in the repository. Under the hood, we use [content-addressable storage](https://en.wikipedia.org/wiki/Content-addressable_storage) to put the objects.


# How it works
![](assets/art-overview.png)


There are three major types of data are stored in the artifact repository
- **Blob**: the actual object to be uploaded to object storage. The path is determined by the sha1 hash of the content. We use the same path strategy as [git](https://git-scm.com/book/zh-tw/v2/Git-Internals-Git-Objects).
- **Commit**: whenever the client pushes a commit to the artifact repository, it creates a commit object to the store. It contains the timestamp, message, and the list of blobs. A commit is also stored at the path of the content hash. It makes it impossible to change the content because the hash would be invalid.
- **Reference**: References are the named tags or time strings to link to a commit. It makes it possible to do the versioning. A special kind of reference `latest` is used by default whenever the client pushes a commit to an artifact repository.

The sample file layout in the artifact repository
```
<artifact repository root>
├── commits
│   ├── 428f898a6792f37cf6805776098387783fdcaa87
│   └── b8d996e998197f8be65f0a0e6ceaa2c7b26bfd11
├── objects
│   ├── 03
│   │   └── cfd743661f07975fa2f1220c5194cbaff48451
│   ├── 23
│   │   └── 715d393e1b723b249642a586597426885f64d1
│   ├── 2b
│   │   └── 77b0d6a5f3bd047fb2c6e44bbb22822f773f94
│   ├── 32
│   │   └── e2eb1b2f4b757fcc3e112d57ca702c65dad526
│   ├── 6d
│   │   └── 7cb1f44cb598a0db93f912791a1c20e3bd7c6b
│   ├── 84
│   │   └── 44869206da2a25b8ee79a8959ed4b34144f2a8
│   ├── ac
│   │   └── 9bcb803e59cd45610d87f2b683319420e7f76b
│   └── d7
│       └── 153b344c84ae187e2a894688310f1956dc45b7
└── refs
    ├── latest
    └── tags
        ├── v1.0.0    
        └── v1.0.1
```

The sample content of a commit
```
{
  "createdAt": "2022-02-06T00:34:45.406454+08:00",
  "messaage": null,
  "blobs": [
    {
      "path": "README.md",
      "hash": "32e2eb1b2f4b757fcc3e112d57ca702c65dad526",
      "mode": 420,
      "size": 6148
    },
    {
      "path": "front.png",
      "hash": "ac9bcb803e59cd45610d87f2b683319420e7f76b",
      "mode": 420,
      "size": 6130505
    },
    {
      "path": "back.png",
      "hash": "d7153b344c84ae187e2a894688310f1956dc45b7",
      "mode": 420,
      "size": 696686
    },
    {
      "path": "tmp",
      "hash": "03cfd743661f07975fa2f1220c5194cbaff48451",
      "mode": 420,
      "size": 4
    },
    {
      "path": "screen.png",
      "hash": "6d7cb1f44cb598a0db93f912791a1c20e3bd7c6b",
      "mode": 420,
      "size": 305686
    },
    {
      "path": "view/1.png",
      "hash": "8444869206da2a25b8ee79a8959ed4b34144f2a8",
      "mode": 420,
      "size": 578961
    },
    {
      "path": "view/2.png",
      "hash": "2b77b0d6a5f3bd047fb2c6e44bbb22822f773f94",
      "mode": 420,
      "size": 214683
    }
  ]
}
```


# FAQs

### Q: Why don't we use object storage directly
Yes. we can use object store to version the artifact using different prefixes. However, we cannot know if there is one file missing when uploading, or if there are any files added to the folder.

### Q: Why don't we use git
Git is not suitable for [versioning big files](https://stackoverflow.com/questions/17888604/git-with-large-files). Whenever we clone a git repository, it will download ALL files for every branch and commit to the local. 

There are good solutions like [git-lfs](https://git-lfs.github.com/) or [dvc](https://dvc.org/) out there to version large files in the git repository. However, we consider that if we require a dependency of some storage like s3 or NFS, why don't we do the versioning just only on the storage?

### Q: Can we do the versioning with multiple datasets or models?
Yes. Just create different repositories with different prefixes in the same s3 bucket or different folders in NFS

### Q: What storage does ArtiV support?
Currently, all are under development. We hope we can support local, NFS, AWS S3, GCS(Google Cloud Storage), Azure Blob Storage in the future.

### Q: How to version a dataset used by the training code versioned by a git repository?
We will provide a file (maybe named `Artfile`) to download artifacts to the speicified folders. Just like the way we manage the package dependency by `requirements.txt` in a python project.

# ArtiV and MLOps

Machine learning application is a data-driven approach to solving problems. The MLOps is a methodology to has a feedback loop from data, model training, evaluation, model publish, deploy, monitoring. There are three core components in an ML application.

- Code
- Datasets
- Models

The code can be training code, application code. Mostly, it is versioned by git and we have been familiar with the way to version it. And we also use git as the single source of truth to drive the whole DevOps lifecycle.

However, for datasets and models, there is still no defacto solution to version them. Usually, these data are stored in cloud object storage, on-premise object storage like MinIO, or NFS. There is still a gap between data storage and version metadata storage. Here is why we would like to build the *art*.

In addition, we are thinking about how to drive the automation when an artifact store event is triggered. In git, we can trigger a job whenever a git event. In the artifact store, we lack the fundamentals to trigger this event. *art* reference the *git* design and provide the *commits* and *references* primitives to make it possible to define a commit or a version is created. It makes it possible to listen to the object storage or file system event to trigger an automation job accordingly.