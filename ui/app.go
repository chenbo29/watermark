package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
)

func RunApp() {
	a := app.NewWithID("com.watermark.ui")
	w := a.NewWindow("文件夹文件查看器")
	w.Resize(fyne.NewSize(800, 500))

	// 文件信息列表
	list := widget.NewMultiLineEntry()
	list.Wrapping = fyne.TextWrapWord
	list.SetPlaceHolder("请选择一个文件夹")

	// 选择文件夹按钮
	btn := widget.NewButton("选择文件夹", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			path := uri.Path()
			showFiles(path, list)
		}, w)
	})

	w.SetContent(container.NewBorder(
		btn, nil, nil, nil,
		container.NewScroll(list),
	))

	w.ShowAndRun()
}

func showFiles(dir string, output *widget.Entry) {
	output.SetText("")
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if path == dir {
			return nil
		}

		line := fmt.Sprintf(
			"%s | %s | %v\n",
			path,
			fileType(info),
			humanize.Bytes(uint64(info.Size())),
			//info.ModTime().Format(time.DateTime),
		)
		output.SetText(output.Text + line)
		return nil
	})
}

func fileType(info os.FileInfo) string {
	if info.IsDir() {
		return "DIR"
	}
	return "FILE"
}
