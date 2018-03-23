package config

import (
	"errors"
	"fmt"
	"net/url"
	"notif2/tools"
	"os"
	"strconv"
	"strings"
)

func StringOption(key string) string {
	val, present := os.LookupEnv(key)
	if !present {
		tools.FailOnError(
			fmt.Sprintf("Configuration key: %s not found", key),
			errors.New("Missing configuration item"),
		)
	}
	return val
}

func IntOption(key string) int {
	valStr := StringOption(key)
	valInt, err := strconv.ParseInt(valStr, 10, 32)
	if err != nil {
		tools.FailOnError(
			fmt.Sprintf(
				"Configuration key: %s value (%s) is not integer",
				key,
				valStr,
			),
			err,
		)
	}
	return int(valInt)
}

func URLOption(key string) *url.URL {
	valStr := StringOption(key)
	res, err := url.Parse(valStr)
	if err != nil {
		tools.FailOnError(
			fmt.Sprintf(
				"Configuration key: %s value (%s) is not a URL",
				key,
				valStr,
			),
			err,
		)
	}
	return res
}

func GetHttpPort(url *url.URL) (int, error) {
	parts := strings.Split(url.Host, ":")
	var validSchemas = map[string]bool{
		"http":  true,
		"https": true,
	}
	if !validSchemas[url.Scheme] {
		return 0, errors.New(fmt.Sprintf("Invalid schema %s", url.Scheme))
	}
	if len(parts) == 1 {
		if url.Scheme == "http" {
			return 80, nil
		} else if url.Scheme == "https" {
			return 443, nil
		} else {
			return 0, errors.New("Not an HTTP URL")
		}
	} else if len(parts) == 2 {
		port, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return 0, err
		}
		return int(port), nil
	} else {
		return 0, errors.New("Bad host part")
	}
}
