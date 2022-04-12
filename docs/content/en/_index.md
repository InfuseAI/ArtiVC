---
title: 
geekdocNav: false
geekdocBreadcrumb: false
geekdocAlign: center
geekdocAnchor: false
---

{{< columns >}}
### ArtiVC (Artifact Version Control) is a version control system for large files.

Do you need to backup your data regularly? Does keeping summarizing and organizing various dataset versions take up your day's productivity? ArtiVC is a handy command-line tool. With only one command, it helps you neatly snapshot your data and tidily switch among different versions of the data. Even better, it seamlessly integrates your existing cloud environment. ArtiVC supports three major cloud providers (AWS S3, Google Cloud Storage, Azure Blob Storage) or stores data in the remote filesystem using SSH. ArtiVC unleashes your performance on your most important jobs with no pain.

<--->
[![asciicast](https://asciinema.org/a/6JEhzpJ5QMiSkiC74s5CyT257.svg)](https://asciinema.org/a/6JEhzpJ5QMiSkiC74s5CyT257?autoplay=1)
{{< /columns >}}

{{< button size="large" relref="usage/getting-started/" >}}Getting Started{{< /button >}}

## Feature Overview

{{< columns >}}
### Data Versioning

Version your data like versioning code. ArtiVC supports commit history, commit message, and version tag. You can diff two commits, and pull data from the speciifc version.

<--->

### Use your own storage

We are used to putting large files in NFS or S3. To use ArtiVC, you can keep putting your files on the same storage without changes.

<--->

### No additional server is required

ArtiVC is a CLI tool. No server or gateway is required to install and operate.

{{< /columns >}}

{{< columns >}}

### Multiple backends support

ArtiVC natively supports local filesystem, remote filesystem (by SSH), AWS S3, Google Cloud Storage, and Azure Blob Storage as backend. And 40+ backends are supported through [Rclone](backends/rclone/) integration. [Learn more](backends/)

<--->

### Painless Configuration

No one like to configure. So we leverage the original configuration as much as possible. Use `.ssh/config` for ssh access, and use `aws configure`, `gcloud auth application-default login`, `az login` for the cloud platforms.

<--->

### Efficient storage and transfer

The file structure of the repository is stored and transferred efficiently by [design](design/how-it-works/). It prevents storing duplicated content and minimum the number of files to upload when pushing a new version. [Learn more](design/benchmark/)


{{< /columns >}}
