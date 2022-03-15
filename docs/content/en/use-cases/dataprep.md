---
title: Dataset Preparation
---

1. Create a dataset folder and use subfolders as image labels
1. Initiate the workspace.
   ```
   art init s3://mybucket/datasets/flowers-classification   
   ```
1. Push your first release
   ```
   art push -m 'first version'
   ```   
1. Clean the dataset, and move the wrong-classified data
1. Push the dataset again
   ```
   # See what data will be pushed
   art push -m 'my second version' --dry-run
   # Push
   art push -m 'my second version' 
   ```
1. If there are new version is pushed by others, sync the data set with remote
   ```
   # Check the difference
   art pull --dry-run
   # Sync with remote
   art pull --sync
   # or use the merge mode
   # art pull --merge
   ```
1. tag the version
   ```
   art push
   art tag v0.1.0
   ```
   and see the change
   ```
   art log
   ```

1. Use the dataset in the other machine
   ```
   art init s3://mybucket/datasets/flowers-classification
   art pull
   ```  