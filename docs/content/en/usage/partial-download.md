---
title: Partial Download
weight: 3
---

By default, ArtiVC download all files of a version. It also supports to download partial of the files in a commit.

You can download partially:

- by specifying multiple files.
- by specifying wildcard-like `.gitignore` pattern.

For more details, please check the [pattern format](https://git-scm.com/docs/gitignore#_pattern_format) in the git document.

## Usage
```shell
# get
avc get -o output repo -- path/to/file1 path/to/file2 data/

# pull
avc pull -- path/to/partial
avc pull v0.1.0 -- path/to/partial ...
```

Note: if you would like to do a partial download, please put the path list after the "--" flag terminator.


## SEE ALSO

* [avc get](/commands/avc_get/)	 - Download data from a repository
* [avc pull](/commands/avc_pull/)	 - Pull data from the repository
