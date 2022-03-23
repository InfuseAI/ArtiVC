---
title: Getting Started
weight: 1
---

<!--more-->

{{< toc >}}


# Install ArtiVC

## Install from Release Page

Download the latest ArtiVC in the [release](https://github.com/InfuseAI/ArtiVC/releases) page

## Install by homebrew

```shell
brew tap infuseai/artivc 
brew install artivc
```

# Quick Start
## Push data
1. Prepare your data. We put data in the folder `/tmp/artivc/workspace`
   ```shell
   mkdir -p /tmp/artivc/workspace
   cd /tmp/artivc/workspace
   ```

   put the first version of data
   ```shell
   echo a > a
   echo b > b
   echo c > c
   ```

1. Select a folder as the repository. Here we use `/tmp/artivc/repo`
   ```shell
   mkdir -p /tmp/artivc/repo
   ```

1. Init the workspace

   ```shell
   # in /tmp/artivc/workspace
   avc init /tmp/artivc/repo
   ```

1. Push the data
   ```shell
   avc push
   ```

1. See the commit log
   ```shell
   avc log
   ```

1. Add more data to your repository

   ```shell
   echo "hello" > hello
   avc status
   ```

   push to the remote
   ```shell
   avc push
   ```

1. Tag a version

   ```shell
   avc tag v0.1.0
   ```

   See the log
   ```shell
   avc log
   ```

## Clone data from exisiting repository
1. Go to the folder to clone repository
   ```shell
   cd /tmp/artivc/
   avc clone /tmp/artivc/repo another-workspace   
   ```
   Then the workspace is created, and the data is downloaded.

1. See the commit log
   ```shell
   cd another-workspace/
   avc log
   ```

## Download data

1. Download the latest version
    ```shell
    avc get -o /tmp/artivc/dl-latest /tmp/artivc/repo
    ```
    
    check the content
    ```shell
    ls /tmp/artivc/dl-latest
    ```

1. Or download the specific version
    ```shell
    avc get -o /tmp/artivc/dl_v0.1.0 /tmp/artivc/repo@v0.1.0
    ```

    check the content
    ```shell
    ls /tmp/artivc/dl_v0.1.0
    ```