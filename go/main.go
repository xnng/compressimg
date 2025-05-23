// Package main 提供图片处理工具
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// 获取子命令
	subCommand := os.Args[1]
	
	// 根据子命令执行相应的功能
	switch subCommand {
	case "heic2jpg":
		// 移除子命令，将剩余参数传递给 heic2jpg 功能
		os.Args = append(os.Args[:1], os.Args[2:]...)
		heic2jpgMain()
	case "compress":
		// 移除子命令，将剩余参数传递给压缩功能
		os.Args = append(os.Args[:1], os.Args[2:]...)
		compressCli()
	case "help":
		printUsage()
	default:
		// 如果第一个参数是文件路径，尝试自动判断功能
		if fileExists(subCommand) {
			ext := strings.ToLower(filepath.Ext(subCommand))
			if ext == ".heic" {
				// 如果是 HEIC 文件，执行 heic2jpg 功能
				heic2jpgMain()
			} else if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				// 如果是支持的图片格式，执行压缩功能
				compressCli()
			} else {
				fmt.Println("error_不支持的文件格式")
				printUsage()
			}
		} else {
			fmt.Println("error_未知命令或文件不存在")
			printUsage()
		}
	}
}

// 打印使用说明
func printUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  图片工具 heic2jpg <输入文件.heic> <输出文件.jpg>")
	fmt.Println("  图片工具 compress [选项] <输入文件>")
	fmt.Println("\n压缩选项:")
	fmt.Println("  -q int")
	fmt.Println("        JPEG 质量 (1-100) (默认 85)")
	fmt.Println("  -pnglevel int")
	fmt.Println("        PNG 压缩级别 (-2 到 9) (默认 0)")
	fmt.Println("  -scale float")
	fmt.Println("        缩放比例 (0-1) (默认 1.0)")
	fmt.Println("  -width int")
	fmt.Println("        最大宽度 (像素) (默认 0，不限制)")
	fmt.Println("  -height int")
	fmt.Println("        最大高度 (像素) (默认 0，不限制)")
	fmt.Println("  -format string")
	fmt.Println("        输出格式 (jpg, png, gif), 默认与输入格式相同")
	fmt.Println("  -o string")
	fmt.Println("        输出文件路径, 默认为 [原文件名]_compressed.[格式]")
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// heic2jpg 主函数在 heic2jpg.go 文件中定义
