---
title: Dry Run
weight: 11
---

Pushing and pulling data is time-consuming. And need to be double-checked before transferring. Dry-run is the feature that allows listing the changeset before sending.


## Push

1. Dry run before pushing
    ```shell
    avc push --dry-run
    ```

1. Do the actual push
    ```
    avc push
    ```

## Pull

1. Dry run before pulling
    ```shell
    avc pull -dry-run
    # or check in delete mode
    # avc pull --delete -dry-run
    ```

1. Do the actual pull

    ```shell
    avc pull
    # avc pull --delete
    ```

