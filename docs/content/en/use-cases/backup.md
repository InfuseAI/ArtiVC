---
title: Data Backup/Snapshot
weight: 1
---

Data backup is one of the most common requirement in different scenario. ArtiVC is a very simple tool to backup, or even snapshot, your data in the cloud storage.

## Snapshot the data

1. Init the repository

    ```shell
   avc init s3://mybucket/mydocuments
   ```
1. Snapshot 

   ```
   avc push
   ```
1. Optionally to tag current snapshot as a version
   ```
   avc tag '2022-Q1'
   ```

## Rollback

1. See the snapshot timeline

    ```
    avc log
    ```

1. Rollback. Use `--delete` to delete local files which are not listed in the snapshot version.

   ```
   avc pull --delete 49175d02
   ```

## Get a file from a version

1. Get a file from a given version

   ```
   avc pull 49175d02 -- path/to/my/file
   ```
