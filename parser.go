package googleplayapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func parseIDs(body []byte) []MetadataResponse {
	const errMsgFormat = "[ERR] googleplayapi.parseIDs(%s...): %v\n"
	var data1 [][]interface{}
	if err := json.Unmarshal(body, &data1); err != nil {
		fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
		return []MetadataResponse{} // TODO handle error
	}

	d := data1[0]

	if d[0] != "wrb.fr" {
		err := fmt.Errorf("the first metadata section element isn't \"wrb.fr\" (%q)", d[0])
		fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
		return []MetadataResponse{}
	}

	if d[1] != "lGYRle" {
		err := fmt.Errorf("the second metadata section element isn't \"lGYRle\" (%q)", d[0])
		fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
		return []MetadataResponse{}
	}

	if d[2] == nil {
		err := fmt.Errorf("the third metadata section element doesn't exist")
		fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
		return []MetadataResponse{}
	}

	var data2 []interface{}
	if err := json.Unmarshal([]byte(d[2].(string)), &data2); err != nil {
		err := fmt.Errorf("unmarshal gp IDs (2): %v", err)
		fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
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
		fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
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
	const errMsgFormat = "googleplayapi.parseMetadata(%s...): %v\n"
	var data1 [][]interface{}
	if err := json.Unmarshal(body, &data1); err != nil {
		err := fmt.Errorf("unmarshal gp metadata: %v", err)
		fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
		return MetadataResponse{}
	}

	var metadataData [][][]interface{}
	var ratingData [][]interface{}
	for _, d := range data1 {
		if d[0].(string) != "wrb.fr" {
			continue
		}

		if d[2] == nil {
			err := fmt.Errorf("the third section element doesn't exist")
			fmt.Fprintf(os.Stderr, "[ERR] %s", err.Error())
			return MetadataResponse{}
		}

		switch d[1].(string) {
		case "jLZZ2e":
			metadataData = processMetadataSection(d[2])
		case "MLWfjd":
			ratingData = processRatingaSection(d[2])
		default:
			err := fmt.Errorf("the second section element not supported (%q)", d[0])
			fmt.Fprintf(os.Stderr, errMsgFormat, body[:10], err)
			return MetadataResponse{} // TODO handle error
		}
	}

	rating := ratingData[0][6].([]interface{})[0].([]interface{})[1].(float64)
	starsCount := ratingData[0][6].([]interface{})[2].([]interface{})[1].(float64)
	stars1Count := ratingData[0][6].([]interface{})[1].([]interface{})[1].([]interface{})[1].(float64)
	stars2Count := ratingData[0][6].([]interface{})[1].([]interface{})[2].([]interface{})[1].(float64)
	stars3Count := ratingData[0][6].([]interface{})[1].([]interface{})[3].([]interface{})[1].(float64)
	stars4Count := ratingData[0][6].([]interface{})[1].([]interface{})[4].([]interface{})[1].(float64)
	stars5Count := ratingData[0][6].([]interface{})[1].([]interface{})[5].([]interface{})[1].(float64)

	return MetadataResponse{
		ArtistName: metadataData[0][12][5].([]interface{})[1].(string),
		// ReleaseDate: data2[0][6][0][1].(float32),
		Rating:      float32(rating),
		StarsCount:  int32(starsCount),
		Stars1Count: int32(stars1Count),
		Stars2Count: int32(stars2Count),
		Stars3Count: int32(stars3Count),
		Stars4Count: int32(stars4Count),
		Stars5Count: int32(stars5Count),
		Title:       metadataData[0][0][0].(string),
		Subtitle:    metadataData[0][10][1].([]interface{})[1].(string),
		Description: metadataData[0][10][0].([]interface{})[1].(string),
		Screenshot1: metadataData[0][12][0].([]interface{})[0].([]interface{})[3].([]interface{})[2].(string),
		Logo:        metadataData[0][12][1].([]interface{})[3].([]interface{})[2].(string),
	}
}

func processMetadataSection(d interface{}) [][][]interface{} {
	var data [][][]interface{}
	if err := json.Unmarshal([]byte(d.(string)), &data); err != nil {
		err := fmt.Errorf("unmarshal gp metadata: %v", err)
		fmt.Fprintf(os.Stderr, "[ERR] %s", err.Error())
		return nil // TODO handle error
	}

	return data
}

func processRatingaSection(d interface{}) [][]interface{} {
	var data [][]interface{}
	if err := json.Unmarshal([]byte(d.(string)), &data); err != nil {
		err := fmt.Errorf("unmarshal gp metadata: %v", err)
		fmt.Fprintf(os.Stderr, "[ERR] %s", err.Error())
		return nil // TODO handle error
	}

	return data
}
