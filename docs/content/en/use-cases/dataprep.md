---
title: Dataset Preparation
---

Dataset Preparation use case is the most commmon use case in ArtiV. You can prepare the unstructured data and push commit to the remote frequently.

1. Create a dataset folder and use subfolders as image labels
1. Initiate the workspace.
   ```shell
   art init s3://mybucket/datasets/flowers-classification   
   ```
1. Push your first release
   ```shell
   art push -m 'first version'
   ```   
1. Clean the dataset, and move the wrong-classified data
1. Push the dataset again
   ```shell
   # See what data will be pushed
   art status
   # Push
   art push -m 'my second version' 
   ```
1. If there are new version is pushed by others, sync the data set with remote
   ```shell
   # Check the difference
   art pull --dry-run
   # Sync with remote
   art pull
   # or use the delete mode
   # art pull --delete --dry-run
   # art pull --delete
   ```
1. tag the version
   ```shell
   art push
   art tag v0.1.0
   ```
   and see the change
   ```shell
   art log
   ```

## Clone by other users

Use the dataset in the other machine

```shell
art clone s3://mybucket/datasets/flowers-classification
cd flowers-classification
```  