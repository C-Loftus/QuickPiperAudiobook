# QuickPiperAudiobook

Create an audiobook for any text content with one command. 
 - Uses any [piper](https://rhasspy.github.io/piper-samples/) model
    - Manages your piper install and associated models
 - Converts [PDFs, epub, txt, mobi, djvu, HTML, docx, and more](https://manual.calibre-ebook.com/generated/en/ebook-convert.html)
    - Can fetch and convert any of the above from remote URLs
 - All conversion is done offline and is entirely private


## Installing

1. Grab a prebuilt [release](https://github.com/C-Loftus/QuickPiperAudiobook/releases/)
    * (Or build from source using `go mod tidy && go build`)

2. Download `ebook-convert` and make sure it is in your PATH. (This is often bundled with [calibre](https://calibre-ebook.com/))

> [!NOTE]  
> You don't need to have piper installed. This program manages piper and the associated models



## Usage 

* Pass in either a local file or a remote URL to generate an audiobook: 
   * i.e. `./QuickPiperAudiobook test.txt`
* For a full list of options use the `--help` flag
   * i.e. `./QuickPiperAudiobook --help`

### Models and Examples


* An example of the default output can be found in [the examples folder](./examples/)
   * Other pretrained models can be listened to at https://rhasspy.github.io/piper-samples/ 
* This program downloads and manages some of the [best quality models](./lib/models.go) for you
   * However, you can use this repo with any Piper model as you have both the `.onnx` and `.onnx.json` file for it. 


### Configuring

A configuration file at `~/.config/QuickPiperAudiobook/config.yml` will be automatically created. You can place a default model and output path so you do not need to specify these args each time.

```yml
output: ~/Audiobooks
model: "en_US-hfc_female-medium.onnx"
```

## Limitations

To my knowledge, piper does not support progress output. Long books (600+ pages) may take a long time (30 min or more) to generate as audio since all computation is being done locally. 