package IdCard

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var CodeData map[string]string

//2017年11月中华人民共和国县以上行政区划代码
//默认已经从下面的网站下载数据保存到data.json
//http://www.mca.gov.cn/article/sj/tjbz/a/2017/201801/201801151447.html
func DownloadCodeData(url string) (map[string]string, error) {
	if url == "" {
		url = "http://www.mca.gov.cn/article/sj/tjbz/a/2017/201801/201801151447.html"
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	// Find the review items
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find(".xl7026226").Text()
		if band != "" {
			m[band[:6]] = band[6:]
		}
	})
	jsonText, err := json.Marshal(m)
	err = ioutil.WriteFile("data.json", jsonText, 0644)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func init() {
	f, err := os.Open("data.json")
	err = json.NewDecoder(f).Decode(&CodeData)
	if err != nil {
		fmt.Println(err)
	}
}

//验证身份号
func Check(id string) bool {
	if CheckCode(id) && CheckBirthday(id) != "" {
		return true
	}
	return false
}

//验证行政区划代码
func CheckAreaCode(id string) bool {
	if CodeData[id[:6]] != "" {
		return true
	}
	return false
}

//验证校验码
func CheckCode(id string) bool {
	id = strings.ToUpper(id)
	codes := []int{1, 0, 88 /* 'X' */, 9, 8, 7, 6, 5, 4, 3, 2}
	var last int
	if id[17:] == "X" {
		last = 88
	} else {
		last, _ = strconv.Atoi(id[17:])
	}
	if codes[checkId(id[:17])] == last {
		return true
	}
	return false
}

//验证生日,返回生日string
func CheckBirthday(id string) string {
	birthday := year(id) + "-" + month(id) + "-" + day(id)
	t, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if t.Unix() > time.Now().Unix() {
		fmt.Println("sorry:不能识别来自未来的人!")
		return ""
	}
	return birthday
}

//获取年龄.
func Age(id string) int {
	birthdayTime, _ := time.Parse("2006-01-02", CheckBirthday(id))
	return timeSub(time.Now(), birthdayTime) / 365
}

//获取地址
func Address(id string) string {
	return province(id) + " " + city(id) + " " + zone(id)
}

//获取性别.
func Sex(id string) string {
	v, _ := strconv.Atoi(id[16:17])
	if v%2 != 0 {
		return "男"
	}
	return "女"
}

//获取星座
func Constellation(id string) string {
	constellation := []string{"水瓶座", "双鱼座", "白羊座", "金牛座", "双子座", "巨蟹座", "狮子座", "处女座", "天秤座", "天蝎座", "射手座", "魔羯座"}
	constellationDays := []int{21, 20, 21, 20, 21, 22, 23, 23, 23, 24, 22, 21}
	month, _ := strconv.Atoi(month(id))
	day, _ := strconv.Atoi(day(id))
	if day < constellationDays[month-1] {
		month--
	}
	if month > 0 {
		return constellation[month-1]
	}
	return constellation[11]
}

//获取属相
func Zodiac(id string) string {
	zodiac := []string{"牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪", "鼠"}
	year, _ := strconv.Atoi(year(id))
	index := (year - 1901) % 12
	return zodiac[index]
}

//通过计算得到最后一位的验证码
func checkId(id17 string) int { // len(id)= 17
	vID := make([]int, 17)
	for i := 0; i < 17; i++ {
		vID[i], _ = strconv.Atoi(id17[i : i+1])
	}
	var wi = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	var res int
	for i := 0; i < 17; i++ {
		res += vID[i] * wi[i]
	}
	return res % 11
}

//获取年.
func year(id string) string {
	return id[6:10]
}

//获取月.
func month(id string) string {
	return id[10:12]
}

//获取日.
func day(id string) string {
	return id[12:14]
}

//2个时间相差天数
func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

//获取省
func province(id string) string {
	return CodeData[id[:2]+"0000"]
}

//获取市
func city(id string) string {
	return CodeData[id[:4]+"00"]
}

//获取区
func zone(id string) string {
	return CodeData[id[:6]]
}
