// Package main 提供图片压缩功能
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// 压缩结果信息
type CompressResult struct {
	OriginalPath     string  // 原始文件路径
	OutputPath       string  // 输出文件路径
	OriginalSize     int64   // 原始文件大小（字节）
	CompressedSize   int64   // 压缩后文件大小（字节）
	CompressionRatio float64 // 压缩率（百分比）
	OriginalWidth    int     // 原始图片宽度
	OriginalHeight   int     // 原始图片高度
	CompressedWidth  int     // 压缩后图片宽度
	CompressedHeight int     // 压缩后图片高度
}

// 压缩选项
type CompressOptions struct {
	Quality     int     // JPEG 质量 (1-100)
	PNGLevel    png.CompressionLevel // PNG 压缩级别
	Scale       float64 // 缩放比例 (0-1)
	MaxWidth    int     // 最大宽度（如果为0则不限制）
	MaxHeight   int     // 最大高度（如果为0则不限制）
	OutputFormat string // 输出格式 (jpg, png, gif)
}

// 默认压缩选项
func DefaultCompressOptions() CompressOptions {
	return CompressOptions{
		Quality:     85,
		PNGLevel:    png.DefaultCompression,
		Scale:       1.0,
		MaxWidth:    0,
		MaxHeight:   0,
		OutputFormat: "",
	}
}

// 压缩图片
func CompressImage(inputPath, outputPath string, options CompressOptions) (*CompressResult, error) {
	// 检查输入文件是否存在
	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("无法访问输入文件: %v", err)
	}

	// 获取原始文件大小
	originalSize := inputInfo.Size()

	// 打开输入文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开输入文件: %v", err)
	}
	defer inputFile.Close()

	// 解码图片
	img, format, err := image.Decode(inputFile)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %v", err)
	}

	// 获取原始图片尺寸
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// 调整图片大小（如果需要）
	resizedImg := img
	if options.Scale < 1.0 && options.Scale > 0 {
		newWidth := int(float64(originalWidth) * options.Scale)
		newHeight := int(float64(originalHeight) * options.Scale)
		resizedImg = imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
	} else if options.MaxWidth > 0 || options.MaxHeight > 0 {
		maxWidth := options.MaxWidth
		maxHeight := options.MaxHeight
		
		// 如果只指定了一个维度，保持宽高比
		if maxWidth == 0 {
			maxWidth = int(float64(originalWidth) * float64(maxHeight) / float64(originalHeight))
		} else if maxHeight == 0 {
			maxHeight = int(float64(originalHeight) * float64(maxWidth) / float64(originalWidth))
		}
		
		// 只有当原始尺寸超过最大尺寸时才调整大小
		if originalWidth > maxWidth || originalHeight > maxHeight {
			resizedImg = imaging.Resize(img, maxWidth, maxHeight, imaging.Lanczos)
		}
	}

	// 获取调整后的尺寸
	resizedBounds := resizedImg.Bounds()
	compressedWidth := resizedBounds.Dx()
	compressedHeight := resizedBounds.Dy()

	// 确定输出格式
	outputFormat := options.OutputFormat
	if outputFormat == "" {
		// 如果未指定输出格式，使用输出文件的扩展名
		outputFormat = strings.ToLower(strings.TrimPrefix(filepath.Ext(outputPath), "."))
		// 如果输出文件没有扩展名，使用输入格式
		if outputFormat == "" {
			outputFormat = format
		}
	}

	// 创建输出目录（如果不存在）
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outputFile.Close()

	// 提取 EXIF 数据（如果是 JPEG 输入）
	var exifData []byte
	if format == "jpeg" || format == "jpg" {
		// 重新打开输入文件以读取 EXIF 数据
		exifFile, err := os.Open(inputPath)
		if err == nil {
			defer exifFile.Close()
			
			// 直接读取文件的前 65536 字节，这通常包含了 EXIF 数据
			buf := make([]byte, 65536)
			n, err := exifFile.Read(buf)
			if err == nil || err == io.EOF {
				// 寻找 EXIF 标记
				for i := 0; i < n-1; i++ {
					if buf[i] == 0xFF && buf[i+1] == 0xE1 {
						// 获取 EXIF 数据长度
						if i+3 < n {
							length := int(buf[i+2])<<8 | int(buf[i+3])
							// 确保我们有足够的数据
							if i+2+length <= n {
								// 复制 EXIF 数据（不包括标记和长度）
								exifData = make([]byte, length-2)
								copy(exifData, buf[i+4:i+2+length])
								break
							}
						}
					}
				}
			}
		}
	}
	
	// 根据输出格式编码图片，保留 EXIF 数据
	switch outputFormat {
	case "jpg", "jpeg":
		if len(exifData) > 0 {
			// 写入 JPEG 文件头
			soi := []byte{0xff, 0xd8}
			if _, err := outputFile.Write(soi); err != nil {
				return nil, fmt.Errorf("写入 JPEG 文件头失败: %v", err)
			}
			
			// 写入 EXIF 数据
			app1Marker := 0xe1
			markerlen := 2 + len(exifData)
			marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
			if _, err := outputFile.Write(marker); err != nil {
				return nil, fmt.Errorf("写入 EXIF 标记失败: %v", err)
			}
			
			if _, err := outputFile.Write(exifData); err != nil {
				return nil, fmt.Errorf("写入 EXIF 数据失败: %v", err)
			}
			
			// 创建跳过 SOI 标记的写入器
			writer := &writerSkipper{outputFile, 2}
			err = jpeg.Encode(writer, resizedImg, &jpeg.Options{Quality: options.Quality})
		} else {
			// 没有 EXIF 数据，正常编码
			err = jpeg.Encode(outputFile, resizedImg, &jpeg.Options{Quality: options.Quality})
		}
	case "png":
		encoder := png.Encoder{CompressionLevel: options.PNGLevel}
		err = encoder.Encode(outputFile, resizedImg)
	case "gif":
		err = gif.Encode(outputFile, resizedImg, nil)
	default:
		return nil, fmt.Errorf("不支持的输出格式: %s", outputFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("编码图片失败: %v", err)
	}

	// 获取压缩后的文件大小
	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("无法获取输出文件信息: %v", err)
	}
	compressedSize := outputInfo.Size()

	// 计算压缩率
	compressionRatio := 100.0 - (float64(compressedSize) / float64(originalSize) * 100.0)

	// 返回压缩结果
	return &CompressResult{
		OriginalPath:     inputPath,
		OutputPath:       outputPath,
		OriginalSize:     originalSize,
		CompressedSize:   compressedSize,
		CompressionRatio: compressionRatio,
		OriginalWidth:    originalWidth,
		OriginalHeight:   originalHeight,
		CompressedWidth:  compressedWidth,
		CompressedHeight: compressedHeight,
	}, nil
}

// 格式化文件大小为人类可读格式
func FormatFileSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// 检查文件格式是否受支持
func IsSupportedFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	default:
		return false
	}
}

// 从文件路径获取格式
func GetFormatFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	return strings.TrimPrefix(ext, ".")
}

// 将图片编码到缓冲区
func EncodeImageToBuffer(img image.Image, format string, options CompressOptions) ([]byte, error) {
	var buf bytes.Buffer
	
	switch format {
	case "jpg", "jpeg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: options.Quality})
		if err != nil {
			return nil, err
		}
	case "png":
		encoder := png.Encoder{CompressionLevel: options.PNGLevel}
		err := encoder.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
	case "gif":
		err := gif.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("不支持的格式: %s", format)
	}
	
	return buf.Bytes(), nil
}

// 从缓冲区解码图片
func DecodeImageFromBuffer(data []byte) (image.Image, string, error) {
	return image.Decode(bytes.NewReader(data))
}

// 从文件读取图片
func ReadImageFromFile(path string) (image.Image, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	
	return image.Decode(file)
}

// 将图片写入文件
func WriteImageToFile(img image.Image, path string, format string, options CompressOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	switch format {
	case "jpg", "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: options.Quality})
	case "png":
		encoder := png.Encoder{CompressionLevel: options.PNGLevel}
		return encoder.Encode(file, img)
	case "gif":
		return gif.Encode(file, img, nil)
	default:
		return fmt.Errorf("不支持的格式: %s", format)
	}
}

// writerSkipper 结构体已移动到 utils.go 文件中

// 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = io.Copy(destFile, sourceFile)
	return err
}
