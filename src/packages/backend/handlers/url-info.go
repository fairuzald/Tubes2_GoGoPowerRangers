package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/models"
)

// URLInfoHandler adalah handler untuk mendapatkan informasi dari URL menggunakan API Wikipedia
func URLInfoHandler(c *gin.Context) {

	// Bind the JSON request body to URLInfoBody struct
	var requestBody models.URLInfoBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON format"})
		return
	}

	// Parse the URL from the request body
	parsedURL, err := url.Parse(requestBody.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URL format"})
		return
	}

	// Unescape special characters in the URL
	unescapedURL, err := url.PathUnescape(parsedURL.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to unescape URL"})
		return
	}

	// Get the page title from the unescaped URL path
	pathParts := strings.Split(unescapedURL, "/")
	pageTitle := pathParts[len(pathParts)-1]

	// Create URL parameters for the Wikipedia API request
	queryParams := url.Values{}
	queryParams.Set("action", "query")
	queryParams.Set("format", "json")
	queryParams.Set("titles", pageTitle)
	queryParams.Set("prop", "pageimages|pageterms")
	queryParams.Set("ppprop", "displaytitle")
	queryParams.Set("piprop", "thumbnail")
	// Set the thumbnail maximum size to 160px
	queryParams.Set("pithumbsize", "160")
	queryParams.Set("wbptterms", "description")

	// Create the URL for the Wikipedia API request
	apiURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?%s", queryParams.Encode())

	// Send the GET request to Wikipedia API
	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data from Wikipedia"})
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response into the struct
	var wikiResponse models.WikipediaResponseURLInfo
	if err := json.NewDecoder(resp.Body).Decode(&wikiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode Wikipedia API response"})
		return
	}

	// Get the first search result from the pages map
	var searchResult models.WikipediaSearchResultURLInfo
	for _, result := range wikiResponse.Query.Pages {
		searchResult = result
		break // Stop after getting the first search result
	}

	// Check if the search result is not empty
	if searchResult.Title != "" {
		formattedResults := gin.H{
			"title": searchResult.Title,
			"url":   requestBody.URL,
			"description": func() string {
				if len(searchResult.Terms.Description) > 0 {
					return searchResult.Terms.Description[0]
				}
				return ""
			}(),

			"image":  searchResult.Image.URL,
			"pageid": searchResult.PageID,
		}

		c.JSON(http.StatusOK, gin.H{"data": formattedResults, "message": "URL information retrieved successfully"})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": nil, "message": "No information found for the provided URL"})
	}
}
