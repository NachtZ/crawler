package main
 
import (
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "errors"
    "strings"
    "time"

    "os"
	//"syscall"
)

type Crawler struct{
    routChan chan byte
    urlChan chan string
    filterChan chan string
    filter map[string]bool
    field string
    link chan []string
    closeChan chan byte
    start time.Time
}

func NewCrawler() *Crawler{
    return &Crawler{
        routChan:make(chan byte,40),
        urlChan:make(chan string,10000),
        filterChan:make(chan string,100000),
        filter:make(map[string]bool),
        closeChan:make(chan byte),
        link: make(chan []string,10000),
        start:time.Now(),
    }
}

func(c * Crawler)writeLink(){
    filename := "./" + c.field + "_links.txt"
    f,err := os.Create(filename)
    if err!=nil{
        fmt.Println(err)
        close(c.closeChan)
        return
    }
    for{
        select{
            case links :=<- c.link:
            fmt.Fprintln(f,links[0]+" :",len(links)-1)
            for _,str := range links[1:]{
                fmt.Fprintln(f,str)
            }
            case <- c.closeChan:
                f.Close()
        }
    }
}

func(c * Crawler)demon(){
    dida := 0
    for{
        c.showState()
        if len(c.routChan) == 0{
            dida ++
        }else{
            dida = 0
        }
        if dida >=5{
            close(c.closeChan)
            return
        }
        time.Sleep(5 * time.Second)
    }
}

func(c * Crawler)filterUrl(){
    filename := "./"+c.field + ".txt"
    f,err := os.Create(filename)
    if err != nil{
        fmt.Println(err)
        close(c.closeChan)
        return
    }
    defer f.Close()
    for{
        select{
            case url := <- c.filterChan:
                if _,ok := c.filter[url]; ok{
                    continue
                }
                //fmt.Println("filter:",url)
                c.filter[url] = true
                f.Write([]byte(url+"\n"))
                //f.Close()
                //f,_ = os.OpenFile(filename,syscall.O_APPEND,777)
                c.urlChan <- url
            
            case <-c.closeChan:
                return
        }
    }
}

func(c* Crawler)showState(){
    fmt.Println("/***************************************/")
    fmt.Println("Crawler in field:",c.field)
    fmt.Println("Routine Now:",len(c.routChan))
    fmt.Println("Url now found:",len(c.filter))
    fmt.Println("Url num in filter channel:",len(c.filterChan))
    fmt.Println("Url num in waiting channel:",len(c.urlChan))
    fmt.Println("Crawler has run:",time.Now().Sub(c.start))
    fmt.Println("/***************************************/")
}

func (c * Crawler)Get(url string) (content string, err error){
    timeout := time.Duration(3 * time.Second)
    client := http.Client{
        Timeout:timeout,
    }
    resp,err := client.Get(url)
    
    if err != nil{
        return 
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil{
        return
    }
    if resp.StatusCode != 200{
        str := fmt.Sprintf("StatusCode is %d.\n",resp.StatusCode)
        err = errors.New(str)
        return
    }
    content = string(data)
    return
}

func (c *Crawler)Analyser(seed,ctx string){
    //fmt.Println("Crawler:",seed)
    reURL := regexp.MustCompile("<a.*?href=\"(.*?)\"")
    tmpMap := make(map[string]bool)
    matches := reURL.FindAllStringSubmatch(ctx,10000)
    links := []string{seed}
    for _,url := range matches{
        if strings.Contains(url[1],c.field){
                c.filterChan <- url[1]
                if _,ok := tmpMap[url[1]];!ok{
                    tmpMap[url[1]] = true
                    links = append(links,url[1])
                }
        }
    }
    c.link <- links
}

func (c *Crawler)Run(seed,field string){
    c.field = field
    ctx, err := c.Get(seed)
    if err != nil{
        fmt.Println(err)
        return
    }
    c.filter[seed] = true
    go c.filterUrl()
    c.Analyser(seed,ctx)
    go c.demon()
    go c.writeLink()
    count := 0
    var url string
    for {
        select{
            case url = <- c.urlChan:
                c.routChan <- 1
                count ++
                go func(url string){
                     
                    html,err1 := c.Get(url)
                    if err1 != nil{
                        fmt.Println(err1)
                        <-c.routChan
                        return 
                    }
                    c.Analyser(url,html)
                    <-c.routChan
                }(url)
                case <-c.closeChan:
                    fmt.Println("Crawler Done!")
                    c.showState()
                    return 

        }
    }
}

func runCrawler(url,field string){
    c := NewCrawler()
    c.Run(url,field)
}