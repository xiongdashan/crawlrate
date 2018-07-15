package models

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/jinzhu/gorm"
	"github.com/otwdev/galaxylib"
)

type Currency struct {
	gorm.Model
	Convertor string  `gorm:"nvarchar(50)"`
	Rate      float64 `gorm:"numberic(30,5)"`
	CrawlTime string
	RawText   string
}

func (c *Currency) FromRemote() {
	address := galaxylib.GalaxyCfgFile.MustValue("data", "url")
	currencyAry := galaxylib.GalaxyCfgFile.MustValueArray("data", "currency", ",")

	currencyBody := ""

	galaxylib.GalaxyDB().OpenDb(func(db *gorm.DB) {

		for _, cy := range currencyAry {
			convertAry := strings.Split(cy, "-")
			param := fmt.Sprintf("From=%s&To=%s", convertAry[0], convertAry[1])
			remote := fmt.Sprintf("%s&%s", address, param)
			regex := regexp.MustCompile(fmt.Sprintf(`1 %s = ([0-9\.]+) %s`, convertAry[0], convertAry[1]))
			convertor := &Currency{
				Convertor: cy,
				CrawlTime: time.Now().Format("2006-01-02"),
			}
			convertor.crawl(remote, regex, db)
			currencyBody = fmt.Sprintf("%s<br/>%s", currencyBody, convertor.RawText)
		}
	})

	n := &Notice{
		Body: currencyBody,
	}

	n.Send()
}

func (c *Currency) crawl(address string, regex *regexp.Regexp, db *gorm.DB) {
	rs, err := http.Get(address)
	if err != nil {
		galaxylib.GalaxyLogger.Error(err)
		return
	}
	doc, _ := goquery.NewDocumentFromResponse(rs)

	section := doc.Find(".uccResultUnit").First()
	c.RawText = section.Text()

	matches := regex.FindAllStringSubmatch(c.RawText, -1)

	c.Rate = galaxylib.DefaultGalaxyConverter.MustFloat(matches[0][1])

	fmt.Println(c.Rate)

	db.Save(c)
}

func (c *Currency) Get() (ret *Currency) {

	var rev []*Currency

	galaxylib.GalaxyDB().OpenDb(func(db *gorm.DB) {
		result := db.Where("convertor=?", c.Convertor).Order("id desc").First(&rev)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
		ret = rev[0]
	})
	return
}
