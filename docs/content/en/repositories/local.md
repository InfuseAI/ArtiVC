---
title: Local Filesystem
weight: 1
---

The local filesystem repository uses a folder of local fileysystem as the repository. The folder can reside in a local disk or a mounted NFS.

## Usage
Init a workspace

```shell
avc init /path/to/mydataset

# Or use relative path
# avc init ../mydataset
```

Clone a repository

```shell
avc clone /path/to/mydataset
cd mydataset/
```