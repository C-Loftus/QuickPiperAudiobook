<h1 style="text-align: center;">QuickPiperAudiobook</h1>

<p align="center">
  <b> English </b> |
  <a href="./README_PL.md">Polski</a>
</p>


Create a natural audiobook for any text content with one command. 

 - Converts [PDFs, epub, txt, mobi, djvu, HTML, docx, and more](https://manual.calibre-ebook.com/generated/en/ebook-convert.html)
 - All conversion is done offline and is entirely private
 - Uses [piper models](https://rhasspy.github.io/piper-samples/); supports many languages

Listen to sample output [ here ](./examples/)

## Installing

1. Grab a prebuilt [release](https://github.com/C-Loftus/QuickPiperAudiobook/releases/)
    * (Or install using `go install github.com/C-Loftus/QuickPiperAudiobook@latest`)
    * (Or build from source using `go mod tidy && go build`)
2. Download `ebook-convert` and make sure it is in your PATH. (This is bundled with [calibre](https://calibre-ebook.com/))
3. _(Optional)_ Download `ffmpeg` for mp3 and chapter support 

> [!NOTE]  
> You don't need to have piper installed. This program manages piper and the associated models

## Usage 

* Pass in either a local file or a remote URL with the proper extension
   * i.e. `./QuickPiperAudiobook test.txt`
* Specify the `--chapters` flag to generate mp3 chapters for epub files
   * i.e. `./QuickPiperAudiobook --chapters test.epub`
* For a full list of options use the `--help` flag
   * i.e. `./QuickPiperAudiobook --help`

### Non-English / UTF-8

* Grab a model for your language of choice (.onnx and .json) from the [piper models](https://rhasspy.github.io/piper-samples/)
  * i.e. `pl_PL-gosia-medium.onnx` and corresponding `pl_PL-gosia-medium.onnx.json` (rename if needed)
* Put them in `~/.config/QuickPiperAudiobook/`
* Use the `--speak-utf-8` and `--model=`  flags to specify you want utf characters to be spoken with a specific model
  * i.e. `./QuickPiperAudiobook --speak-utf-8 --model=pl_PL-gosia-medium.onnx MaszynaTuringa_Wikipedia.pdf`

> [!NOTE]  
> Consider specifying this model as the default in the configuration file if you plan to use it frequently

### Configuring

* You can create a config file at `~/.config/QuickPiperAudiobook/` to specify preferred values if you do not want to specify these as cli args each time
  * i.e. you can use any arbitrary model by putting the associated `.onnx` and `.onnx.json` file for it in `~/.config/QuickPiperAudiobook/`
  * A full example config can be found [here](./examples/config.yaml)

```yml
# An example for `~/.config/QuickPiperAudiobook/config.yaml`

# the default output directory to use if the user does not specify --output in the cli args
output: ~/Audiobooks
# the default model to use if the user does not specify --model in the cli args
model: "en_US-hfc_female-medium.onnx"
# output the audiobook as an mp3 file (requires ffmpeg in your PATH)
mp3: false
# generate chapter metadata when outputting mp3s (requires an epub input and ffmpeg in your PATH)
chapters: false
```

## Notes

- Piper does not support progress output. Long audiobooks may take a long time to generate since all computation is being done locally. 
- This repo has only been tested on Linux.

## Support

Thank you for considering supporting this project.

I accept donations on Github or Paypal. If you would like to sponsor this project or reach out to me for business reasons, you can contact me via [email](mailto:github@colton.place) or [my website](https://colton.place/contact/)