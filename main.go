package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	illusts_url        = "https://www.pixiv.net/artworks/"
	punctuation_string = `!"#$%&'()*+,-./:;<=>?@*【】[\]^_{|}~`
	dirLoc             = "D:/Pixiv/"
)

var (
	u, err            = url.Parse("http://127.0.0.1:7890") // 代理（这里用的是自己的代理
	profile_url       = "https://www.pixiv.net/ajax/user/%s/profile/all"
	Illustrators_maps = make(map[string]interface{})
	Illustrators      = []string{"2353373", "21479436", "20728711", "30236169", "145944"} //仅仅测试
	illusts_id        = make([]string, 0)

	// 使用代理
	client = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
	}
)

func main() {

	var wg sync.WaitGroup

	for index := range Illustrators {
		// fmt.Println(Illustrators[index])
		profile_url = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/profile/all", Illustrators[index])
		res, username := found_illusts_id(illusts_id, profile_url) //得到作品ID的数组切片
		// remove_punctuation_string(username)
		// println(username)
		username = remove_punctuation_string(username) //移除特殊字符
		check_DirIsExist(username)
		time.Sleep(5 * time.Second)
		var length int = len(res)
		wg.Add(1)
		go Designated_Illustrator_Download(string(Illustrators[index]), length, res, &wg, username)
	}

	wg.Wait()
	fmt.Println("爬取完成")
}

// 指定画师ID下载
func Designated_Illustrator_Download(Illustrator_id string, Artwork_Length int,
	res []string, wg *sync.WaitGroup, username string) {

	defer wg.Done()
	for i := Artwork_Length - 1; i != 0; i-- {
		resp := get_illusts_info(illusts_url, res[i]) //获取插图的body返回信息
		// fmt.Println(res[i], i)
		regexp, _ := regexp.Compile(`"original":"(.*?)"`)
		img_link := regexp.FindString(resp)        //获取插图原图的链接信息
		img_url := strings.Split(img_link, `"`)[3] //获取插图链接
		title := get_illusts_title(resp)           //获得插图的标题
		title = remove_punctuation_string(title)
		time.Sleep(1 * time.Second)
		downloadImg(img_url, title, username) //下载图片
	}
}

// 去除特殊字符
func remove_punctuation_string(context string) string {
	for index := range punctuation_string {
		remove := string(punctuation_string[index])
		context = strings.Replace(context, remove, "", -1)
	}
	return context
}

// 查看目标目录是否存在并创建
func check_DirIsExist(DirName string) {
	if _, err := os.Stat(dirLoc + DirName); os.IsNotExist(err) {
		DirName = remove_punctuation_string(DirName)
		os.Mkdir(dirLoc+DirName, 0755)
		time.Sleep(1 * time.Second)
	}
}

// 设置请求头爬取
func http_Client_SetAndDo(AimUrl string) (http_resp *http.Response) {
	req, _ := http.NewRequest("GET", AimUrl, nil)
	req.Header.Add("referer", "https://www.pixiv.net/")
	resp, _ := client.Do(req)
	return resp
}

// 下载图片 需要传入图片的URL链接 以及标题名
func downloadImg(img_url string, title string, username string) {
	title = strings.Trim(title, `"`)
	fmt.Println(title, username)
	resp := http_Client_SetAndDo(img_url) //
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	// 写出数据
	err = ioutil.WriteFile(dirLoc+username+"/"+title+".png", bytes, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

// 获取插图的标题
func get_illusts_title(html_body string) string {
	regexp, _ := regexp.Compile(`"illustTitle":"(.*?)"`)
	illustTitle := regexp.FindString(html_body)
	title := strings.Split(illustTitle, ":")[1]
	return title
}

// 获取插图地址的网页信息
func get_illusts_info(aimurl string, id string) string {
	aimurl = aimurl + id

	resp := http_Client_SetAndDo(aimurl)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))
	return string(body)
}

// 寻找画师ID
func found_illusts_id(illusts_id []string, aimurl string) ([]string, string) {
	resp := http_Client_SetAndDo(aimurl)
	time.Sleep(1 * time.Second)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	ctx := StrToMap(string(body))
	// fmt.Println(ctx)
	username := ctx["body"].(map[string]interface{})["pickup"].([]interface{})[0].(map[string]interface{})["userName"].(string)

	ctx2 := ctx["body"].(map[string]interface{})["illusts"]

	for key := range ctx2.(map[string]interface{}) {
		illusts_id = append(illusts_id, key)
	}
	return illusts_id, username
}

// string转map
func StrToMap(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}
