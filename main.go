package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 五十音图行定义
type gojuonLine struct {
	romaji   string
	hiragana []string
	katakana []string
}

var gojuon = []gojuonLine{
	{"a", []string{"あ", "い", "う", "え", "お"}, []string{"ア", "イ", "ウ", "エ", "オ"}},
	{"ka", []string{"か", "き", "く", "け", "こ"}, []string{"カ", "キ", "ク", "ケ", "コ"}},
	{"sa", []string{"さ", "し", "す", "せ", "そ"}, []string{"サ", "シ", "ス", "セ", "ソ"}},
	{"ta", []string{"た", "ち", "つ", "て", "と"}, []string{"タ", "チ", "ツ", "テ", "ト"}},
	{"na", []string{"な", "に", "ぬ", "ね", "の"}, []string{"ナ", "ニ", "ヌ", "ネ", "ノ"}},
	{"ha", []string{"は", "ひ", "ふ", "へ", "ほ"}, []string{"ハ", "ヒ", "フ", "ヘ", "ホ"}},
	{"ma", []string{"ま", "み", "む", "め", "も"}, []string{"マ", "ミ", "ム", "メ", "モ"}},
	{"ya", []string{"や", "ゆ", "よ"}, []string{"ヤ", "ユ", "ヨ"}},
	{"ra", []string{"ら", "り", "る", "れ", "ろ"}, []string{"ラ", "リ", "ル", "レ", "ロ"}},
	{"wa", []string{"わ", "を", "ん"}, []string{"ワ", "ヲ", "ン"}},
}

var kanaToRomaji = map[string][]string{
	// Hiragana
	"あ": {"a"}, "い": {"i"}, "う": {"u"}, "え": {"e"}, "お": {"o"},
	"か": {"ka"}, "き": {"ki"}, "く": {"ku"}, "け": {"ke"}, "こ": {"ko"},
	"さ": {"sa"}, "し": {"shi", "si"}, "す": {"su"}, "せ": {"se"}, "そ": {"so"},
	"た": {"ta"}, "ち": {"chi", "ti"}, "つ": {"tsu", "tu"}, "て": {"te"}, "と": {"to"},
	"な": {"na"}, "に": {"ni"}, "ぬ": {"nu"}, "ね": {"ne"}, "の": {"no"},
	"は": {"ha"}, "ひ": {"hi"}, "ふ": {"fu", "hu"}, "へ": {"he"}, "ほ": {"ho"},
	"ま": {"ma"}, "み": {"mi"}, "む": {"mu"}, "め": {"me"}, "も": {"mo"},
	"や": {"ya"}, "ゆ": {"yu"}, "よ": {"yo"},
	"ら": {"ra"}, "り": {"ri"}, "る": {"ru"}, "れ": {"re"}, "ろ": {"ro"},
	"わ": {"wa"}, "を": {"o", "wo"}, "ん": {"n"},

	// Katakana
	"ア": {"a"}, "イ": {"i"}, "ウ": {"u"}, "エ": {"e"}, "オ": {"o"},
	"カ": {"ka"}, "キ": {"ki"}, "ク": {"ku"}, "ケ": {"ke"}, "コ": {"ko"},
	"サ": {"sa"}, "シ": {"shi", "si"}, "ス": {"su"}, "セ": {"se"}, "ソ": {"so"},
	"タ": {"ta"}, "チ": {"chi", "ti"}, "ツ": {"tsu", "tu"}, "テ": {"te"}, "ト": {"to"},
	"ナ": {"na"}, "ニ": {"ni"}, "ヌ": {"nu"}, "ネ": {"ne"}, "ノ": {"no"},
	"ハ": {"ha"}, "ヒ": {"hi"}, "フ": {"fu", "hu"}, "ヘ": {"he"}, "ホ": {"ho"},
	"マ": {"ma"}, "ミ": {"mi"}, "ム": {"mu"}, "メ": {"me"}, "モ": {"mo"},
	"ヤ": {"ya"}, "ユ": {"yu"}, "ヨ": {"yo"},
	"ラ": {"ra"}, "リ": {"ri"}, "ル": {"ru"}, "レ": {"re"}, "ロ": {"ro"},
	"ワ": {"wa"}, "ヲ": {"o", "wo"}, "ン": {"n"},
}

// 用户选择的假名（已选定的平假名字符）
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

func main() {
	rand.Seed(time.Now().UnixNano())
	a := app.New()
	w := a.NewWindow("五十音学习助手 (FiftySound)")

	modeSelect := widget.NewSelect([]string{"模式一: 假名 => 罗马音", "模式二: 罗马音 => 假名手写"}, func(string) {})
	modeSelect.PlaceHolder = "请点击下拉框，选择你想要的模式"

	hiraganaCheck := widget.NewCheck("平假名", nil)
	hiraganaCheck.SetChecked(true)
	katakanaCheck := widget.NewCheck("片假名", nil)
	katakanaCheck.SetChecked(true)

	statsLabel := widget.NewLabel("当前正确率: 0.00%")

	selectKanaBtn := widget.NewButton("选择假名范围", func() {
		showKanaSelectionDialog(a, w)
	})

	startBtn := widget.NewButton("开始", func() {
		mode := modeSelect.Selected
		if mode == "" {
			dialog.ShowInformation("提示", "请选择一个模式后再开始。", w)
			return
		}
		targets := getTargets(hiraganaCheck.Checked, katakanaCheck.Checked)
		if len(targets) < 2 {
			dialog.ShowInformation("提示", "请选择至少2个假名进行练习。", w)
			return
		}
		stats := &Stats{}
		if mode == "模式一: 假名 => 罗马音" {
			showModeOne(a, targets, w, stats, statsLabel, hiraganaCheck.Checked, katakanaCheck.Checked)
		} else {
			showModeTwo(a, targets, w, hiraganaCheck.Checked, katakanaCheck.Checked)
		}
	})

	content := container.NewVBox(
		widget.NewLabel("五十音学习助手 (FiftySound)"),
		modeSelect,
		hiraganaCheck,
		katakanaCheck,
		selectKanaBtn,
		startBtn,
		statsLabel,
	)
	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}

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

// 根据用户选择的hiragana/katakana和selectedChars确定最终targets
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

// 显示选择假名的对话框
func showKanaSelectionDialog(a fyne.App, parent fyne.Window) {
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

		// 行选中变化时更新本行所有项
		currentLineCharChecks := lineCharChecks
		lineCheck.OnChanged = func(checked bool) {
			for _, c := range currentLineCharChecks {
				c.SetChecked(checked)
			}
		}

		// 字项变化时更新行选中状态
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

	// 恢复时，如果没有全选，则保持当前状态
	// 如果用户之前选了整行，应当此时也更新行选择框的状态
	for i, line := range gojuon {
		// 对每行检查是否已全选
		lineCheck := lineChecks[i]
		// 找出该行所有假名
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
			// 找到该行对应的字check并更新
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
				// 恢复初始状态：不自动清空用户之前勾选，因为用户可能想重新调整
				// 这里不强制取消全选，以保持用户的状态
				// 如果想要恢复全部未选中，可以启用下面代码：
				// lc.SetChecked(false)
			}
			for _, item := range checks {
				item.check.Enable()
				// 同理这里不强制取消勾选，因为用户可能想在不全随机的情况下修改
				// 如需强制清空可加:
				// item.check.SetChecked(false)
			}
		}
	}

	scroll := container.NewScroll(grid)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	dialogWin := a.NewWindow("选择假名范围")
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

// KanaPool 用于在练习中轮流出现不重复的题目，
// 用完之后重新洗牌，以确保每个假名都出现后才重复。
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

// 模式一： 假名 => 罗马音（使用 KanaPool 确保每个假名轮询出现）
func showModeOne(a fyne.App, targets []string, mainWin fyne.Window, stats *Stats, statsLabel *widget.Label, hSelected, kSelected bool) {
	w := a.NewWindow("模式一")
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

// 模式二： 罗马音 => 假名手写（使用 KanaPool 确保每个假名轮询出现）
func showModeTwo(a fyne.App, targets []string, mainWin fyne.Window, hSelected, kSelected bool) {
	w := a.NewWindow("模式二")

	question := widget.NewLabel("")
	feedback := widget.NewLabel("正确答案: ")
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
			currentRomaji = val[rand.Intn(len(val))]
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

// 自定义绘图widget
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
