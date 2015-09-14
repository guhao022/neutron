package storage

import (
    "os"
    "io/ioutil"
    "encoding/json"
    "bufio"
)

func New(storpath, filename string) (*Storage, error) {

    // 检测文件夹是否存在   若不存在  创建文件夹
    if _, err := os.Stat(storpath); err != nil {

        if os.IsNotExist(err) {

            err = os.Mkdir(storpath, os.ModePerm)

            if err != nil {
                return nil, err
            }
        } else {
            return nil, err
        }
    }

    return &Storage{storpath: storpath, name: filename}, nil
}

type Storage struct {
    storpath    string
    name        string
}

// 获取文件信息
func (sto Storage) Get(value interface{}) error {
    return read(sto.name + ".json", value)
}

// 缓存文件
func (sto Storage) Store(value interface{}) error {
    return write(sto.name + ".json", value)
}


func getFile(path string) (*os.File, error) {
    f, err := os.OpenFile(path, os.O_RDWR, 0666)
    if err != nil {
        if _, ok := err.(*os.PathError); ok {
            return os.Create(path)
        }
        return nil, err
    }
    return f, nil
}

func read(path string, value interface{}) error {
    f, err := getFile(path)
    defer f.Close()

    if err != nil {
        return err
    }

    return json.NewDecoder(bufio.NewReader(f)).Decode(&value)
}

func write(path string, value interface{}) error {
    content, err := json.Marshal(value)

    if err != nil {
        return err
    }
    return ioutil.WriteFile(path, content, os.ModePerm)
}

