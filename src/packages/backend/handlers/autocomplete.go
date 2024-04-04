package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"backend/models"

	"github.com/gin-gonic/gin"
)

func AutoCompleteHandler(c *gin.Context) {
	// Get search term and limit from query parameter
	searchTerm := c.Query("search")
	limit := c.Query("limit")

	// Create query parameters to send API request to Wikipedia
	queryParams := url.Values{}
	queryParams.Set("action", "query")
	queryParams.Set("format", "json")
	queryParams.Set("gpssearch", searchTerm)
	queryParams.Set("generator", "prefixsearch")
	queryParams.Set("prop", "pageprops|pageimages|pageterms")
	queryParams.Set("redirects", "")
	queryParams.Set("ppprop", "displaytitle")
	queryParams.Set("piprop", "thumbnail")
	// Set the thumbnail maximum size to 160px
	queryParams.Set("pithumbsize", "160")
	queryParams.Set("wbptterms", "description")
	queryParams.Set("gpsnamespace", "0")

	// Set gpslimit only if limit is provided and is not null
	if limit != "" {
		// Set gpslimit only
		queryParams.Set("gpslimit", limit)
	}

	queryParams.Set("origin", "*")

	// Create URL to Wikipedia API and add the params
	url := fmt.Sprintf("https://en.wikipedia.org/w/api.php?%s", queryParams.Encode())

	// Shoot the request to Wikipedia API
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data from Wikipedia"})
		return
	}
	defer resp.Body.Close()

	// Decode JSON response from Wikipedia API based on model
	var result models.AutoComplete
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode JSON response"})
		return
	}

	// Validation if there is no result
	if len(result.Query.Pages) == 0 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	// Format the result to be more readable
	var formattedResults []gin.H
	for _, page := range result.Query.Pages {
		// Make sure the description is not empty and only take the first one
		description := ""
		if len(page.Terms.Description) > 0 {
			description = page.Terms.Description[0]
		}

		// Format the result
		formattedResult := gin.H{
			"pageid":      page.PageID,
			"title":       page.Title,
			"description": description,
			"image": gin.H{
				"url":    page.Thumbnail.Source,
				"width":  page.Thumbnail.Width,
				"height": page.Thumbnail.Height,
			},
		}
		formattedResults = append(formattedResults, formattedResult)
	}

	c.JSON(http.StatusOK, formattedResults)
}
