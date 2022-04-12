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

# Configuration
Here we describe how to configure credentials to access the remote backend. The principle of ArtiVC is "Use your tool's config". This allows you to access among the toolchains painlessly.

{{<tabs "Configuration">}}
{{<tab "Local Filesystem">}}
No configuration required
{{</tab>}}

{{<tab "Remote Filesystem (SSH)">}}
1. Configure the `~/.ssh/config`
   ```bash
   Host myserver
      HostName myserver.hosts
      User myname
      IdentityFile ~/.ssh/id_ed25519
   ```    
1. Check if you can access the SSH server successfully
   ```
   ssh myserver
   ```

For more information, please see the [Remote Filesystem (SSH) backend](../../backends/ssh)   
{{</tab>}}

{{<tab "AWS S3">}}
1. Install the [AWS CLI](https://aws.amazon.com/cli/)
2. Configure the AWS CLI
   ```
   aws configure
   ```
3. Check current config
   ```
   aws configure list      
   ```

For more information, please see the [AWS S3 backend](../../backends/s3)

{{</tab>}}

{{<tab "Google Cloud Storage">}}
1. Install the [glcoud CLI](https://cloud.google.com/sdk/gcloud)
2. Login the application default credential
   ```
   gcloud auth application-default login  
   ```
3. Check the current credential is available
   ```
   gcloud auth application-default print-access-token 
   ```

For more information, please see the [Google Cloud Storage backend](../../backends/gcs)
{{</tab>}}

{{<tab "Azure Blob Storage">}}
1. Install the [Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli)
2. Login the Azure CLI
   ```
   az login
   ```
3. Check the login status
   ```
   az account show
   ```

For more information, please see the [Azure Blob Storage](../../backends/azureblob)   
{{</tab>}}
{{</tabs>}}

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

   {{<tabs "avc init">}}
   {{<tab "Local Filesystem">}}
   ```shell
   # in /tmp/artivc/workspace
   avc init /tmp/artivc/repo
   ```
   {{</tab>}}

   {{<tab "Remote Filesystem (SSH)">}}   
   ```shell
   # in /tmp/artivc/workspace
   avc init <host>:path/to/repo
   ```   
   {{</tab>}}

   {{<tab "AWS S3">}}
   ```shell
   # in /tmp/artivc/workspace
   avc init s3://<bucket>/path/to/repo
   ```
   {{</tab>}}

   {{<tab "Google Cloud Storage">}}
   ```shell
   # in /tmp/artivc/workspace
   avc init gs://<bucket>/path/to/repo
   ```
   {{</tab>}}

   {{<tab "Azure Blob Storage">}}
   ```shell
   # in /tmp/artivc/workspace
   avc init https://<storageaccount>.blob.core.windows.net/<container>/path/to/repo
   ```
   {{</tab>}}

   {{</tabs>}}

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
   
   {{<tabs "avc clone">}}
   {{<tab "Local Filesystem">}}
   ```shell
   cd /tmp/artivc/
   avc clone /tmp/artivc/repo another-workspace   
   ```
   {{</tab>}}

   {{<tab "Remote Filesystem (SSH)">}}   
   ```shell
   cd /tmp/artivc/
   avc clone <host>:path/to/repo
   ```
   {{</tab>}}

   {{<tab "AWS S3">}}
   ```shell
   cd /tmp/artivc/
   avc clone s3://<bucket>/path/to/repo
   ```
   {{</tab>}}

   {{<tab "Google Cloud Storage">}}
   ```shell
   cd /tmp/artivc/
   avc clone gs://<bucket>/path/to/repo
   ```
   {{</tab>}}

   {{<tab "Azure Blob Storage">}}
   ```shell
   cd /tmp/artivc/
   avc clone https://<storageaccount>.blob.core.windows.net/<container>/path/to/repo
   ```
   {{</tab>}}

   {{</tabs>}}

   Then the workspace is created, and the data is downloaded.

1. See the commit log
   ```shell
   cd another-workspace/
   avc log
   ```

## Download data

1. Download the latest version

   {{<tabs "avc get">}}
   {{<tab "Local Filesystem">}}
   ```shell
   avc get -o /tmp/artivc/dl-latest /tmp/artivc/repo
   ```
   {{</tab>}}

   {{<tab "Remote Filesystem (SSH)">}}   
   ```shell
   avc get -o /tmp/artivc/dl-latest <host>:path/to/repo
   ```
   {{</tab>}}

   {{<tab "AWS S3">}}
   ```shell
   avc get -o /tmp/artivc/dl-latest s3://<bucket>/path/to/repo
   ```
   {{</tab>}}

   {{<tab "Google Cloud Storage">}}
   ```shell
   avc get -o /tmp/artivc/dl-latest gs://<bucket>/path/to/repo
   ```
   {{</tab>}}

   {{<tab "Azure Blob Storage">}}
   ```shell
   avc get -o /tmp/artivc/dl-latest https://<storageaccount>.blob.core.windows.net/<container>/path/to/repo
   ```
   {{</tab>}}

   {{</tabs>}}
    
    check the content
    ```shell
    ls /tmp/artivc/dl-latest
    ```

1. Or download the specific version

    {{<tabs "avc get version">}}
   {{<tab "Local Filesystem">}}
   ```shell
   avc get -o /tmp/artivc/dl-v0.1.0 /tmp/artivc/repo@v0.1.0
   ```
   {{</tab>}}

   {{<tab "Remote Filesystem (SSH)">}}   
   ```shell
   avc get -o /tmp/artivc/dl-v0.1.0 <host>:path/to/repo@v0.1.0
   ```
   {{</tab>}}

   {{<tab "AWS S3">}}
   ```shell
   avc get -o /tmp/artivc/dl-v0.1.0 s3://<bucket>/path/to/repo@v0.1.0
   ```
   {{</tab>}}

   {{<tab "Google Cloud Storage">}}
   ```shell
   avc get -o /tmp/artivc/dl-v0.1.0 gs://<bucket>/path/to/repo@v0.1.0
   ```
   {{</tab>}}

   {{<tab "Azure Blob Storage">}}
   ```shell
   avc get -o /tmp/artivc/dl-v0.1.0 https://<storageaccount>.blob.core.windows.net/<container>/path/to/repo@v0.1.0
   ```
   {{</tab>}}

   {{</tabs>}}

    check the content
    ```shell
    ls /tmp/artivc/dl-v0.1.0
    ```