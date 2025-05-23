// heic2jpg.js
const { execFileSync } = require('child_process');
const path = require('path');

/**
 * 将 HEIC 图片转换为 JPG 格式
 * @param {string} inputPath - HEIC 文件路径
 * @param {string} outputPath - JPG 输出路径
 * @returns {Object} 转换结果对象 {success, message, duration}
 */
function convertHeicToJpg(inputPath, outputPath) {
  // 可执行文件路径
  const exePath = path.resolve(__dirname, 'heic2jpg');

  const result = {
    success: false,
    message: '',
    duration: 0
  };

  const start = Date.now();

  try {
    // 执行转换
    const stdout = execFileSync(exePath, [inputPath, outputPath], { encoding: 'utf-8' });
    const end = Date.now();

    // 记录耗时
    result.duration = (end - start) / 1000;

    // 解析输出结果，检查前缀
    const output = stdout.trim();

    if (output.startsWith('success_')) {
      // 成功消息
      result.success = true;
      result.message = output.substring(8);
    } else {
      // 错误消息
      result.success = false;
      result.message = output.substring(6);
    }
    console.log(result)
  } catch (error) {
    console.error(error);
  }
}

convertHeicToJpg("./input22.heic", "./output.jpg");