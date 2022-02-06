package util

import (
	"net/url"
)

func GetSchemeFromUrl(urlString string) string {
	url, err := url.Parse(urlString)

	if err != nil {
		return ""
	}

	return url.Scheme
}

func GetUrlFromString(urlString string) (*url.URL, error) {
	return url.Parse(urlString)
}