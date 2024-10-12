# QuickPiperAudiobook

Create a natural audiobook for any text content with one command. 

 - Converts [PDFs, epub, txt, mobi, djvu, HTML, docx, and more](https://manual.calibre-ebook.com/generated/en/ebook-convert.html)
 - All conversion is done offline and is entirely private
 - Uses [piper models](https://rhasspy.github.io/piper-samples/); supports many languages

Listen to sample output [ here ](./examples/)

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


### Configuring

* A configuration file at `~/.config/QuickPiperAudiobook/config.yml` will be automatically created. 
* You can specify a default model and output path so you do not need to specify these args each time.
* You can use any arbitrary model by putting the associated `.onnx` and `.onnx.json` file for it in `~/.config/QuickPiperAudiobook/`


```yml
# An example for `~/.config/QuickPiperAudiobook/config.yml`

# the default output directory to use if the user does not specify --output in the cli args
output: ~/Audiobooks
# the default model to use if the user does not specify --model in the cli args
model: "en_US-hfc_female-medium.onnx"
```

## Notes

Piper does not support progress output. Long audiobooks may take a long time to generate since all computation is being done locally. 

This repo has only been tested on Linux.
