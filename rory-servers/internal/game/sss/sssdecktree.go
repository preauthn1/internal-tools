package sss

import (
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// 13水紀錄牌組的 tree 結構
type SSSDeckTree struct {
	leafLimit int // 生長限制大小
	leafCount int // 目前葉子數量

	isStop bool // 是否停止生長

	source    *poker.CardGroup         // 原始13張牌
	record    map[int]([]*SSSDeckNode) // key: layer  value: node slices
	sortNodes []*SSSDeckNode           // 照大小牌組排列

	root *SSSDeckNode
}

// NewDeckTree 建構函式: 最大牌型設為 root, leafCount 控制生長大小
func (t *SSSDeckTree) NewDeckTree(table *SSSTable, cards *poker.CardGroup, leafCount int) *SSSDeckTree {
	t.isStop = false
	t.leafCount = 0
	t.leafLimit = leafCount

	t.source = poker.NewCardGroup(cards.Cards)

	var rootDeck = new(SSSDeck)
	rootDeck.SetMaxDeck(table, t.source)
	t.root = &SSSDeckNode{data: rootDeck, left: nil, middle: nil, right: nil, layer: 0, nodeType: SSSDeckNodeType_Root, position: 0}

	t.record = make(map[int]([]*SSSDeckNode), 0)
	var layerNodes = make([]*SSSDeckNode, 0)
	layerNodes = append(layerNodes, t.root)
	t.record[0] = layerNodes

	t.sortNodes = make([]*SSSDeckNode, 0)
	t.sortNodes = append(t.sortNodes, t.root)

	return t
}

func (t *SSSDeckTree) Root() *SSSDeckNode {
	return t.root
}

func (t *SSSDeckTree) Print() {

	sort.Sort(t.source)
	fmt.Println(t.source)

	// for layer, nodes := range t.record {
	// 	fmt.Println("[", layer, "]")
	// 	for _, node := range nodes {
	// 		fmt.Println(node)
	// 	}
	// }

	for _, node := range t.sortNodes {
		fmt.Println(node)
	}
}

// GetSortNodes
func (t *SSSDeckTree) GetSortNodes() []*SSSDeckNode {
	return t.sortNodes
}

// GrowUp 讓樹長大(產生root底下的可能推薦牌型, 直到 leaf 超過 leafLimit 停止)
func (t *SSSDeckTree) GrowUp(s *SSSTable) {

	// 看看 tree 有沒有下一層的支葉
	var layer = 0 //從 root 開始
	var hasNextLayer bool = true
	for !t.isStopGrowUp() && hasNextLayer {
		hasNextLayer = false
		var layerNodes = t.record[layer]
		var nextLayerNodes = make([]*SSSDeckNode, 0)
		for _, node := range layerNodes {
			// 不同牌組搜尋策略(依照牌型，依照牌組)
			var extraDeck, rightDeck, middleDeck, leftDeck *SSSDeck = nil, nil, nil, nil
			var extraNode, rightNode, middleNode, leftNode *SSSDeckNode = nil, nil, nil, nil
			var extraResult, rightResult, middleResult, leftResult = false, false, false, false

			// extra
			extraDeck = node.data.NextDeckByExtra(s, t.source)

			if extraDeck != nil && !t.existDeck(extraDeck) { //如果找的到牌組，建立支葉
				extraResult, extraNode = t.addExtraLeaf(node, extraDeck)
				if extraResult {
					nextLayerNodes = append(nextLayerNodes, extraNode)
					t.record[extraNode.layer] = nextLayerNodes

					t.sortNodes = append(t.sortNodes, extraNode)
					if t.isStopGrowUp() {
						return
					}
					hasNextLayer = true
				}
			}

			// 右邊
			rightDeck = node.data.NextDeckByFoot(s, t.source)

			if rightDeck != nil && !t.existDeck(rightDeck) { //如果找的到牌組，建立右邊支葉
				rightResult, rightNode = t.addRightLeaf(node, rightDeck)
				if rightResult {
					// fmt.Println("[", rightNode.layer, "]", "右", k)
					nextLayerNodes = append(nextLayerNodes, rightNode)
					t.record[rightNode.layer] = nextLayerNodes

					t.sortNodes = append(t.sortNodes, rightNode)
					if t.isStopGrowUp() {
						return
					}
					hasNextLayer = true
				}
			}

			// 中間
			middleDeck = node.data.NextDeckByHead(s, t.source)
			if middleDeck != nil && !t.existDeck(middleDeck) { //如果找的到牌組，建立中間支葉
				middleResult, middleNode = t.addMiddleLeaf(node, middleDeck)
				if middleResult {
					nextLayerNodes = append(nextLayerNodes, middleNode)
					t.record[middleNode.layer] = nextLayerNodes
					t.sortNodes = append(t.sortNodes, middleNode)
					if t.isStopGrowUp() {
						return
					}
					hasNextLayer = true
				}
			}

			// 左邊
			leftDeck = node.data.NextDeckByBody(s, t.source)

			if leftDeck != nil && !t.existDeck(leftDeck) { //如果找的到牌組，建立左邊支葉
				leftResult, leftNode = t.addLeftLeaf(node, leftDeck)
				if leftResult {
					nextLayerNodes = append(nextLayerNodes, leftNode)
					t.record[leftNode.layer] = nextLayerNodes
					t.sortNodes = append(t.sortNodes, leftNode)
					if t.isStopGrowUp() {
						return
					}
					hasNextLayer = true
				}
			}
		}

		if hasNextLayer {
			layer++
		}
	}
}

// existDeck 是否存在牌組
func (t *SSSDeckTree) existDeck(data *SSSDeck) bool {
	for _, nodes := range t.record {
		for _, node := range nodes {
			if node.data.IsEqualDeck(data) {
				// fmt.Println("@@", node.data)
				return true
			}
		}
	}
	return false
}

// addExtraLeaf 新增節點的extra葉子
func (t *SSSDeckTree) addExtraLeaf(node *SSSDeckNode, data *SSSDeck) (result bool, currentNode *SSSDeckNode) {
	if node == nil || data == nil {
		return false, nil
	}

	if t.leafCount >= t.leafLimit {
		return false, nil
	}

	t.leafCount++
	node.extra = node.extra.bindData(data, node.layer+1, SSSDeckNodeType_Extra, t.leafCount)
	return true, node.extra
}

// addRightLeaf 新增節點的右邊葉子
func (t *SSSDeckTree) addRightLeaf(node *SSSDeckNode, data *SSSDeck) (result bool, currentNode *SSSDeckNode) {
	if node == nil || data == nil {
		return false, nil
	}

	if t.leafCount >= t.leafLimit {
		return false, nil
	}

	t.leafCount++
	node.right = node.right.bindData(data, node.layer+1, SSSDeckNodeType_Right, t.leafCount)
	return true, node.right
}

// addMiddleLeaf 新增節點的中間葉子
func (t *SSSDeckTree) addMiddleLeaf(node *SSSDeckNode, data *SSSDeck) (result bool, currentNode *SSSDeckNode) {
	if node == nil || data == nil {
		return false, nil
	}

	if t.leafCount >= t.leafLimit {
		return false, nil
	}

	t.leafCount++
	node.middle = node.middle.bindData(data, node.layer+1, SSSDeckNodeType_Middle, t.leafCount)
	return true, node.middle
}

// addLeftLeaf 新增節點的左邊葉子
func (t *SSSDeckTree) addLeftLeaf(node *SSSDeckNode, data *SSSDeck) (result bool, currentNode *SSSDeckNode) {
	if node == nil || data == nil {
		return false, nil
	}

	if t.leafCount >= t.leafLimit {
		return false, nil
	}

	t.leafCount++
	node.left = node.left.bindData(data, node.layer+1, SSSDeckNodeType_Left, t.leafCount)
	return true, node.left
}

// checkGrowUp 是否停止生長
func (t *SSSDeckTree) isStopGrowUp() bool {
	if t.leafCount >= t.leafLimit {
		t.isStop = true
		return true
	}
	return t.isStop
}
