package model

import "encoding/xml"

type DIDLLite struct {
	XMLName    xml.Name    `xml:"DIDL-Lite"`
	Text       string      `xml:",chardata"`
	Dc         string      `xml:"dc,attr"`
	Upnp       string      `xml:"upnp,attr"`
	Xmlns      string      `xml:"xmlns,attr"`
	Containers []Container `xml:"container,omitempty"`
	Items      []Item      `xml:"item,omitempty"`
}

type Container struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id,attr"`
	ParentID    string `xml:"parentID,attr"`
	Restricted  string `xml:"restricted,attr"`
	ChildCount  string `xml:"childCount,attr"`
	Title       string `xml:"title"`
	Class       string `xml:"class"` // one of object.container.storageFolder, object.container.album.musicAlbum
	StorageUsed string `xml:"storageUsed"`
}
type Item struct {
	Text                string `xml:",chardata"`
	ID                  string `xml:"id,attr"`
	ParentID            string `xml:"parentID,attr"`
	Restricted          string `xml:"restricted,attr"`
	RefID               string `xml:"refID,attr"`
	Title               string `xml:"title"`
	Class               string `xml:"class"`
	Description         string `xml:"description"`
	Creator             string `xml:"creator"`
	Date                string `xml:"date"`
	Artist              string `xml:"artist"`
	Album               string `xml:"album"`
	Genre               string `xml:"genre"`
	OriginalTrackNumber string `xml:"originalTrackNumber"`
	Res                 struct {
		Text            string `xml:",chardata"`
		Size            string `xml:"size,attr"`
		Duration        string `xml:"duration,attr"`
		Bitrate         string `xml:"bitrate,attr"`
		SampleFrequency string `xml:"sampleFrequency,attr"`
		NrAudioChannels string `xml:"nrAudioChannels,attr"`
		ProtocolInfo    string `xml:"protocolInfo,attr"`
	} `xml:"res"`
}
