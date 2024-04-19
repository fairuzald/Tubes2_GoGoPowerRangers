package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"backend/models"

	"github.com/gin-gonic/gin"
)

// AutoCompleteHandler handles autocomplete requests to fetch data from the Wikipedia API
func AutoCompleteHandler(c *gin.Context) {
	// Get search term and limit from query parameter
	searchTerm := c.Query("search")
	limit := c.Query("limit")

	// Create query parameters for the Wikipedia API request
	queryParams := url.Values{}
	queryParams.Set("action", "query")
	queryParams.Set("format", "json")
	queryParams.Set("gpssearch", searchTerm)
	queryParams.Set("generator", "prefixsearch")
	queryParams.Set("prop", "pageprops|pageimages|pageterms|info")
	queryParams.Set("inprop", "url")
	queryParams.Set("redirects", "")
	queryParams.Set("ppprop", "displaytitle")
	queryParams.Set("piprop", "thumbnail")
	queryParams.Set("pilimit", "max")
	// Set the thumbnail maximum size to 160px
	queryParams.Set("pithumbsize", "160")
	queryParams.Set("wbptterms", "description")
	queryParams.Set("gpsnamespace", "0")

	// Set gpslimit only if limit is provided and is not null
	if limit != "" {
		queryParams.Set("gpslimit", limit)
	}

	queryParams.Set("origin", "*")

	// Create the URL for the Wikipedia API request
	apiURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?%s", queryParams.Encode())

	// Send the GET request to Wikipedia API
	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data from Wikipedia"})
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response from Wikipedia API into the AutoComplete struct
	var result models.AutoComplete
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode JSON response"})
		return
	}

	// Handle the case when there are no results
	if len(result.Query.Pages) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": nil, "message": "No results found"})
		return
	}

	// Format and append the results to allResults
	var formattedResults []gin.H

	for _, page := range result.Query.Pages {
		description := ""
		if len(page.Terms.Description) > 0 {
			description = page.Terms.Description[0]
		}

		formattedResult := gin.H{
			"pageid":      page.PageID,
			"title":       page.Title,
			"description": description,
			"image":       page.Thumbnail.Source,
			"url":         page.FullURL,
		}
		formattedResults = append(formattedResults, formattedResult)
	}

	// Return the combined formatted results as JSON
	if len(formattedResults) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": nil, "message": "No results found"})

	} else {
		c.JSON(http.StatusOK, gin.H{"data": formattedResults, "message": "Results retrieved successfully"})
	}
}
