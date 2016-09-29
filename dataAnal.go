package main

import(
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

func readFile(path string, link string) map[string]int{
    m := make(map[string]int)
    file,err := ioutil.ReadFile(path)
    if err!= nil{
        fmt.Println("Error in read file:",path)
        return nil
    }
    for i,str := range strings.Split(string(file),"\n"){
        m[str] = i
    }
    file, err = ioutil.ReadFile(link)
    if err != nil{
        fmt.Println("Error in read file:",link)
        return nil
    }
    count := 0
    cur := 0
    node,err := os.Create("linknode.txt")
    if err!= nil{
        fmt.Println("Error in read file: linknode.txt")
        return nil
    }
    for _,str := range strings.Split(string(file),"\n"){
        tmp := ""
        if count == 0{
            fmt.Sscanf(str,"%s :%d",&tmp,&count)
            cur = m[tmp]
            fmt.Println(tmp, cur,count)
        }else{
            fmt.Sscanf(str,"%s",&tmp)
            fmt.Fprintf(node,"%d %d\n",cur,m[tmp])
            fmt.Println(tmp, cur)
            count --
        }
    }
    node.Close()
    return m
}

func runTotal(seed,field string){
    runCrawler(seed,field)
    m := readFile(field+".txt",field+"_links.txt")
    runPageRank(m,"linknode.txt",field+"PGResult.txt")
    runRing("linknode.txt",field+"ringResult.txt")
}

func main(){
    runTotal("http://www.bupt.edu.cn","bupt.edu.cn")
}