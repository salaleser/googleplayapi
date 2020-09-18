package googleplayapi

// MetadataResponse is a vitalina's Application's metadata structure.
type MetadataResponse struct { // TODO add more fields
	Title       string
	Link        string
	AppID       string
	ArtistName  string
	Rating      float32
	StarsCount  int32
	Stars1Count int32
	Stars2Count int32
	Stars3Count int32
	Stars4Count int32
	Stars5Count int32
	ReleaseDate string
	Subtitle    string
	Description string
	Screenshot1 string // TODO add array
	Logo        string
}
