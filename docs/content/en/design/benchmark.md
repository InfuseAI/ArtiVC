---
title: Benchmark
weight: 2
---

{{<toc>}}

This document shows the performance benchmark of ArtiV and alternatives. We choose these three tools to compare

- [AWS CLI](https://aws.amazon.com/cli/) (2.4.18): AWS CLI is the AWS official tool to upload/download and manage s3. We use it as the baseline to compare with other tools.
- [Rclone](https://rclone.org/) (v1.57.0) RClone is a command line program to manage files on cloud storage. It is just like rsync, but for cloud storage.
- [DVC](https://dvc.org/) (v2.9.5): DVC is an open-source tool for data science and machine learning projects. It provides Git-like commands to version large files.

We upload data to the same bucket and path. The following are the example commands for each tool.


{{<tabs "tool">}}

{{<tab "AWS CLI">}}
```bash
# upload
aws s3 cp --recursive . s3://art-vcs/benchmark/

# download
aws s3 cp --recursive s3://art-vcs/benchmark/ .
```
{{< /tab >}}

{{<tab "Rclone">}}
```bash
# upload
rclone copy . s3:art-vcs

# download
rclone copy s3:art-vcs .
```
{{< /tab >}}

{{<tab "DVC">}}
```bash
# init
git init
dvc init
dvc remote add -d s3 s3://art-vcs/benchmark/

# upload
dvc add data
dvc push

# download
dvc pull
```
{{< /tab >}}

{{<tab "ArtiV">}}
```bash
# init
art init s3://art-vcs/benchmark/

# upload
art push

# download
art pull
```
{{< /tab >}}

{{< /tabs >}}

Testing Environment

{{<tabs "environment">}}

{{<tab "Hardware">}}

```shell
Hardware Overview:
  Model Name:	MacBook Pro
  Model Identifier:	MacBookPro18,3
  Chip:	Apple M1 Pro
  Total Number of Cores:	10 (8 performance and 2 efficiency)
  Memory:	32 GB
```  
{{< /tab >}}

{{<tab "Network">}}

```shell
$ speedtest
   Speedtest by Ookla

     Server: Chunghwa Mobile - Hsinchu (id = 18450)
        ISP: Chunghwa Telecom
    Latency:     5.25 ms   (0.59 ms jitter)
   Download:   224.54 Mbps (data used: 198.9 MB)
     Upload:    79.58 Mbps (data used: 136.3 MB)
Packet Loss: Not available.
 Result URL: https://www.speedtest.net/result/c/4dd6b2ec-61d1-4b7b-a179-811951412997
```
{{< /tab >}}

{{<tab "S3">}}
The region is `ap-northeast-1`
{{< /tab >}}

{{< /tabs >}}

# Test Plans

We test these three scenarios

1. **[Transfer large size files](#transfer-large-size-files)**
2. **[Transfer large amount of files](#transfer-large-amount-of-files)**
3. **[Determining which files to upload and download](#determining-which-files-to-upload)**

## Transfer large size files

In this test case, we create 10 x 100mb files under `data/` folder. S3 supports [multipart upload](https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html). If the client supports multipart upload and supports parallel upload and download, the result should be similar. We also try to transfer twice, the second trial is to test if the tool can know there is existing content at the same location.

Prepare Data
```bash
# create 10 x 100mb files in data/
seq 10 | xargs -I {} dd if=/dev/urandom of=data/100m_{} bs=100m count=1
```

Commands
{{<tabs "test1">}}
{{<tab "AWS CLI">}}

```bash
# upload
time aws s3 cp --recursive . s3://art-vcs/benchmark/

# download
time aws s3 cp --recursive s3://art-vcs/benchmark/ .
```
{{</tab>}}

{{<tab "Rclone">}}

```bash
# upload
time rclone copy --progress . s3:art-vcs/benchmark/

# download
time rclone copy --progress s3:art-vcs/benchmark/ .
```
{{</tab>}}

{{<tab "DVC">}}


```bash
# init 
git init
dvc init
dvc remote add -d s3 s3://art-vcs/benchmark/

# track data and compute hash
dvc add data

# upload
time dvc push 

# download
rm -rf .dvc/cache
time dvc pull
```
{{</tab>}}

{{<tab "ArtiV">}}


```bash
# init
art init s3://art-vcs/benchmark/large-files/artiv/

# upload
time art push

# download
time art pull
```
{{</tab>}}
{{</tabs>}}

Result

![](../images/benchmark1.svg)

|  | Upload (1st / 2nd)| Download (1st / 2nd) |
| --- | --- | --- |
| AWS CLI | 1m43s / 1m43s | 1m53s / 1m52s |
| Rclone | 1m59s / 2.6s | 2m22s / 2.6s |
| DVC | 1m44s / 1s | 1m27s / <1s |
| ArtiV | 1m43s / <1s | 1m30s / <1s |


For the first trial of download, all four tools have a similar result. But for the second trial, all tools except AWS CLI would check the content hash from local to remote. ArtiV has the best performance because the file list and hashes are stored in ONE commit object. That is, there is only one API request required for this test case.

## Transfer large amount of files

In this test case, we create 1,000 x 1kb files under `data/` folder. It is quite common for a dataset to contain many small files. Because the file size is small, if the tool can reduce the number of API requests, it would get better performance.

Prepare data

```bash
# create 1000 files with 1k size
seq 1000 | xargs -I {} dd if=/dev/urandom of=data/1k_{} bs=1k count=1
```


Commands
{{<tabs "test2">}}

{{<tab "AWS CLI">}}

```bash
# upload
time aws s3 cp --recursive . s3://art-vcs/benchmark/

# download
time aws s3 cp --recursive s3://art-vcs/benchmark/ .
```
{{</tab>}}

{{<tab "Rclone">}}

```bash
# upload
time rclone copy --progress . s3:art-vcs/benchmark/

# download
time rclone copy --progress s3:art-vcs/benchmark/ .
```
{{</tab>}}

{{<tab "DVC">}}

```bash
# init 
git init
dvc init
dvc remote add -d s3 s3://art-vcs/benchmark/

# track data and compute hash
dvc add data

# upload
time dvc push

# download
rm -rf .dvc/cache
time dvc pull
```
{{</tab>}}

{{<tab "ArtiV">}}

```bash
# init
art init s3://art-vcs/benchmark/

# upload
time art push

# download
time art pull
```
{{</tab>}}
{{</tabs>}}

Result

![](../images/benchmark2.svg)

|  | Upload (1st / 2nd) | Download (1st / 2nd) |
| --- | --- | --- |
| aws cli | 16s / 16s | 10s / 10s |
| rclone | 51s / 12s | 20s / 12s |
| dvc | 20s / 1s | 18s / <1s |
| artiv | 12s / <1s | 8s / <1s |


For the first trial, RClone gets the worst performance because it own have four threads (transfers) by default. We can use `--transfers 10` option to increase the transfer speed. ArtiV has the best performance for the first trial because there is an optimization for the first push if there is no commit in the remote repository. It will upload all the files without a content check.

For the second trial, ArtiV and DVC have an efficient way to know no transfer is required. RClone would check all the content hash one by one.

## Determining which files to upload

To transfer files is costly. If we can know the files are existing in local or remote in advance, we can prevent unnecessary transfer. This test is to understand how much time to take for difference checking. This test is inspired by the article [Cloud Data Sync Methods and Benchmark: DVC vs Rclone](https://dvc.org/blog/dvc-vs-rclone).

To compare the local and remote files, there are two methods

1. Check the content hash with remote for each local file. The method would get better performance if there are few files in local and a lot of files in the remote.
2. List the remote files and compare the remote hash with local files. The performance of the method is linearly related to the amount of data on the remote.

There are three scenarios to test

1. 50000 local files and no remote files
1. 500 local files and 50000 remote files
1. 1 local file and 50000 remote files


Prepare data

```bash
# create 50,000 files with 1k size
seq 50000 | xargs -I {} dd if=/dev/urandom of=data/1k_{} bs=1k count=1
```

Commands
{{<tabs "test3">}}
{{<tab "Rclone">}}


```bash
time rclone copy --dry-run s3:art-vcs/benchmark/ .
```
{{</tab>}}
{{<tab "Rclone (no traverse)">}}


```bash
time rclone copy --dry-run —no-traverse s3:art-vcs/benchmark/ .
```
{{</tab>}}
{{<tab "DVC">}}

```bash
time dvc status -c
```
{{</tab>}}
{{<tab "ArtiV">}}

```bash
time art push --dry-run
```
{{</tab>}}
{{</tabs>}}

Result

![](../images/benchmark3.svg)
    |  | local 50000 / remote 0 | local 500 / remote 50000 | local 1 / remote 50000 |
    | --- | --- | --- | --- |
    | rclone | 3s | 16s | 11s |
    | rclone (—no-traverse) | 7m48s | 6.6s | 2s |
    | dvc | 3s | 6.6s | 2.2s |
    | artiv | 1.1s | <1s | <1s |


Just like the [blog post](https://dvc.org/blog/dvc-vs-rclone), DVC uses an adaptive method to query data from remote, but ArtiV still outperforms DVC, why? The reason is DVC uses a `data.dvc` file holding the content hash of a version of the folder, and the file list and md5 hashes are stored in a `<hash>.dir` file. Even though there is the file list stored in the `<hash>.dir`, it still cannot guarantee that all the files are available in the remote. DVC still needs to use one of the two methods to synchronize the status between local and remote.

ArtiV uses another way. ArtiV is a centralized version control system and the commit object is stored in this repository. All push commands should guarantee that all files should be successfully uploaded to the repository and then the commit object can be uploaded to the repository. So if we can get the commit object from the repo, we can say that all the files listed in the commit objects are available in the repo. There is no additional need to check the existence one by one.


# Conclusions

From the benchmark, we know ArtiV has a similar performance as AWS CLI while downloading and uploading the data. Using the commit object, we can easily manipulate the changeset with only one API call, no matter how many objects in local or remote.