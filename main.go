package main

import (
	"Go-Pixiv/pixiv"
	"Go-Pixiv/utils"
	"fmt"
	"sync"
	"time"
)

var (
	illusts_id = make([]string, 0)
)

func main() {
	Illustrators_Config_map, _, RootDir_Path := utils.Config_init()
	var wg sync.WaitGroup
	for key, _ := range Illustrators_Config_map {
		// fmt.Println(Illustrators[index])
		profile_url := fmt.Sprintf(pixiv.Profile_api, Illustrators_Config_map[key])
		res := pixiv.Found_illusts_id(illusts_id, profile_url) //得到作品ID的数组切片
		username := utils.Remove_Punctuation_String(key)       //移除特殊字符
		time.Sleep(1 * time.Second)
		var length int = len(res)
		wg.Add(1)
		go pixiv.Designated_llustrator_Download(RootDir_Path, Illustrators_Config_map[key], length, res, &wg, username)
	}

	wg.Wait()
	fmt.Println("爬取完成")
}

// // 指定画师ID下载
// func Designated_Illustrator_Download(Illustrator_id string, Artwork_Length int,
// 	res []string, wg *sync.WaitGroup, username string) {

// 	defer wg.Done()
// 	for i := Artwork_Length - 1; i != 0; i-- {
// 		resp := get_illusts_info(illusts_url, res[i]) //获取插图的body返回信息
// 		// fmt.Println(res[i], i)
// 		regexp, _ := regexp.Compile(`"original":"(.*?)"`)
// 		img_link := regexp.FindString(resp)        //获取插图原图的链接信息
// 		img_url := strings.Split(img_link, `"`)[3] //获取插图链接
// 		title := get_illusts_title(resp)           //获得插图的标题
// 		title = utils.Remove_Punctuation_String(title)
// 		time.Sleep(1 * time.Second)
// 		downloadImg(img_url, title, username) //下载图片
// 	}
// }
