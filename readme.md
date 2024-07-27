# QuickPiperAudiobook

Create an audiobook for any text content with one command. 
 - Uses any [piper](https://rhasspy.github.io/piper-samples/) model
 - Generates all audio offline
 - Can take in URLs or local files which contain valid PDF, txt, mobi or epub



## Installing

Grab a prebuilt [release](https://github.com/C-Loftus/QuickPiperAudiobook/releases/)

Build from source using 

```sh
go mod tidy && go build
```

## Running 

```
./QuickPiperAudiobook --help
```

## Dependencies

- `ebook-convert` needs to be in your PATH. (This is often bundled with [calibre](https://calibre-ebook.com/))

> [!NOTE]  
> You don't need to have piper installed. This program manages piper and the associated models

## Limitations

To my knowledge, piper does not support progress output. Long book (600+ pages) may take a long time (30 min or more) to generate as audio since all computation is being done locally. 