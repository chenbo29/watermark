package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/imagex/v2"
)

type SegmentedUploadImagesT struct {
	FilePath string
	StoreKey string
}

// SegmentedUploadImages 上传文件
func SegmentedUploadImages(imagesTS []SegmentedUploadImagesT) (*imagex.CommitUploadImageResult, error) {
	// 默认 ImageX 实例为 `cn-north-1`，如果您想使用其他区域的实例，请使用 `imagex.NewInstanceWithRegion(区域名)` 显式指定区域
	instance := imagex.DefaultInstance

	instance.SetCredential(base.Credentials{
		AccessKeyID:     os.Getenv("VOLC_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("VOLC_SECRET_ACCESS_KEY"),
	})

	params := &imagex.ApplyUploadImageParam{
		ServiceId: os.Getenv("VOLC_SERVICE_ID"), // 服务 ID
		// UploadHost:     "",						//指定上传域名
		// ContentTypes:   []string{"image/jpg"},	//指定Content-Type
		// StorageClasses: []string{"ARCHIVE"},		//指定存储类型
	}

	files := make([]io.Reader, 0)
	sizeArr := make([]int64, 0)
	for _, v := range imagesTS {
		err := addFile(params, &files, &sizeArr, v.StoreKey, v.FilePath)
		if err != nil {
			log.Fatal(v.FilePath, err)
		}
	}

	return instance.SegmentedUploadImages(context.Background(), params, files, sizeArr)
}

func addFile(params *imagex.ApplyUploadImageParam, files *[]io.Reader, sizeArr *[]int64, storeKey string, filePath string) error {
	var err error
	var file *os.File
	var fileInfo os.FileInfo

	// 读取文件
	file, err = os.Open(filePath)
	if err != nil {
		return err
	}
	fileInfo, err = file.Stat()
	if err != nil {
		return err
	}

	params.StoreKeys = append(params.StoreKeys, storeKey)
	*files = append(*files, file)
	*sizeArr = append(*sizeArr, fileInfo.Size())
	return nil
}
