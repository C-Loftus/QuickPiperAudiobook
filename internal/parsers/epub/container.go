// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package epub

// Representation of the container.xml at the root of the book
// This file defines the version of xml used, the encoding
// and links to the full-path of the content.opf which defines
// the actual structure of the book's internal contents
/*
Example:

<?xml version='1.0' encoding='utf-8'?>
<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">
  <rootfiles>
    <rootfile media-type="application/oebps-package+xml" full-path="EPUB/content.opf"/>
  </rootfiles>
</container>
*/
type Container struct {
	Rootfile Rootfile `xml:"rootfiles>rootfile" json:"rootfile"`
}

// Rootfile root file
type Rootfile struct {
	Path string `xml:"full-path,attr" json:"path"`
	Type string `xml:"media-type,attr" json:"type"`
}
