/*
用map来存储节点。
map中存在的node如下：
map[idx]*node
再用map[string]bool来存储已经查找到的环。
在检索的时候，使用dfs来进行检索，在dfs的时候记录下环的路径。（string或者数组）
*/

package Crawler

import (
    "fmt"
    "strconv"
    "strings"
    "os"
)

type node struct{
    idx int
    next []*node
    state int
}

type RF struct{
    nodes []*node
    maxnode int
    size int
    file *os.File
    count int
}

func NewRF(size int)(rf *RF){
    rf = &RF{
        nodes:make([]*node,size),
        size:size,
        maxnode:-1,
    }
    return 
}

func (rf *RF)insert(idx int, links []int){
    if idx >= rf.size{
        t := make([]*node,idx)
        rf.size += idx
        rf.nodes = append(rf.nodes,t...)
        rf.maxnode = idx
    }
    if idx >rf.maxnode{
        rf.maxnode = idx
    }
    if rf.nodes[idx] == nil{
        rf.nodes[idx] = &node{
            idx:idx,
            next:make([]*node,len(links)),
            state:0,
        }
    }else{
        if len(rf.nodes[idx].next)!= 0{
            return
        }
        rf.nodes[idx].next = make([]*node,len(links))   
    }
    tmp := rf.nodes[idx].next
    for i:=0;i<len(links);i++{
        if links[i] >= rf.size{
            t := make([]*node,links[i])
            rf.size += links[i]
            rf.nodes = append(rf.nodes,t...)
            rf.maxnode = links[i]
        }
        if links[i] >= rf.maxnode{
            rf.maxnode = links[i]
        }
        if rf.nodes[links[i]] == nil{
            rf.nodes[links[i]] = &node{
                idx:links[i],
                next:nil,
                state:0,
            }
        }
        //fmt.Println(idx,links[i])
        tmp[i] = rf.nodes[links[i]]
    }
}

func (rf *RF)dfs(idx int,path string){
    if rf.nodes[idx] == nil{
        return
    }
    if rf.nodes[idx].state !=0{

        if path == ""||rf.nodes[idx].state == 1{
            return 
        }
        //here need to cut the trim the ring.
        //first, we need to know that this visited node is the tail and the head of the ring
        //so we should first find the start of the ring, to cut the fist n nodes until the path's first node is this node.
        //then add the tail of ring in the end of path.
        //after that, we get a ring, we check whether the ring is found before and then insert it to the map.
     //   if _,ok := rf.rings[path];ok == true{
     //       fmt.Println("in dfs, 1");
     //       return ;//can't happend in fact
     //   }
        s := "#" + strconv.Itoa(idx) +"#"
        idx := strings.Index(path,s)
        if idx == -1{
            return
        }
        //path = strings.TrimLeft(path,s)
        path = path[idx:] + "#"
        rf.count ++
        fmt.Fprintln(rf.file,rf.count,":",path)
    //    rf.rings[path] = true
        return
    }
    rf.nodes[idx].state = -1
    tmp := rf.nodes[idx].next
    for i:=0;i<len(tmp);i++{
        if tmp[i] == nil{
            break
        }
        rf.dfs(tmp[i].idx,path + "#" + strconv.Itoa(tmp[i].idx))
    }
    rf.nodes[idx].state = 1
}


func (rf *RF)ReadFile(path string){
    file,err := os.Open(path)
    if err != nil{
        fmt.Println(err)
        return
    }
    defer file.Close()
    x,y := -1,0
    n,err :=fmt.Fscanf(file,"%d %d\n",&x,&y)
    if n !=2 || err!=nil{
        fmt.Println(n,err)
        return
    }
    idx :=x
    nums := []int{y}
    for {
        n,err =fmt.Fscanf(file,"%d %d\n",&x,&y)
        if n !=2 || err!=nil{
            fmt.Println("Read Done!",n,err)
            break
        }
        if x != idx{
            rf.insert(idx,nums)
            idx = x
            nums = []int{y}
            if idx%10000 == 0{
                fmt.Println("Read",idx,"Now!")
            }
        }else{
            nums = append(nums,y)
        }
    }
    rf.insert(idx,nums)
}

func(rf *RF)SearchRing(path string){
    file,err := os.Create(path)
    if err !=nil{
        return 
    }
    defer file.Close()
    rf.file = file
    for i:=0;i<=rf.maxnode;i++{
        if rf.nodes[i]!=nil{
            rf.dfs(i,"#" + strconv.Itoa(i))
            rf.nodes[i].state = 1
        }
    }
    /*
    count := 0    
    for key,_ := range rf.rings{
        count ++
        fmt.Fprintln(file,count, ":",key)
    }*/
}
/*
func main2(){
    if len(os.Args) <=1{
        fmt.Println("Need 1 arguments.")
        return
    }
    rf := NewRF(5)
    fmt.Println("Begin read ",os.Args[1]," now.")
    rf.ReadFile(os.Args[1])
    fmt.Println("Read done.")
    rf.SearchRing()
    fmt.Println("Searching done.")
}*/
func runRing(path1 string,path2 string){

    rf := NewRF(5)
    fmt.Println("Begin read ",path1," now.")
    rf.ReadFile(path1)
    fmt.Println("Read done.")
    rf.SearchRing(path2)
    fmt.Println("Searching done.")
}