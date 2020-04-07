package api

// 恕我直言，小学生编程班的学生都比你们会写前端
// 工地上可能随便找个打扫卫生的老奶奶都知道写omitempty

// TtRoomData is TT_ROOM_DATA
type TtRoomData struct {
	Type      string `json:"type"`
	State     string `json:"state"`
	IsOn      bool   `json:"isOn"`
	IsOff     bool   `json:"isOff"`
	IsReplay  bool   `json:"isReplay"`
	IsPayRoom bool   `json:"isPayRoom"`
	//IsSecret        int    `json:"isSecret"`
	// RoomPayPassword int    `json:"roomPayPassword"`
	// ID              int64 `json:"id"`
	// SID             int64 `json:"sid"`
	// Channel         int64  `json:"channel"`
	// LiveChannel     int64  `json:"liveChannel"`
	// LiveID          int64  `json:"liveId"`
	// ShortChannel    int64  `json:"shortChannel"`
	// IsBluRay
	GameFullName string `json:"gameFullName"`
	GameHostName string `json:"gameHostName"`
	// ScreenType   int    `json:"screenType"`
	StartTime  interface{} `json:"startTime"`
	TotalCount interface{} `json:"totalCount"`
	// CameraOpen int   `json:"cameraOpen"`
	// LiveCompatibleFlag
	// BussType int `json:"bussType"`
	// IsPlatinum
	ScreenShot string `json:"screenShot"`
	// PreviewURL
	// GameID
	// LiveSourceType  int    `json:"liveSourceType"`
	PrivateHost string `json:"privateHost"`
	ProfileRoom string `json:"profileRoom"`
	// RecommendStatus int    `json:"recommendStatus"`
	// Popular
	// GID                  int64  `json:"gid"`
	Introduction string `json:"introduction"`
	// IsRedirectHuya       int    `json:"isRedirectHuya"`
	IsShowMmsProgramList bool `json:"isShowMmsProgramList"`
}

// TtProfileInfo is TT_PROFILE_INFO
type TtProfileInfo struct{}

// Stream is stream
type Stream struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   []struct {
		GameLiveInfo struct {
			UID                string `json:"uid"`
			Sex                string `json:"sex"`
			GameFullName       string `json:"gameFullName"`
			GameHostName       string `json:"gameHostName"`
			StartTime          string `json:"startTime"`
			ActivityID         string `json:"activityId"`
			Level              string `json:"level"`
			TotalCount         string `json:"totalCount"`
			RoomName           string `json:"roomName"`
			IsSecret           string `json:"isSecret"`
			CameraOpen         string `json:"cameraOpen"`
			LiveChannel        string `json:"liveChannel"`
			BussType           string `json:"bussType"`
			YYID               string `json:"yyid"`
			ScreenShot         string `json:"screenshot"`
			ShortChannel       string `json:"shorChannel"`
			Avatar180          string `json:"avatar180"`
			GID                string `json:"gid"`
			Channel            string `json:"channel"`
			Introduction       string `json:"introduction"`
			ProfileHomeHost    string `json:"profileHomeHost"`
			LiveSourceType     string `json:"liveSourceType"`
			ScreenType         string `json:"screenType"`
			BitRate            string `json:"bitRate"`
			GameType           string `json:"gameType"`
			AttendeeCount      string `json:"attendeeCount"`
			MultiStreamFlag    string `json:"multiStreamFlag"`
			CodecType          string `json:"codecType"`
			LiveCompatibleFlag string `json:"liveCompatibleFlag"`
			ProfileRoom        string `json:"profileRoom"`
			LiveID             string `json:"liveId"`
			RecommendTagName   string `json:"recommendTagName"`
			ContentIntro       string `json:"contentIntro"`
		} `json:"gameLiveInfo"`
		GameStreamInfoList []struct {
			SCdnType            string `json:"sCdnType"`
			IIsMaster           int    `json:"iIsMaster"`
			LChannelID          int64  `json:"lchannelId"`
			LSubChannelID       int64  `json:"lSubChannelId"`
			LPresenterUID       int64  `json:"lPresenterUid"`
			SStreamName         string `json:"sStreamName"`
			SFlvURL             string `json:"sFlvUrl"`
			SFlvURLSuffix       string `json:"sFlvUrlSuffix"`
			SFlvAntiCode        string `json:"sFlvAntiCode"`
			SHlsURL             string `json:"sHlsUrl"`
			SHlsURLSuffix       string `json:"sHlsUrlSuffix"`
			SHlsAntiCode        string `json:"sHlsAntiCode"`
			ILineIndex          int    `json:"iLineIndex"`
			IIsMultiStream      int    `json:"iIsMultiStream"`
			IPCPriorityRate     int    `json:"iPCPriorityRate"`
			IWebPriorityRate    int    `json:"iWebPriorityRate"`
			IMobilePriorityRate int    `json:"iMobilePriorityRate"`
			// VFlvIPList
			IIsP2PSupport   int    `json:"iIsP2PSupport"`
			SP2pURL         string `json:"sP2pUrl"`
			Sp2pURLSuffix   string `json:"sP2pUrlSuffix"`
			LFreeFlag       int    `json:"lFreeFlag"`
			NewCFlvAntiCode string `json:"newCFlvAntiCode"`
		} `json:"gameStreamInfoList"`
	} `json:"data"`
	Count            int `json:"count"`
	VMultiStreamInfo []struct {
		SDisplayName string `json:"sDisplayName"`
		IBitRate     int    `json:"iBitRate"`
	} `json:"vMultiStreamInfo"`
	IWebDefaultBitRate int `json:"iWebDefaultBitRate"`
}
