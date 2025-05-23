// Package main 提供图片压缩命令行工具
package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func compressCli() {
	// 定义命令行参数
	var (
		quality      = flag.Int("q", 85, "JPEG 质量 (1-100)")
		pngLevel     = flag.Int("pnglevel", int(png.DefaultCompression), "PNG 压缩级别 (-2 到 9)")
		scale        = flag.Float64("scale", 1.0, "缩放比例 (0-1)")
		maxWidth     = flag.Int("width", 0, "最大宽度 (像素)")
		maxHeight    = flag.Int("height", 0, "最大高度 (像素)")
		outputFormat = flag.String("format", "", "输出格式 (jpg, png, gif), 默认与输入格式相同")
		outputPath   = flag.String("o", "", "输出文件路径, 默认为 [原文件名]_compressed.[格式]")
	)

	// 解析命令行参数
	flag.Parse()

	// 获取输入文件路径
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("error_参数不足，用法: compress [选项] <输入文件>")
		flag.PrintDefaults()
		return
	}

	inputPath := args[0]

	// 确保输入文件存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Println("error_输入文件不存在")
		return
	}

	// 确保输入文件格式受支持
	if !IsSupportedFormat(inputPath) {
		fmt.Println("error_不支持的输入文件格式，支持的格式: jpg, jpeg, png, gif")
		return
	}

	// 如果未指定输出路径，生成默认输出路径
	output := *outputPath
	if output == "" {
		ext := filepath.Ext(inputPath)
		baseName := strings.TrimSuffix(inputPath, ext)
		
		// 确定输出格式
		outputExt := ext
		if *outputFormat != "" {
			outputExt = "." + *outputFormat
		}
		
		output = baseName + "_compressed" + outputExt
	}

	// 设置压缩选项
	options := DefaultCompressOptions()
	options.Quality = *quality
	options.PNGLevel = png.CompressionLevel(*pngLevel)
	options.Scale = *scale
	options.MaxWidth = *maxWidth
	options.MaxHeight = *maxHeight
	options.OutputFormat = *outputFormat

	// 执行压缩
	result, err := CompressImage(inputPath, output, options)
	if err != nil {
		fmt.Printf("error_%s\n", err.Error())
		return
	}

	// 输出压缩结果
	fmt.Println("success_压缩成功")
	fmt.Printf("原始文件: %s\n", result.OriginalPath)
	fmt.Printf("输出文件: %s\n", result.OutputPath)
	fmt.Printf("原始大小: %s (%dx%d)\n", FormatFileSize(result.OriginalSize), result.OriginalWidth, result.OriginalHeight)
	fmt.Printf("压缩大小: %s (%dx%d)\n", FormatFileSize(result.CompressedSize), result.CompressedWidth, result.CompressedHeight)
	fmt.Printf("压缩率: %.2f%%\n", result.CompressionRatio)
}
