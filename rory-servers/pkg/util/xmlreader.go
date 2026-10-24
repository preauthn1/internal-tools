/*
将XML读入为map
*/

package util

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type result struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

// Item ?
type Item struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

// ReadXMLFile 讀取xml文件
func ReadXMLFile(fileName string) []Item {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var ret result
	err = xml.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	return ret.Items
}

// ReadXMLFileMap 讀取xml文件到map中
func ReadXMLFileMap(fileName string) map[string]string {
	data := ReadXMLFile(fileName)

	m := make(map[string]string)
	for _, v := range data {
		m[v.Key] = v.Value
	}

	return m
}

// ReadXMLStringMap 讀取XML字符串到map中
func ReadXMLStringMap(data string) map[string]string {
	btdata := []byte(data)
	return ReadXMLMap(btdata)
}

// ReadXMLMap 讀取xml數據到map中
func ReadXMLMap(data []byte) map[string]string {

	var ret result
	err := xml.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	m := make(map[string]string)
	for _, v := range ret.Items {
		m[v.Key] = v.Value
	}

	return m
}
