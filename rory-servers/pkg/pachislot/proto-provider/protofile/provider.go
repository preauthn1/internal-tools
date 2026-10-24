package protofile

import (
	"fmt"
	"io/ioutil"
	"os"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
)

// TODO
// "SPRule.bytes":               new(ppb.SPRule_ARRAY),
// "SPCondition.bytes":          new(ppb.SPCondition_ARRAY),
// "Condition.bytes":            new(ppb.Condition_ARRAY),

// Provider 透過 Excel 轉出的 .bytes 檔案轉換
type Provider struct {
	path         string        // 路徑 Folder
	transformers []Transformer // 轉換器
}

// NewProvider create Provider
func NewProvider(path string) (*Provider, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("NewProvider open base path error, %v", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("NewProvider open base path error, %v", err)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("NewProvider error, %s is not a foler", path)
	}

	p := new(Provider)
	p.path = path

	return p, nil
}

// RegisterTransformer 註冊轉換器
func (p *Provider) RegisterTransformer(t Transformer) error {
	path := p.path + "/" + t.File()

	// 檢查檔案是否存在
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("%s is not a file", path)
	}

	p.transformers = append(p.transformers, t)

	return p.load(t)
}

// Load 讀取檔案
func (p *Provider) load(t Transformer) error {
	path := p.path + "/" + t.File()
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(b, t.Item())
	if err != nil {
		return err
	}
	return nil
}

// Load 讀取檔案
func (p *Provider) Load() error {
	for _, t := range p.transformers {
		err := p.load(t)
		if err != nil {
			return err
		}
	}
	return nil
}

// Proto 取得 Proto
func (p *Provider) Proto() (*pachislot.Proto, error) {
	_p := new(pachislot.Proto)
	for _, t := range p.transformers {
		err := t.Transform(_p)
		if err != nil {
			return nil, err
		}
	}
	return _p, nil
}
