package storage

import (
    "time"
    "crypto/rand"
    r "math/rand"
    "bufio"
    "errors"
    "io"
    "os"
    "path/filepath"
    "regexp"
)


func RandomCreateBytes(n int, alphabets ...byte) []byte {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    var randby bool
    if num, err := rand.Read(bytes); num != n || err != nil {
        r.Seed(time.Now().UnixNano())
        randby = true
    }
    for i, b := range bytes {
        if len(alphabets) == 0 {
            if randby {
                bytes[i] = alphanum[r.Intn(len(alphanum))]
            } else {
                bytes[i] = alphanum[b%byte(len(alphanum))]
            }
        } else {
            if randby {
                bytes[i] = alphabets[r.Intn(len(alphabets))]
            } else {
                bytes[i] = alphabets[b%byte(len(alphabets))]
            }
        }
    }
    return bytes
}

func SelfPath() string {
    path, _ := filepath.Abs(os.Args[0])
    return path
}

func SelfDir() string {
    return filepath.Dir(SelfPath())
}

func FileExists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func SearchFile(filename string, paths ...string) (fullpath string, err error) {
    for _, path := range paths {
        if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
            return
        }
    }
    err = errors.New(fullpath + " not found in paths")
    return
}

// like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
    re, err := regexp.Compile(patten)
    if err != nil {
        return
    }

    fd, err := os.Open(filename)
    if err != nil {
        return
    }
    lines = make([]string, 0)
    reader := bufio.NewReader(fd)
    prefix := ""
    isLongLine := false
    for {
        byteLine, isPrefix, er := reader.ReadLine()
        if er != nil && er != io.EOF {
            return nil, er
        }
        if er == io.EOF {
            break
        }
        line := string(byteLine)
        if isPrefix {
            prefix += line
            continue
        } else {
            isLongLine = true
        }

        line = prefix + line
        if isLongLine {
            prefix = ""
        }
        if re.MatchString(line) {
            lines = append(lines, line)
        }
    }
    return lines, nil
}
