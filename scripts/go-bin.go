package main

import (
  "bufio"
  "bytes"
  "compress/bzip2"
  "flag"
  "io"
  "log"
  "strconv"
  "strings"
  "os"
  "os/exec"
  "path/filepath"
)

var storedFiles =map[string][]byte {}

func GetStored(fn string)(r io.Reader){
  bb:=bytes.NewBuffer(storedFiles[fn])
  r= bzip2.NewReader(bb)
  return
}

func printCode(fn string){
  fd, err:=os.Create(fn)
  if err!=nil {
    log.Fatal(err)
  }
  defer fd.Close()
  io.WriteString(fd, 
`package main
import("io";"bytes";"compress/bzip2")
var storedFiles = map[string][]byte {
`)
  for k,v:= range storedFiles{
    io.WriteString(fd,`"`+k+`":{`)
    strs:=make([]string,len(v))
    for i,s:=range v{
      strs[i]=strconv.Itoa(int(s))
    }
    io.WriteString(fd,strings.Join(strs,","))
    io.WriteString(fd,`},`+"\n")
  }
 io.WriteString(fd,` }

func GetStored(fn string)(r io.Reader){
  bb:=bytes.NewBuffer(storedFiles[fn])
  r= bzip2.NewReader(bb)
  return
}
`)
}

func readFile(fn string, info os.FileInfo, err error) error{
    if ! info.Mode().IsRegular() {
      return nil
    }
    log.Println("retrieving: "+fn)
    var bb bytes.Buffer
    cmd:=exec.Command("bzip2","-c",fn)
    cmd.Stdout=bufio.NewWriter(&bb)
    cmd.Run()
    storedFiles[fn]=bb.Bytes()
    return nil
}

func main(){
  outFile:=flag.String("o","output.go","output go file for embedded resources")
  flag.Parse()
  for _,fn := range flag.Args() {
    filepath.Walk(fn,readFile)
  }
  printCode(*outFile)
}
