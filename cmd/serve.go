package cmd

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"r1wallet/config"
	"r1wallet/handler"
	"r1wallet/repositories"
	"r1wallet/routes"
	"r1wallet/services"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve r1 wallet application",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCMD.AddCommand(serveCmd)
}

func serve() {
	ca := config.InitializeConfig()
	rep := repositories.NewRepository(ca.DB, ca.RDB)
	ser := services.NewServices(rep)
	hndl := handler.NewBaseHandler(ser)
	go initializeStreamServer(ser, ca, hndl)
	initializeHttpServer(hndl)
}

func initializeHttpServer(handler *handler.BaseHandler) {
	e := echo.New()
	e.HideBanner = true
	p := prometheus.NewPrometheus("r1wallet", nil)
	p.Use(e)
	routes.RegisterRoutes(e, handler)
	e.Logger.Fatal(e.Start(":1323"))
}

func initializeStreamServer(service *services.Services, config *config.ConfiguredApp, handler *handler.BaseHandler) {
	ch := make(chan string)
	go service.Consumer.Consume(ch, config.Config.App.IncreaseChannelName)
	handler.Credit.HandleIncreaseRequestFromChannel(ch)
}
