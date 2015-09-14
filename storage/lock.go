package storage
import (
    "os"
    "path"
    "io/ioutil"
    "log"
    "errors"
    "io"
    "crypto/md5"
    "fmt"
)

type Lock interface {
    Lock()      error   // 锁定文件
    Unlock()    error   // 解锁
    Locked()    bool    // 检测是否加锁
    Recovery()  error   // 恢复
}

func NewLock(lckpath, lckname, needlckfile string) (*lockfile, error) {
    // 检测锁文件是否存在
    lckfile := path.Join(lckpath, lckname)

    _, err := os.Stat(lckfile)

    if os.IsExist(err) {

        var lc lckcontent

        if err = read(lckfile, &lc); err != nil {

            return nil, nil
        }

        return &lockfile{

            name: lckname,

            path: lckpath,

            needlckfi: needlckfile,

            content: lc,

        }, nil

        //return &l, nil
    }

    // 读取需要锁定的文件内容

    _, err = os.Stat(needlckfile)

    if os.IsNotExist(err) {
        panic("file is not exist")
    }

    bytes, err := ioutil.ReadFile(needlckfile)

    if err != nil {
        fmt.Printf("read file error: %s\n", err)
        return nil, err
    }

    fileMd5, err := FileMd5(needlckfile)

    if err != nil {
        fmt.Printf("file md5 error: %s\n", err)
        return nil, err
    }

    var lckcont = lckcontent{bytes,fileMd5}

    return &lockfile{

        name: lckname,

        path: lckpath,

        needlckfi: needlckfile,

        content: lckcont,

    }, nil
}

func NewLckStor(stor Storage, lockpath ...string) (*lockfile, error) {

    needlckfile := path.Join(stor.storpath, stor.name + ".json")

    lckname := stor.name

    lckpath := stor.storpath

    if len(lockpath) > 0 {
        lckpath = lockpath[0]
    }

    // 检测锁文件是否存在
    lckfile := path.Join(lckpath, lckname)

    _, err := os.Stat(lckfile)

    if os.IsExist(err) {

        var lc lckcontent

        if err = read(lckfile, &lc); err != nil {

            return nil, nil
        }

        return &lockfile{

            name: lckname,

            path: lckpath,

            needlckfi: needlckfile,

            content: lc,

        }, nil

        //return &l, nil
    }

    // 读取需要锁定的文件内容

    _, err = os.Stat(needlckfile)

    if os.IsNotExist(err) {
        panic("file is not exist")
    }

    bytes, err := ioutil.ReadFile(needlckfile)

    if err != nil {
        fmt.Printf("read file error: %s\n", err)
        return nil, err
    }

    fileMd5, err := FileMd5(needlckfile)

    if err != nil {
        fmt.Printf("file md5 error: %s\n", err)
        return nil, err
    }

    var lckcont = lckcontent{bytes,fileMd5}

    return &lockfile{

        name: lckname,

        path: lckpath,

        needlckfi: needlckfile,

        content: lckcont,

    }, nil
}

type lckcontent struct {
    filebytes   []byte
    md5         []byte
}

type lockfile struct {
    name        string
    path        string
    needlckfi   string
    content     lckcontent
}

func (l *lockfile) Lock() error {
    // 检测是否存在 lck 文件，如果存在，先解锁
    if l.Locked() {
        // 为了文件安全，先检测文件 MD5 是否相同

        same, err := l.checkMd5()

        if err !=nil {
            panic(err)
        }

        if !same {
            return errors.New("file is locked!")
        }

        l.Unlock()
    }

    var lckc lckcontent
    err := write(l.lckfile(), lckc)

    return err
}

func (l *lockfile) Unlock() error {
    return os.Remove(l.lckfile())
}

func (l *lockfile) Locked() bool {

    _, err := os.Stat(l.lckfile())

    return os.IsExist(err)
}

func (l *lockfile) Recovery() error {
    // 检测文件 MD5 是否相同

    same, err := l.checkMd5()

    if err != nil {
        panic(err)
    }

    if !same {
        // 获取文件内容
        var lc lckcontent
        content := lc.filebytes

        // 把内容写入文件
        err = ioutil.WriteFile(l.needlckfi, content, os.ModePerm)

        return err
    }

    return nil
}

// 获取 lock 文件路径
func (l *lockfile) lckfile() string {
    lckfile := path.Join(l.path, l.name + ".lck")
    return lckfile
}

// 检测文件 MD5
func (l *lockfile) checkMd5() (bool, error) {
    // 读取 lck 文件

    var lc lckcontent
    err := read(l.lckfile(), &lc)

    if err != nil {
        log.Printf("read lckfile error: %s\n", err)
        return false, err
    }

    //然后计算需要加锁文件 MD5
    filemd5, err := FileMd5(l.needlckfi)

    if err != nil {
        return false, err
    }

    if string(lc.md5) != string(filemd5) {
        return false, nil
    }

    return true, nil
}

// 获取文件 Md5 值
func FileMd5(fi string) ([]byte, error) {
    file, err := os.Open(fi)
    if err != nil {
        log.Printf("open file error : %s\n", err)
        return nil, err
    }
    md5f := md5.New()
    io.Copy(md5f, file)
    return md5f.Sum(nil), nil
}