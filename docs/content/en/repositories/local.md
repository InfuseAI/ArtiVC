---
title: Local Filesystem
weight: 1
---

The local filesystem repository uses a folder of local fileysystem as the repository. The folder can reside in a local disk or a mounted NFS.

## Usage
Init a workspace

```shell
art init /path/to/mydataset

# Or use relative path
# art init ../mydataset
```

Clone a repository

```shell
art clone /path/to/mydataset
cd mydataset/
```