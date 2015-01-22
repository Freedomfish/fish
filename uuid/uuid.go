package uuid

import (
    "fmt"
    "time"
    "math/rand"
    "crypto/md5"
    "strconv"
    "io"
)

func NewUuid() string {
    nano := time.Now().UnixNano()
    rand.Seed(nano)
    rndNum := rand.Int63()
    sessionId := Md5(Md5(strconv.FormatInt(nano, 10))+Md5(strconv.FormatInt(rndNum, 10)))
    return sessionId
}

func Md5(text string) string {
    hashMd5 := md5.New()
    io.WriteString(hashMd5, text)
    return fmt.Sprintf("%x", hashMd5.Sum(nil))
}
