package main
import (
    "fmt"
    "math"
    "os"
    //"io/ioutil"
    "time"
)
type Node struct{
    Right *Node
    Down *Node
    Col int
    Row int
}

type Matrix struct{
    colHead *Node
    rowHead *Node
    colTail *Node
    num []int
    val []float64
    maxNode int
    size int
} 

func NewMatrix(max int) *Matrix{
    m := &Matrix{
        colHead:new(Node),
        rowHead:new(Node),
        maxNode:-1,
        num:make([]int,max),
        size:max,
    }
    m.colHead.Col = -1
    m.colHead.Row = -1
    m.rowHead.Col = -1
    m.rowHead.Row = -1
    m.colTail = m.colHead
    return m
}

func(m *Matrix)Insert(x,y int){
    now := &Node{
        Right:nil,
        Down:nil,
        Col:x,
        Row:y,
    }
    node,before := m.colHead,m.colHead
    if x>m.maxNode ||y>m.maxNode{
        if x>y{
            m.maxNode = x
        }else{
            m.maxNode = y
        }
        if m.maxNode >=m.size{
            m.size += m.maxNode/10
            t := make([]int,m.size+1-len(m.num))
            m.num = append(m.num,t...)
        }
    }
    m.num[y]++
    if m.colTail.Col <= x{
        if m.colTail.Col == x{
            if m.colTail.Row < y{
                node = m.colTail
            }else{
                m.colTail = now
            }
        }else{
            node = m.colTail
            m.colTail = now
        }
    }
    for node!=nil && node.Right!=nil{
        if node.Col == x{
            if node.Row >y{
                before.Right,now.Down,now.Right = now,before.Right,before.Right.Right
                return
            }
            break
        }
        if node.Col < x && node.Right.Col>x{
            node.Right,now.Right = now,node.Right
            return 
        }
        before = node
        node = node.Right
    }
    if node.Col!= x && node.Right == nil{
        node.Right = now

        return
    }
    
    if node.Down == nil{
        node.Down = now
        return 
    }
    for node!=nil &&node.Down!=nil{
        if node.Row == y{
            return
        }
        if node.Row <y && node.Down.Row >y{
            node.Down,now.Down = now,node.Down
            return
        }
        node = node.Down
    }
    node.Down = now
}

func(m * Matrix)muti(v []float64) (res[]float64){
    res = make([]float64,len(v))
    alpha := 0.15
    alN := alpha/float64(m.maxNode+1)
    idx :=0
    colNode,rowNode := m.colHead,m.colHead
    for colNode = m.colHead.Right;colNode!=nil;colNode = colNode.Right{
        idx = colNode.Col
        for rowNode = colNode;rowNode!=nil;rowNode = rowNode.Down{
            res[idx] += v[rowNode.Row]*m.val[rowNode.Row]
        }
    }
    for i:=0;i<len(res);i++{
        res[i] = alpha*res[i] +alN
    }
    return res
}

func(m * Matrix)check(){
    colNode,rowNode := m.colHead,m.colHead
    for colNode = m.colHead.Right;colNode!=nil;colNode = colNode.Right{
        fmt.Println(colNode.Col,":")
        for rowNode = colNode;rowNode!=nil;rowNode = rowNode.Down{
            fmt.Println(rowNode.Col,rowNode.Row,m.val[rowNode.Row])
        }
    }
    
}

func (m * Matrix)BuildGM(){
    m.val = make([]float64,m.maxNode+1)
    for i:=0;i<=m.maxNode;i++{
        m.val[i] = 1/float64(m.num[i])
    }
}

func (m * Matrix)calVector(p float64)(res []float64){
    v := make([]float64,m.maxNode+1)
    for i:=0;i<=m.maxNode;i++{
        v[i] = 1
    }
    pNow := 1.0
    for pNow >= p{
        pNow = 0
        res = m.muti(v)
        for i:=0;i<len(res);i++{
            pNow += math.Pow(res[i] - v[i],2)
        }
        pNow = math.Sqrt(pNow)
        v = res
    }
    pNow = 0
    for i:=0;i<len(res);i++{
        pNow += res[i] 
    }
    for i:=0;i<len(res);i++{
        res[i]/=pNow
    }
    return res
}

type pairs struct{
    idx int
    val float64
}

func sort(maps map[int]string,val []float64,path string){
    p := []*pairs{} 
    total := 0.0
    for i:=0;i<len(val);i++{
        tmp := &pairs{
            idx:i,
            val:val[i],
        }
        p = append(p,tmp)
    }
    for i:=0;i<len(val);i++{
        for j:=i+1;j<len(val);j++{
            if p[i].val<p[j].val{
                p[i],p[j] = p[j],p[i]
            }
        }
    }
    for i:=0;i<len(p);i++{
        //fmt.Println(p[i].idx,p[i].val)
        total += p[i].val
    }
    //fmt.Println(total)
    if path == ""{
        for i:=0;i<len(p);i++{
        fmt.Println(p[i].idx,p[i].val)
        }
        fmt.Println(total)
    }else{
        file,err := os.Create(path)
        if err!=nil{
            fmt.Println(err)
            return
        }
        for i:=0;i<len(p);i++{
            fmt.Fprintln(file,maps[p[i].idx],p[i].val)
        }
        fmt.Fprintln(file,total)
    }
}

func (m* Matrix)ReadFile(path string){
    file,err := os.Open(path)
    count :=0
    if err!=nil{
        fmt.Println(err)
        return
    }
    defer file.Close()
    x,y := 0,0
    for ;;{
        n,err:=fmt.Fscanf(file,"%d %d",&x,&y)
        if err!=nil ||n != 2{
            fmt.Println(n,err)
            break
        }
        m.Insert(y,x)
        count ++
        if count % 10000 == 0{
            fmt.Println("Read",count,"Data")
        }
        n,err = fmt.Fscanf(file,"%d %d",&x,&y)
    }
    fmt.Println("Read",count,"Data")
}

func runPageRank(maps map[string]int,path1 string, path2 string){
    start := time.Now()
    //arg_num := len(os.Args)
    newMap := make(map[int]string)
    for k,v := range maps{
        newMap[v] = k
    }
    m := NewMatrix(1000)
    //if arg_num <=2{
      //  fmt.Println("Not enough args.")
        //return
    //}
    //fmt.Println(os.Args)
    m.ReadFile(path1)
    m.BuildGM()
    //m.check()
    res := m.calVector(0.00000001)
    sort(newMap,res,path2)
    fmt.Println(time.Now().Sub(start))
}