package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var debug = false

// AppIDs returns application IDs by a keyword.
func AppIDs(keyword string, gl string, hl string) []MetadataResponse {
	const errMsg = "[ERR] scraper.GpAppIDs(%s,%s,%s): %v\n"
	const baseURL = "https://play.google.com/_/PlayStoreUi/data/batchexecute"
	uri, err := url.Parse(baseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, keyword, gl, hl, err)
		return []MetadataResponse{}
	}

	query, err := url.ParseQuery(uri.RawQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, keyword, gl, hl, err)
		return []MetadataResponse{}
	}
	query.Add("gl", gl)
	query.Add("hl", hl)
	uri.RawQuery = query.Encode()

	value := fmt.Sprintf("[[[lGYRle,'[[null,[[10,[10,%d]],true,null,[96,27,4,8,57,30,110,79,11,16,49,1,3,9,12,104,55,56,51,10,34,77]],[%s],4,null,null,null,[]]]',null,%s]]]", 5, keyword, keyword)

	data := url.Values{}
	data.Add("f.req", value)

	resp, err := http.PostForm(baseURL, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, keyword, gl, hl, err)
		return []MetadataResponse{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, keyword, gl, hl, err)
		return []MetadataResponse{}
	}

	return parseIDs(body[5:])
}

// App returns an Application's metadata by ID packageName.
func App(packageName string, gl string, hl string) MetadataResponse {
	const errMsg = "[ERR] scraper.AppGp(%s,%s,%s): %v\n"
	const baseURL = "https://play.google.com/_/PlayStoreUi/data/batchexecute"
	uri, err := url.Parse(baseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, packageName, gl, hl, err)
		return MetadataResponse{}
	}

	query, err := url.ParseQuery(uri.RawQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, packageName, gl, hl, err)
		return MetadataResponse{}
	}
	query.Add("gl", gl)
	query.Add("hl", hl)
	uri.RawQuery = query.Encode()

	// v1 := fmt.Sprintf("[d5UeYe,'[[%s,7]]',null,%s]", appID, appID)        // price
	// v3 := fmt.Sprintf("[MLWfjd,'[[%s,7]]',null,%s]", appID, appID)        // rating
	// v4 := fmt.Sprintf("[IoIWBc,'[[null,[%s,7]]]',null,%s]", appID, appID) // version
	// v5 := fmt.Sprintf("[k8610b,'[[null,[%s,7]]]',null,%s]", appID, appID) // short rating
	// v6 := fmt.Sprintf("[BQ0Ly,'[[null,[%s,7]]]',null,%s]", appID, appID)  // downloads

	value := fmt.Sprintf("[[[jLZZ2e,'[[%s,7],2]',null,%s]]]",
		packageName, packageName)

	data := url.Values{}
	data.Add("f.req", value)

	resp, err := http.PostForm(baseURL, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, packageName, gl, hl, err)
		return MetadataResponse{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg, packageName, gl, hl, err)
		return MetadataResponse{}
	}

	return parseMetadata(body[5:])
}
