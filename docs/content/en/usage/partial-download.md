---
title: Partial Download
weight: 3
---

ArtiVC supports not only whole repository download, but also partial files download in the repository.

You can download partially:
- by specifying multiple files.
- by specifying wildcard like `.gitignore` pattern.

For more details, please check the [pattern format](https://git-scm.com/docs/gitignore#_pattern_format) in the git document.


## Usage
```shell
# get
avc get -o output repo -- path/to/file1 path/to/file2 data/

# pull
avc pull -- path/to/partial
avc pull v0.1.0 -- path/to/partial ...
```

Note: if you would like to do partial download with `avc pull`, please specify the "--" flag terminator.


## SEE ALSO

* [avc get](/commands/avc_get/)	 - Download data from a repository
* [avc pull](/commands/avc_pull/)	 - Pull data from the repository
