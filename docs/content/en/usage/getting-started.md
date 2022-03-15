---
title: Getting Started
weight: 1
---

This page tells you how to get started with the Geekdoc theme, including installation and basic configuration.

<!--more-->

{{< toc >}}


# Install ArtiV

## Install from Release Page

Download the latest ArtiV in the [release](https://github.com/InfuseAI/ArtiV/releases) page

## Install by homebrew

```shell
brew tap infuseai/artiv 
brew install artiv
```

# Quick Start
## Push data
1. Prepare your data. We put data in the folder `/tmp/artiv/workspace`
   ```shell
   mkdir -p /tmp/artiv/workspace
   cd /tmp/artiv/workspace
   ```

   put the first version of data
   ```shell
   echo a > a
   echo b > b
   echo c > c
   ```

1. Select a folder as the repository. Here we use `/tmp/artiv/repo`
   ```shell
   mkdir -p /tmp/artiv/repo
   ```

1. Init the workspace

   ```shell
   # in /tmp/artiv/workspace
   art init /tmp/artiv/repo
   ```

1. Push the data
   ```shell
   art push
   ```

1. See the commit log
   ```shell
   art log
   ```

1. Add more data to your repository

   ```shell
   echo "hello" > hello
   art status
   ```

   push to the remote
   ```shell
   art push
   ```

1. Tag a version

   ```shell
   art tag v0.1.0
   ```

   See the log
   ```shell
   art log
   ```

## Clone data from exisiting repository
1. Go to the folder to clone repository
   ```shell
   cd /tmp/artiv/
   art clone /tmp/artiv/repo another-workspace   
   ```
   Then the workspace is created, and the data is downloaded.

1. See the commit log
   ```shell
   cd another-workspace/
   art log
   ```

## Download data

1. Download the latest version
    ```shell
    art get -o /tmp/artiv/dl-latest /tmp/artiv/repo
    ```
    
    check the content
    ```shell
    ls /tmp/artiv/dl-latest
    ```

1. Or download the specific version
    ```shell
    art get -o /tmp/artiv/dl_v0.1.0 /tmp/artiv/repo@v0.1.0
    ```

    check the content
    ```shell
    ls /tmp/artiv/dl_v0.1.0
    ```