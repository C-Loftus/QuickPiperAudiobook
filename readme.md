# QuickPiperAudiobook

Create an audiobook for any text content with one command. 
 - Uses any [piper](https://rhasspy.github.io/piper-samples/) model
 - Manages your piper install and associated models
 - Can convert [PDFs, epub, txt, mobi, djvu, HTML, docx, and more](https://manual.calibre-ebook.com/generated/en/ebook-convert.html)
 - Can fetch content from remote URLs
 - All processing on your text file is done offline and is entirely private



## Installing

1. Grab a prebuilt [release](https://github.com/C-Loftus/QuickPiperAudiobook/releases/)
    * (Or build from source using `go mod tidy && go build`)

2. Download `ebook-convert` and make sure it is in your PATH. (This is often bundled with [calibre](https://calibre-ebook.com/))

> [!NOTE]  
> You don't need to have piper installed. This program manages piper and the associated models



## Running 

```
./QuickPiperAudiobook --help
```

## Limitations

To my knowledge, piper does not support progress output. Long books (600+ pages) may take a long time (30 min or more) to generate as audio since all computation is being done locally. 