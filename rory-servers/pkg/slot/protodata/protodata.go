package protodata

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"example.invalid/golang/protobuf/proto"
)

// DefaultPaths 預設讀取路徑
var DefaultPaths = []string{"./protodata", "../../protodata", "./protoconfig", "../../protoconfig"}

// ProtoOption 各個表的選項
type ProtoOption int

// ProtoOptions
const (
	OptALL ProtoOption = iota
	OptSlotPattern
	OptMainGameTable
	OptSlotPayLine
	OptSlotPayTable
	OptSlotSymbolDistribution
	OptSymbolAttributes
	OptJackpotTable
	OptBonusTable
	OptPickGameTable
	OptRespinTable
	OptFreeSpinTable
	OptSelectionTable
	OptExtraRule
	OptReelBgDistribution
	OptEnd
)

// ProtoData for config data from xls
type ProtoData struct {
	basicPath              string
	opts                   map[ProtoOption]bool
	PatternID              int32
	SlotPattern            []*SlotPattern
	SymbolAttributes       []*SymbolAttributes
	MainGameTable          []*MainGameTable
	SlotSymbolDistribution []*SlotSymbolDistribution
	JackpotTable           []*JackpotTable
	SlotPayLine            []*SlotPayLine
	SlotPayTable           []*SlotPayTable
	BonusTable             []*BonusTable
	PickGameTable          []*PickGameTable
	RespinTable            []*RespinTable
	FreeSpinTable          []*FreeSpinTable
	SelectionTable         []*SelectionTable
	ExtraRule              []*ExtraRule
	ReelBgDistribution     []*ReelBgDistribution
}

// NewProtoData create ProtoData by given path and the proto option what you want to load
func NewProtoData(path string, opts ...ProtoOption) (*ProtoData, error) {
	// Check path folder exist
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a folder", path)
	}

	data := new(ProtoData)
	data.basicPath = path
	data.opts = make(map[ProtoOption]bool)
	for _, opt := range opts {
		data.opts[opt] = true
	}
	if data.opts[OptALL] {
		// for i := 0; i < int(OptEnd); i++ {
		// 	data.opts[ProtoOption(i)] = true
		// }
		// 只獨有檔案的
		f, _ := os.Open(path)
		files, err := f.Readdir(0)
		if err != nil {
			return nil, err
		}
		for _, _fi := range files {
			switch _fi.Name() {
			case "SlotPattern.bytes":
				data.opts[OptSlotPattern] = true
			case "SlotPayLine.bytes":
				data.opts[OptSlotPayLine] = true
			case "SlotPayTable.bytes":
				data.opts[OptSlotPayTable] = true
			case "MainGameTable.bytes":
				data.opts[OptMainGameTable] = true
			case "SymbolAttributes.bytes":
				data.opts[OptSymbolAttributes] = true
			case "JackpotTable.bytes":
				data.opts[OptJackpotTable] = true
			case "FreeSpinTable.bytes":
				data.opts[OptFreeSpinTable] = true
			case "SlotSymbolDistribution.bytes":
				data.opts[OptSlotSymbolDistribution] = true
			case "BonusTable.bytes":
				data.opts[OptBonusTable] = true
			case "PickGameTable.bytes":
				data.opts[OptPickGameTable] = true
			case "SelectionTable.bytes":
				data.opts[OptSelectionTable] = true
			case "RespinTable.bytes":
				data.opts[OptRespinTable] = true
			case "ReelBgDistribution.bytes":
				data.opts[OptReelBgDistribution] = true
			case "ExtraRule.bytes":
				data.opts[OptExtraRule] = true
			}
		}
	}

	return data, nil
}

// Reload from config file
func (pd *ProtoData) Reload() error {

	// SlotPattern
	if val, ok := pd.opts[OptSlotPattern]; val && ok {
		i := new(SlotPattern_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/SlotPattern.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.SlotPattern = i.GetItems()
	}

	// SymbolAttributes
	if val, ok := pd.opts[OptSymbolAttributes]; val && ok {
		i := new(SymbolAttributes_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/SymbolAttributes.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.SymbolAttributes = i.GetItems()
	}

	// MainGameTable
	if val, ok := pd.opts[OptMainGameTable]; val && ok {
		i := new(MainGameTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/MainGameTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.MainGameTable = i.GetItems()
	}

	// SlotSymbolDistribution
	if val, ok := pd.opts[OptSlotSymbolDistribution]; val && ok {
		i := new(SlotSymbolDistribution_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/SlotSymbolDistribution.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.SlotSymbolDistribution = i.GetItems()
	}

	// JackpotTable
	if val, ok := pd.opts[OptJackpotTable]; val && ok {
		i := new(JackpotTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/JackpotTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.JackpotTable = i.GetItems()
	}

	// SlotPayLine
	if val, ok := pd.opts[OptSlotPayLine]; val && ok {
		i := new(SlotPayLine_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/SlotPayLine.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.SlotPayLine = i.GetItems()
	}

	// SlotPayTable
	if val, ok := pd.opts[OptSlotPayTable]; val && ok {
		i := new(SlotPayTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/SlotPayTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.SlotPayTable = i.GetItems()
	}

	// BonusTable
	if val, ok := pd.opts[OptBonusTable]; val && ok {
		i := new(BonusTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/BonusTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.BonusTable = i.GetItems()
	}

	// PickGameTable
	if val, ok := pd.opts[OptPickGameTable]; val && ok {
		i := new(PickGameTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/PickGameTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.PickGameTable = i.GetItems()
	}

	// RespinTable
	if val, ok := pd.opts[OptRespinTable]; val && ok {
		i := new(RespinTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/RespinTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.RespinTable = i.GetItems()
	}

	// FreeSpinTable
	if val, ok := pd.opts[OptFreeSpinTable]; val && ok {
		i := new(FreeSpinTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/FreeSpinTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.FreeSpinTable = i.GetItems()
	}

	// SelectionTable
	if val, ok := pd.opts[OptSelectionTable]; val && ok {
		i := new(SelectionTable_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/SelectionTable.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.SelectionTable = i.GetItems()
	}

	// ExtraRule
	if val, ok := pd.opts[OptExtraRule]; val && ok {
		i := new(ExtraRule_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/ExtraRule.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.ExtraRule = i.GetItems()
	}

	// ReelBgDistribution
	if val, ok := pd.opts[OptReelBgDistribution]; val && ok {
		i := new(ReelBgDistribution_ARRAY)
		if err := pd.Load(fmt.Sprintf("%s/ReelBgDistribution.bytes", pd.basicPath), i); err != nil {
			return err
		}
		pd.ReelBgDistribution = i.GetItems()
	}

	return nil
}

// Print proto config content
func (pd *ProtoData) Print() {

	// SlotPattern
	{
		fmt.Println("=== SlotPattern ===", len(pd.SlotPattern))
		for _, row := range pd.SlotPattern {
			fmt.Println(row.GetMyComment(), row)
		}
	}

	// SymbolAttributes
	{
		fmt.Println("=== SymbolAttributes ===", len(pd.SymbolAttributes))
		for _, row := range pd.SymbolAttributes {
			fmt.Println(row)
		}
	}

	// MainGameTable
	{
		fmt.Println("=== MainGameTable ===", len(pd.MainGameTable))
		for _, row := range pd.MainGameTable {
			fmt.Println(row)
		}
	}

	// SlotSymbolDistribution
	{
		fmt.Println("=== SlotSymbolDistribution ===", len(pd.SlotSymbolDistribution))
		for _, row := range pd.SlotSymbolDistribution {
			fmt.Println(row)
		}
	}

	// JackpotTable
	{
		fmt.Println("=== JackpotTable ===", len(pd.JackpotTable))
		for _, row := range pd.JackpotTable {
			fmt.Println(row)
		}
	}

	// SlotPayLine
	{
		fmt.Println("=== SlotPayLine ===", len(pd.SlotPayLine))
		for _, row := range pd.SlotPayLine {
			fmt.Println(row)
		}
	}

	// SlotPayTable
	{
		fmt.Println("=== SlotPayTable ===", len(pd.SlotPayTable))
		for _, row := range pd.SlotPayTable {
			fmt.Println(row)
		}
	}

	// BonusTable
	{
		fmt.Println("=== BonusTable ===", len(pd.BonusTable))
		for _, row := range pd.BonusTable {
			fmt.Println(row)
		}
	}

	// PickGameTable
	{
		fmt.Println("=== PickGameTable ===", len(pd.PickGameTable))
		for _, row := range pd.PickGameTable {
			fmt.Println(row)
		}
	}

	// RespinTable
	{
		fmt.Println("=== RespinTable ===", len(pd.RespinTable))
		for _, row := range pd.RespinTable {
			fmt.Println(row)
		}
	}

	// FreeSpinTable
	{
		fmt.Println("=== FreeSpinTable ===", len(pd.FreeSpinTable))
		for _, row := range pd.FreeSpinTable {
			fmt.Println(row)
		}
	}

	// SelectionTable
	{
		fmt.Println("=== SelectionTable ===", len(pd.SelectionTable))
		for _, row := range pd.SelectionTable {
			fmt.Println(row)
		}
	}

	// ExtraRule
	{
		fmt.Println("=== ExtraRule ===", len(pd.ExtraRule))
		for _, row := range pd.ExtraRule {
			fmt.Println(row)
		}
	}

	// ReelBgDistribution
	{
		fmt.Println("=== ReelBgDistribution ===", len(pd.ReelBgDistribution))
		for _, row := range pd.ReelBgDistribution {
			fmt.Println(row)
		}
	}

}

// Load proto data by path & type
func (pd *ProtoData) Load(_path string, i interface{}) error {

	var b64 bool

	dir := path.Dir(_path)
	base := path.Base(_path)

	if _, err := os.Stat(_path); err != nil {
		b64 = true
	}

	// check default file exist
	if b64 {
		buf, err := ioutil.ReadFile(dir + "/b64-" + base)
		if err != nil {
			return err
		}

		b, err := base64.StdEncoding.DecodeString(string(buf))
		if err != nil {
			return err
		}

		if err := proto.Unmarshal(b, i.(proto.Message)); err != nil {
			return err
		}

		return nil
	}

	b, err := ioutil.ReadFile(_path)
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(b, i.(proto.Message)); err != nil {
		return err
	}

	return nil
}
