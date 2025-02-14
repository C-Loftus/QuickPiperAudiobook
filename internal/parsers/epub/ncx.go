package epub

//Ncx OPS/toc.ncx
type Ncx struct {
	Points []NavPoint `xml:"navMap>navPoint" json:"points"`
}

//NavPoint nav point
type NavPoint struct {
	Text    string     `xml:"navLabel>text" json:"text"`
	Content Content    `xml:"content" json:"content"`
	Points  []NavPoint `xml:"navPoint" json:"points"`
}

//Content nav-point content
type Content struct {
	Src string `xml:"src,attr" json:"src"`
}
