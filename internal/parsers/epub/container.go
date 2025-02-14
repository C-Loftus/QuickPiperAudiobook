package epub

// Container META-INF/container.xml file
type Container struct {
	Rootfile Rootfile `xml:"rootfiles>rootfile" json:"rootfile"`
}

// Rootfile root file
type Rootfile struct {
	Path string `xml:"full-path,attr" json:"path"`
	Type string `xml:"media-type,attr" json:"type"`
}
