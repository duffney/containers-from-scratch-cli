# containers-from-scratch-cli
A simple Go CLI that builds Linux containers. Based on Liz Rice's [containers-from-scratch](https://github.com/lizrice/containers-from-scratch/tree/master).


Download Alpine for chroot:

```bash
mkdir alpine
cd alpine
curl -o alpine.tar.gz http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.0-x86_64.tar.gz
tar xvf alpine.tar.gz
rm alpine.tar.gz
touch ALPINE_CONTAINER_ROOT
```
