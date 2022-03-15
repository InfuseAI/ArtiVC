---
weight: 3
---

- **Object storage or NFS:**  Create folders or prefixes to add version semantics. It is straightforward to store different versions of data. But it lacks commit message, metadata, and history tracking. And we cannot know which one is the latest version.
- **S3 versioning**: [S3 Versioning](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Versioning.html) can achieve the object-level versioning. We can get the latest, but possible to roll back to the previous version.
- **Git LFS**: [Git LFS](https://git-lfs.github.com/) is an open-source Git extension for versioning large files developed by Github.
   - [Github](https://docs.github.com/en/repositories/working-with-files/managing-large-files/about-git-large-file-storage), [Gitlab](https://docs.gitlab.com/ee/topics/git/lfs/index.html), [Huggingface](https://huggingface.co/docs/transformers/model_sharing#repository-features) provides Git LFS feature with limited quota.
   - If you want to put data on your own storage, you need to install a [standalone Git LFS server](https://github.com/git-lfs/lfs-test-server)
- **DVC**: [DVC](https://dvc.org/) is built to make ML models shareable and reproducible. It is designed to handle large files, data sets, machine learning models, and metrics as well as code.
    - use `git` command to version small files or metadata, use `dvc` to manage large files.
    - you need to know both git and dvc. In the workflow, the two commands should switch back and forth. See the [dvc tutorial](https://dvc.org/doc/use-cases/versioning-data-and-model-files/tutorial)
- **LakeFS**: [LakeFS](https://lakefs.io/) provides a multi-server solution to make s3 as git-like repositories
    - the [architecture](https://docs.lakefs.io/understand/architecture.html) is much heavier than ArtiV due to an extra database for metadata storage and S3 configuration. (e.g., S3 gateway)