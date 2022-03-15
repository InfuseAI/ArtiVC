---
title: Configuration
weight: -10
---

{{< toc >}}

# Ignore file
Just like git, you can put a `.artignore` file at the root of workspace to define the excluding list. The rule is a regular expression of path. Here is the example

```
# Each line defines a regular expression rule

# Ignore files
^test$
^path/to/my/file$
\.DS_Store$

# Ignore folders. Use a forward slash at the end
^build/
^path/to/my/folder/
/build/

# Ignore all file with extension '.py'
\.py$
````

# S3 repository
Prepare the `~/.aws/credentials` to access the s3 repository. Please see the [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)
