package main

import (
	"net/http"

	"github.com/jasonlvhit/gocron"

	"github.com/labstack/echo"
	"github.com/otwdev/crawlrate/models"
	"github.com/otwdev/galaxylib"
)

func main() {

	galaxylib.DefaultGalaxyConfig.InitConfig()
	galaxylib.DefaultGalaxyLog.ConfigLogger()

	crawlData()

	e := echo.New()

	e.GET("/currency/:convert", func(c echo.Context) error {
		convertor := c.Param("convert")
		currency := &models.Currency{
			Convertor: convertor,
		}
		ret := currency.Get()

		data := &struct {
			Convert string
			Rate    float64
			Date    string
		}{
			Convert: ret.Convertor,
			Rate:    ret.Rate,
			Date:    ret.CrawlTime,
		}

		return c.JSON(http.StatusOK, data)
	})

	port := galaxylib.GalaxyCfgFile.MustValue("data", "port")

	e.Start(port)

}

func crawlData() {

	crawlTime := galaxylib.GalaxyCfgFile.MustValue("data", "crawlTime")

	go func() {

		gocron.Every(1).Days().At(crawlTime).Do(func() {
			c := &models.Currency{}

			c.FromRemote()
		})

		<-gocron.Start()
	}()
}
