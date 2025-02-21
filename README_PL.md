<h1 align=center>QuickPiperAudiobook</h1>

<p align="center">
  <a href="./README.md">English</a> |
  <b> Polski </b>
</p>

Stwórz naturalny audiobook dla dowolnej treści tekstowej za pomocą jednego polecenia.

- Konwertuje [PDF, epub, txt, mobi, djvu, HTML, docx i inne](https://manual.calibre-ebook.com/generated/en/ebook-convert.html)
- Cała konwersja odbywa się offline, na twoim komputerze i jest całkowicie prywatna
- Wykorzystuje [modele Piper](https://rhasspy.github.io/piper-samples/); obsługuje wiele języków w tym polski

Posłuchaj przykładowego nagrania [tutaj](./examples/)

## Instalacja

1. Pobierz wstępnie zbudowane [wydanie](https://github.com/C-Loftus/QuickPiperAudiobook/releases/)
   * (lub zbuduj ze źródła za pomocą `go mod tidy && go build`)

2. Pobierz `ebook-convert` i upewnij się, że znajduje się w twojej ścieżce PATH. (Często jest on dołączony do [Calibre](https://calibre-ebook.com/))

> [!NOTE]
> Nie musisz mieć zainstalowanego Piper. Ten program zarządza Piper i powiązanymi modelami.

## Użycie

* Podaj lokalny plik lub zdalny URL, aby wygenerować audiobook:
  * np. `./QuickPiperAudiobook test.txt`
* Aby uzyskać pełną listę opcji, użyj flagi `--help`
  * np. `./QuickPiperAudiobook --help`

### Języki inne niż angielski oraz obsługa kodowania UTF-8

* Pobierz model dla wybranego języka (.onnx i .json) z [piper models](https://rhasspy.github.io/piper-samples/)
  * np. `pl_PL-gosia-medium.onnx` i odpowiadający mu `pl_PL-gosia-medium.onnx.json` (zmień nazwę, jeśli to konieczne)
* Umieść je w katalogu `~/.config/QuickPiperAudiobook/`
* Użyj flag `--speak-utf-8` i `--model=`
  * np. `./QuickPiperAudiobook --speak-utf-8 --model=pl_PL-gosia-medium.onnx MaszynaTuringa_Wikipedia.pdf`

> [!NOTE]
> Pomyśl o dodaniu wybranego, podstawowego modelu do pliku konfiguracyjnego jeżeli zamierzasz go często używać

### Konfiguracja

* Plik konfiguracyjny `~/.config/QuickPiperAudiobook/config.yml` zostanie automatycznie utworzony.
* Możesz określić domyślny model oraz ścieżkę wyjściową, aby nie musieć za każdym razem podawać tych argumentów.
* Możesz używać dowolnego modelu, umieszczając powiązane pliki `.onnx` i `.onnx.json` w katalogu `~/.config/QuickPiperAudiobook/`.

```yml
# Przykład dla `~/.config/QuickPiperAudiobook/config.yml`

# domyślny katalog wyjściowy, jeśli użytkownik nie określi --output w argumentach CLI
output: ~/Audiobooks
# domyślny model, jeśli użytkownik nie określi --model w argumentach CLI
model: "en_US-hfc_female-medium.onnx"
