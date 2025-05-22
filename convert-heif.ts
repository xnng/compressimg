import * as fs from 'fs';
import * as util from 'util';
import heicConvert from 'heic-convert';
import sharp from 'sharp';

const readFile = util.promisify(fs.readFile);
const writeFile = util.promisify(fs.writeFile);

/**
 * 输出格式类型
 */
type OutputFormat = 'JPEG' | 'PNG';

/**
 * HEIF/HEIC 图片转换选项
 */
interface ConvertOptions {
  /** 输出格式 ('JPEG' 或 'PNG') */
  format?: OutputFormat;
  /** 图片质量 (1-100) */
  quality?: number;
  /** 最大宽度 (像素) */
  maxWidth?: number;
  /** 最大高度 (像素) */
  maxHeight?: number;
}

/**
 * 转换耗时统计
 */
interface ConversionTimings {
  /** 总耗时（毫秒） */
  total: number;
  /** 格式转换耗时（毫秒） */
  conversion: number;
  /** 压缩处理耗时（毫秒） */
  compression: number;
}

/**
 * 转换成功结果
 */
interface ConversionSuccess {
  /** 标记成功 */
  success: true;
  /** 原始文件大小（字节） */
  originalSize: number;
  /** 转换后文件大小（字节） */
  convertedSize: number;
  /** 压缩率（百分比） */
  compressionRatio: number;
  /** 图片宽度（像素） */
  width: number;
  /** 图片高度（像素） */
  height: number;
  /** 耗时统计 */
  timings: ConversionTimings;
}

/**
 * 转换失败结果
 */
interface ConversionFailure {
  /** 标记失败 */
  success: false;
  /** 错误信息 */
  error: string;
}

/**
 * 转换结果联合类型
 */
type ConversionResult = ConversionSuccess | ConversionFailure;

/**
 * 将 HEIF/HEIC 图片转换为指定格式
 * @param inputPath HEIF/HEIC 图片路径
 * @param outputPath 输出图片路径
 * @param options 转换选项
 * @returns 转换结果信息
 */
export async function convertHeifToImage(
  inputPath: string,
  outputPath: string,
  options: ConvertOptions = {}
): Promise<ConversionResult> {
  const { format = 'JPEG', quality = 90, maxWidth, maxHeight } = options;

  try {
    // 记录开始时间
    const startTime = performance.now();
    
    // 读取 HEIF/HEIC 文件
    const inputBuffer = await readFile(inputPath);

    // 转换为指定格式
    const conversionStartTime = performance.now();
    const convertedBuffer = await heicConvert({
      buffer: inputBuffer,
      format: format
    });
    const conversionEndTime = performance.now();
    const conversionTime = conversionEndTime - conversionStartTime;

    // 使用 sharp 处理图片质量和尺寸
    // 记录压缩处理开始时间
    const compressionStartTime = performance.now();
    let sharpInstance = sharp(convertedBuffer);

    // 调整图片大小 (如果指定了最大宽度或高度)
    if (maxWidth || maxHeight) {
      const resizeOptions: sharp.ResizeOptions = {
        fit: 'inside',
        withoutEnlargement: true
      };

      if (maxWidth) resizeOptions.width = maxWidth;
      if (maxHeight) resizeOptions.height = maxHeight;

      sharpInstance = sharpInstance.resize(resizeOptions);
    }

    // 应用质量设置并输出为指定格式
    let outputBuffer: Buffer;
    if (format === 'JPEG') {
      outputBuffer = await sharpInstance
        .jpeg({
          quality,
          mozjpeg: true
        })
        .toBuffer();
    } else {
      outputBuffer = await sharpInstance
        .png({
          quality: Math.min(Math.floor(quality / 10), 9),  // 将 1-100 转换为 0-9
        })
        .toBuffer();
    }

    // 记录压缩处理结束时间
    const compressionEndTime = performance.now();
    const compressionTime = compressionEndTime - compressionStartTime;
    
    // 获取处理后的图像信息
    const processedInfo = await sharp(outputBuffer).metadata();

    // 写入输出文件
    await writeFile(outputPath, outputBuffer);
    
    // 记录结束时间
    const endTime = performance.now();
    const totalTime = endTime - startTime;

    // 计算压缩率
    const originalSize = inputBuffer.length;
    const convertedSize = outputBuffer.length;
    const compressionRatio = 100 - (convertedSize / originalSize * 100);

    return {
      success: true,
      originalSize,
      convertedSize,
      compressionRatio,
      width: processedInfo.width || 0,
      height: processedInfo.height || 0,
      timings: {
        total: Math.round(totalTime),
        conversion: Math.round(conversionTime),
        compression: Math.round(compressionTime)
      }
    } as ConversionSuccess;
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : String(error)
    } as ConversionFailure;
  }
}

// 使用示例
async function example() {
  const result = await convertHeifToImage(
    './input.heic',
    './output.png',
    {
      format: 'PNG',
      quality: 90
    }
  );

  if (result.success) {
    console.log('转换成功!');
    console.log(`原始大小: ${(result.originalSize / 1024).toFixed(2)} KB`);
    console.log(`转换后大小: ${(result.convertedSize / 1024).toFixed(2)} KB`);
    console.log(`压缩率: ${result.compressionRatio.toFixed(2)}%`);
    console.log(`尺寸: ${result.width}x${result.height}`);
    console.log(`耗时统计:`);
    console.log(`  - 总耗时: ${result.timings.total} ms`);
    console.log(`  - 格式转换: ${result.timings.conversion} ms`);
    console.log(`  - 压缩处理: ${result.timings.compression} ms`);
  } else {
    console.error(`转换失败: ${result.error}`);
  }
}

// 如果直接运行此文件，则执行示例
if (require.main === module) {
  example().catch(console.error);
}