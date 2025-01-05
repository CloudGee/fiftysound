package vocabulary

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ==================================================
// 0. 数据结构 & 常量
// ==================================================

// WordItem 表示单词结构
type WordItem struct {
	Kana   string   `json:"假名"`
	Kanji  string   `json:"日本汉字"`
	Chines []string `json:"中文释义"`
}

func sameWord(a, b WordItem) bool {
	return a.Kana == b.Kana && a.Kanji == b.Kanji
}

const githubZipURL = "https://github.com/CloudGee/JapaneseVocabulary/archive/refs/heads/main.zip"

// 保存每个节点的选中状态
var selectedPaths = make(map[string]bool)

// DirEntry 表示目录/文件节点
type DirEntry struct {
	Name     string
	FullPath string
	IsDir    bool
	Content  []byte
	Parent   *DirEntry
	Children []*DirEntry
}

// 全局
var (
	rootDir      *DirEntry
	vocabDir     *DirEntry
	nodeIndex    map[string]*DirEntry
	fileContents map[string][]byte

	// 选中的单词
	selectedWords []WordItem
)

// ==================================================
// 1. 模块主页面：ShowVocabularyMainPage
//    先加载词库，然后进入“模块主页面”
// ==================================================

// ShowVocabularyMainPage 主界面改用下拉框选择模式
// ShowVocabularyMainPage 主界面改用下拉框选择模式
func ShowVocabularyMainPage(myApp fyne.App, parent fyne.Window) {
	progress := dialog.NewProgressInfinite("提示", "正在从GitHub请求最新词库...", parent)
	progress.Show()

	var err error
	maxRetries := 4
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		err = loadZipAndInit()
		if err == nil {
			fmt.Println("Operation succeeded")
			break
		}

		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	progress.Hide()

	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	mainWin := myApp.NewWindow("新标日语单词练习 - 模块主页面")

	// 下拉框选择模式
	modeSelect := widget.NewSelect([]string{
		"模式1: 中文 => 假名&汉字",
		"模式2: 假名(汉字) => 中文",
		"模式3: 背单词",
	}, nil)
	modeSelect.PlaceHolder = "请点击下拉框，选择你想要的模式"

	// 开始按钮
	startBtn := widget.NewButton("开始", func() {
		if modeSelect.Selected == "" {
			dialog.ShowInformation("提示", "请先选择模式", mainWin)
			return
		}
		if len(selectedWords) == 0 {
			dialog.ShowInformation("提示", "请先选择单词", mainWin)
			return
		}

		switch modeSelect.Selected {
		case "模式1: 中文 => 假名&汉字":
			showModeOneWords(myApp, mainWin, selectedWords)
		case "模式2: 假名(汉字) => 中文":
			showModeTwoWords(myApp, mainWin, selectedWords)
		case "模式3: 背单词":
			showModeThreeWords(myApp, mainWin, selectedWords)
		}
	})

	// 选择文件或目录按钮
	selBtn := widget.NewButton("请先选择需要练习的单元", func() {
		showSelectTree(myApp, mainWin)
	})

	mainWin.SetContent(container.NewVBox(
		widget.NewLabel("新标日语单词练习 (模块主页面)"),
		widget.NewLabel("请选择操作："),
		selBtn,
		modeSelect,
		startBtn, // 替换为开始按钮
	))
	mainWin.Resize(fyne.NewSize(400, 300))
	mainWin.Show()
}

func showSelectTree(myApp fyne.App, parent fyne.Window) {
	selWin := myApp.NewWindow("选择需要练习的单元")
	if vocabDir == nil {
		dialog.ShowInformation("提示", "未找到vocabularyLib目录", parent)
		return
	}

	myTree := widget.NewTree(
		func(uid string) []string {
			if uid == "" {
				return []string{vocabDir.FullPath}
			}
			nd := nodeIndex[uid]
			if nd == nil {
				return nil
			}
			res := make([]string, len(nd.Children))
			for i, c := range nd.Children {
				res[i] = c.FullPath
			}
			return res
		},
		func(uid string) bool {
			if uid == "" {
				return true
			}
			nd := nodeIndex[uid]
			return nd != nil && nd.IsDir
		},
		func(branch bool) fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),      // 用于显示文件或目录名
				widget.NewCheck("", nil), // 用于文件的选择框
			)
		},
		func(uid string, branch bool, obj fyne.CanvasObject) {
			nd := nodeIndex[uid]
			hbox := obj.(*fyne.Container)
			label := hbox.Objects[0].(*widget.Label)
			check := hbox.Objects[1].(*widget.Check)

			if nd != nil {
				// 去除 `.json` 后缀
				displayName := nd.Name
				if strings.HasSuffix(strings.ToLower(displayName), ".json") {
					displayName = strings.TrimSuffix(displayName, ".json")
				}
				label.Text = displayName
				label.Refresh()

				if nd.IsDir {
					// 目录节点：禁用选择框，隐藏选择框
					check.SetChecked(false)
					check.Disable()
					check.Hide()
				} else {
					// 文件节点：启用选择框
					check.SetChecked(selectedPaths[uid]) // 恢复之前的选择状态
					check.Enable()
					check.Show()
					check.OnChanged = func(checked bool) {
						selectedPaths[uid] = checked // 更新选择状态
					}
				}
			}
		},
	)

	// 默认展开所有节点
	myTree.OpenAllBranches()

	// 确认按钮逻辑
	confirmBtn := widget.NewButton("确认", func() {
		var combined []WordItem
		for fullPath, checked := range selectedPaths {
			if !checked {
				continue
			}
			ws, _ := loadJSON(fullPath)
			combined = append(combined, ws...)
		}

		if len(combined) == 0 {
			dialog.ShowInformation("提示", "没有选到任何 JSON 文件", parent)
			return
		}

		selectedWords = combined
		parent.Show()  // 确保主界面保持打开
		selWin.Close() // 关闭选择文件或目录的窗口
	})

	// 取消按钮逻辑
	cancelBtn := widget.NewButton("取消", func() {
		parent.Show()
	})

	selWin.SetContent(container.NewBorder(
		nil,
		container.NewHBox(confirmBtn, cancelBtn),
		nil, nil,
		myTree,
	))
	selWin.Resize(fyne.NewSize(600, 500))
	selWin.Show()
}

// 更新所有子节点的选中状态（异步处理）
func updateChildrenChecked(n *DirEntry, checked bool) {
	go func() {
		var stack []*DirEntry
		stack = append(stack, n)

		for len(stack) > 0 {
			// 获取栈顶节点
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// 更新当前节点状态
			selectedPaths[current.FullPath] = checked

			// 将子节点加入栈中
			stack = append(stack, current.Children...)
		}
	}()
}

// 更新父节点的选中状态
func updateParentChecked(n *DirEntry) {
	if n == nil || n.Parent == nil {
		return
	}

	parent := n.Parent
	allChecked := true
	for _, child := range parent.Children {
		if !selectedPaths[child.FullPath] {
			allChecked = false
			break
		}
	}

	selectedPaths[parent.FullPath] = allChecked
	updateParentChecked(parent)
}

// loadZipAndInit 下载 ZIP 到内存并解析
func loadZipAndInit() error {
	data, err := downloadZip(githubZipURL)
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("ZIP解压失败: %w", err)
	}

	// 构造 rootDir
	rootDir = &DirEntry{
		Name:     "ROOT",
		FullPath: "ROOT",
		IsDir:    true,
	}
	nodeIndex = make(map[string]*DirEntry)
	nodeIndex["ROOT"] = rootDir
	fileContents = make(map[string][]byte)

	// 目录优先
	files := make([]*zip.File, len(zr.File))
	copy(files, zr.File)
	sort.Slice(files, func(i, j int) bool {
		if files[i].FileInfo().IsDir() != files[j].FileInfo().IsDir() {
			return files[i].FileInfo().IsDir()
		}
		return files[i].Name < files[j].Name
	})
	for _, f := range files {
		createDirEntry(rootDir, f)
	}

	// 找到 vocabularyLib
	vocabDir = findDirByName(rootDir, "vocabularyLib")
	if vocabDir != nil {
		vocabDir.Name = "标准日本语第二版"
	}

	return nil
}

func downloadZip(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP状态码=%d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func createDirEntry(root *DirEntry, zf *zip.File) {
	parts := strings.Split(zf.Name, "/")
	curr := root
	for i, p := range parts {
		if p == "" {
			continue
		}

		found := false
		var child *DirEntry
		for _, c := range curr.Children {
			if c.Name == p {
				child = c
				found = true
				break
			}
		}

		if !found {
			isDir := i < len(parts)-1 || zf.FileInfo().IsDir()
			newFullPath := curr.FullPath + "/" + p
			child = &DirEntry{
				Name:     p,
				FullPath: newFullPath,
				IsDir:    isDir,
				Parent:   curr,
			}
			curr.Children = append(curr.Children, child)
			nodeIndex[newFullPath] = child

			// 目录优先再排序
			sort.Slice(curr.Children, func(i, j int) bool {
				ci, cj := curr.Children[i], curr.Children[j]
				if ci.IsDir != cj.IsDir {
					return ci.IsDir
				}
				return ci.Name < cj.Name
			})
		}

		curr = child

		// 如果是文件 => 读内容
		if i == len(parts)-1 && !zf.FileInfo().IsDir() {
			rc, err := zf.Open()
			if err == nil {
				bs, _ := io.ReadAll(rc)
				rc.Close()
				curr.Content = bs
				fileContents[curr.FullPath] = bs
			}
		}
	}
}

func findDirByName(dir *DirEntry, name string) *DirEntry {
	if dir.Name == name && dir.IsDir {
		return dir
	}
	for _, c := range dir.Children {
		if c.IsDir {
			if found := findDirByName(c, name); found != nil {
				return found
			}
		}
	}
	return nil
}

// 递归向下全选/全取消
func checkAllChildren(n *DirEntry, checked bool, sel map[string]bool) {
	sel[n.FullPath] = checked
	for _, c := range n.Children {
		sel[c.FullPath] = checked
		if c.IsDir {
			checkAllChildren(c, checked, sel)
		}
	}
}

// 若某子节点取消 => 父目录取消
func uncheckParentIfNeeded(n *DirEntry, sel map[string]bool) {
	if n == nil || n.Parent == nil {
		return
	}
	parent := n.Parent
	// 看看 parent 的所有子是否都选中
	allChecked := true
	for _, c := range parent.Children {
		if !sel[c.FullPath] {
			allChecked = false
			break
		}
	}
	if !allChecked {
		sel[parent.FullPath] = false
		uncheckParentIfNeeded(parent, sel)
	}
}

// 若同级都选 => 父也选
func checkParentIfAllSiblingsChecked(n *DirEntry, sel map[string]bool) {
	if n == nil || n.Parent == nil {
		return
	}
	parent := n.Parent
	// 检查同级
	for _, sibling := range parent.Children {
		if !sel[sibling.FullPath] {
			return
		}
	}
	// 同级全选 => 父选
	sel[parent.FullPath] = true
	checkParentIfAllSiblingsChecked(parent, sel)
}

// 递归收集 json
func collectJSON(n *DirEntry, files *[]string) {
	if !n.IsDir {
		if strings.HasSuffix(strings.ToLower(n.Name), ".json") {
			*files = append(*files, n.FullPath)
		}
		return
	}
	for _, c := range n.Children {
		collectJSON(c, files)
	}
}

func loadJSON(pathStr string) ([]WordItem, error) {
	data, ok := fileContents[pathStr]
	if !ok {
		return nil, errors.New("no content")
	}
	var arr []WordItem
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

// ==================================================
// 3. 三种模式：showModeOneWords, showModeTwoWords, ...
// ==================================================

type WordPool struct {
	items []WordItem
	index int
	last  *WordItem
}

func newWordPool(words []WordItem) *WordPool {
	p := &WordPool{
		items: make([]WordItem, len(words)),
	}
	copy(p.items, words)
	rand.Seed(time.Now().UnixNano())
	p.shuffle()
	return p
}

func (p *WordPool) shuffle() {
	rand.Shuffle(len(p.items), func(i, j int) {
		p.items[i], p.items[j] = p.items[j], p.items[i]
	})
	p.index = 0
}

func (p *WordPool) nextWord() WordItem {
	if p.index >= len(p.items) {
		p.shuffle()
	}
	item := p.items[p.index]
	p.index++

	// 避免连续相同
	if p.last != nil && sameWord(item, *p.last) && len(p.items) > 1 {
		p.shuffle()
		item = p.items[0]
		p.index = 1
	}
	p.last = &item
	return item
}

// 模式1: "中文" => 假名&汉字
func showModeOneWords(myApp fyne.App, parent fyne.Window, words []WordItem) {
	win := myApp.NewWindow("模式1: 中文 => 假名&汉字")

	pool := newWordPool(words)
	question := widget.NewLabel("")
	kanaEntry := widget.NewEntry()
	kanjiEntry := widget.NewEntry()
	feedback := widget.NewLabel("")

	var current WordItem

	var refresh = func() {
		kanaEntry.SetText("")
		kanjiEntry.SetText("")
		feedback.SetText("")
		current = pool.nextWord()
		question.SetText("中文释义: " + strings.Join(current.Chines, "/"))
	}

	judgeBtn := widget.NewButton("判题", func() {
		k := strings.TrimSpace(kanaEntry.Text)
		j := strings.TrimSpace(kanjiEntry.Text)
		if k == current.Kana && j == current.Kanji {
			feedback.SetText("正确！")
		} else {
			feedback.SetText(fmt.Sprintf("错误，正确答案: %s / %s", current.Kana, current.Kanji))
		}
	})

	nextBtn := widget.NewButton("下一题", func() {
		refresh()
	})

	closeBtn := widget.NewButton("关闭", func() {
		win.Close()
	})

	win.SetContent(container.NewVBox(
		question,
		widget.NewLabel("假名："), kanaEntry,
		widget.NewLabel("汉字："), kanjiEntry,
		container.NewHBox(judgeBtn, nextBtn),
		feedback,
		closeBtn,
	))
	win.Resize(fyne.NewSize(400, 300))
	refresh()
	win.Show()
}

// 模式2: "假名(汉字)" => 中文
func showModeTwoWords(myApp fyne.App, parent fyne.Window, words []WordItem) {
	win := myApp.NewWindow("模式2: 假名(汉字) => 中文")

	pool := newWordPool(words)
	question := widget.NewLabel("")
	answerEntry := widget.NewEntry()
	feedback := widget.NewLabel("")

	var current WordItem

	var refresh = func() {
		answerEntry.SetText("")
		feedback.SetText("")
		current = pool.nextWord()
		question.SetText(fmt.Sprintf("请填写中文: %s (%s)", current.Kana, current.Kanji))
	}

	judgeBtn := widget.NewButton("判题", func() {
		ans := strings.TrimSpace(answerEntry.Text)
		correct := false
		for _, c := range current.Chines {
			if c == ans {
				correct = true
				break
			}
		}
		if correct {
			feedback.SetText("正确！")
		} else {
			feedback.SetText("错误！正确答案: " + strings.Join(current.Chines, "/"))
		}
	})

	nextBtn := widget.NewButton("下一题", func() {
		refresh()
	})

	closeBtn := widget.NewButton("关闭", func() {
		win.Close()
	})

	win.SetContent(container.NewVBox(
		question,
		answerEntry,
		container.NewHBox(judgeBtn, nextBtn),
		feedback,
		closeBtn,
	))
	win.Resize(fyne.NewSize(400, 300))
	refresh()
	win.Show()
}

// 模式3: 背单词 (显示中文+假名+汉字)
func showModeThreeWords(myApp fyne.App, parent fyne.Window, words []WordItem) {
	win := myApp.NewWindow("模式3: 背单词")

	pool := newWordPool(words)
	wordLabel := widget.NewLabel("")

	var showOne = func() {
		w := pool.nextWord()
		wordLabel.SetText(fmt.Sprintf("[中文] %s\n[假名] %s\n[汉字] %s",
			strings.Join(w.Chines, "/"), w.Kana, w.Kanji))
	}

	nextBtn := widget.NewButton("下一词", func() {
		showOne()
	})
	closeBtn := widget.NewButton("关闭", func() {
		win.Close()
	})

	win.SetContent(container.NewVBox(
		wordLabel,
		nextBtn,
		closeBtn,
	))
	win.Resize(fyne.NewSize(400, 300))
	showOne()
	win.Show()
}
