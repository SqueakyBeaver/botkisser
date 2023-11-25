package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetPrideFlag(name string) (string, error) {
	// Return link to flag file

	// Get page name
	// https://www.lgbtqia.wiki/1.40/api.php?action=query&format=json&list=search&meta=&formatversion=2&srsearch=Demiboy&srwhat=nearmatch

	// Get flag url
	// https://www.lgbtqia.wiki/1.40/api.php?action=query&format=json&prop=pageimages%7Cimages&list=&meta=&titles=Demiboy&formatversion=latest&piprop=original
	fmt.Printf("Getting flag for: %v\n", name)
	pageTitle := fetchPageTitle(name, "new")
	if pageTitle != "" {

		flagUrl := fetchFlagUrl(pageTitle, "new")

		if flagUrl != "" {
			return flagUrl, nil
		}
	}
	// assume we need to fetch with the old api
	fmt.Printf("Using old api for query: %s\n", name)
	pageTitle = fetchPageTitle(name, "old")

	flagUrl := fetchFlagUrl(pageTitle, "old")
	if flagUrl != "" {
		return flagUrl, nil
	} else {
		return "", fmt.Errorf("no results found for \"%s\"", name)
	}
}

func fetchPageTitle(name string, version string) string {
	var vstr string
	var prefix string
	if version == "new" {
		vstr = "w"
		prefix = "new"
	} else {
		vstr = "1.40"
		prefix = "www"
	}
	respUrl := fmt.Sprintf("https://%s.lgbtqia.wiki/%s/api.php?action=query&format=json&list=search&meta=&formatversion=2&srsearch=%s&redirects=1&srwhat=nearmatch", prefix, vstr, name)

	searchResp, err := http.Get(respUrl)
	if err != nil {
		log.Panicf("Error fetching pride search: %v", err)
	}

	defer searchResp.Body.Close()

	searchBody, err := io.ReadAll(searchResp.Body)
	if err != nil {
		log.Panicf("Error parsing pride search response: %v", err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(searchBody, &data)
	if err != nil {
		log.Panicf("Error reading pride search response json: %v\n%v", err, data)
	}

	// I have to assert the type each time
	// Same thing as data["query"]["search"][0]
	searchData := fmt.Sprint(data["query"])
	if strings.Contains(searchData, "title") {
		return fmt.Sprint(data["query"].(map[string]interface{})["search"].([]interface{})[0].(map[string]interface{})["title"])
	}
	return ""
}

func fetchFlagUrl(pageTitle string, version string) string {
	if pageTitle == "" {
		return ""
	}

	var vstr string
	var prefix string
	if version == "new" {
		vstr = "w"
		prefix = "new"
	} else {
		vstr = "1.40"
		prefix = "www"
	}

	respUrl := fmt.Sprintf("https://%s.lgbtqia.wiki/%s/api.php?action=query&format=json&prop=pageimages|images&list=&meta=&titles=%s&redirects=1&formatversion=latest&piprop=original", prefix, vstr, pageTitle)

	flagResp, err := http.Get(respUrl)
	if err != nil {
		log.Panicf("Error fetching pride flag: %v", err)
	}

	defer flagResp.Body.Close()

	flagBody, err := io.ReadAll(flagResp.Body)
	if err != nil {
		log.Panicf("Error parsing flag response: %v", err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(flagBody, &data)
	if err != nil {
		log.Panicf("Error reading flag response json: %v\n%v", err, data)
	}

	pageData := fmt.Sprint(data["query"].(map[string]interface{})["pages"].([]interface{})[0])
	if strings.Contains(pageData, "original") {
		return fmt.Sprint(data["query"].(map[string]interface{})["pages"].([]interface{})[0].(map[string]interface{})["original"].(map[string]interface{})["source"])
	}

	return ""
}

func main() {
	println(GetPrideFlag("trans"))
	println(GetPrideFlag("bi"))
	println(GetPrideFlag("genderqueer"))
	println(GetPrideFlag("demigirl"))
}
