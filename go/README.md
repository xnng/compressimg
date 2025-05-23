# 图片处理工具

这是一个多功能图片处理工具，支持图片压缩和 HEIC 转 JPG 功能。

## 功能特点

### 图片压缩功能

- **支持多种输入格式**：jpg、png、gif
- **支持多种输出格式**：jpg、png、gif
- **保留 EXIF 数据**：保留原始图片的方向信息和其他元数据
- **灵活的压缩选项**：
  - 质量设置（对于 JPG）
  - 压缩级别（对于 PNG）
  - 缩放比例
  - 最大宽度/高度限制
- **详细的压缩信息**：
  - 原始文件大小和尺寸
  - 压缩后文件大小和尺寸
  - 压缩率

### HEIC 转 JPG 功能

- 将 HEIC 格式图片转换为 JPG 格式
- 保留原始 EXIF 数据（包括图片方向信息）

## 安装

### 从源码编译

1. 克隆仓库：

```bash
git clone https://github.com/xnng/compressimg.git
cd compressimg/go
```

2. 编译：

```bash
go build -o imgtools
```

3. 将编译好的可执行文件移动到系统路径下（可选）：

```bash
sudo mv imgtools /usr/local/bin/
```

## 使用方法

### 图片压缩

基本用法：

```bash
./imgtools compress [选项] <输入文件>
```

可用选项：

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `-q` | JPEG 质量 (1-100) | 85 |
| `-pnglevel` | PNG 压缩级别 (-2 到 9) | 0 |
| `-scale` | 缩放比例 (0-1) | 1.0 |
| `-width` | 最大宽度 (像素) | 0（不限制） |
| `-height` | 最大高度 (像素) | 0（不限制） |
| `-format` | 输出格式 (jpg, png, gif) | 与输入格式相同 |
| `-o` | 输出文件路径 | [原文件名]_compressed.[格式] |

#### 示例

1. 基本压缩（使用默认参数）：

```bash
./imgtools compress image.jpg
```

2. 指定 JPEG 质量：

```bash
./imgtools compress -q 75 image.jpg
```

3. 调整图片尺寸（按比例缩放）：

```bash
./imgtools compress -scale 0.8 image.jpg
```

4. 限制最大宽度（保持宽高比）：

```bash
./imgtools compress -width 1024 image.jpg
```

5. 转换格式：

```bash
./imgtools compress -format png image.jpg
```

6. 组合多个选项：

```bash
./imgtools compress -q 80 -scale 0.7 -o output.jpg image.png
```

### HEIC 转 JPG

基本用法：

```bash
./imgtools heic2jpg <输入文件.heic> <输出文件.jpg>
```

示例：

```bash
./imgtools heic2jpg photo.heic photo.jpg
```

## 压缩效果示例

以下是一些压缩效果的示例：

1. 高质量压缩（质量 = 85，不调整尺寸）：
   - 原始大小：1.2 MB (3000x2000)
   - 压缩大小：800 KB (3000x2000)
   - 压缩率：33%

2. 中等质量压缩（质量 = 75，缩放 = 0.8）：
   - 原始大小：1.2 MB (3000x2000)
   - 压缩大小：450 KB (2400x1600)
   - 压缩率：62.5%

3. 低质量压缩（质量 = 60，缩放 = 0.5）：
   - 原始大小：1.2 MB (3000x2000)
   - 压缩大小：150 KB (1500x1000)
   - 压缩率：87.5%

## 技术细节

- 使用 Go 语言开发
- 依赖库：
  - `github.com/disintegration/imaging`：用于图像处理和调整大小
  - `github.com/jdeng/goheif`：用于 HEIC 格式处理

## 常见问题

1. **为什么有时压缩后的文件反而变大了？**
   
   这可能是因为原始图片已经被高度压缩，或者使用的质量设置太高。尝试降低质量设置或调整图片尺寸。

2. **如何获得最佳压缩效果？**
   
   通常，组合使用质量调整和尺寸调整可以获得最佳效果。例如，使用 `-q 75 -scale 0.8` 可以在保持较好图片质量的同时获得显著的压缩效果。

3. **支持批量处理吗？**
   
   当前版本不支持直接的批量处理。但你可以使用脚本来批量处理多个文件。

## 许可证

MIT
