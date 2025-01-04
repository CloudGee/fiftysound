package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	// 假设你的五十音模块在这个路径
	"FiftySound/modules/fifty_sounds"

	// 我们要用到 ShowVocabularyPractice
	"FiftySound/modules/vocabulary"
)

func main() {
	myApp := app.New()
	myWin := myApp.NewWindow("日语学习 - 主菜单")

	// 五十音按钮
	btnFiftySounds := widget.NewButton("五十音练习", func() {
		// 这里传入两个参数：myApp, myWin
		// 具体实现你在 fifty_sounds.ShowFiftySounds 里写
		fifty_sounds.ShowFiftySounds(myApp, myWin)
	})

	// 新标日语单词练习按钮
	btnVocabulary := widget.NewButton("新标日语单词练习", func() {
		// 把 myApp 和 myWin 一起传
		vocabulary.ShowVocabularyMainPage(myApp, myWin)
	})

	myWin.SetContent(container.NewVBox(
		widget.NewLabel("请选择要进入的功能："),
		btnFiftySounds,
		btnVocabulary,
	))
	myWin.Resize(fyne.NewSize(400, 300))
	myWin.ShowAndRun()
}
