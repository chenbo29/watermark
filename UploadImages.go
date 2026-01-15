package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/imagex/v2"
)

type SegmentedUploadImagesT struct {
	FilePath string
	StoreKey string
}

// UploadImages 上传文件
func UploadImages(imagesTS []SegmentedUploadImagesT) (*imagex.CommitUploadImageResult, error) {
	var (
		err         error
		ret         *imagex.CommitUploadImageResult
		imagesTSNew = make([]SegmentedUploadImagesT, 0)
	)
	for k, v := range imagesTS {
		imagesTSNew = append(imagesTSNew, v)
		if k%10 == 0 {
			ret, err = UploadImagesSingle(imagesTSNew)
			if err != nil {
				return ret, err
			}
			imagesTSNew = imagesTSNew[:0]
		}
	}
	return ret, err
}

func UploadImagesSingle(imagesTS []SegmentedUploadImagesT) (*imagex.CommitUploadImageResult, error) {
	// 默认 ImageX 实例为 `cn-north-1`，如果您想使用其他区域的实例，请使用 `imagex.NewInstanceWithRegion(区域名)` 显式指定区域
	instance := imagex.DefaultInstance

	instance.SetCredential(base.Credentials{
		AccessKeyID:     os.Getenv("VOLC_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("VOLC_SECRET_ACCESS_KEY"),
	})

	params := &imagex.ApplyUploadImageParam{
		ServiceId: "wpk84y9dot", // 服务 ID
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
			params.StoreKeys = append(params.StoreKeys, file.Name())
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
	params.ContentTypes = append(params.ContentTypes, GetImageContentType(file))
	*files = append(*files, file)
	*sizeArr = append(*sizeArr, fileInfo.Size())
	return nil
}
