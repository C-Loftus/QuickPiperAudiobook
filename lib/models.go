package lib

// Piper has hundreds of pretrained models on the sample Website
// These are some of the best ones for English. However, as long
// as you have both the .onnx and .onnx.json files locally, you
// can use any model you want or even train your own.
var ModelToURL = map[string]string{
	"en_US-hfc_male-medium.onnx":              "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/hfc_male/medium/en_US-hfc_male-medium.onnx",
	"en_US-hfc_female-medium.onnx":            "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/hfc_female/medium/en_US-hfc_female-medium.onnx",
	"en_US-lessac-medium.onnx":                "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx",
	"en_GB-northern_english_male-medium.onnx": "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_GB/northern_english_male/medium/en_GB-northern_english_male-medium.onnx",
}
