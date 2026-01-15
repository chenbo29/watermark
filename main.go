package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dromara/carbon/v2"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/volcengine/volc-sdk-golang/service/imagex/v2"
)

func main() {
	a := app.NewWithID("com.watermark")
	w := a.NewWindow("图片处理")
	config := LoadConfig()

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
			showFiles(config, path, list)
		}, w)
	})

	w.SetContent(container.NewBorder(
		btn, nil, nil, nil,
		container.NewScroll(list),
	))
	w.ShowAndRun()
}

func showFiles(config *Config, dir string, output *widget.Entry) {
	var (
		err    error
		images = make([]UploadImagesT, 0)
	)
	output.SetText(fmt.Sprintf("开始读取文件：%v\n", carbon.Now().ToDateTimeString()))
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
		images = append(images, UploadImagesT{FilePath: path, StoreKey: fmt.Sprintf("%v%v", uuid.NewString(), filepath.Ext(path))})
		return nil
	})
	output.SetText(output.Text + fmt.Sprintf("文件数量：%d\n", len(images)))
	if len(images) > int(config.MaxFileNum) {
		output.SetText(output.Text + fmt.Sprintf("文件数量不能超过%d", config.MaxFileNum))
	}
	var ret *imagex.CommitUploadImageResult
	ret, err = NewImageX().UploadImages(images)
	if err != nil {
		log.Println(fmt.Sprintf("上传出错：%v", err))
	} else {
		log.Println(fmt.Sprintf("上传结果：%+v", ret))
	}

}

func fileType(info os.FileInfo) string {
	if info.IsDir() {
		return "DIR"
	}
	return "FILE"
}
