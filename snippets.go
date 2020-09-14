package main
import("io";"bytes";"compress/bzip2")
var storedFiles = map[string][]byte {
"snippets/gossip.cfg":{66,90,104,57,49,65,89,38,83,89,172,3,246,129,0,0,55,217,128,0,18,80,3,195,144,46,70,222,0,32,0,84,37,79,84,100,211,67,67,4,244,130,73,53,61,65,232,128,0,194,114,172,156,41,1,1,50,77,110,65,64,93,4,88,44,6,8,96,114,78,22,206,21,104,207,10,24,53,50,38,141,202,219,135,159,55,44,40,245,153,26,35,241,119,36,83,133,9,10,192,63,104,16},
 }

func GetStored(fn string)(r io.Reader){
  bb:=bytes.NewBuffer(storedFiles[fn])
  r= bzip2.NewReader(bb)
  return
}
