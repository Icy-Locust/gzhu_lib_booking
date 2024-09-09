package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/storage"
	"github.com/chromedp/chromedp"
)

var (
	Un           string
	Pd           string
	Cookies      string
	Cookies_path string = "./cookie.txt"
	UserAgent    string = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36`
)

var seatmap = map[string]int{
	"101": 101266684,
	"3A":  100588305,
	"206": 100586975,
}

func Calc_seat(room string, no int) (string, error) {
	devid, valid := seatmap[room]
	if !valid {
		return "", errors.New("calc_seat: unknown room number")
	}
	return strconv.Itoa(devid + no - 1), nil
}

func Get_cookies() error {
	// create chrome instance
	/* Without headless
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background())
	defer cancel()
	*/
	/* With remote debug
	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:9222")
	defer cancel()
	*/
	if Get_room_info() {
		log.Printf("Save cookie: %s", Cookies)
		err := os.WriteFile(Cookies_path, []byte(Cookies), 0644)
		if err != nil {
			log.Println(err)
		}
		return nil
	}
	dat, err := os.ReadFile(Cookies_path)
	if err != nil {
		log.Println(err)
	} else {
		Cookies = string(dat)
		if Get_room_info() {
			log.Printf("Use saved cookie: %s", Cookies)
			return nil
		}
	}
	if Un == "" || Pd == "" {
		return errors.New("get_cookies: empty user name or password")
	}
	log.Println("Start browser login")

	// use allocator
	//ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var all_cookies []*network.Cookie

	err = chromedp.Run(ctx,
		// login
		chromedp.Navigate(`https://newcas.gzhu.edu.cn/cas/login?service=http://libbooking.gzhu.edu.cn/#/ic/home`),
		// wait for page loading
		chromedp.WaitVisible(`#un`),
		// enter user name and password
		chromedp.SendKeys(`#un`, Un),
		chromedp.SendKeys(`#pd`, Pd),
		// click login button
		chromedp.Click(`#index_login_btn`, chromedp.NodeVisible),
		// wait for jumping
		chromedp.WaitVisible(`.footer`),
		// get cookie
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies_list, err := storage.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			all_cookies = cookies_list
			return nil
		}),
	)
	if err != nil {
		return err
	}
	// parse cookie
	ret := ""
	for _, c := range all_cookies {
		if c.Domain == "libbooking.gzhu.edu.cn" {
			ret += c.Name + "=" + c.Value + "; "
		}
	}
	// cut last '; '
	ret = ret[:len(ret)-2]
	Cookies = ret
	log.Printf("Get cookie: %s", Cookies)
	err = os.WriteFile(Cookies_path, []byte(Cookies), 0644)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func Get_room_info() bool {
	if Cookies == "" {
		return false
	}
	// 101 roomid
	roomid := "100647013"
	// current status
	date := time.Now().Format("20060102")
	// constructing request
	url := fmt.Sprintf("http://libbooking.gzhu.edu.cn/ic-web/reserve?roomIds=%s&resvDates=%s&sysKind=8", roomid, date)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	// set header
	req.Header.Set(`Cookie`, Cookies)
	req.Header.Add(`User-Agent`, UserAgent)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)
	ret, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	// get return code
	var dat map[string]interface{}
	if err := json.Unmarshal(ret, &dat); err != nil {
		log.Panic(err)
	}
	return (dat["code"]) == float64(0)
}
