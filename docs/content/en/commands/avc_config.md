## avc config

Configure the workspace

### Synopsis

Configure the workspace. The config file is stored at ".avc/config".

```
avc config [<key> [<value>]]
```

### Examples

```
  # List the config
  avc config

  # Get the config
  avc config repo.url

  # Set the config
  avc config repo.url s3://your-bucket/data
```

### Options

```
  -h, --help   help for config
```

### SEE ALSO

* [avc](/commands/avc/)	 - ArtiVC is a version control system for large files

###### Auto generated by spf13/cobra on 24-Mar-2022