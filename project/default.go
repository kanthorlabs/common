package project

import (
	"os"
)

var (
	DefaultRegion    = "sg"
	DefaultNamespace = "kanthor"
	DefaultTier      = "default"
)

func Region() string {
	region := os.Getenv("PROJECT_REGION")
	if region != "" {
		return region
	}
	return DefaultRegion
}

func Namespace() string {
	ns := os.Getenv("PROJECT_NAMESPACE")
	if ns != "" {
		return ns
	}
	return DefaultNamespace
}

func Tier() string {
	tier := os.Getenv("PROJECT_TIER")
	if tier != "" {
		return tier
	}
	return DefaultTier
}
