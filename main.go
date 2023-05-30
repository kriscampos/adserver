package main

import (
	"github.com/kriscampos/adserver/internal/ad_engine"
	"github.com/kriscampos/adserver/internal/router"
)

func main() {
	adEngine := ad_engine.NewAdEngine()
	adEngine.Start()
	defer adEngine.Stop()

	router.SetupRouter(adEngine).Run()
}
