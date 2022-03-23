---
title: Expose the dataset
weight: 3
---

ArtiVC repository can be exposed as a http endpoint. In S3, we can just make the bucket and give the data consumer the http endpiont of the repository. In this way, we can download data through CDN or other reverse proxies.

1. [Make your S3 bucket public](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/?nc1=h_ls)
1. Copy the public URL of your repository. For example
   ```shell
   https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```
1. Download the data
   ```shell
   avc get -o /tmp/output https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```   
