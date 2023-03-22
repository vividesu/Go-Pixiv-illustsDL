package pixiv

import (
	"Go-Pixiv/utils"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 寻找画师ID
func Found_illusts_id(illusts_id []string, aimurl string) []string {
	resp := utils.Http_Client_SetAndDo(aimurl)
	defer resp.Body.Close()
	time.Sleep(1 * time.Second)
	body, _ := io.ReadAll(resp.Body)
	ctx := utils.StrToMap(string(body))
	// username := ctx["body"].(map[string]interface{})["pickup"].([]interface{})[0].(map[string]interface{})["userName"].(string)
	ctx2 := ctx["body"].(map[string]interface{})["illusts"]
	for key := range ctx2.(map[string]interface{}) {
		illusts_id = append(illusts_id, key)
	}

	return illusts_id
}

// 获取插图地址的网页信息
func Get_illusts_info(aimurl string, id string) string {
	aimurl = aimurl + id
	resp := utils.Http_Client_SetAndDo(aimurl)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

// 获取插图的标题
func Get_illusts_title(html_body string) string {
	regexp, _ := regexp.Compile(`"illustTitle":"(.*?)"`)
	illustTitle := regexp.FindString(html_body)
	title := strings.Split(illustTitle, ":")[1]
	return title
}

// 下载图片 需要传入图片的URL链接 以及标题名
func downloadImg(RootDir string, img_url string, title string, username string) {
	title = strings.Trim(title, `"`)
	fmt.Println(title, username)
	resp := utils.Http_Client_SetAndDo(img_url) //
	defer resp.Body.Close()
	// 写出数据
	f, err := os.Create(RootDir + username + "/" + title + ".png")
	if err != nil {
		fmt.Println("os.create err %V \n", err)
	}
	io.Copy(f, resp.Body)

}

// 指定画师ID下载
func Designated_llustrator_Download(RootDir string, Illustrator_id string, Artwork_Length int, res []string, wg *sync.WaitGroup, username string) {
	defer wg.Done()
	for i := Artwork_Length - 1; i != 0; i-- {
		resp := Get_illusts_info(Illusts_api, res[i]) //获插图的body返回信息
		// fmt.Println(res[i], i)
		regexp, _ := regexp.Compile(`"original":"(.*?)"`)
		img_link := regexp.FindString(resp)        //获取插图原图的接信息
		img_url := strings.Split(img_link, `"`)[3] //获插图链接
		title := Get_illusts_title(resp)           //获得插图的标题
		title = utils.Remove_Punctuation_String(title)
		time.Sleep(1 * time.Second)
		downloadImg(RootDir, img_url, title, username) //下载图片
	}
}
