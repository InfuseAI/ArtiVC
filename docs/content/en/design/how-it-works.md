---
weight: 1
---
There are three major types of data are stored in the artifact repository
- **Blob**: the actual object to be uploaded to object storage. The path is determined by the **SHA1** hash of the content. We use the same path strategy as [git](https://git-scm.com/book/zh-tw/v2/Git-Internals-Git-Objects).
- **Commit**: whenever the client pushes a commit to the repository, it creates a commit object to the store. It contains the timestamp, message, and the list of blobs. A commit is also stored at the path of the content hash. It makes it impossible to change the content because the hash would be invalid. The content of a commit is a gzip-compressed JSON content.
- **Reference**: References are the named tags or time strings to link to a commit. It makes it possible to do the versioning. A special kind of reference `latest` is used by default whenever the client pushes a commit to an artifact repository.

![](../images/artiv-overview.png)



The sample file layout in the artifact repository
```
<artifact repository root>
├── commits
│   ├── 428f898a6792f37cf6805776098387783fdcaa87
│   └── b8d996e998197f8be65f0a0e6ceaa2c7b26bfd11
├── objects
│   ├── 03
│   │   └── cfd743661f07975fa2f1220c5194cbaff48451
│   ├── 23
│   │   └── 715d393e1b723b249642a586597426885f64d1
│   ├── 2b
│   │   └── 77b0d6a5f3bd047fb2c6e44bbb22822f773f94
│   ├── 32
│   │   └── e2eb1b2f4b757fcc3e112d57ca702c65dad526
│   ├── 6d
│   │   └── 7cb1f44cb598a0db93f912791a1c20e3bd7c6b
│   ├── 84
│   │   └── 44869206da2a25b8ee79a8959ed4b34144f2a8
│   ├── ac
│   │   └── 9bcb803e59cd45610d87f2b683319420e7f76b
│   └── d7
│       └── 153b344c84ae187e2a894688310f1956dc45b7
└── refs
    ├── latest
    └── tags
        ├── v1.0.0    
        └── v1.0.1
```

The sample content of a commit
```
{
  "createdAt": "2022-02-06T00:34:45.406454+08:00",
  "messaage": null,
  "blobs": [
    {
      "path": "README.md",
      "hash": "32e2eb1b2f4b757fcc3e112d57ca702c65dad526",
      "mode": 420,
      "size": 6148
    },
    {
      "path": "front.png",
      "hash": "ac9bcb803e59cd45610d87f2b683319420e7f76b",
      "mode": 420,
      "size": 6130505
    },
    {
      "path": "back.png",
      "hash": "d7153b344c84ae187e2a894688310f1956dc45b7",
      "mode": 420,
      "size": 696686
    },
    {
      "path": "tmp",
      "hash": "03cfd743661f07975fa2f1220c5194cbaff48451",
      "mode": 420,
      "size": 4
    },
    {
      "path": "screen.png",
      "hash": "6d7cb1f44cb598a0db93f912791a1c20e3bd7c6b",
      "mode": 420,
      "size": 305686
    },
    {
      "path": "view/1.png",
      "hash": "8444869206da2a25b8ee79a8959ed4b34144f2a8",
      "mode": 420,
      "size": 578961
    },
    {
      "path": "view/2.png",
      "hash": "2b77b0d6a5f3bd047fb2c6e44bbb22822f773f94",
      "mode": 420,
      "size": 214683
    }
  ]
}
```