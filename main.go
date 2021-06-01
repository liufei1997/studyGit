package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"git.in.codoon.com/backend/common/clog"
	"git.in.codoon.com/backend/common/thirdutil"
	"git.in.codoon.com/backend/system_service/models"
	"git.in.codoon.com/third/pinyin/pinyin"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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

type ActiveMainSt struct {
	Id             int64     `gorm:"primary_key;column:id"    sql:"type:int(11);auto_increment;not null" json:"id"`
	CreateTime     time.Time `gorm:"column:create_time"       sql:"type:datetime;not null" json:"create_time"`
	TpId           int       `gorm:"column:tp_id"             sql:"type:int(11);not null;default:0" json:"tp_id"`
	StTime         time.Time `gorm:"column:st_time"           sql:"type:datetime" json:"st_time"`
	EtTime         time.Time `gorm:"column:et_time"           sql:"type:datetime;not null" json:"et_time"`
	Name           string    `gorm:"column:name"              sql:"type:varchar(25)" json:"name"`
	Placard        string    `gorm:"column:placard"           sql:"type:varchar(200)" json:"placard"`
	Details        string    `gorm:"column:details"           sql:"type:varchar(2000)" json:"details"`
	Fee            int       `gorm:"column:fee"               sql:"type:int(11)" json:"fee"`
	DeadLine       time.Time `gorm:"column:deadline"          sql:"type:datetime" json:"deadline"`
	JoinNum        int       `gorm:"column:join_num"          sql:"type:int(11)" json:"join_num"`
	ActualJoinNum  int       `gorm:"column:actual_join_num"   sql:"type:int(11)" json:"actual_join_num"`
	IsNeedApproval int8      `gorm:"column:is_need_approval"  sql:"type:int(6)" json:"is_need_approval"`
	IsNeedPhone    int8      `gorm:"column:is_need_phone"     sql:"type:int(6)" json:"is_need_phone"`
	City           string    `gorm:"column:city"              sql:"type:varchar(20)" json:"city"`
	IngAndLat      string    `gorm:"column:lng_and_lat"       sql:"type:varchar(50)" json:"lng_and_lat"`
	Address        string    `gorm:"column:address"           sql:"type:varchar(50)" json:"address"`
	CreateUserId   string    `gorm:"column:create_user_id"    sql:"type:varchar(36)" json:"create_user_id"`
	IsCodoon       int8      `gorm:"column:is_codoon"         sql:"type:int(6)" json:"is_codoon"`
	SortKey        int       `gorm:"column:sort_key"          sql:"type:int(11)" json:"sort_key"`
	//0:正常 1:举报 2:删除 3:敏感 4取消
	IsDelete        int8      `gorm:"column:is_delete"         sql:"type:int(6) json:"is_delete"`
	FindCity        int       `gorm:"column:find_city"         sql:"type:int(11)" json:"find_city"`
	CancelReason    string    `gorm:"column:cancel_reason"     sql:"type:varchar(300)" json:"cancel_reason"`
	ThirdLink       string    `gorm:"column:third_link"        sql:"type:varchar(100)" json:"third_link"`
	GroupId         int       `gorm:"column:group_id"          sql:"type:int(11)" json:"group_id"`
	NeedGroupMember int       `gorm:"column:need_group_member" sql:"type:int(11)" json:"need_group_member"`
	UpdateTime      time.Time `gorm:"column:update_time"       sql:"type:datetime" json:"update_time"`
	DetailsUrl      string    `gorm:"column:details_url"       sql:"type:varchar(300)" json:"details_url"`
	BrowseCount     int       `gorm:"column:browse_count"      sql:"type:int(11)" json:"browse_count"`
	IsOfficePub     int8      `gorm:"column:is_office_pub"     sql:"type:int(6)" json:"is_office_pub"`
	UserInfoFields  string    `gorm:"column:user_info_fields"  sql:"type:varchar(1024)" json:"user_info_fields"`
	GpsSign         int8      `gorm:"column:gps_sign"          sql:"type:int(6)" json:"gps_sign"`
	GpsSignTime     time.Time `gorm:"column:gps_sign_time"     sql:"type:datetime" json:"gps_sign_time"`
	GpsIngAndLat    string    `gorm:"column:gps_lng_and_lat"   sql:"type:varchar(50)" json:"gps_lng_and_lat"`
	//Phone           string    `gorm:"column:phone"             sql:"type:varchar(100)" json:"phone"`
	PayType int `gorm:"column:pay_type"          sql:"type:int(11)" json:"pay_type"`
	//是否支持退款 0不支持, 1支持
	RefundType  int   `gorm:"column:refund_type"       sql:"type:int(11)" json:"refund_type"`
	PayFee      int64 `gorm:"column:pay_fee"          sql:"type:bigint(20)" json:"pay_fee"`
	ReportCount int   `gorm:"column:report_count"          sql:"type:int(11)" json:"report_count"`
	FrozenType  int   `gorm:"column:frozen_type"          sql:"type:int(11)" json:"frozen_type"`
	//1:基于位置签到 0 不受位置影响签到
	IsSignWithLocation int8      `gorm:"column:is_sign_with_location"    sql:"type:int(6)" json:"is_sign_with_location"`
	FromType           int8      `gorm:"column:from_type"         sql:"type:int(6)" json:"from_type"`
	FromVersion        string    `gorm:"column:from_version"  sql:"type:varchar(16)" json:"from_version"`
	DetailsRnModel     string    `gorm:"column:details_rn_model"           sql:"type:varchar(3000)" json:"details_rn_model"`
	ActiveCityCode     string    `gorm:"column:active_city_code"           sql:"type:varchar(16)" json:"active_city_code"`
	ActiveLoseTime     time.Time `gorm:"active_lose_time"                  sql:"type:datetime"    json:"active_lost_time"`

	//新的报名规则字段
	UserApplyRule     string `json:"user_apply_rule"`
	ExtData           string `json:"ext_data"`
	GpsLiveInfo       string `gorm:"column:gps_live_info"          sql:"type:varchar(1024)" json:"gps_live_info"`
	EnableRanking     bool   `gorm:"not null;default:0" json:"enable_ranking"`      // 开启排名
	ActivateUserCount int    `gorm:"not null;default:0" json:"activate_user_count"` // 总运动人数
	TotalDistance     int    `gorm:"not null;default:0" json:"total_distance"`      // 总数据米

}

type CsvFile struct {
	Name string
	Data []byte
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
	//var name float64
	//name = 81578
	//accLength := name / 1000
	//
	//fmt.Println(float64(int(accLength*100)) / 100)

	//f1, err := os.Open("a")
	//if err != nil {
	//
	//}
	//defer f1.Close()
	//f2, err := os.Open("b")
	//if err != nil {
	//
	//}
	//defer f2.Close()
	//f3, err := os.Open("c")
	//if err != nil {
	//
	//}
	//defer f3.Close()
	//var files = []*os.File{f1, f2, f3}
	//dest := "testCompress666.zip"
	//err = Compress(files, dest)
	//if err != nil {
	//
	//} else {
	//	fmt.Println("压缩成功")
	//}

	//f1, _ :=os.Create("cs.csv")
	//var files = []*os.File{f1}
	//dest := "压缩cs.zip"
	//err := Compress(files, dest)
	//if err != nil {
	//
	//} else {
	//	fmt.Println("压缩成功")
	//}

	//var aaa [][]string
	//aaa =append(aaa,[]string{"111","111","222"})
	//CreateCSV("aaa.csv",aaa)





	//m3u8_dir_path := "13429235688@01@05@09"
	//data := strings.Split(m3u8_dir_path, "@")
	////需要在index=1的地方插入时间戳，方便以“id time emo val arou”的顺序写入csv文件
	////将index=1后面的数据保存到一个临时的切片（其实是对index后面数据的复制）
	//tmp := append([]string{}, data[1:]...)
	//
	////拼接插入的时间戳
	//data = append(data[0:1], time.Now().Format("20060102150405"))
	//
	////与临时切片再组合得到最终的需要的切片
	//data = append(data, tmp...)
	//
	//fmt.Printf("data[0]: %s\n", data[0])
	//fmt.Printf("data[1]: %s\n", data[1])
	//fmt.Printf("data[2]: %s\n", data[2])
	//fmt.Printf("data[3]: %s\n", data[3])
	//fmt.Printf("data[4]: %s\n", data[4])
	//
	//csvName := data[0] + ".csv"
	//file, er := os.Open(csvName)
	//defer file.Close()
	//
	//// 如果文件不存在，创建文件
	//if er != nil && os.IsNotExist(er) {
	//
	//	file, er := os.Create(csvName)
	//	if er != nil {
	//		panic(er)
	//	}
	//	defer file.Close()
	//
	//	// 写入字段标题
	//	w := csv.NewWriter(file) //创建一个新的写入文件流
	//	title := []string{"user_id", "time", "emo", "val", "arou"}
	//
	//	// 这里必须刷新，才能将数据写入文件。
	//	w.Write(title)
	//	w.Write(data)
	//	w.Flush()
	//	fmt.Printf("if end")
	//} else {
	//	// 如果文件存在，直接加在末尾
	//	txt, err := os.OpenFile(csvName, os.O_APPEND|os.O_RDWR, 0666)
	//	defer txt.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//	w := csv.NewWriter(txt) //创建一个新的写入文件流
	//	w.Write(data)
	//	w.Flush()
	//	fmt.Printf("else end")
	//}
	//

	//files := make([]CsvFile, 0)
	//// 文件名， 每行的内容， 主题
	//group := "csvfile1"
	//
	//var content [][]string
	//content = append(content,[]string{"aaa","bbb","ccc"})
	//content = append(content,[]string{"a2","22","c222"})
	//
	//title := []string{"第一列","第二列","第三列"}
	//data, err := GenCSV(group, content, title)
	//if err != nil {
	//}
	//file := CsvFile{
	//	Name: fmt.Sprintf("%v.csv", group),
	//	Data: data,
	//}
	//files = append(files, file)
	//zipData, err := BytesZip(files)
	//fileName := "测试文件1"
	//
	//fileUrl, err := PutZipToOss(zipData, fileName)
	//fmt.Println(fileUrl)


	title := []string{"1","2", "3"}
	var rows [][]string
	rows = append(rows,[]string{"a","b","c"})
	rows = append([][]string{title}, rows...)
	fmt.Println(rows)

	generateCSV("csv111",rows)



}

func PutZipToOss(data []byte, fileName string) (string, error) {
	filUrl, err := thirdutil.DefaultOssClient.PutU(fileName, data)
	return filUrl, err
}


func BytesZip(files []CsvFile) ([]byte, error) {
	buff := new(bytes.Buffer)
	w := zip.NewWriter(buff)
	for _, v := range files {
		ww, err := w.Create(v.Name)
		if err != nil {
			return nil, err
		}
		_, err = ww.Write(v.Data)
		if err != nil {
			return nil, err
		}
	}
	err := w.Close()
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func GenCSV(groupName string, rows [][]string, title []string) ([]byte, error) {
	if groupName == "" {
		return nil, errors.New("sheet name is nil")
	}
	if title != nil && len(title) > 0 {
		rows = append([][]string{title}, rows...)
	}
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	err := w.Write([]string{"\xEF\xBB\xBF"})
	if err != nil {
		return nil, err
	}
	err = w.WriteAll(rows)
	return buf.Bytes(), err
}

func generateCSV(csvName string, rows [][]string) ([]byte, error) {

	if csvName == "" {
		return nil, errors.New("sheet name is nil")
	}
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	err := w.Write([]string{"\xEF\xBB\xBF"})
	if err != nil {
		return nil, err
	}
	err = w.WriteAll(rows)
	return buf.Bytes(), err
}

//生成csv文件
func CreateCSV(txtname string,title [][]string)  {
	f , err := os.Create(txtname)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF")
	w:=csv.NewWriter(f)
	w.WriteAll(title)
	w.Flush()
}

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = prefix + "/" + header.Name
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	file.Close()
	if err != nil {
		return err
	}
	return nil
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
