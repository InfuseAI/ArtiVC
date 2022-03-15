---
title: Expose the dataset
weight: 3
---
1. [Make your S3 bucket public](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/?nc1=h_ls)
1. Copy the public URL of your repository. For example
   ```
   https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```
1. Download the data
   ```
   art get -o /tmp/output https://mybucket.s3.ap-northeast-1.amazonaws.com/datasets/flowers-classification
   ```   
