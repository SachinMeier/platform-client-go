package main

import (
	"context"
	"fmt"

	platform "github.com/SachinMeier/platform-client-go/platform"
)

func main() {
	ctx := context.Background()
	url, account_id, api_key, err := platform.LoadEnv()
	if err != nil {
		fmt.Printf("ENV failed to load. Missing variables: %s", err.Error())
	}
	pc := platform.NewPlatformClient(
		ctx,
		url,
		account_id,
		api_key,
	)
	Getaccount(pc)
}
