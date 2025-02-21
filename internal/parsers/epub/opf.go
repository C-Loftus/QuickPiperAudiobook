package epub

/*
In the epub spec, the .opf file
defines the spine of the book which
lists the content files in the book
as well as the order in which they
should be read.

Example (abbreviated):

<?xml version='1.0' encoding='utf-8'?>
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="bookid" version="3.0" prefix="rendition: http://www.idpf.org/vocab/rendition/#">

	<metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
	  <meta property="dcterms:modified">2022-01-05T10:44:50Z</meta>
	  <meta name="cover" content="cover-image"/>
	  <dc:title>Writing An Interpreter In Go</dc:title>
	  <dc:creator opf:role="aut" opf:file-as="Ball, Thorsten">Ball, Thorsten</dc:creator>
	  <dc:identifier id="bookid">urn:uuid:273fd756-62f2-4858-8d67-99e08f24bba9</dc:identifier>
	  <dc:identifier opf:scheme="ASIN">B01N2T1VD2</dc:identifier>
	  <dc:contributor opf:file-as="CompanyName" opf:role="own">Epubor</dc:contributor>
	  <dc:contributor opf:file-as="PersonalName" opf:role="own">Ultimate</dc:contributor>
	  <dc:contributor opf:file-as="eCore" opf:role="bkp">eCore v0.9.11.728 [ http://www.epubor.com/ecore.html ]</dc:contributor>
	  <dc:contributor opf:file-as="SiteURL" opf:role="own">http://www.epubor.com</dc:contributor>
	  <dc:date>2016-12-10T16:00:00+00:00</dc:date>
	  <dc:publisher>cj5_7929</dc:publisher>
	  <dc:language>en</dc:language>
	</metadata>
	<manifest>
	  <item href="text00000.html" id="id_1" media-type="application/xhtml+xml"/>
	</manifest>
	<spine toc="ncx">
	  <itemref idref="id_1"/>
	</spine>
	<guide>
	  <reference type="toc" title="Table of contents" href="text00001.html#toc"/>
	</guide>

</package>
*/
type Opf struct {
	Metadata Metadata   `xml:"metadata" json:"metadata"`
	Manifest []Manifest `xml:"manifest>item" json:"manifest"`
	Spine    Spine      `xml:"spine" json:"spine"`
}

// The metadata section of the .opf file
type Metadata struct {
	Title       []string     `xml:"title" json:"title"`
	Language    []string     `xml:"language" json:"language"`
	Identifier  []Identifier `xml:"identifier" json:"identifier"`
	Creator     []Author     `xml:"creator" json:"creator"`
	Subject     []string     `xml:"subject" json:"subject"`
	Description []string     `xml:"description" json:"description"`
	Publisher   []string     `xml:"publisher" json:"publisher"`
	Contributor []Author     `xml:"contributor" json:"contributor"`
	Date        []Date       `xml:"date" json:"date"`
	Type        []string     `xml:"type" json:"type"`
	Format      []string     `xml:"format" json:"format"`
	Source      []string     `xml:"source" json:"source"`
	Relation    []string     `xml:"relation" json:"relation"`
	Coverage    []string     `xml:"coverage" json:"coverage"`
	Rights      []string     `xml:"rights" json:"rights"`
	Meta        []Metafield  `xml:"meta" json:"meta"`
}

// The identifier section within the opf metadata
type Identifier struct {
	Data   string `xml:",chardata" json:"data"`
	ID     string `xml:"id,attr" json:"id"`
	Scheme string `xml:"scheme,attr" json:"scheme"`
}

// The schema for author metadata in the opf metadata
type Author struct {
	Data   string `xml:",chardata" json:"author"`
	FileAs string `xml:"file-as,attr" json:"file_as"`
	Role   string `xml:"role,attr" json:"role"`
}

// The schema for date metadata in the opf metadata
type Date struct {
	Data  string `xml:",chardata" json:"data"`
	Event string `xml:"event,attr" json:"event"`
}

// The schema for metafield metadata in the opf metadata
type Metafield struct {
	Name    string `xml:"name,attr" json:"name"`
	Content string `xml:"content,attr" json:"content"`
}

// The manifest section of the .opf file
type Manifest struct {
	ID           string `xml:"id,attr" json:"id"`
	Href         string `xml:"href,attr" json:"href"`
	MediaType    string `xml:"media-type,attr" json:"type"`
	Fallback     string `xml:"media-fallback,attr" json:"fallback"`
	Properties   string `xml:"properties,attr" json:"properties"`
	MediaOverlay string `xml:"media-overlay,attr" json:"overlay"`
}

// The Spine section of the .opf file.
// Contains all spine items and thus
// is an overview of the book's structure
type Spine struct {
	ID              string      `xml:"id,attr" json:"id"`
	Toc             string      `xml:"toc,attr" json:"toc"`
	PageProgression string      `xml:"page-progression-direction,attr" json:"progression"`
	Items           []SpineItem `xml:"itemref" json:"items"`
}

// SpineItem is a specific item in the spine that essentially
// defines a chapter or logical section of the book
type SpineItem struct {
	// An idref links the spine item to an item in the manifest
	IDref      string `xml:"idref,attr" json:"id_ref"`
	Linear     string `xml:"linear,attr" json:"linear"`
	ID         string `xml:"id,attr" json:"id"`
	Properties string `xml:"properties,attr" json:"properties"`
}
