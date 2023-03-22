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
