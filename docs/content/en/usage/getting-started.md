---
title: Getting Started
weight: -20
---

This page tells you how to get started with the Geekdoc theme, including installation and basic configuration.

<!--more-->

{{< toc >}}


## Install ArtiV

**Install from Release Page**

Download the latest ArtiV in the [release](https://github.com/InfuseAI/ArtiV/releases) page

**Install by homebrew**

```
brew tap infuseai/artiv 
brew install artiv
```

## Push data
1. Select the bucket and path to put your versioned data. Here we use `s3://mybucket/path/to/repo`. You can use local repository instead (e.g. `/path/to/repo`)
1. Prepare the `~/.aws/credentials` if you use the S3 repository
1. Push your data to the repository
   ```
   cd /path/to/my/data
   art init s3://mybucket/path/to/repo
   art push
   ```

1. See the commit log
   ```
   art log
   ```

1. Add more data to your repository

   ```
   echo "hello" > hello
   art push --dry-run
   art push
   ```

1. Tag a version
   ```
   art tag v0.1.0
   art log
   ```


## Pull data
1. Pull data from the repository
   ```
   cd /path/to/download
   art init s3://mybucket/path/to/repo
   art pull
   ```

1. or download data with single `get` command
    ```
    art get -o /path/to/download s3://mybucket/path/to/repo    
    ```
    or download the specific version
    ```
    art get -o /path/to/download s3://mybucket/path/to/repo@v0.1.0
