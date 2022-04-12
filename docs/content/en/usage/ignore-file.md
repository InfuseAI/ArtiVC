---
title: Ignore File
weight: 12
---

Just like git, you can put a `.avcignore` file at the root of workspace to define the excluding list. The rule is the same as `.gitignore`. For more details, please check the [pattern format](https://git-scm.com/docs/gitignore#_pattern_format) in the git document.

Here is the example:

```shell
# Ignore files
test
path/to/my/file
.DS_Store

# Ignore folders. Use a forward slash at the end
build/
path/to/my/folder/
/build/

# Ignore all file with extension '.py'
*.py
````


