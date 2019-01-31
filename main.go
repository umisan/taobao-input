package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/umisan/taobao/config"
	service "github.com/umisan/taobao/service/config"
)

/*
最初に既存のitem.jsonをitem_listとして読み出す
すでに存在するURLか判断するためにURLのmapを作成
CSVの読み込み(sampleのフォーマットと仮定)
mapを更新しながら新しい商品をitem_listに追加
最後にitem_listを保存して終了
*/

func wait() {
	var temp string
	fmt.Println("終了するにはqを入力してください")
	fmt.Scan(&temp)
}

func main() {
	//item.jsonの読み出し
	db := "item.json"
	var item_list config.ItemList
	item_list.Generate(db)

	//URLのmapの作成
	var url_map map[string]bool
	url_map = make(map[string]bool)
	for _, item := range item_list {
		_, ok := url_map[item.Link]
		if !ok {
			url_map[item.Link] = true
		}
	}

	//csvの読み込み
	fmt.Println("利用するcsvファイル名を入力してください")
	var input_file string
	fmt.Scan(&input_file)
	fmt.Println("検索間隔を入力してください(秒)")
	var duration int
	fmt.Scan(&duration)
	if duration == 0 {
		duration = 1
	}
	csv_byte, err := ioutil.ReadFile(input_file)
	if err != nil {
		log.Println(err)
		wait()
		return
	}
	csv_str := string(csv_byte[:])
	csv_reader := strings.NewReader(csv_str)
	reader := csv.NewReader(csv_reader)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		wait()
		return
	}

	//item_listへの追加
	var index uint = 0
	if len(item_list) != 0 {
		index = item_list[len(item_list)-1].Id + 1
	}
	for i := 1; i < len(records); i++ {
		var new_item config.Item
		new_item.Maker = records[i][0]
		new_item.Number = records[i][1]
		new_item.Name = records[i][2]
		new_item.Stock = "0"
		new_item.Link = records[i][3]
		if _, ok := url_map[new_item.Link]; !ok {
			//アイテム追加処理
			fmt.Println("追加： ", new_item.Link)
			url_map[new_item.Link] = true
			new_item_list, err := service.GenerateNewItems(new_item)
			if err != nil {
				log.Println(err)
				wait()
				return
			}
			for i, _ := range new_item_list {
				new_item_list[i].Id = index
				index++
			}
			item_list = append(item_list, new_item_list...)
		}
		time.Sleep(time.Duration(duration) * time.Second)
	}

	//書き込んで保存
	item_list.WriteData(db)
	fmt.Println("アイテムの追加は正常に終了しました")
	wait()
}
