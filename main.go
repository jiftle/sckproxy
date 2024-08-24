package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/jiftle/sckproxy/version"

	"github.com/jiftle/sckproxy/cmd"
)

func main() {
	g.Log().Infof(context.Background(), "Name: %v, Version: %v, BuildTime: %v, starting...", version.Name, version.Version, version.BuildTime)
	cmd.Execute()
}
