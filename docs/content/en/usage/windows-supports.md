---
title: Windows Support
weight: 100
---

ArtiVC is supported on Windows platforms by running on the WSL.


> WSL Prerequisites:
>
> You must be running Windows 10 version 2004 and higher (Build 19041 and higher) or Windows 11.


## What's WSL
WSL is Windows Subsystem for Linux. It lets users work on a Linux environment on Windows directly.

Check [official WSL site](https://docs.microsoft.com/en-us/windows/wsl/) for more information.


## How To
You could download ArtiVC linux amd64/arm64 version binary from our [github release page](https://github.com/InfuseAI/ArtiVC/releases) depends on your platform architecture.

After extracting ArtiVC binary from the archive, what you need to do is executing it on linux command line interface as usual.

- extract the binary from archive
```shell
$ tar zxvf ArtiVC-v0.4.0-linux-arm64.tar.gz
```

- execution
```shell
$ ./avc version
version.BuildInfo{"Version":"v0.4.0","GitCommit":"febfe3bc579d77130570ba7d12fcf404326b0f7d","GitTreeState":"clean","GoVersion":"go1.17.8"}
```
