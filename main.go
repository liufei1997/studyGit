package main

import (
	"errors"
	"fmt"
	"git.in.codoon.com/backend/common/clog"
	"git.in.codoon.com/backend/system_service/models"
	"git.in.codoon.com/third/pinyin/pinyin"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type PublishRecord struct {
	//发布日期
	Date string `json:"date"`
	//链接
	Link string `json:"link"`
}

type Province struct {
	Code   int    `json:"code"`
	Name   string `json:"name"`
	Link   string `json:"-"`
	Cities []City `json:"cities"`
}

type City struct {
	Code     int      `json:"code"`
	Name     string   `json:"name"`
	Link     string   `json:"-"`
	Counties []County `json:"counties"`
}

type County struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Link string `json:"-"`
}

type PROVINCE_CITY_REGION_MODEL struct {
	ID                int    `gorm:"column:id" sql:"type:int(11)" json:"id"`
	ProvinceCode      int    `gorm:"column:province_code" sql:"type:int(11)" json:"province_code"`
	ProvinceName      string `gorm:"column:province_name" sql:"type:varchar(128)" json:"province_name"`
	ProvinceNamePy    string `gorm:"column:province_name_py" sql:"type:varchar(128)" json:"province_name_py"`
	CityCode          int    `gorm:"column:city_code" sql:"type:int(11)" json:"city_code"`
	CityName          string `gorm:"column:city_name" sql:"type:varchar(128)" json:"city_name"`
	CityNamePy        string `gorm:"column:city_name_py" sql:"type:varchar(128)" json:"city_name_py"`
	RegionCode        int    `gorm:"column:region_code" sql:"type:int(11)" json:"region_code"`
	RegionName        string `gorm:"column:region_name" sql:"type:varchar(128)" json:"region_name"`
	RegionNamePy      string `gorm:"column:region_name_py" sql:"type:varchar(128)" json:"region_name_py"`
	CityCodeTelephone string `gorm:"column:city_code_telephone" sql:"type:varchar(8)" json:"city_code_telephone"`
	Area              string `gorm:"column:area" sql:"type:varchar(64)" json:"area"`
}

func main() {
	//data, _ := GetProvinceUrlAndData("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/")
	//fmt.Println(data)
	//
	//res, _ := json.Marshal(data)
	//WriteWithIoutil("全部数据",res)

	//publishRecords, _ := GetPublishRecordV2()
	//fmt.Println(publishRecords)
	//publishRecordsData,_ :=json.Marshal(publishRecords)
	//WriteWithIoutil("所有年份链接及发布日期",publishRecordsData)

	//data, _ := GetProvinceUrlAndData("http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/")
	//fmt.Println(data)
	HasLetter("11")


}

func HasLetter(s string) {
	reg := regexp.MustCompile(`[a-zA-Z]`)
	matchString := reg.MatchString(s)
	fmt.Println(matchString)

}

// 将数据写入文件
func WriteWithIoutil(fileName string, data []byte) {
	if ioutil.WriteFile(fileName, data, 0644) == nil {
		fmt.Println("写入文件成功:")
	}
}

////UpdateAmapData
//func UpdateAmapDataV2(c *gin.Context) {
//	clog.Logger.Info("auto update map data")
//
//	var updatedAt, prefixUrl, msg string
//	publishRecords, err := GetPublishRecord()
//	if err != nil {
//		msg = "GetPublishRecord failed"
//		clog.Logger.Error("GetPublishRecord error:%v", err)
//		goto notice
//	}
//	if len(publishRecords) == 0 {
//		msg = "publishRecords 'length is zero "
//		goto notice
//	}
//	updatedAt = publishRecords[0].Date
//	prefixUrl = publishRecords[0].Link
//
//	clog.Logger.Info("UpdateMapData get info,time:%s ", updatedAt)
//
//	if updatedAt != "" {
//		//从redis获取上一次更新时间 一样则不处理
//		lastUpdateTime := models.GetLastMapDataUpdateTime()
//		if lastUpdateTime == updatedAt {
//			clog.Logger.Info("UpdateMapData time equal,so don't update, time:%s", lastUpdateTime)
//			msg = "最新的一条记录没变，无需更新"
//			goto notice
//		}
//		go UpdateData(prefixUrl, updatedAt)
//	} else {
//		clog.Logger.Error("publishRecords[0].Date  err")
//		msg = "publishRecords[0].Date  是空字符串"
//		goto notice
//	}
//
//notice:
//	fmt.Println("")
//
//}

func UpdateData(prefixUrl string, updateAt string) error {
	provinces, err := GetProvinceUrlAndData(prefixUrl)
	if err != nil {
		return err
	}

	areaList := models.PROVINCE_CITY_REGION_MODEL_LIST{}

	err = areaList.GetAll()
	if err != nil {
		clog.Logger.Error("areaList.GetAll error:%v", err)
		return err
	}

	codeMap := make(map[int]models.PROVINCE_CITY_REGION_MODEL, 0)

	idMap := make(map[int]int, 0)

	for _, area := range areaList {

		//用于删除数据
		idMap[area.ID] = 0
		if area.CityCode == 0 && area.RegionCode == 0 {
			//省
			codeMap[area.ProvinceCode] = area
		} else if area.CityCode != 0 && area.RegionCode == 0 {
			//市
			codeMap[area.CityCode] = area
		} else if area.CityCode != 0 && area.RegionCode != 0 {
			//区
			codeMap[area.RegionCode] = area
		}
	}

	allData := prepareData(provinces)
	for _, data := range allData {
		p := models.PROVINCE_CITY_REGION_MODEL{}
		p.ProvinceCode = data.ProvinceCode
		p.ProvinceName = data.ProvinceName
		p.ProvinceNamePy = data.ProvinceNamePy
		p.CityCode = data.CityCode
		p.CityName = data.CityName
		p.CityNamePy = data.CityNamePy
		p.RegionCode = data.RegionCode
		p.RegionName = data.RegionName
		p.RegionNamePy = data.RegionNamePy
		p.Area = data.Area

		// 省级数据
		if data.CityCode == 0 {
			//判断是否存在
			if area, ok := codeMap[data.ProvinceCode]; !ok {
				saveProvinceErr := p.Save()
				if saveProvinceErr != nil {
					clog.Logger.Error("saveProvinceErr error:%v", err)
					return saveProvinceErr
				}
				clog.Logger.Info("UpdateMapData province v%", p)
			} else {
				//判断各个值是否相等
				if area.Equal(&p) {
					idMap[area.ID] = 1
				} else {
					clog.Logger.Info("UpdateMapData update data 省 [old:%+v] [new:%+v]", area, p)
					saveProvinceErr := p.Save()
					if saveProvinceErr != nil {
						clog.Logger.Error("saveProvinceErr error:%v", saveProvinceErr)
						return saveProvinceErr
					}
				}
			}
			continue
		}
		// 市级数据
		if data.RegionCode == 0 {
			//判断是否存在
			if area, ok := codeMap[data.CityCode]; !ok {
				saveCityErr := p.Save()
				if saveCityErr != nil {
					clog.Logger.Error("if block saveCityErr error:%v", saveCityErr)
					return saveCityErr
				}
				clog.Logger.Info("UpdateMapData city v%", p)
			} else {
				//判断各个值是否相等
				if area.Equal(&p) {
					idMap[area.ID] = 1
				} else {
					clog.Logger.Info("UpdateAmapData update data 市 [old:%+v] [new:%+v]", area, p)
					saveCityErr := p.Save()
					if saveCityErr != nil {
						clog.Logger.Error("else block saveCityErr error:%v", saveCityErr)
						return saveCityErr
					}
				}
			}
			continue
		} else {
			// 区级数据
			//判断是否存在
			if area, ok := codeMap[data.RegionCode]; !ok {
				saveCountyErr := p.Save()
				if saveCountyErr != nil {
					clog.Logger.Error("if block saveCountyErr error:%v", saveCountyErr)
					return saveCountyErr
				}
				clog.Logger.Info("UpdateMapData Region v%", p)
			} else {
				//判断各个值是否相等
				if area.Equal(&p) {
					idMap[area.ID] = 1
				} else {
					clog.Logger.Info("UpdateMapData update data 区 [old:%+v] [new:%+v]", area, p)
					saveCountyErr := p.Save()
					if saveCountyErr != nil {
						clog.Logger.Error("else block saveCountyErr error:%v", saveCountyErr)
						return saveCountyErr
					}
				}
			}
		}
	}

	//检测是否有需要删除的项
	for k, v := range idMap {
		if v == 1 {
			continue
		}
		clog.Logger.Info("UpdateMapData delete k:%d", k)

		//删除
		model := models.PROVINCE_CITY_REGION_MODEL{ID: k}

		deleteErr := model.Delete()

		if deleteErr != nil {
			clog.Logger.Info("UpdateMapData delete data id:%d error:%v", k, err)
		}
	}

	// Todo  存入redis
	models.SetLastMapDataUpdateTime(updateAt)
	return nil
}

// 获取  http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/  这个页面的数据
// 根据ul[class='center_list_contlist'] 获取所有记录的更新日期及其链接地址
func GetPublishRecord() (publishRecords []PublishRecord, err error) {
	fetchUrl := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/"
	tempPublishRecords := make([]PublishRecord, 0)
	defer func() {
		if err == nil {
			publishRecords = tempPublishRecords
		}
	}()
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	// 设置gbk解码，防止乱码
	c.DetectCharset = true
	// 禁用 cookies
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
	c.OnHTML("div[class='center'] div[class='center_list'] ul[class='center_list_contlist']", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("ul li a ", func(i int, element *colly.HTMLElement) bool {
			hrefValue := element.Attr("href")
			if hrefValue == "" {
				err = fmt.Errorf("cant find herf value")
				return false
			}
			recordUrl := strings.ReplaceAll(hrefValue, filepath.Base(hrefValue), "")
			recordsUpdateTime := element.DOM.Find("span font[class='cont_tit02']").Text()
			if recordsUpdateTime == "" {
				err = fmt.Errorf("cant publish time value")
				return false
			}
			record := PublishRecord{
				Date: recordsUpdateTime,
				Link: recordUrl,
			}
			tempPublishRecords = append(tempPublishRecords, record)
			return true
		})
	})

	c.OnError(func(response *colly.Response, er error) {
		err = fmt.Errorf("visit %s OnError:%v", fetchUrl, er)
		return
	})
	if er := c.Visit(fetchUrl); er != nil {
		err = fmt.Errorf("visit %s error:%v", fetchUrl, er)
		return
	}
	return
}

// 获取所有省份对应的链接地址及省级数据
func GetProvinceUrlAndData(prefixUrl string) (provinces []Province, err error) {
	provs := make([]Province, 0)
	defer func() {
		if err == nil {
			provinces = provs
		}
	}()

	//Todo
	c := colly.NewCollector(colly.CacheDir("./缓存"))
	extensions.RandomUserAgent(c)
	// 设置gbk解码，防止乱码
	c.DetectCharset = true
	// 禁用 cookies
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
	//省级列表
	c.OnHTML("tr[class='provincetr']", func(e *colly.HTMLElement) {
		//遍历每一行
		e.ForEachWithBreak("tr[class='provincetr'] > td", func(i int, item *colly.HTMLElement) bool {
			// 省份名称
			provinceName := item.ChildText("a")
			href := item.ChildAttr("a", "href")
			// 每个省份对应的链接地址, 最后一条td 里面没有 a 标签，排除这个td
			provinceHref := prefixUrl + href
			if href != "" && len(href) > 2 {
				provinceCode := href[:2]
				code, atoiErr := strconv.Atoi(provinceCode)
				if atoiErr != nil {
					return false
				}
				p := Province{
					Code:   code,
					Name:   provinceName,
					Link:   provinceHref,
					Cities: nil,
				}
				provs = append(provs, p)
			} else {
				err = errors.New("href 数据错误")
				return false
			}
			return true
		})
	})

	c.OnError(func(response *colly.Response, er error) {
		err = fmt.Errorf("visit %s OnError:%v", prefixUrl, er)
		return
	})

	err = c.Visit(prefixUrl)
	if err != nil {
		clog.Logger.Info("visit %s error: %v: ", prefixUrl, err)
		return
	} else {
		clog.Logger.Info("start visit %s", prefixUrl)
	}

	for i := 0; i < len(provs); i++ {
		cities, getCityErr := GetCityNameAndCode(prefixUrl, provs[i].Link)
		if getCityErr == nil {
			provs[i].Cities = cities
		} else {
			clog.Logger.Info("visit %s error: %v: ", prefixUrl, err)
			err = getCityErr
			return
		}

		for j := 0; j < len(provs[i].Cities); j++ {

			if provs[i].Cities[j].Name == "东莞市" || provs[i].Cities[j].Name == "中山市" {
				towns, GetTownOfDonguanAndhongshanErr := GetTownOfDonguanAndhongshan(prefixUrl, (provs[i].Cities[j]).Link)
				if GetTownOfDonguanAndhongshanErr == nil {
					(provs[i].Cities[j]).Counties = towns
					continue
				} else {
					err = GetTownOfDonguanAndhongshanErr
					return
				}
			} else {
				counties, GetCountyNameAndCodeErr := GetCountyNameAndCode(prefixUrl, (provs[i].Cities[j]).Link)
				if GetCountyNameAndCodeErr == nil {
					(provs[i].Cities[j]).Counties = counties
				} else {
					err = GetCountyNameAndCodeErr
					return
				}
			}
		}
	}

	// 单独增加下面三个地区
	xianggang := Province{
		Code:   81,
		Name:   "香港特别行政区",
		Link:   "",
		Cities: nil,
	}
	aomeng := Province{
		Code:   82,
		Name:   "澳门特别行政区",
		Link:   "",
		Cities: nil,
	}
	taiwang := Province{
		Code:   71,
		Name:   "台湾省",
		Link:   "",
		Cities: nil,
	}

	provs = append(provs, xianggang, aomeng, taiwang)
	return
}

// 获取所有市的链接及市级数据
func GetCityNameAndCode(prefixUrl string, provinceUrl string) (cities []City, err error) {

	cts := make([]City, 0)
	defer func() {
		if err == nil {
			cities = cts
		}
	}()
	// Todo
	c := colly.NewCollector(colly.CacheDir("./缓存"))
	extensions.RandomUserAgent(c)
	// 设置gbk解码
	c.DetectCharset = true
	// 禁用 cookies
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
	//市级列表
	c.OnHTML(".citytable tbody", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("tr[class='citytr']", func(i int, item *colly.HTMLElement) bool {
			// 城市地址
			hrefValue, exists := item.DOM.Find("tr[class='citytr'] td >a").Attr("href")
			if !exists {
				clog.Logger.Info("hrefValue  don't exist")
				return false
			}
			cityUrl := prefixUrl + hrefValue
			text := item.Text
			// 城市code
			if len(text) < 13 {
				clog.Logger.Info("获取市 len(text) < 13 ,数据有问题")
				return false
			}
			cityCode := text[:4]
			code, atoiErr := strconv.Atoi(cityCode)
			if atoiErr != nil {
				clog.Logger.Error("strconv.Atoi(cityCode) error:%v", atoiErr)
				return false
			}
			// 城市名称
			cityName := text[12:]
			city := City{
				Code:     code,
				Name:     cityName,
				Link:     cityUrl,
				Counties: nil,
			}
			cts = append(cts, city)
			return true
		})
	})

	c.OnError(func(response *colly.Response, er error) {
		err = fmt.Errorf("visit %s OnError:%v", provinceUrl, er)
		return
	})

	if er := c.Visit(provinceUrl); er != nil {
		err = fmt.Errorf("visit %s error:%v", provinceUrl, er)
		return
	}
	return
}

// 存储所有区县的链接及其对应的名称和区划代码
func GetCountyNameAndCode(prefixUrl string, cityUrl string) (counties []County, err error) {

	couns := make([]County, 0)
	defer func() {
		if err == nil {
			counties = couns
		}
	}()
	// Todo
	c := colly.NewCollector(colly.CacheDir("./缓存"))
	extensions.RandomUserAgent(c)
	// 设置gbk解码
	c.DetectCharset = true
	// 禁用 cookies
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
	//区县列表
	c.OnHTML(".countytable tbody", func(e *colly.HTMLElement) {
		//遍历每一行
		e.ForEachWithBreak("tr[class='countytr']", func(i int, item *colly.HTMLElement) bool {
			// 获取每个区县的url
			text := item.Text
			if len(text) < 13 {
				clog.Logger.Info("获取区 len(text) < 13 ,数据有问题")
				return false
			}
			topTwo := text[:2]
			threeToFour := text[2:4]
			topSix := text[0:6]
			// 区县代码
			countyCode, _ := strconv.Atoi(topSix)
			// 区县的地址
			countyUrl := prefixUrl + topTwo + "/" + threeToFour + "/" + topSix + ".html"
			// 区县名称
			countyName := text[12:]
			county := County{
				Code: countyCode,
				Name: countyName,
				Link: countyUrl,
			}
			couns = append(couns, county)
			return true
		})
	})

	c.OnError(func(response *colly.Response, er error) {
		err = fmt.Errorf("visit %s OnError:%v", cityUrl, er)
		return
	})
	if er := c.Visit(cityUrl); er != nil {
		err = fmt.Errorf("visit %s error:%v", cityUrl, er)
		return
	}
	return
}

// 获取东莞市和中山市下属的所有镇名称和区划代码
func GetTownOfDonguanAndhongshan(prefixUrl string, cityUrl string) (counties []County, err error) {

	towns := make([]County, 0)
	defer func() {
		if err == nil {
			counties = towns
		}
	}()

	c := colly.NewCollector(colly.CacheDir("./缓存"))
	extensions.RandomUserAgent(c)
	// 设置gbk解码
	c.DetectCharset = true
	// 禁用 cookies
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
	//镇列表
	c.OnHTML(".towntable tbody", func(e *colly.HTMLElement) {
		//遍历每一行
		e.ForEachWithBreak("tr[class='towntr']", func(i int, item *colly.HTMLElement) bool {
			// 获取每个镇的url
			text := item.Text
			if len(text) < 13 {
				clog.Logger.Info("获取镇 len(text) < 13 ,数据有问题")
				return false
			}
			topTwo := text[:2]
			threeToFour := text[2:4]
			topNine := text[0:9]
			// 镇的地址
			townUrl := prefixUrl + topTwo + "/" + threeToFour + "/" + topNine + ".html"
			// 镇的code
			townCode, _ := strconv.Atoi(text[:6])
			// 镇名称
			townName := text[12:]
			// 镇级数据
			town := County{
				Code: townCode,
				Name: townName,
				Link: townUrl,
			}
			towns = append(towns, town)
			return true
		})
	})
	c.OnError(func(response *colly.Response, er error) {
		err = fmt.Errorf("visit %s OnError:%v", cityUrl, er)
		return
	})
	if er := c.Visit(cityUrl); er != nil {
		err = fmt.Errorf("visit %s error:%v", cityUrl, er)
		return
	}
	return
}

// 将抓取到的数据处理成对应数据库表的形式
func prepareData(provinces []Province) []PROVINCE_CITY_REGION_MODEL {
	regions := make([]PROVINCE_CITY_REGION_MODEL, 0)
	dict := pinyin.NewDict()
	for _, p := range provinces {
		pro := PROVINCE_CITY_REGION_MODEL{
			ID:                0,
			ProvinceCode:      p.Code,
			ProvinceName:      p.Name,
			ProvinceNamePy:    dict.Sentence(p.Name).None(),
			CityCode:          0,
			CityName:          "",
			CityNamePy:        "",
			RegionCode:        0,
			RegionName:        "",
			RegionNamePy:      "",
			CityCodeTelephone: "",
			Area:              "",
		}
		regions = append(regions, pro)
		for _, city := range p.Cities {
			cty := PROVINCE_CITY_REGION_MODEL{
				ID:                0,
				ProvinceCode:      p.Code,
				ProvinceName:      p.Name,
				ProvinceNamePy:    dict.Sentence(p.Name).None(),
				CityCode:          city.Code,
				CityName:          city.Name,
				CityNamePy:        dict.Sentence(city.Name).None(),
				RegionCode:        0,
				RegionName:        "",
				RegionNamePy:      "",
				CityCodeTelephone: "",
				Area:              "",
			}
			regions = append(regions, cty)
			for _, county := range city.Counties {
				region := PROVINCE_CITY_REGION_MODEL{
					ID:                0,
					ProvinceCode:      p.Code,
					ProvinceName:      p.Name,
					ProvinceNamePy:    dict.Sentence(p.Name).None(),
					CityCode:          city.Code,
					CityName:          city.Name,
					CityNamePy:        dict.Sentence(city.Name).None(),
					RegionCode:        county.Code,
					RegionName:        county.Name,
					RegionNamePy:      dict.Sentence(county.Name).None(),
					CityCodeTelephone: "",
					Area:              "",
				}
				regions = append(regions, region)
			}
		}
	}
	return regions
}
