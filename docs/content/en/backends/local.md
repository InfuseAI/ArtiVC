---
title: Local Filesystem
weight: 1
---

Use the local filesystem as the repository backend. The folder can reside in a local disk or a mounted NFS.

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