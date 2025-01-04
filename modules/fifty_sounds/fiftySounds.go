package fifty_sounds

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ======================= 五十音图行定义 =======================
type gojuonLine struct {
	romaji   string
	hiragana []string
	katakana []string
}

var gojuon = []gojuonLine{
	// 五十音 (Unvoiced Sounds)
	{
		romaji:   "a",
		hiragana: []string{"あ", "い", "う", "え", "お"},
		katakana: []string{"ア", "イ", "ウ", "エ", "オ"},
	},
	{
		romaji:   "ka",
		hiragana: []string{"か", "き", "く", "け", "こ"},
		katakana: []string{"カ", "キ", "ク", "ケ", "コ"},
	},
	{
		romaji:   "sa",
		hiragana: []string{"さ", "し", "す", "せ", "そ"},
		katakana: []string{"サ", "シ", "ス", "セ", "ソ"},
	},
	{
		romaji:   "ta",
		hiragana: []string{"た", "ち", "つ", "て", "と"},
		katakana: []string{"タ", "チ", "ツ", "テ", "ト"},
	},
	{
		romaji:   "na",
		hiragana: []string{"な", "に", "ぬ", "ね", "の"},
		katakana: []string{"ナ", "ニ", "ヌ", "ネ", "ノ"},
	},
	{
		romaji:   "ha",
		hiragana: []string{"は", "ひ", "ふ", "へ", "ほ"},
		katakana: []string{"ハ", "ヒ", "フ", "ヘ", "ホ"},
	},
	{
		romaji:   "ma",
		hiragana: []string{"ま", "み", "む", "め", "も"},
		katakana: []string{"マ", "ミ", "ム", "メ", "モ"},
	},
	{
		romaji:   "ya",
		hiragana: []string{"や", "ゆ", "よ"},
		katakana: []string{"ヤ", "ユ", "ヨ"},
	},
	{
		romaji:   "ra",
		hiragana: []string{"ら", "り", "る", "れ", "ろ"},
		katakana: []string{"ラ", "リ", "ル", "レ", "ロ"},
	},
	{
		romaji:   "wa",
		hiragana: []string{"わ", "を", "ん"},
		katakana: []string{"ワ", "ヲ", "ン"},
	},

	// 浊音・半濁音 (Voiced and Semi-Voiced Sounds)
	{
		romaji:   "ga",
		hiragana: []string{"が", "ぎ", "ぐ", "げ", "ご"},
		katakana: []string{"ガ", "ギ", "グ", "ゲ", "ゴ"},
	},
	{
		romaji:   "za",
		hiragana: []string{"ざ", "じ", "ず", "ぜ", "ぞ"},
		katakana: []string{"ザ", "ジ", "ズ", "ゼ", "ゾ"},
	},
	{
		romaji:   "da",
		hiragana: []string{"だ", "ぢ", "づ", "で", "ど"},
		katakana: []string{"ダ", "ヂ", "ヅ", "デ", "ド"},
	},
	{
		romaji:   "ba",
		hiragana: []string{"ば", "び", "ぶ", "べ", "ぼ"},
		katakana: []string{"バ", "ビ", "ブ", "ベ", "ボ"},
	},
	{
		romaji:   "pa",
		hiragana: []string{"ぱ", "ぴ", "ぷ", "ぺ", "ぽ"},
		katakana: []string{"パ", "ピ", "プ", "ペ", "ポ"},
	},

	// 拗音 (Contracted Sounds)
	{
		romaji:   "kya",
		hiragana: []string{"きゃ", "きゅ", "きょ"},
		katakana: []string{"キャ", "キュ", "キョ"},
	},
	{
		romaji:   "gya",
		hiragana: []string{"ぎゃ", "ぎゅ", "ぎょ"},
		katakana: []string{"ギャ", "ギュ", "ギョ"},
	},
	{
		romaji:   "sha",
		hiragana: []string{"しゃ", "しゅ", "しょ"},
		katakana: []string{"シャ", "シュ", "ショ"},
	},
	{
		romaji:   "ja",
		hiragana: []string{"じゃ", "じゅ", "じょ"},
		katakana: []string{"ジャ", "ジュ", "ジョ"},
	},
	{
		romaji:   "cha",
		hiragana: []string{"ちゃ", "ちゅ", "ちょ"},
		katakana: []string{"チャ", "チュ", "チョ"},
	},
	{
		romaji:   "nya",
		hiragana: []string{"にゃ", "にゅ", "にょ"},
		katakana: []string{"ニャ", "ニュ", "ニョ"},
	},
	{
		romaji:   "hya",
		hiragana: []string{"ひゃ", "ひゅ", "ひょ"},
		katakana: []string{"ヒャ", "ヒュ", "ヒョ"},
	},
	{
		romaji:   "bya",
		hiragana: []string{"びゃ", "びゅ", "びょ"},
		katakana: []string{"ビャ", "ビュ", "ビョ"},
	},
	{
		romaji:   "pya",
		hiragana: []string{"ぴゃ", "ぴゅ", "ぴょ"},
		katakana: []string{"ピャ", "ピュ", "ピョ"},
	},
	{
		romaji:   "mya",
		hiragana: []string{"みゃ", "みゅ", "みょ"},
		katakana: []string{"ミャ", "ミュ", "ミョ"},
	},
	{
		romaji:   "rya",
		hiragana: []string{"りゃ", "りゅ", "りょ"},
		katakana: []string{"リャ", "リュ", "リョ"},
	},
}

// ======================= kanaToRomaji 定义 =======================
var kanaToRomaji = map[string][]string{
	// 平假名
	"あ": {"a"}, "い": {"i"}, "う": {"u"}, "え": {"e"}, "お": {"o"},
	"か": {"ka"}, "き": {"ki"}, "く": {"ku"}, "け": {"ke"}, "こ": {"ko"},
	"が": {"ga"}, "ぎ": {"gi"}, "ぐ": {"gu"}, "げ": {"ge"}, "ご": {"go"},
	"さ": {"sa"}, "し": {"shi", "si"}, "す": {"su"}, "せ": {"se"}, "そ": {"so"},
	"ざ": {"za"}, "じ": {"ji", "zi"}, "ず": {"zu"}, "ぜ": {"ze"}, "ぞ": {"zo"},
	"た": {"ta"}, "ち": {"chi", "ti"}, "つ": {"tsu", "tu"}, "て": {"te"}, "と": {"to"},
	"だ": {"da"}, "ぢ": {"ji(di)", "ji", "di"}, "づ": {"zu(du)", "zu", "du"}, "で": {"de"}, "ど": {"do"},
	"な": {"na"}, "に": {"ni"}, "ぬ": {"nu"}, "ね": {"ne"}, "の": {"no"},
	"は": {"ha"}, "ひ": {"hi"}, "ふ": {"fu", "hu"}, "へ": {"he"}, "ほ": {"ho"},
	"ば": {"ba"}, "び": {"bi"}, "ぶ": {"bu"}, "べ": {"be"}, "ぼ": {"bo"},
	"ぱ": {"pa"}, "ぴ": {"pi"}, "ぷ": {"pu"}, "ぺ": {"pe"}, "ぽ": {"po"},
	"ま": {"ma"}, "み": {"mi"}, "む": {"mu"}, "め": {"me"}, "も": {"mo"},
	"や": {"ya"}, "ゆ": {"yu"}, "よ": {"yo"}, "ゃ": {"ya"}, "ゅ": {"yu"}, "ょ": {"yo"},
	"ら": {"ra"}, "り": {"ri"}, "る": {"ru"}, "れ": {"re"}, "ろ": {"ro"},
	"わ": {"wa"}, "を": {"o(wo)", "o", "wo"}, "ん": {"n"},
	"きゃ": {"kya"}, "きゅ": {"kyu"}, "きょ": {"kyo"},
	"ぎゃ": {"gya"}, "ぎゅ": {"gyu"}, "ぎょ": {"gyo"},
	"しゃ": {"sha"}, "しゅ": {"shu"}, "しょ": {"sho"},
	"じゃ": {"ja"}, "じゅ": {"ju"}, "じょ": {"jo"},
	"ちゃ": {"cha"}, "ちゅ": {"chu"}, "ちょ": {"cho"},
	"にゃ": {"nya"}, "にゅ": {"nyu"}, "にょ": {"nyo"},
	"ひゃ": {"hya"}, "ひゅ": {"hyu"}, "ひょ": {"hyo"},
	"びゃ": {"bya"}, "びゅ": {"byu"}, "びょ": {"byo"},
	"ぴゃ": {"pya"}, "ぴゅ": {"pyu"}, "ぴょ": {"pyo"},
	"みゃ": {"mya"}, "みゅ": {"myu"}, "みょ": {"myo"},
	"りゃ": {"rya"}, "りゅ": {"ryu"}, "りょ": {"ryo"},

	// 片假名
	"ア": {"a"}, "イ": {"i"}, "ウ": {"u"}, "エ": {"e"}, "オ": {"o"},
	"カ": {"ka"}, "キ": {"ki"}, "ク": {"ku"}, "ケ": {"ke"}, "コ": {"ko"},
	"ガ": {"ga"}, "ギ": {"gi"}, "グ": {"gu"}, "ゲ": {"ge"}, "ゴ": {"go"},
	"サ": {"sa"}, "シ": {"shi", "si"}, "ス": {"su"}, "セ": {"se"}, "ソ": {"so"},
	"ザ": {"za"}, "ジ": {"ji", "zi"}, "ズ": {"zu"}, "ゼ": {"ze"}, "ゾ": {"zo"},
	"タ": {"ta"}, "チ": {"chi", "ti"}, "ツ": {"tsu", "tu"}, "テ": {"te"}, "ト": {"to"},
	"ダ": {"da"}, "ヂ": {"ji(di)", "ji", "di"}, "ヅ": {"zu(du)", "zu", "du"}, "デ": {"de"}, "ド": {"do"},
	"ナ": {"na"}, "ニ": {"ni"}, "ヌ": {"nu"}, "ネ": {"ne"}, "ノ": {"no"},
	"ハ": {"ha"}, "ヒ": {"hi"}, "フ": {"fu", "hu"}, "ヘ": {"he"}, "ホ": {"ho"},
	"バ": {"ba"}, "ビ": {"bi"}, "ブ": {"bu"}, "ベ": {"be"}, "ボ": {"bo"},
	"パ": {"pa"}, "ピ": {"pi"}, "プ": {"pu"}, "ペ": {"pe"}, "ポ": {"po"},
	"マ": {"ma"}, "ミ": {"mi"}, "ム": {"mu"}, "メ": {"me"}, "モ": {"mo"},
	"ヤ": {"ya"}, "ユ": {"yu"}, "ヨ": {"yo"}, "ャ": {"ya"}, "ュ": {"yu"}, "ョ": {"yo"},
	"ラ": {"ra"}, "リ": {"ri"}, "ル": {"ru"}, "レ": {"re"}, "ロ": {"ro"},
	"ワ": {"wa"}, "ヲ": {"o(wo)", "o", "wo"}, "ン": {"n"},
	"キャ": {"kya"}, "キュ": {"kyu"}, "キョ": {"kyo"},
	"ギャ": {"gya"}, "ギュ": {"gyu"}, "ギョ": {"gyo"},
	"シャ": {"sha"}, "シュ": {"shu"}, "ショ": {"sho"},
	"ジャ": {"ja"}, "ジュ": {"ju"}, "ジョ": {"jo"},
	"チャ": {"cha"}, "チュ": {"chu"}, "チョ": {"cho"},
	"ニャ": {"nya"}, "ニュ": {"nyu"}, "ニョ": {"nyo"},
	"ヒャ": {"hya"}, "ヒュ": {"hyu"}, "ヒョ": {"hyo"},
	"ビャ": {"bya"}, "ビュ": {"byu"}, "ビョ": {"byo"},
	"ピャ": {"pya"}, "ピュ": {"pyu"}, "ピョ": {"pyo"},
	"ミャ": {"mya"}, "ミュ": {"myu"}, "ミョ": {"myo"},
	"リャ": {"rya"}, "リュ": {"ryu"}, "リョ": {"ryo"},
}

// ======================= 其它数据 & 函数 =======================
var selectedChars []string

type Stats struct {
	Total   int
	Correct int
}

func (s Stats) Accuracy() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Correct) / float64(s.Total) * 100.0
}

// ======================= 对外暴露的入口函数 =======================
// 你需要在 main.go 中这样调用：ShowFiftySounds(myApp, myWin)
func ShowFiftySounds(myApp fyne.App, w fyne.Window) {
	// 1) 新开子窗口
	newWin := myApp.NewWindow("五十音 - 子窗口")

	rand.Seed(time.Now().UnixNano())

	// 2) 下拉选择模式
	modeSelect := widget.NewSelect([]string{"模式一: 假名 => 罗马音", "模式二: 罗马音 => 假名手写"}, func(string) {})
	modeSelect.PlaceHolder = "请点击下拉框，选择你想要的模式"

	// 3) 平假名、片假名复选框
	hiraganaCheck := widget.NewCheck("平假名", nil)
	hiraganaCheck.SetChecked(true)
	katakanaCheck := widget.NewCheck("片假名", nil)
	katakanaCheck.SetChecked(true)

	// 4) 统计标签
	statsLabel := widget.NewLabel("当前正确率: 0.00%")

	// 5) 选择假名范围按钮
	selectKanaBtn := widget.NewButton("选择假名范围", func() {
		// 打开“选择假名范围”的新窗口或对话框
		showKanaSelectionDialog(myApp, newWin)
	})

	// 6) “开始”按钮
	startBtn := widget.NewButton("开始", func() {
		mode := modeSelect.Selected
		if mode == "" {
			dialog.ShowInformation("提示", "请选择一个模式后再开始。", newWin)
			return
		}
		targets := getTargets(hiraganaCheck.Checked, katakanaCheck.Checked)
		if len(targets) < 3 {
			dialog.ShowInformation("提示", "请选择至少2个假名进行练习。", newWin)
			return
		}
		stats := &Stats{}
		if mode == "模式一: 假名 => 罗马音" {
			// 这里使用 newWin 作为父窗口，或者也可以继续使用 main 窗口
			showModeOne(myApp, targets, newWin, stats, statsLabel, hiraganaCheck.Checked, katakanaCheck.Checked)
		} else {
			showModeTwo(myApp, targets, newWin, hiraganaCheck.Checked, katakanaCheck.Checked)
		}
	})

	// 7) 布局并设置到 newWin
	content := container.NewVBox(
		widget.NewLabel("五十音学习助手 (FiftySound)"),
		modeSelect,
		hiraganaCheck,
		katakanaCheck,
		selectKanaBtn,
		startBtn,
		statsLabel,
	)
	newWin.SetContent(content)
	newWin.Resize(fyne.NewSize(400, 300))

	// 8) 显示该子窗口
	newWin.Show()
}

// ======================= 选择假名范围 =======================
func showKanaSelectionDialog(myApp fyne.App, parent fyne.Window) {
	var lineChecks []*widget.Check
	var checks []struct {
		check *widget.Check
		char  string
		line  *widget.Check
	}

	allRandomCheck := widget.NewCheck("全部随机(包含所有五十音)", nil)
	grid := container.NewVBox(allRandomCheck)

	for _, line := range gojuon {
		lineCheck := widget.NewCheck(line.romaji+"行:", nil)
		lineChecks = append(lineChecks, lineCheck)

		rowItems := []fyne.CanvasObject{lineCheck}
		var lineCharChecks []*widget.Check

		for _, hchar := range line.hiragana {
			c := widget.NewCheck(hchar, nil)
			if contains(selectedChars, hchar) {
				c.SetChecked(true)
			}
			lineCharChecks = append(lineCharChecks, c)

			checks = append(checks, struct {
				check *widget.Check
				char  string
				line  *widget.Check
			}{c, hchar, lineCheck})

			rowItems = append(rowItems, c)
		}

		// 行选中 => 更新本行所有项
		currentLineCharChecks := lineCharChecks
		lineCheck.OnChanged = func(checked bool) {
			for _, c := range currentLineCharChecks {
				c.SetChecked(checked)
			}
		}

		// 字项 => 更新行是否全选
		for _, c := range currentLineCharChecks {
			c2 := c
			c2.OnChanged = func(_ bool) {
				allSelected := true
				for _, cc := range currentLineCharChecks {
					if !cc.Checked {
						allSelected = false
						break
					}
				}
				lineCheck.OnChanged = nil
				lineCheck.SetChecked(allSelected)
				lineCheck.OnChanged = func(checked bool) {
					for _, ch := range currentLineCharChecks {
						ch.SetChecked(checked)
					}
				}
			}
		}

		row := container.NewHBox(rowItems...)
		grid.Add(row)
	}

	// 恢复时，如果没有全选则保持当前状态
	for i, line := range gojuon {
		lineCheck := lineChecks[i]
		hiraganas := line.hiragana
		allSelected := true
		for _, hchar := range hiraganas {
			if !contains(selectedChars, hchar) {
				allSelected = false
				break
			}
		}
		lineCheck.OnChanged = nil
		lineCheck.SetChecked(allSelected)
		lineCheck.OnChanged = func(checked bool) {
			var lineCharChecks []*widget.Check
			for _, c := range checks {
				if c.line == lineCheck && contains(hiraganas, c.char) {
					lineCharChecks = append(lineCharChecks, c.check)
				}
			}
			for _, c := range lineCharChecks {
				c.SetChecked(checked)
			}
		}
	}

	allRandomCheck.OnChanged = func(checked bool) {
		if checked {
			for _, lc := range lineChecks {
				lc.SetChecked(true)
				lc.Disable()
			}
			for _, item := range checks {
				item.check.SetChecked(true)
				item.check.Disable()
			}
		} else {
			for _, lc := range lineChecks {
				lc.Enable()
			}
			for _, item := range checks {
				item.check.Enable()
			}
		}
	}

	scroll := container.NewScroll(grid)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	dialogWin := myApp.NewWindow("选择假名范围") // 原先是 a.NewWindow("选择假名范围")
	dialogWin.SetContent(container.NewVBox(
		scroll,
		widget.NewButton("确认", func() {
			if allRandomCheck.Checked {
				selectedChars = getAllHiragana()
				dialogWin.Close()
				return
			}
			var chosen []string
			for _, item := range checks {
				if item.check.Checked {
					chosen = append(chosen, item.char)
				}
			}
			selectedChars = chosen
			dialogWin.Close()
		}),
		widget.NewButton("取消", func() {
			dialogWin.Close()
		}),
	))
	dialogWin.Resize(fyne.NewSize(700, 500))
	dialogWin.Show()
}

// ======================= 辅助函数 =======================
func contains(arr []string, t string) bool {
	for _, a := range arr {
		if a == t {
			return true
		}
	}
	return false
}

func getAllHiragana() []string {
	var all []string
	for _, line := range gojuon {
		all = append(all, line.hiragana...)
	}
	return all
}

// 根据用户选择的 hiragana/katakana + selectedChars，确定最终 targets
func getTargets(h, k bool) []string {
	var result []string
	if len(selectedChars) == 0 {
		return result
	}

	type indexPos struct {
		lineIndex int
		colIndex  int
	}
	hiraganaIndexMap := make(map[string]indexPos)
	for li, line := range gojuon {
		for ci, c := range line.hiragana {
			hiraganaIndexMap[c] = indexPos{li, ci}
		}
	}

	for _, ch := range selectedChars {
		pos, ok := hiraganaIndexMap[ch]
		if !ok {
			continue
		}
		lineData := gojuon[pos.lineIndex]

		if h && pos.colIndex < len(lineData.hiragana) {
			result = append(result, lineData.hiragana[pos.colIndex])
		}
		if k && pos.colIndex < len(lineData.katakana) {
			result = append(result, lineData.katakana[pos.colIndex])
		}
	}
	return result
}

// ======================= KanaPool (真随机，不重复一轮) =======================
type KanaPool struct {
	items []string
	index int
}

func newKanaPool(targets []string) *KanaPool {
	p := &KanaPool{
		items: make([]string, len(targets)),
	}
	copy(p.items, targets)
	p.shuffle()
	return p
}

func (p *KanaPool) shuffle() {
	rand.Shuffle(len(p.items), func(i, j int) {
		p.items[i], p.items[j] = p.items[j], p.items[i]
	})
	p.index = 0
}

func (p *KanaPool) next() string {
	if p.index >= len(p.items) {
		p.shuffle()
	}
	k := p.items[p.index]
	p.index++
	return k
}

// ======================= 模式1： 假名 => 罗马音 =======================
func showModeOne(myApp fyne.App, targets []string, mainWin fyne.Window, stats *Stats, statsLabel *widget.Label, hSelected, kSelected bool) {
	w := myApp.NewWindow("模式一")
	question := widget.NewLabel("")
	answerEntry := widget.NewEntry()
	feedback := widget.NewLabel("")

	pool := newKanaPool(targets)
	var currentKana string

	nextQuestion := func() {
		answerEntry.SetText("")
		feedback.SetText("")
		currentKana = pool.next()
		prompt := ""
		if isHiragana(currentKana) {
			prompt = fmt.Sprintf("请输入平假名 %s 的罗马音：", currentKana)
		} else {
			prompt = fmt.Sprintf("请输入片假名 %s 的罗马音：", currentKana)
		}
		question.SetText(prompt)
	}

	judgeBtn := widget.NewButton("判断", func() {
		q := currentKana
		ans := strings.TrimSpace(answerEntry.Text)
		stats.Total++
		if checkRomaji(q, ans) {
			feedback.SetText("正确")
			stats.Correct++
		} else {
			feedback.SetText("错误，正确答案: " + strings.Join(kanaToRomaji[q], "/"))
		}
		statsLabel.SetText(fmt.Sprintf("当前正确率: %.2f%%", stats.Accuracy()))
	})

	nextBtn := widget.NewButton("下一题", func() {
		nextQuestion()
	})

	backBtn := widget.NewButton("返回主界面", func() {
		w.Close()
	})

	w.SetContent(container.NewVBox(
		question,
		answerEntry,
		container.NewHBox(judgeBtn, nextBtn),
		feedback,
		backBtn,
	))
	w.Resize(fyne.NewSize(400, 300))
	nextQuestion()
	w.Show()
}

// ======================= 模式2： 罗马音 => 假名手写 =======================
func showModeTwo(myApp fyne.App, targets []string, mainWin fyne.Window, hSelected, kSelected bool) {
	w := myApp.NewWindow("模式二")

	question := widget.NewLabel("")
	feedback := widget.NewLabel("用鼠标在下方空白处写出假名后，点击显示答案以进行比较。正确答案: ")
	drawingArea := newDrawingWidget()
	clearBtn := widget.NewButtonWithIcon("清空画布", theme.ContentClearIcon(), func() {
		drawingArea.Clear()
	})

	pool := newKanaPool(targets)
	var currentRomaji string
	var currentKana string

	nextQuestion := func() {
		drawingArea.Clear()
		feedback.SetText("正确答案: ")
		currentKana = pool.next()

		if val, ok := kanaToRomaji[currentKana]; ok && len(val) > 0 {
			currentRomaji = val[0]
		} else {
			currentRomaji = "a"
			currentKana = "あ"
		}

		if isHiragana(currentKana) {
			question.SetText("请写出平假名 " + currentRomaji)
		} else {
			question.SetText("请写出片假名 " + currentRomaji)
		}
	}

	showAnswerBtn := widget.NewButton("显示答案", func() {
		feedback.SetText("正确答案: " + currentKana)
	})

	nextBtn := widget.NewButton("下一题", func() {
		nextQuestion()
	})

	backBtn := widget.NewButton("返回主界面", func() {
		w.Close()
	})

	w.SetContent(container.NewBorder(
		container.NewVBox(question, container.NewHBox(showAnswerBtn, nextBtn), feedback),
		container.NewHBox(backBtn, clearBtn),
		nil, nil,
		drawingArea,
	))
	w.Resize(fyne.NewSize(600, 400))
	nextQuestion()
	w.Show()
}

// ======================= 自定义绘图widget =======================
type drawingWidget struct {
	widget.BaseWidget
	lines    []*canvas.Line
	lastPos  *fyne.Position
	bg       *canvas.Rectangle
	renderer *drawingRenderer
}

func newDrawingWidget() *drawingWidget {
	d := &drawingWidget{}
	d.ExtendBaseWidget(d)
	return d
}

func (d *drawingWidget) CreateRenderer() fyne.WidgetRenderer {
	d.bg = canvas.NewRectangle(color.White)
	d.bg.SetMinSize(fyne.NewSize(400, 300))
	d.renderer = &drawingRenderer{
		d:    d,
		objs: []fyne.CanvasObject{d.bg},
	}
	return d.renderer
}

func (d *drawingWidget) Dragged(event *fyne.DragEvent) {
	if d.lastPos == nil {
		d.lastPos = &event.Position
		return
	}
	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 2
	line.Position1 = *d.lastPos
	line.Position2 = event.Position
	d.lines = append(d.lines, line)
	d.renderer.objs = append(d.renderer.objs, line)
	canvas.Refresh(d)
	d.lastPos = &event.Position
}

func (d *drawingWidget) DragEnd() {
	d.lastPos = nil
}

func (d *drawingWidget) Tapped(_ *fyne.PointEvent) {
	// 不处理点击事件
}

func (d *drawingWidget) Clear() {
	d.lines = nil
	d.renderer.objs = []fyne.CanvasObject{d.bg}
	canvas.Refresh(d)
}

type drawingRenderer struct {
	d    *drawingWidget
	objs []fyne.CanvasObject
}

func (r *drawingRenderer) Layout(size fyne.Size) {
	r.objs[0].Resize(size)
}

func (r *drawingRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 300)
}

func (r *drawingRenderer) Refresh() {
	canvas.Refresh(r.d)
}

func (r *drawingRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *drawingRenderer) Objects() []fyne.CanvasObject {
	return r.objs
}

func (r *drawingRenderer) Destroy() {}

// ======================= 辅助：判定罗马音、反找假名、判断平假名 =======================
func checkRomaji(kana, ans string) bool {
	if val, ok := kanaToRomaji[kana]; ok {
		for _, correct := range val {
			if strings.EqualFold(ans, correct) {
				return true
			}
		}
	}
	return false
}

func reverseFindKana(romaji string) string {
	for k, v := range kanaToRomaji {
		for _, vv := range v {
			if strings.EqualFold(vv, romaji) {
				return k
			}
		}
	}
	return "null"
}

func isHiragana(c string) bool {
	for _, line := range gojuon {
		for _, h := range line.hiragana {
			if h == c {
				return true
			}
		}
	}
	return false
}
