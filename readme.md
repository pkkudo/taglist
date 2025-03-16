**no longer works with the deprecated api on docker hub**

# taglist

A go script which returns the latest tag which includes x.y.z version format from the specified docker hub repository.

# build/install

<https://go.dev/doc/tutorial/compile-install>

```shell
go build
# ./taglist will be created
```

# usage

```shell
$ ./taglist --help
Usage of ./taglist:
  -all
        Write the list of tags found in the ./output.txt file
  -exclude string
        Exclude pattern
  -filter string
        Filter to use to look for the latest tag including the specified pattern
  -repo string
        Docker Hub Repository such as alpine, busybox, and jupyter/base-notebook

$ ./taglist -repo python
3.13.0a6-windowsservercore-ltsc2022
$ ./taglist -repo python -exclude windows
3.13.0a6-slim-bullseye

$ ./taglist -repo jupyter/base-notebook
x86_64-notebook-7.0.6

$ ./taglist -repo alpine -all
3.19.1
$ cat tags.txt
latest
edge
20240329
20240315
3.19.1
3.19
3.18.6
3.18
3.17.7
3.17
3.16.9
3.16
3
20231219
3.19.0
3.18.5
3.17.6
3.16.8
3.15.11
3.15
3.18.4
20230901
3.18.3
3.17.5
3.16.7
3.15.10
3.18.2
3.17.4
3.16.6
3.15.9
```
