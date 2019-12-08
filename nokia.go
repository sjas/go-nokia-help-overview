package main

import "crypto/tls"
import . "fmt"
import "io/ioutil"
import "net/http"
import "log"
import "github.com/antchfx/htmlquery"
import "strings"
import "regexp"
import "bufio"
import "github.com/Workiva/go-datastructures/set"
import "reflect"

func main(){
        url:="http://documentation.nokia.com/html/"
        doc,err:=htmlquery.LoadURL(url);if err!=nil{log.Fatal(err)}
        list:=htmlquery.Find(doc,"//a");if err!=nil{log.Fatal(err)}
        urls:=[]string{}
        for _,i:=range list{
                tmp:=Sprintf("%s%s",url,strings.TrimSpace(htmlquery.InnerText(i)))
                urls=append(urls,tmp)
        }
        jsps:=set.New()
        for _,url:=range urls{
                doc,err:=htmlquery.LoadURL(url);if err!=nil{log.Fatal(err)}
                html:=htmlquery.OutputHTML(doc,true)
                res,err:=regexp.MatchString(`jsp`,html);if err!=nil{log.Fatal(err)}
                if res{
                        sc:=bufio.NewScanner(strings.NewReader(html))
                        re:=regexp.MustCompile(`http.*jsp*`)
                        for sc.Scan(){
                                res:=re.FindAllString(sc.Text(),1)
                                if len(res)==1{jsps.Add(res[0])}
                        }
                }
        }
        endresult:=make(map[string]string)
        jsplist:=jsps.Flatten()
        Println(reflect.TypeOf(jsplist))
        for _,jsp:=range jsplist{
                c:=&http.Client{Transport:&http.Transport{TLSClientConfig:&tls.Config{InsecureSkipVerify:true,},},}
                html,err:=c.Get(jsp.(string));if err!=nil{log.Fatal(err)}
                website,err:=ioutil.ReadAll(html.Body);if err!=nil{log.Fatal(err)}
                html.Body.Close()
                doc,err:=htmlquery.Parse(strings.NewReader(string(website)));if err!=nil{log.Fatal(err)}
                titles:=htmlquery.Find(doc,"//title")
                for _,title:=range titles{
                        titlestring:=htmlquery.InnerText(title)
                        endresult[titlestring]=url
                        Println(titlestring)
                        Println(jsp)
                }
        }
        //Print(endresult)
}
