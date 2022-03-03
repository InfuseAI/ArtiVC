# ArtiV

ArtiV (**Arti**fact **V**ersions) is a version control system for large files. 

To store and share large files, we may use NFS or object storage (e.g. s3, MinIO). However, if we would like to do versioning on top of them, it is not a trivial thing. ArtiV is a CLI tool to enable you to version files on your storage without pain. You don't need to install any additional server or gateway and we turn your storage into the versioned repository.

[![asciicast](https://asciinema.org/a/XP7AWyrCszebm1KuF3PcYfnRK.svg)](https://asciinema.org/a/XP7AWyrCszebm1KuF3PcYfnRK?autoplay=1)

# Quick Start

## Install ArtiV
Download the latest ArtiV in the [release](https://github.com/InfuseAI/ArtiV/releases) page

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

# Usage
ArtiV is a version control system for large files. It treats a folder of the local filesystem or a prefix of object store as a repository. Unlike other version control systems which have a version control server to control the versions, there is no server in ArtiV. We interact with the file system or object store to do all version control operations.

Here is the basic concept of ArtiV
- **Repository**: the remote storage to store the versioned data. Currently, we support local and s3 repositories. Just like the remote repository in git.
- **Workspace**: the local directory to tackle data. Just like the local workspace in git.
- **Commit**: a snapshot of a set of data. It contains the commit message, timestamp, and the list of blobs.
- **Blob**:  a single object of data.
- **Tag**: an alias of a commit.


## Commands

```
AritV is a version control system for large files

Usage:
  art [command]

Examples:
  # Push data to the repository
  cd /path/to/my/data
  art init s3://mybucket/path/to/repo
  art push

  # Pull data from the repository
  cd /path/to/download
  art init s3://mybucket/path/to/repo
  art pull

  # Download by quick command
  art get -o /path/to/download s3://mybucket/path/to/repo

  # Show command help
  art <command> -h

  For more information, please check https://github.com/infuseai/artiv

Basic Commands:
  init        Initiate a workspace
  config      Configure the workspace
  pull        Pull data from the repository
  push        Push data to the repository
  tag         List or manage tags
  list        List files of a commit
  log         Log commits
  diff        Diff workspace/commits/references

Quick Commands (Download or upload without a workspace):
  get         Download data from a repository
  put         Upload data to a repository

Other Commands:
  version     Print the version information
  docs        Generate docs
  help        Help about any command
  completion  Generate the autocompletion script for the specified shell

Flags:
  -h, --help   help for art
```

Please see [the command reference](./docs/commands/art.md)

## Cheatsheet
![](./docs/assets/cheatsheet.png)


## Configurations
### S3 repository
Prepare the `~/.aws/credentials` to access the s3 repository. Please see the [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

# Use Cases

## Case 1: Dataset Preparation of an Image Classification Project

1. Create a dataset folder and use subfolders as image labels
1. Initiate the workspace.
   ```
   art init s3://mybucket/datasets/flowers-classification   
   ```
1. Push your first release
   ```
   art push -m 'first version'
   ```   
1. Clean the dataset, and move the wrong-classified data
1. Push the dataset again
   ```
   # See what data will be pushed
   art push -m 'my second version' --dry-run
   # Push
   art push -m 'my second version' 
   ```
1. If there are new version is pushed by others, sync the data set with remote
   ```
   # Check the difference
   art pull --dry-run
   # Sync with remote
   art pull --sync
   # or use the merge mode
   # art pull --merge
   ```
1. tag the version
   ```
   art push
   art tag v0.1.0
   ```
   and see the change
   ```
   art log
   ```

1. Use the dataset in the other machine
   ```
   art init s3://mybucket/datasets/flowers-classification
   art pull
   ```   

## Case 2: Run an Experiment
Here we use three repositories
- Dataset for training
- A base model for transfer learning of fine-tuning
- Experiment output. Versioned by the timestamp of an experiment.
In this use case, we use `get` and `put` commands to simplify the commands for programmatic use cases.

1. Clone the training code
   ```
   git clone https://github.com/my-org/my-ml-project.git
   cd my-ml-project
   ```
1. Download the dataset
   ```
   art get -o dataset/ s3://mybucket/datasets/flowers-classification@v0.1.0
   ```
1. Download the base model   
   ```
   art get -o base/ s3://mybucket/models/my-base-model@v0.3.0
   ```   
1. Train and output your training result (trained model, experiment log, hyperparams, etc) to `artifacts/` folder
   ```
   python ./train.py
   ``` 
1. Upload the artifacts 
   ```
   art put artifacts/ s3://mybucket/experiments/project1@202220303-100504
   ```


## Case 3: Expose the dataset by HTTP
1. [Make your S3 bucket public](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/?nc1=h_ls)
1. Copy the public URL of your repository. For example
   ```
   https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```
1. Download the data
   ```
   art get -o /tmp/output https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```   

# How it works
![](./docs/assets/art-overview.png)

There are three major types of data are stored in the artifact repository
- **Blob**: the actual object to be uploaded to object storage. The path is determined by the **SHA1** hash of the content. We use the same path strategy as [git](https://git-scm.com/book/zh-tw/v2/Git-Internals-Git-Objects).
- **Commit**: whenever the client pushes a commit to the repository, it creates a commit object to the store. It contains the timestamp, message, and the list of blobs. A commit is also stored at the path of the content hash. It makes it impossible to change the content because the hash would be invalid. The content of a commit is a gzip-compressed JSON content.
- **Reference**: References are the named tags or time strings to link to a commit. It makes it possible to do the versioning. A special kind of reference `latest` is used by default whenever the client pushes a commit to an artifact repository.

The sample file layout in the artifact repository
```
<artifact repository root>
├── commits
│   ├── 428f898a6792f37cf6805776098387783fdcaa87
│   └── b8d996e998197f8be65f0a0e6ceaa2c7b26bfd11
├── objects
│   ├── 03
│   │   └── cfd743661f07975fa2f1220c5194cbaff48451
│   ├── 23
│   │   └── 715d393e1b723b249642a586597426885f64d1
│   ├── 2b
│   │   └── 77b0d6a5f3bd047fb2c6e44bbb22822f773f94
│   ├── 32
│   │   └── e2eb1b2f4b757fcc3e112d57ca702c65dad526
│   ├── 6d
│   │   └── 7cb1f44cb598a0db93f912791a1c20e3bd7c6b
│   ├── 84
│   │   └── 44869206da2a25b8ee79a8959ed4b34144f2a8
│   ├── ac
│   │   └── 9bcb803e59cd45610d87f2b683319420e7f76b
│   └── d7
│       └── 153b344c84ae187e2a894688310f1956dc45b7
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


# ArtiV and MLOps

Machine learning application is a data-driven approach to solving problems. The MLOps is a methodology to has a feedback loop from data, model training, evaluation, model publish, deploy, monitoring. There are three core components in an ML application.

- Code
- Datasets
- Models

The code can be training code, application code. Mostly, it is versioned by git and we have been familiar with the way to version it. And we also use git as the single source of truth to drive the whole DevOps lifecycle.

However, for datasets and models, there is still no defacto solution to version them. Usually, these data are stored in cloud object storage, on-premise object storage like MinIO, or NFS. There is still a gap between data storage and version metadata storage. Here is why we would like to build the *ArtiV*.

In addition, we are thinking about how to drive the automation when an artifact store event is triggered. In git, we can trigger a job whenever a git event. In the artifact store, we lack the fundamentals to trigger this event. *art* reference the *git* design and provide the *commits* and *references* primitives to make it possible to define a commit or a version that is created. It makes it possible to listen to the object storage or the file system event to trigger an automation job accordingly.

# Alternatives

- **Object storage or NFS:**  Create folders or prefixes to add version semantics. It is straightforward to store different versions of data. But it lacks commit message, metadata, and history tracking. And we cannot know which one is the latest version.
- **S3 versioning**: [S3 Versioning](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Versioning.html) can achieve the object-level versioning. We can get the latest, but possible to roll back to the previous version.
- **Git LFS**: [Git LFS](https://git-lfs.github.com/) is an open-source Git extension for versioning large files developed by Github.
   - [Github](https://docs.github.com/en/repositories/working-with-files/managing-large-files/about-git-large-file-storage), [Gitlab](https://docs.gitlab.com/ee/topics/git/lfs/index.html), [Huggingface](https://huggingface.co/docs/transformers/model_sharing#repository-features) provides Git LFS feature with limited quota.
   - If you want to put data on your own storage, you need to install a [standalone Git LFS server](https://github.com/git-lfs/lfs-test-server)
- **DVC**: [DVC](https://dvc.org/) is built to make ML models shareable and reproducible. It is designed to handle large files, data sets, machine learning models, and metrics as well as code.
    - use `git` command to version small files or metadata, use `dvc` to manage large files.
    - you need to know both git and dvc. In the workflow, the two commands should switch back and forth. See the [dvc tutorial](https://dvc.org/doc/use-cases/versioning-data-and-model-files/tutorial)
- **LakeFS**: [LakeFS](https://lakefs.io/) provides a multi-server solution to make s3 as git-like repositories
    - the [architecture](https://docs.lakefs.io/understand/architecture.html) is much heavier than ArtiV due to an extra database for metadata storage and S3 configuration. (e.g., S3 gateway)


