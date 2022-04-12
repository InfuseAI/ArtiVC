---
title: Remote Filesystem (SSH)
weight: 2
---

{{< toc >}}

Use remote filesystem through SSH as the repository backend.

## Features

- Use the original `~/.ssh/config` without pain
- Support password authentication
- Support public key authentication
- Support SSH agent
- Support SSH proxy by `ProxyCommand` or `ProxyJump`
- Support host key checking through `~/.ssh/known_hosts`
- Concurrent uploading and downloading

## Configuration
To use SSH backend, you have to configure `~/.ssh/config`.  

1. Configure the `~/.ssh/config`
    ```
    Host myserver
        HostName myserver.hosts
        User myname
        IdentityFile ~/.ssh/id_ed25519
    ```

1. Check if you can access the SSH server successfully

    ```
    ssh myserver
    ```

For more information, please see [ssh config manual](https://man.openbsd.org/ssh_config)

## Usage

Init a workspace

```shell
avc init myserver:path/to/mydataset
```

Clone a repository

```shell
avc clone myserver:path/to/mydataset
cd mydataset/
```

## SSH Proxy

There are two ways to connect to the destination server through bastion (proxy) server.

### ProxyCommand

1. Configure `~/.ssh/config`
    ```
    host bastion
        hostname mybastion.hosts
        user myname
        identityfile ~/id_ed25519

    host myserver
        hostname myserver.hosts
        user myname
        identityfile ~/id_ed25519
        ProxyCommand ssh bastion -W %h:%p
    ```
1. Init the workspace
    ```
    avc init myserver:path/to/mydataset
    ```

### ProxyJump


1. Configure `~/.ssh/config`
    ```
    host bastion
        hostname mybastion.hosts
        user myname
        identityfile ~/id_ed25519

    host myserver
        hostname myserver.hosts
        user myname
        identityfile ~/id_ed25519
        ProxyJump bastion
    ```

1. Init the workspace
    ```
    avc init myserver:path/to/mydataset
    ```

## Environment Variables


| Name | Description | Default value |
| --- | --- | --- |
| `SSH_USER` | The login user | The current user. |
| `SSH_PASSWORD` | The passowrd to be used for password authentication |  |
| `SSH_PORT` | The port of the ssh server | 22 |
| `SSH_IDENTITY_FILE` | The identity file to be used for pubkey authentication |  |
| `SSH_KEY_PASSPHRASE` | The passphrase of the identity key  |  |
| `SSH_STRICT_HOST_KEY_CHECKING` | Set the value to `no` to disable the key checking | |

## Supported Directives for SSH config

- [Port](https://man.openbsd.org/ssh_config#Port)
- [User](https://man.openbsd.org/ssh_config#User)
- [IdentityFile](https://man.openbsd.org/ssh_config#IdentityFile)
- [StrictHostKeyChecking](https://man.openbsd.org/ssh_config#StrictHostKeyChecking)
- [ProxyCommand](https://man.openbsd.org/ssh_config#ProxyCommand)
- [ProxyJump](https://man.openbsd.org/ssh_config#ProxyJump)
