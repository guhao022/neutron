package storage

import (
    "time"
    "os"
    "path"
    "errors"
)

func New(filename, author string) (f *File) {
    f.Id = string(RandomCreateBytes(16))
    f.Path = "/"
    f.Name = filename
    f.Extension = "json"
    f.Author = author
    f.Locked = true
    f.LckName = filename + ".lck"
    f.CreateAt = time.Now().Unix()
    f.UpdateAt = time.Now().Unix()
    return
}

type File struct {
    Id          string
    Path        string
    Name        string
    Extension   string
    ContentType string
    Author      string
    Size        int64
    Locked      bool
    LckName     string
    CreateAt    int64
    UpdateAt    int64
}

//新建文件 默认 json
func (this *File) CreateFile(text string) error{
    file_full_name := this.Path + this.Name + "." + this.Extension

    path := path.Dir(file_full_name)
    fi, err := os.Stat(path)
    if os.IsNotExist(err) {
        if err := os.MkdirAll(path, os.ModePerm); err != nil {
            return err
        }
    }
    if !fi.IsDir() {
        return errors.New("path is not a directory")
    }
    tmp, err := os.OpenFile(file_full_name, os.O_WRONLY|os.O_CREATE, os.ModePerm)
    if err != nil {
        return err
    }
    defer tmp.Close()
    if _, err := tmp.Write([]byte(text)); err != nil {
        return err
    }
    return nil
}

func (this *File) unlock() error {
    return os.Remove(this.LckName)
}

func (this *File) lock() error {
    path := path.Dir(this.LckName)
    fi, err := os.Stat(path)
    if os.IsNotExist(err) {
        if err := os.MkdirAll(path, os.ModePerm); err != nil {
            return err
        }
    }
    if !fi.IsDir() {
        return errors.New("path is not a directory")
    }
    lck, err := os.Create(this.LckName)
    if err != nil {
        return err
    }
    lck.Close()
    return nil
}

func (this *File) locked() bool {
    _, err := os.Stat(this.LckName)
    return !os.IsNotExist(err)
}