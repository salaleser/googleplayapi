package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

// MetadataResponse is a vitalina's Application's metadata structure.
type MetadataResponse struct { // TODO add more fields
	Title       string
	Link        string
	AppID       string
	ArtistName  string
	Rating      float32
	ReleaseDate string
	Subtitle    string
	Description string
	Screenshot1 string // TODO add array
	Logo        string
}

func parseIDs(body []byte) []MetadataResponse {
	const errMsg = "[ERR] scraper.parseGpIDs(%s...): %v\n"
	var data1 [][]interface{}
	if err := json.Unmarshal(body, &data1); err != nil {
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return []MetadataResponse{} // TODO handle error
	}

	d := data1[0]

	if d[0] != "wrb.fr" {
		err := fmt.Errorf("the first metadata section element isn't \"wrb.fr\" (%q)", d[0])
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return []MetadataResponse{}
	}

	if d[1] != "lGYRle" {
		err := fmt.Errorf("the second metadata section element isn't \"lGYRle\" (%q)", d[0])
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return []MetadataResponse{}
	}

	if d[2] == nil {
		err := fmt.Errorf("the third metadata section element doesn't exist")
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return []MetadataResponse{}
	}

	var data2 []interface{}
	if err := json.Unmarshal([]byte(d[2].(string)), &data2); err != nil {
		err := fmt.Errorf("unmarshal gp IDs (2): %v", err)
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return []MetadataResponse{}
	}

	i1 := data2[0].([]interface{})
	if i1 == nil {
		i1JSON, _ := json.Marshal(data2[0])
		log.Printf("cast interface 1: %q", errors.New(string(i1JSON)))
		return []MetadataResponse{} // TODO handle error
	}

	i1_1 := i1[1]
	if i1_1 == nil {
		i1_1JSON, _ := json.Marshal(i1)
		log.Printf("cast interface 1.1: %q", errors.New(string(i1_1JSON)))
		return []MetadataResponse{} // TODO handle error
	}

	i2 := i1_1.([]interface{})
	if i2 == nil {
		i2JSON, _ := json.Marshal(i1_1)
		log.Printf("cast interface 2: %q", errors.New(string(i2JSON)))
		return []MetadataResponse{} // TODO handle error
	}

	i3 := i2[0].([]interface{})
	if i3 == nil {
		i3JSON, _ := json.Marshal(i2)
		log.Printf("cast interface 3: %q", errors.New(string(i3JSON)))
		return []MetadataResponse{} // TODO handle error
	}

	i4 := i3[0].([]interface{})
	if i4 == nil {
		i4JSON, _ := json.Marshal(i3)
		log.Printf("cast interface 4: %q", errors.New(string(i4JSON)))
		return []MetadataResponse{} // TODO handle error
	}

	i5 := i4[0].([]interface{})
	if i5 == nil {
		i5JSON, _ := json.Marshal(i4)
		log.Printf("cast interface 5: %q", errors.New(string(i5JSON)))
		return []MetadataResponse{} // TODO handle error
	}

	// FIXME
	if len(i5) < 2 {
		err := fmt.Errorf("len check: < 2")
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return []MetadataResponse{}
	}

	// FIXME interfaces
	metadatas := make([]MetadataResponse, 0)
	for _, d := range i5 {
		metadata := MetadataResponse{
			Title:  d.([]interface{})[2].(string),
			AppID:  d.([]interface{})[12].([]interface{})[0].(string),
			Rating: -1,
		}

		metadatas = append(metadatas, metadata)
	}

	return metadatas
}

func parseMetadata(body []byte) MetadataResponse {
	const errMsg = "[ERR] scraper.parseGpMetadata(%s...): %v\n"
	var data1 [][]interface{}
	if err := json.Unmarshal(body, &data1); err != nil {
		err := fmt.Errorf("unmarshal gp metadata: %v", err)
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return MetadataResponse{}
	}

	d := data1[0]

	if d[0] != "wrb.fr" {
		err := fmt.Errorf("the first metadata section element isn't \"wrb.fr\" (%q)", d[0])
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return MetadataResponse{} // TODO handle error
	}

	if d[1] != "jLZZ2e" {
		err := fmt.Errorf("the second metadata section element isn't \"jLZZ2e\" (%q)", d[0])
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return MetadataResponse{} // TODO handle error
	}

	if d[2] == nil {
		err := fmt.Errorf("the third metadata section element doesn't exist")
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return MetadataResponse{} // TODO handle error
	}

	var data2 [][][]interface{}
	if err := json.Unmarshal([]byte(d[2].(string)), &data2); err != nil {
		err := fmt.Errorf("unmarshal gp metadata: %v", err)
		fmt.Fprintf(os.Stderr, errMsg, body[:10], err)
		return MetadataResponse{} // TODO handle error
	}

	return MetadataResponse{
		// AppID: appID,
		ArtistName: data2[0][12][5].([]interface{})[1].(string),
		// ReleaseDate: data2[0][6][0][1].(float32),
		// Rating:      data2[0][0][0].(string),
		Title:       data2[0][0][0].(string),
		Subtitle:    data2[0][10][1].([]interface{})[1].(string),
		Description: data2[0][10][0].([]interface{})[1].(string),
		Screenshot1: data2[0][12][0].([]interface{})[0].([]interface{})[3].([]interface{})[2].(string),
		Logo:        data2[0][12][1].([]interface{})[3].([]interface{})[2].(string),
	}
}
