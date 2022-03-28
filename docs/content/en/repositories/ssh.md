---
title: SSH
weight: 40
---

The SSH repository allows you to use the file system of a remote ssh server as a repository. Under the hood, it uses the `ssh` command to communicate with remote.

## Configuration
To use SSH repository, you have to configure `~/.ssh/config`.  

1. Check `ssh` command is available

   ```
   ssh -V
   ```

1. Configure the `~/.ssh/config`
    ```
    # ~/.ssh/config
    Host ml-vm
        HostName <vm address>
        User myname
        IdentityFile ~/.ssh/id_myname

    Host ml-vm2
        HostName <vm address>
        User myname
        Port 2222
        IdentityFile ~/.ssh/id_myname
    ```

1. Check if you can access the remote without password

    ```
    ssh mydataset
    ```

For more information, please see [ssh config manual](https://linux.die.net/man/5/ssh_config)

## Usage

Init a workspace

```shell
# avc init ssh://<host>/<path>
avc init ssh://ml-vm/home/myname/datasets
```

Clone a repository

```shell
avc clone ssh://ml-vm/home/myname/mydataset
cd mydataset/
```

> Note: the path has to be an absolute path. It is not a relative path to your home directory.