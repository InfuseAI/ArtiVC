---
title: Benchmark
weight: 2
---

{{<toc>}}

This document shows the performance of ArtiV. We select s3 as the backend and upload data to the same bucket and path.


We choose the four tools to compare

- aws cli
- rclone
- dvc
- artiv


All four tools support to use the credentials in `~/.aws/credentials`


# Tools

{{<tabs "tool">}}

{{<tab "AWS CLI (2.4.18)">}}

[AWS CLI](https://aws.amazon.com/cli/) is the AWS official tool to upload/download and manage s3. We use it as the baseline to compare with other tools.

The basic commands to upload and download files

```
# upload
aws s3 cp --recursive . s3://art-vcs/benchmark/

# download
aws s3 cp --recursive s3://art-vcs/benchmark/ .
```
{{< /tab >}}

{{<tab "Rclone (v1.57.0)">}}

[Rclone](https://rclone.org/) is a command line program to manage files on cloud storage. It is just like rsync, but for cloud storage.

The basic commands to upload and download files

```
rclone copy . s3:art-vcs
rclone copy s3:art-vcs .
```
{{< /tab >}}

{{<tab "DVC (v2.9.5)">}}

[DVC](https://dvc.org/) is an open-source tool for data science and machine learning projects. It provides Git-like commands to version large files.


To init a dvc repository

```
git init
dvc init
dvc remote add -d s3 s3://art-vcs/benchmark/
```

The basic commands to upload and download files

```
# upload
dvc add data
dvc push

# download
dvc pull
```
{{< /tab >}}

{{<tab "ArtiV (v0.5.0)">}}

[ArtiV](https://artiv.netlify.app/) is a version control tool to manage large files.

To init a artiv repository

```
art init s3://art-vcs/benchmark/
```

The basic commands to upload and download files

```
# upload
art push

# download
art pull
```
{{< /tab >}}

{{< /tabs >}}

# Environment

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

We will test the three metrics

1. **[Transfer large size files](#transfer-large-size-files)**: To see if it support the multipart upload feature in s3. And if it can benefit from parallel download/upload
2. **[Transfer large amount of files](#transfer-large-amount-of-files)**: Some datasets contain thousands or millions of small files. This test is to see the impact of the number of files.
3. **[Determining which files to upload and download](#determining-which-files-to-upload)**: To transfer files is costly, if we can know the files are existing in local or remote, we can prevent unnecessary transfer. This test is to understand how much time to take for difference checking. This test is inspired by the article [Cloud Data Sync Methods and Benchmark: DVC vs Rclone](https://dvc.org/blog/dvc-vs-rclone).



## Transfer large size files

In this test case, we create 10 x 100mb files under `data/` folder

```bash
# create 10 files with 100M size
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
# Init 
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
# Init
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

## Transfer large amount of files

In this test case, we create 1,000 x 1kb files under `data/` folder

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
# Init 
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
# Init
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

## Determining which files to upload
In this test case, we test the performance to DRY RUN which files should be uploaded and no actual content would be transferred. 

To compare the local and remote files, there are two methods

1. Check the 


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
# Init 
git init
dvc init
dvc remote add -d s3 s3://art-vcs/benchmark/

# track data and compute hash
dvc add data

# upload
dvc push

# status
time dvc status -c
```
{{</tab>}}
{{<tab "ArtiV">}}

```bash
art init s3://art-vcs/benchmark/
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

# Conclusions