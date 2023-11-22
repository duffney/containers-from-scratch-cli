# containers-from-scratch-cli

Goal: Write 1 paragraph per day.

**Who is this repo for?** _If you've written a Dockerfile, built an image, run a container and now want to know what's going on under the hood, this repo is for you!_ 


Let's run some containers. Or should I say processes. 

> "Containers are just isolated groups of proccess running on a single host" --



## What's docker build doing?

Download Alpine for chroot:

```bash
mkdir alpine
cd alpine
curl -o alpine.tar.gz http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.0-x86_64.tar.gz
tar xvf alpine.tar.gz
rm alpine.tar.gz
touch ALPINE_CONTAINER_ROOT
```

Download Ubuntu for chroot:

```bash
sudo apt install debootstrap
sudo debootstrap jammy ./ubuntu-rootfs http://archive.ubuntu.com/ubuntu/
```

stress mem and cpu

```bash
./stress --vm 1 --vm-bytes 100M
```
