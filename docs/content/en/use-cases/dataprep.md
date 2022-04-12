---
title: Dataset Preparation
weight: 2
---

Organizing dataset can be a hassle, especially as data is constantly evolving. ArtiVC is the most suitable tool to organize the dataset. There are the following benefits.

- No need to transfer files with the existing content. Even you rename or copy to different folder. ArtiVC knows they are the same content. It is common to move or keep the same images, videos when the dataset is evolving.
- Version tagging. If there is a stable version of dataset, we can tag a commit as the human-readable version.

## Prepare a dataset

Here are the common steps to prepare a dataset

1. Create a dataset folder and use subfolders as image labels
1. Initiate the workspace.
   ```shell
   avc init s3://mybucket/datasets/flowers-classification   
   ```
1. Push your first release
   ```shell
   avc push -m 'first version'
   ```   
1. Clean the dataset, and move the wrong-classified data
1. Push the dataset again
   ```shell
   # See what data will be pushed
   avc status
   # Push
   avc push -m 'my second version' 
   ```
1. If there are new versions is pushed by others, sync the data set with remote
   ```shell
   # Check the difference
   avc pull --dry-run
   # Sync with remote
   avc pull
   # or use the delete mode
   # avc pull --delete --dry-run
   # avc pull --delete
   ```
1. tag the version
   ```shell
   avc push
   avc tag v0.1.0
   ```
   and see the change
   ```shell
   avc log
   ```

## Clone the dataset

Use the dataset in the other machine

```shell
avc clone s3://mybucket/datasets/flowers-classification
cd flowers-classification
```  