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

ArtiVC natively supports local filesystem, remote filesystem (by SSH), S3, GCS as backend. And 40+ backends are supported through [Rclone](backends/rclone/) integration. [Learn more](backends/)

<--->

### Expose your data to public

To serve a repository as a public HTTP endpoint, the repository turn to a http repository right away. Then the data consumer can download your data with an one-line command. [Learn more](use-cases/expose/)

<--->

### Efficient storage and transfer

The file structure of repository is storage and transfer effiecntly by [design](design/how-it-works/). It prevent from storing duplicated content and minimum the round-trip time to determine change set to transfer. [Learn more](design/benchmark/)


{{< /columns >}}
