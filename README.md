# containers-from-scratch-cli

**Who is this repo for?** If you've written a Dockerfile, built an image, run a container and now want to know what's going on under the hood, this repo is for you! 

In this repo, is a Go CLI that mimic the Docker run and build commands. After building the cli and executing its commands, I encourage you to poke around the code base or watch the livestreams to indulge your curiosity about how containers work. :)

Step 1: Build the containerCLI

```bash
go build . -o containercli
```

Step 2: Run a container with


A few notes:
- Container isolation is achieved by the use of a Linux Kernel feature called Namespaces.
- Resource constraint is achieved by the use of cgroups.




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
