// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package epub

// Ncx OPS/toc.ncx
type Ncx struct {
	NavPoints []NavPoint `xml:"navMap>navPoint" json:"points"`
}

type NavPoint struct {
	NavLabel  NavLabel `xml:"navLabel" json:"navLabel"`
	Id        string   `xml:"id,attr" json:"id"`
	PlayOrder int      `xml:"playOrder,attr" json:"playOrder"`
}

// NavPoint nav point
type NavLabel struct {
	Text    string  `xml:"navLabel>text" json:"text"`
	Content Content `xml:"content" json:"content"`
}

// Content nav-point content
type Content struct {
	Src string `xml:"src,attr" json:"src"`
}
