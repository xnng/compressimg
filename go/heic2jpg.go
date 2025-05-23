// Package main 提供简单的 HEIC 转 JPG 功能
package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"github.com/jdeng/goheif"
)

// 将 HEIC 文件转换为 JPG 文件
// 保留原始 EXIF 数据（包括图片方向信息）
func convertHeicToJpg(inputPath, outputPath string) error {
	// 打开输入文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %v", err)
	}
	defer inputFile.Close()

	// 提取 EXIF 数据
	exif, err := goheif.ExtractExif(inputFile)
	if err != nil {
		log.Printf("警告: 无法从 %s 提取 EXIF 数据: %v\n", inputPath, err)
		// 继续处理，即使没有 EXIF 数据
	}

	// 重置文件指针
	_, err = inputFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("重置文件指针失败: %v", err)
	}

	// 解码 HEIC 图片
	img, err := goheif.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("解码 HEIC 图片失败: %v", err)
	}

	// 创建输出目录（如果不存在）
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outputFile.Close()

	// 准备写入器，保留 EXIF 数据
	var writer interface {
		Write([]byte) (int, error)
	} = outputFile

	if exif != nil {
		// 写入 JPEG 文件头
		soi := []byte{0xff, 0xd8}
		if _, err := outputFile.Write(soi); err != nil {
			return fmt.Errorf("写入 JPEG 文件头失败: %v", err)
		}

		// 写入 EXIF 数据
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := outputFile.Write(marker); err != nil {
			return fmt.Errorf("写入 EXIF 标记失败: %v", err)
		}

		if _, err := outputFile.Write(exif); err != nil {
			return fmt.Errorf("写入 EXIF 数据失败: %v", err)
		}

		// 创建跳过 SOI 标记的写入器
		writer = &writerSkipper{outputFile, 2}
	}

	// 编码为 JPG
	err = jpeg.Encode(writer, img, nil)
	if err != nil {
		return fmt.Errorf("编码 JPG 失败: %v", err)
	}

	return nil
}
// writerSkipper 结构体已移动到 utils.go 文件中

// heic2jpgMain 是 heic2jpg 功能的主函数
func heic2jpgMain() {
	// 检查命令行参数
	if len(os.Args) < 3 {
		fmt.Println("error_参数不足，用法: heic2jpg <输入文件.heic> <输出文件.jpg>")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// 确保输入文件存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Println("error_输入文件不存在")
		return
	}

	// 确保输入文件是 HEIC 格式
	ext := filepath.Ext(inputPath)
	if ext != ".heic" && ext != ".HEIC" {
		fmt.Println("error_输入文件不是 HEIC 格式")
		return
	}

	// 确保输出文件是 JPG 格式
	ext = filepath.Ext(outputPath)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".JPG" && ext != ".JPEG" {
		fmt.Println("error_输出文件不是 JPG 格式")
		return
	}

	// 执行转换
	err := convertHeicToJpg(inputPath, outputPath)
	if err != nil {
		fmt.Printf("error_%s\n", err.Error())
		return
	}

	fmt.Println("success_转换成功")
}
