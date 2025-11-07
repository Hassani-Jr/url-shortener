package validator

import (
	"net/url"
	"strings"

	"github.com/Hassani-Jr/url-shortener/pkg/logger/apperror"
)

const MaxURLLength = 2000

func ValidateURL(rawURL string) (string, error){
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == ""{
		return "",apperror.BadRequest("URL is required",nil)
	}

	if len(rawURL) > MaxURLLength {
		return "", apperror.BadRequest("URL is too long",nil)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "",apperror.BadRequest("Invalid URL format",nil)
	}

	if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https"){
		return  "", apperror.BadRequest("URL must use http or https",nil)
	}

	if parsedURL.Host == ""{
		return "", apperror.BadRequest("URL must have a valid host", nil)
	}

	return rawURL,nil
}