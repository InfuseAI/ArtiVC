<p align="center">
  <img src="docs/assets/ArtiVC_workspace.svg" width="40%" />
</p>

<p align="center">
	<a href="https://raw.githubusercontent.com/infuseai/artivc/main/LICENSE" >
		<img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="Apache License" /></a>
	<a href="https://github.com/InfuseAI/ArtiVC/actions?query=workflow%3A%22Go+test%22+branch%3Amain++">
		<img src="https://github.com/infuseai/artivc/workflows/Go%20test/badge.svg?branch=main" alt="Go tests status" /></a>
</p>


# ArtiVC

[ArtiVC](https://artivc.io/) (**Arti**facts **V**ersion **C**ontrol) is a handy command-line tool for data versioning on cloud storage. With only one command, it helps you neatly snapshot your data and Switch data between versions. Even better, it seamlessly integrates your existing cloud environment. ArtiVC supports three major cloud providers (AWS S3, Google Cloud Storage, Azure Blob Storage) and the remote filesystem using SSH.

[![asciicast](https://asciinema.org/a/6JEhzpJ5QMiSkiC74s5CyT257.svg)](https://asciinema.org/a/6JEhzpJ5QMiSkiC74s5CyT257?autoplay=1)

Try it out from the [Getting Started](https://artivc.io/usage/getting-started/) guide

# Features

- **Data Versioning**: Version your data like versioning code. ArtiVC supports commit history, commit message, and version tag. You can diff two commits, and pull data from the specific version.
- **Use your own storage**: We are used to putting large files in NFS or S3. To use ArtiVC, you can keep putting your files on the same storage without changes.
- **No additional server is required**: ArtiVC is a CLI tool. No server or gateway is required to install and operate.
- **Multiple backends support**:  ArtiVC natively supports local filesystem, remote filesystem (by SSH), AWS S3, Google Cloud Storage, and Azure Blob Storage as backend. And 40+ backends are supported through [Rclone](https://artivc.io/backends/rclone/) integration. [Learn more](https://artivc.io/backends/)
- **Painless Configuration**:  No one like to configure. So we leverage the original configuration as much as possible. Use `.ssh/config` for ssh access, and use `aws configure`, `gcloud auth application-default login`, `az login` for the cloud platforms.
- **Efficient storage and transfer**:  The file structure of the repository is stored and transferred efficiently by [design](https://artivc.io/design/how-it-works/). It prevents storing duplicated content and minimum the number of files to upload when pushing a new version. [Learn more](https://artivc.io/design/benchmark/)

# Documentation

For more detail, please read the [ArtiVC documentation](https://artivc.io/usage/getting-started/)
