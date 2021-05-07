package main

import ("fmt"

"github.com/djimenez/iconv-go"
"github.com/gocolly/colly"
"github.com/gocolly/colly/debug"
)





func getP() {
	var preUrl string = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/"
	var provinceHref string

	// http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"), colly.MaxDepth(6), colly.Debugger(&debug.LogDebugger{}))
	//省级列表
	c.OnHTML("table[class='provincetable'] > tbody", func(e *colly.HTMLElement) {
		//遍历每一行
		e.ForEach("tr[class='provincetr'] > td", func(i int, item *colly.HTMLElement) {
			// 省份
			//text := item.ChildText("a")
			href := item.ChildAttr("a", "href")
			fmt.Println(href)
			provinceHref = preUrl + href
			{
				c.OnHTML("table[class='citytable'] > tbody", func(e *colly.HTMLElement) {

					//遍历每一行
					e.ForEach("tr[class='citytr'] > td", func(i int, item *colly.HTMLElement) {
						// 城市名称
						cityName := item.ChildText("a")

						fmt.Println(cityName)

						//href := item.ChildAttr("a", "href")
						//href = preUrl + href
						//fmt.Println(href)
						//c.Visit(href)

					})
				})

				err := c.Visit(provinceHref)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			//fmt.Println(href)
			//c.Visit(href)

		})
	})

	err := c.Visit("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

}
