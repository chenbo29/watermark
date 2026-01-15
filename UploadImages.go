package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/imagex/v2"
)

type UploadImagesT struct {
	FilePath string
	StoreKey string
}

type ImageX struct {
	ctx    context.Context
	config *Config
}

func NewImageX() *ImageX {
	return &ImageX{
		ctx:    context.Background(),
		config: LoadConfig(),
	}
}

// UploadImages 上传文件
func (i *ImageX) UploadImages(imagesTS []UploadImagesT) (*imagex.CommitUploadImageResult, error) {
	var (
		err         error
		ret         *imagex.CommitUploadImageResult
		imagesTSNew = make([]UploadImagesT, 0)
	)
	for k, v := range imagesTS {
		imagesTSNew = append(imagesTSNew, v)
		if k%10 == 0 {
			ret, err = i.uploadImagesSingle(imagesTSNew)
			if err != nil {
				return ret, err
			}
			imagesTSNew = imagesTSNew[:0]
		}
	}
	ret, err = i.uploadImagesSingle(imagesTSNew)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (i *ImageX) uploadImagesSingle(imagesTS []UploadImagesT) (*imagex.CommitUploadImageResult, error) {
	// 默认 ImageX 实例为 `cn-north-1`，如果您想使用其他区域的实例，请使用 `imagex.NewInstanceWithRegion(区域名)` 显式指定区域
	instance := imagex.DefaultInstance

	instance.SetCredential(base.Credentials{
		AccessKeyID:     i.config.VolcAccessKeyId,
		SecretAccessKey: i.config.VolcSecretAccessKey,
	})

	params := &imagex.ApplyUploadImageParam{
		ServiceId: i.config.VolcServiceId, // 服务 ID
		//UploadHost:     "quickok-blog.tos-cn-shanghai.volces.com", //指定上传域名
		//ContentTypes:   []string{"image/jpg"}, //指定Content-Type
		//StorageClasses: []string{"ARCHIVE"},   //指定存储类型
		Overwrite: false,
	}

	files := make([][]byte, 0)
	for _, v := range imagesTS {
		var (
			file *os.File
		)
		dat, err := os.ReadFile(v.FilePath)
		if err != nil {
			log.Println(fmt.Sprintf("文件：%v，读取失败：%v", v.FilePath, err))
		} else {
			files = append(files, dat)
			// 读取文件
			file, _ = os.Open(v.FilePath)
			params.StoreKeys = append(params.StoreKeys, filepath.Base(v.FilePath))
			params.ContentTypes = append(params.ContentTypes, GetImageContentType(file))
			_ = file.Close()
		}
	}
	ret, err := instance.UploadImages(params, files)
	if err != nil {
		log.Println(fmt.Sprintf("参数信息：%+v", params))
	}
	return ret, err
}

func GetImageContentType(r *os.File) string {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return ""
	}
	return http.DetectContentType(buf[:n])
}
