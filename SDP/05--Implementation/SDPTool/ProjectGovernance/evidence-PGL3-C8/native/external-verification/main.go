package main
import("fmt";"os"; lab "pgi3pilot")
func main(){for _,v:=range [][3]string{{"Work","\t Élan\t Case\u00a0","Work: Élan\t Case"},{"Prefix","\n\u2003",""},{"","A",": A"},{"  P ","  MiXeD  ","  P : MiXeD"},{"🚜"," 作業 ","🚜: 作業"}}{if got:=lab.FormatLabel(v[0],v[1]);got!=v[2]{fmt.Printf("FAIL: %q != %q\n",got,v[2]);os.Exit(1)}};fmt.Println("PASS: five external behavior cases")}
