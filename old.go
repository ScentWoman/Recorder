package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	_ "net/http/pprof"
)

type LiveInfo struct {
	isLive    bool
	title     string
	hostname  string
	streamUrl map[string]string
}

type counter struct {
	c  int
	mu sync.Mutex
}

type pipe struct {
	resp       *http.Response
	ctx        context.Context
	subcounter *counter
}

var (
	Url = flag.String("url", "", "liveroom url")

	jsonReg        = regexp.MustCompile(`"stream".*?}   `)
	SplitInterval  = 3600 * time.Second
	TickerInterval = 30 * time.Second
	//srv           *drive.Service
	//workDirID     = ""
	mainDirID      = "Save/"
	backDirID      = "Save/Backup/"
	report         counter
	txFlag, alFlag = &atomic.Value{}, &atomic.Value{}
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("127.0.0.1:8888", nil))
	}()
	flag.Parse()
	// srv = initSrv()
	// checkWorkDir()

	if *Url == "" {
		log.Fatal("Url is empty! Please check the usage using -help.")
	}

	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt)

	alFlag.Store(true)
	txFlag.Store(false)

	go work()
	go func() {
		log.Println("logger")
		t := time.NewTicker(TickerInterval)
		defer t.Stop()
		for {
			<-t.C
			report.mu.Lock()
			if TickerInterval <= time.Second {
				fmt.Printf("\rSpeed: %7.1f KB/s.", float64(report.c)/1024.0/float64(TickerInterval/time.Second))
			} else {
				fmt.Printf("Speed: %7.1f KB/s.\n", float64(report.c)/1024.0/float64(TickerInterval/time.Second))
			}
			report.c = 0
			report.mu.Unlock()
		}
	}()

	<-c
}

func work() {
	for {
		//checkDateDir()
		if !alFlag.Load().(bool) && !txFlag.Load().(bool) {
			time.Sleep(time.Second)
		} else {
			info := parseLiveInfo(*Url)
			if info.isLive {
				log.Println()
				log.Println(info.hostname+":", info.title, "is on live.")

				suffix := info.hostname + "_" + info.title

				if info.streamUrl["AL"] != "" && alFlag.Load().(bool) {
					req, _ := http.NewRequest("GET", info.streamUrl["AL"], nil)
					req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36")
					req.Header.Add("referer", *Url)
					req.Header.Add("Sec-Fetch-Mode", "cors")
					alFlag.Store(false)
					go record(req, nameRegularize(suffix+"_"), mainDirID, alFlag)
				}

				if info.streamUrl["TX"] != "" && txFlag.Load().(bool) {
					req, _ := http.NewRequest("GET", info.streamUrl["TX"], nil)
					req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36")
					req.Header.Add("referer", *Url)
					req.Header.Add("Sec-Fetch-Mode", "cors")
					txFlag.Store(false)
					go record(req, nameRegularize(suffix+"_"), backDirID, txFlag)
				}
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// func checkDateDir() {
// 	dirList, err := srv.Files.List().
// 		Q("'" + workDirID + "' in parents and name='" + time.Now().Add(-6*time.Hour).Format("2006-01-02") + "'").Do()
// 	if err != nil {
// 		log.Println("get dir", err)
// 		return
// 	}
// 	if len(dirList.Files) == 0 {
// 		log.Println("Create:", time.Now().Add(-6*time.Hour).Format("2006-01-02"))
// 		f, fe := srv.Files.Create(&drive.File{
// 			Name:     time.Now().Add(-6 * time.Hour).Format("2006-01-02"),
// 			MimeType: "application/vnd.google-apps.folder",
// 			Parents:  []string{workDirID},
// 		}).Do()
// 		if fe != nil {
// 			log.Println("create dir:", fe)
// 		}
// 		mainDirID = f.Id
// 	} else {
// 		mainDirID = dirList.Files[0].Id
// 		log.Println(dirList.Files[0].Name)
// 	}

// 	dirList, err = srv.Files.List().
// 		Q("'" + mainDirID + "' in parents and name='Backup'").Do()
// 	if err != nil {
// 		log.Println("get Backup dir", err)
// 		return
// 	}
// 	if len(dirList.Files) == 0 {
// 		log.Println("Create:", time.Now().Add(-6*time.Hour).Format("2006-01-02"))
// 		f, fe := srv.Files.Create(&drive.File{
// 			Name:     "Backup",
// 			MimeType: "application/vnd.google-apps.folder",
// 			Parents:  []string{mainDirID},
// 		}).Do()
// 		if fe != nil {
// 			log.Println("create Backup dir:", fe)
// 		}
// 		backDirID = f.Id
// 	} else {
// 		backDirID = dirList.Files[0].Id
// 		log.Println(dirList.Files[0].Name)
// 	}
// }

// func checkWorkDir() {
// 	list, _ := srv.Files.List().Q("'root' in parents and name='Record'").Do()
// 	log.Println("Find:", list.Files[0].Name)
// 	workDirID = list.Files[0].Id
// }

func record(req *http.Request, suffix string, dir string, flag *atomic.Value) {
	ctx, cancel := context.WithTimeout(context.Background(), SplitInterval)
	defer func() {
		cancel()
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	defer flag.Store(true)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	p := pipe{
		resp:       resp,
		ctx:        ctx,
		subcounter: &counter{},
	}

	go func() {
		time.Sleep(10 * time.Second)
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			<-t.C
			p.subcounter.mu.Lock()
			if p.subcounter.c == 0 {
				cancel()
				p.subcounter.mu.Unlock()
				return
			}
			p.subcounter.c = 0
			p.subcounter.mu.Unlock()
		}
	}()

	f, e := os.Create(dir + suffix + time.Now().Format("15-04-05") + ".flv")
	if e != nil {
		log.Println("Create file:", e)
		return
	}
	_, err = io.Copy(f, p)

	// _, err = srv.Files.Create(&drive.File{
	// 	Name:    suffix + time.Now().Format("15-04-05") + ".flv",
	// 	Parents: []string{dir},
	// }).Media(p).Do()
	if err != nil {
		log.Println(err)
		return
	}
}

func parseLiveInfo(url string) (info LiveInfo) {
	info = LiveInfo{
		streamUrl: make(map[string]string),
	}
	defer func() { _ = recover() }()

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch url:", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	text := string(body)

	info.isLive = strings.Contains(text, "<body class=\"liveStatus-on")

	if strings.Contains(text, "哎呀，虎牙君找不到这个主播，要不搜索看看？") {
		log.Fatal("Liveroom does not exist!")
	}

	if strings.Contains(text, `"isOn":false`) {
		//log.Fatal("Live is off.")
		return
	}
	jsonPart := jsonReg.FindString(text)

	var ijson interface{}
	_ = json.Unmarshal([]byte(html.UnescapeString(jsonPart[10:len(jsonPart)-3])), &ijson)
	liveInfoList := ((ijson.(map[string]interface{})["data"]).([]interface{})[0]).(map[string]interface{})["gameLiveInfo"].(map[string]interface{})
	streamInfoList := ((ijson.(map[string]interface{})["data"]).([]interface{})[0]).(map[string]interface{})["gameStreamInfoList"].([]interface{})

	for _, v := range streamInfoList {
		m := v.(map[string]interface{})
		surl := itoS(m["sFlvUrl"]) + "/" + itoS(m["sStreamName"]) + "." + itoS(m["sFlvUrlSuffix"]) + "?t=100&sv=1910112100&" + itoS(m["sFlvAntiCode"])
		surl = strings.ReplaceAll(surl, "http://", "https://")
		info.streamUrl[m["sCdnType"].(string)] = surl
	}
	info.hostname = itoS(liveInfoList["nick"])
	info.title = itoS(liveInfoList["introduction"])

	return
}

func (p pipe) Read(b []byte) (n int, e error) {
	select {
	case <-p.ctx.Done():
		e = io.EOF
		return
	default:
		n, e = p.resp.Body.Read(b)
		report.mu.Lock()
		report.c += n
		report.mu.Unlock()
		p.subcounter.mu.Lock()
		p.subcounter.c += n
		p.subcounter.mu.Unlock()
		return
	}
}

// func initSrv() *drive.Service {
// 	b, err := ioutil.ReadFile("credentials.json")
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	// If modifying these scopes, delete your previously saved token.json.
// 	config, err := google.ConfigFromJSON(b, drive.DriveScope)
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}
// 	client := getClient(config)

// 	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Drive client: %v", err)
// 	}

// 	return srv
// }

// // Retrieve a token, saves the token, then returns the generated client.
// func getClient(config *oauth2.Config) *http.Client {
// 	// The file token.json stores the user's access and refresh tokens, and is
// 	// created automatically when the authorization flow completes for the first
// 	// time.
// 	tokFile := "token.json"
// 	tok, err := tokenFromFile(tokFile)
// 	if err != nil {
// 		tok = getTokenFromWeb(config)
// 		saveToken(tokFile, tok)
// 	}
// 	return config.Client(context.Background(), tok)
// }

// // Request a token from the web, then returns the retrieved token.
// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
// 	fmt.Printf("Go to the following link in your browser then type the "+
// 		"authorization code: \n%v\n", authURL)

// 	var authCode string
// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		log.Fatalf("Unable to read authorization code %v", err)
// 	}

// 	tok, err := config.Exchange(context.TODO(), authCode)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve token from web %v", err)
// 	}
// 	return tok
// }

// // Retrieves a token from a local file.
// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(tok)
// 	return tok, err
// }

// // Saves a token to a file path.
// func saveToken(path string, token *oauth2.Token) {
// 	fmt.Printf("Saving credential file to: %s\n", path)
// 	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		log.Fatalf("Unable to cache oauth token: %v", err)
// 	}
// 	defer f.Close()
// 	_ = json.NewEncoder(f).Encode(token)
// }

func nameRegularize(name string) string {
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "*", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "|", "_")
	name = strings.ReplaceAll(name, "\n", "_")
	name = strings.ReplaceAll(name, "\r", "_")
	name = strings.ReplaceAll(name, " ", "_")
	if len(name) > 255 {
		name = name[:255]
	}
	return name
}

func itoS(v interface{}) (s string) {
	defer func() { _ = recover() }()
	s = v.(string)
	return
}
