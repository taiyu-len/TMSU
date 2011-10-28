package vfs

import (
    "fmt"
    "os"
    "github.com/hanwen/go-fuse/fuse"
)

type FuseVfs struct {
    fuse.DefaultFileSystem
    state *fuse.MountState
}

// factories

func Mount(path string) (*FuseVfs, os.Error) {
    fuseVfs := FuseVfs{}

    state, _, error := fuse.MountPathFileSystem(path, &fuseVfs, nil)
    if error != nil { return nil, error }

    fuseVfs.state = state

    return &fuseVfs, nil
}

// api

func (this *FuseVfs) Unmount() {
    this.state.Unmount()
}

func (this *FuseVfs) Loop() {
    this.state.Loop()
}

func (this *FuseVfs) GetAttr(name string, context *fuse.Context) (*os.FileInfo, fuse.Status) {
    fmt.Println("GetAttr", name)

    switch (name) {
        case "file.txt":
            return &os.FileInfo{
                                   Mode: fuse.S_IFREG | 0644,
                                   Size: int64(len(name)),
                               }, fuse.OK
        case "":
            return &os.FileInfo{
                                   Mode: fuse.S_IFDIR | 0755,
                               }, fuse.OK
    }

    return nil, fuse.ENOENT
}

func (this *FuseVfs) OpenDir(name string, context *fuse.Context) (chan fuse.DirEntry, fuse.Status) {
    fmt.Println("Open dir", name)

    switch name {
        case "": return this.topDirectories()
        case "tags": return this.tagDirectories()
        case "untagged": return this.untaggedFiles()
    }

    return nil, fuse.ENOENT
}

func (this *FuseVfs) Open(name string, flags uint32, context *fuse.Context) (fuse.File, fuse.Status) {
    fmt.Println("Open", name)

    if name != "file.txt" { return nil, fuse.ENOENT }

    if flags & fuse.O_ANYWRITE != 0 { return nil, fuse.EPERM }

    return fuse.NewDataFile([]byte(name)), fuse.OK
}

// implementation

func (this *FuseVfs) topDirectories() (chan fuse.DirEntry, fuse.Status) {
    fmt.Println("topDirectories")

    channel := make(chan fuse.DirEntry, 2)
    channel <- fuse.DirEntry{ Name: "tags", Mode: fuse.S_IFDIR }
    channel <- fuse.DirEntry{ Name: "untagged", Mode: fuse.S_IFDIR }
    close(channel)

    fmt.Println("/topDirectories")

    return channel, fuse.OK
}

func (this *FuseVfs) tagDirectories() (chan fuse.DirEntry, fuse.Status) {
    fmt.Println("tagDirectories")

    channel := make(chan fuse.DirEntry, 1)
    //TODO query db
    //channel <- fuse.DirEntry{ Name: "test1", Mode: fuse.S_IFDIR }
    //channel <- fuse.DirEntry{ Name: "test2", Mode: fuse.S_IFDIR }
    close(channel)

    fmt.Println("/tagDirectories")

    return channel, fuse.OK
}

func (this *FuseVfs) untaggedFiles() (chan fuse.DirEntry, fuse.Status) {
    fmt.Println("untaggedFiles")

    channel := make(chan fuse.DirEntry, 1)
    //TODO query db
    close(channel)

    fmt.Println("/untaggedFiles")

    return channel, fuse.OK
}
