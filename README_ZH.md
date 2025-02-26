
<h1 align=center>QuickPiperAudiobook</h1>
<p align="center">
  <a href="./README.md">English</a>
  |
  <a href="./README_PL.md">Polski</a>
  |
  <b> 简体中文 </b>
</p>

只需一条命令即可为任何文本内容生成自然语音的有声读物。
 - 支持转换 [PDF、EPUB、TXT、MOBI、DJVU、HTML、DOCX 等多种格式](https://manual.calibre-ebook.com/generated/en/ebook-convert.html)
 - 所有转换都在本地完成，保证数据隐私
 - 基于 [Piper 模型](https://rhasspy.github.io/piper-samples/)；支持多种语言
在[此处](./examples/)听取示例音频
## 安装
1. 下载预构建的[发行版](https://github.com/C-Loftus/QuickPiperAudiobook/releases/)
    * (或通过 `go mod tidy && go build` 从源码构建)
2. 安装 `ebook-convert` 并确保它在系统 PATH 中（该工具包含在 [Calibre](https://calibre-ebook.com/) 中）
3. *(可选)* 安装 `ffmpeg` 以支持 MP3 和章节功能
> [!NOTE]  
> 无需单独安装 Piper。本程序会自动管理 Piper 及其相关模型
## 使用方法
* 输入本地文件路径或带有正确扩展名的远程 URL
   * 示例：`./QuickPiperAudiobook test.txt`
* 使用 `--chapters` 参数为 EPUB 文件生成带章节的 MP3
   * 示例：`./QuickPiperAudiobook --chapters test.epub`
* 使用 `--help` 参数查看所有可用选项
   * 示例：`./QuickPiperAudiobook --help`
### 非英语/UTF-8 支持
* 从 [Piper 模型库](https://rhasspy.github.io/piper-samples/)下载您所需语言的模型文件（.onnx 和 .json）
  * 例如：`zh_CN-huayan-medium.onnx` 及对应的 `zh_CN-huayan-medium.onnx`（如需要请重命名）
* 将它们放置在 `~/.config/QuickPiperAudiobook/` 目录中
* 使用 `--speak-utf-8` 和 `--model=` 参数指定用特定模型处理 UTF-8 字符
  * 示例：`./QuickPiperAudiobook --speak-utf-8 --model=zh_CN-huayan-medium.onnx test_chinese_data.pdf`
> [!提示]  
> 如果您经常使用某个模型，建议在配置文件中将其设为默认模型
### 配置设置
* 可以在 `~/.config/QuickPiperAudiobook/` 创建配置文件，以避免每次都需要在命令行指定参数
  * 例如：可以将任意模型的 `.onnx` 和 `.onnx.json` 文件放在 `~/.config/QuickPiperAudiobook/` 目录中
  * 完整的配置示例请参考[这里](./examples/config.yaml)
```yml
# `~/.config/QuickPiperAudiobook/config.yaml` 示例
# 用户未指定 --output 参数时使用的默认输出目录
output: ~/Audiobooks
# 用户未指定 --model 参数时使用的默认模型
model: "en_US-hfc_female-medium.onnx"
# 是否输出为 MP3 文件（需要 ffmpeg）
mp3: false
# 是否在输出 MP3 时生成章节元数据（需要 EPUB 输入和 ffmpeg）
chapters: false
```
## 注意事项
- Piper 不提供进度输出。由于所有处理都在本地进行，生成长篇有声读物可能需要较长时间。
- 本项目仅在 Linux 环境下测试过。
## 支持项目
感谢您考虑支持本项目。
您可以通过 GitHub 或 PayPal 进行捐赠。如需赞助本项目或出于商业目的联系我，请通过[电子邮件](mailto:github@colton.place)或[我的网站](https://colton.place/contact/)与我联系。
