## avc put

Upload data to a repository

```
avc put [-m <message>] <dir> <repository>[@<tag>]
```

### Examples

```
  # Upload the latest version
  avc put ./folder/ /path/to/mydataset

  # Upload the specific version
  avc put ./folder/ /path/to/mydataset@v1.0.0
```

### Options

```
  -h, --help             help for put
  -m, --message string   Commit meessage
```

### SEE ALSO

* [avc](/commands/avc/)	 - ArtiVC is a version control system for large files

###### Auto generated by spf13/cobra on 24-Mar-2022
