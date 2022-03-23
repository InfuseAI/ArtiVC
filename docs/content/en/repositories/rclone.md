---
title: Rclone
weight: 50
---

[Rclone](https://rclone.org/) is a command-line program to manage files on cloud storage. As a analogy, **Rclone** is the **rsync** for cloud storage. Rclone supports [40+ providers](https://rclone.org/#providers). The Rclone repository is to use **Rclone** command to communicate with the remote backend.


## Configuration

1. [Install the Rclone](https://rclone.org/install/)
1. Configure the Rclone
    ```
    rclone config
    ```
1. Check if the Rclone commmand can access your remote backend and path.
    ```
    rclone lsd <remote>:<path>
    ```

## Usage

Init a workspace

```shell
# avc init rclone://<remote>/<path>
avc init rclone://myremote/path/to/mydataset
```

Clone a repository

```shell
avc clone rclone://myremote/path/to/mydataset
cd mydataset/
```