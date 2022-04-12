---
title: 
geekdocNav: false
geekdocBreadcrumb: false
geekdocAlign: center
geekdocAnchor: false
---

{{< columns >}}
### ArtiVC (Artifact Version Control) is a version control system for large files.


**rsync** is an ssh-based tool that provides fast incremental file transfer.<br>
**Rclone** is a rsync-like tool for cloud storage.<br>
**ArtiVC** is like Git for files versioning and like Rclone for cloud storage.

<--->
[![asciicast](https://asciinema.org/a/6JEhzpJ5QMiSkiC74s5CyT257.svg)](https://asciinema.org/a/6JEhzpJ5QMiSkiC74s5CyT257?autoplay=1)
{{< /columns >}}

{{< button size="large" relref="usage/getting-started/" >}}Getting Started{{< /button >}}

## Feature Overview

{{< columns >}}
### Data Versioning

Version your data like versioning code. ArtiVC supports commmit history, commit message, version tag. You can diff two commits, pull data from speciifc version.

<--->

### Use your own storage

We are used to putting large files in NFS or S3. To use ArtiVC, you can keep put your files on the same storage without changes.

<--->

### No additional server required

ArtiVC is a CLI tool. No server or gateway is required to install and operate.

{{< /columns >}}

{{< columns >}}

### Multiple backends support

ArtiVC natively supports local filesystem, remote filesystem (by SSH), AWS S3, Google Cloud Storage, Azure Blob Storage as backend. And 40+ backends are supported through [Rclone](backends/rclone/) integration. [Learn more](backends/)

<--->

### Painless Configuration

No one like to configure. So we leverage the original configuraion as much as possible. Use `.ssh/config` for ssh access, and use `aws configure`, `gcloud auth application-default login`, `az login` for the cloud platforms.

<--->

### Efficient storage and transfer

The file structure of repository is storage and transfer effiecntly by [design](design/how-it-works/). It prevents from storing duplicated content and minimum the number of files to upload when pushing a new version. [Learn more](design/benchmark/)


{{< /columns >}}
