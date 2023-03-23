package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Punctuation_string = `!"#$%&'()*+,-./:;<=>?@*【】[\]^_{|}~`
)

// 去除特殊字符
func Remove_Punctuation_String(context string) string {
	for index := range Punctuation_string {
		remove := string(Punctuation_string[index])
		context = strings.Replace(context, remove, "", -1)
	}
	return context
}

// 查看目标目录是否存在并创建
func Check_Root_DirIsExist(DirName string) {
	if _, err := os.Stat(DirName); os.IsNotExist(err) {
		os.Mkdir(DirName, 0755)
		time.Sleep(2 * time.Second)
	} else {
		return
	}
}

func Check_Illustrators_DirIsExist(RootDir string, Illustrators_DirName string) {
	if _, err := os.Stat(RootDir + Illustrators_DirName); os.IsNotExist(err) {
		Illustrators_DirName = Remove_Punctuation_String(Illustrators_DirName)
		os.Mkdir(RootDir+Illustrators_DirName, 0755)
		time.Sleep(2 * time.Second)
	} else {
		return
	}
}

// 读写具体结构不清楚的JSON
func Read_NotSure_Json(FileName string, json_data map[string]interface{}) {
	// 打开文件
	file, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// 解析JSON数据
	errors := json.Unmarshal([]byte(data), &json_data)
	if errors != nil {
		panic(errors)
	}
}

func Read_list_Json(FileName string) (Unmarshal_Data map[string][]string) {
	var json_data map[string][]string
	file, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	errors := json.Unmarshal([]byte(data), &json_data)
	if errors != nil {
		panic(errors)
	}
	return json_data
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

// 判断获取到的数组与本地读出来的数组不同处 产生新切片用于下载
func Compare_exits_PID(illustsId_res []string, loacl_illustsId []string) (UndowloadFile_slice []string) {
	s := make(map[string]bool)
	New_slice := make([]string, 0)
	for _, v := range loacl_illustsId {
		s[v] = true
	}
	for _, v := range illustsId_res {
		if s[v] == false {
			New_slice = append(New_slice, v)
		} else if s[v] {
			continue
		}
	}
	return
}

// 判断画师名称是否存在 不存在则写入 "name" : []的数据格式
func compare_exits_Illustrator(IllustratorName string, storaged_illusts_json map[string][]string) (write_to_loacljson map[string][]string) {
	_, ok := storaged_illusts_json[IllustratorName]
	if ok {
		return nil
	} else {
		//添加到待写入的map里
		write_to_loacljson[IllustratorName] = make([]string, 0)
	}
	return
}

// 创建目录 以及整合Illustrators的map[画师名][画师ID]数据和画师名称数组 返回
func Config_init() (map[string]string, []string, string) {
	Illustrators_Config_map := make(map[string]string)
	Illustrators_Config_List := make([]string, 0)
	config_value := Read_Config_Json()
	RootDir := config_value["data"].(map[string]interface{})["storage path"].(string)
	Illustrators_Info := config_value["data"].(map[string]interface{})["Illustrators"].(map[string]interface{})
	for Name := range Illustrators_Info {
		Illustrators_Config_map[Name] = Illustrators_Info[Name].(string)
		Illustrators_Config_List = append(Illustrators_Config_List, Name)
	}
	Check_Pixiv_DirIsExist(RootDir, Illustrators_Config_List)

	return Illustrators_Config_map, Illustrators_Config_List, RootDir
}

// 读取Config.json配置文件数据
func Read_Config_Json() map[string]interface{} {
	content, err := ioutil.ReadFile("configs\\config.json")
	if err != nil {
		fmt.Println("文件打开失败 [Err:%s]", err)
	}
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

// 检查Pixiv根目录以及 按画师名的目录是否均创建
func Check_Pixiv_DirIsExist(RootDir string, IllustratorsList []string) {
	if _, err := os.Stat(RootDir); os.IsNotExist(err) {
		os.Mkdir(RootDir, 0755)
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println("指定根目录已创建")
	}

	for _, SubDir := range IllustratorsList {
		SubDir = Remove_Punctuation_String(SubDir)
		if _, err := os.Stat(RootDir + SubDir); os.IsNotExist(err) {
			os.Mkdir(RootDir+SubDir, 0755)
			time.Sleep(2 * time.Second)
		}
	}
	fmt.Println("目录检查 完")
}
