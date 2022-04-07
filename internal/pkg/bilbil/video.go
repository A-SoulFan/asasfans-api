package bilbil

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	webVideoSearchURL  = "https://api.bilibili.com/x/web-interface/search/type?context=&search_type=video&page=%d&order=pubdate&keyword=%s&duration=0&category_id=&tids_2=&__refresh__=true&_extra=&tids=0&highlight=1&single_column=0"
	webVideoInfoURL    = "https://api.bilibili.com/x/web-interface/view?bvid=%s"
	webVideoTagInfoURL = "https://api.bilibili.com/x/web-interface/view/detail/tag?aid=%s"
)

type SDK struct {
	logger *zap.Logger
}

func NewSDK(logger *zap.Logger) *SDK {
	return &SDK{logger: logger}
}

type ResponseBasic struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Ttl     int         `json:"ttl"`
	Data    interface{} `json:"data"`
}

type VideoSearchInfo struct {
	Type         string        `json:"type"`
	Id           int           `json:"id"`
	Author       string        `json:"author"`
	Mid          int           `json:"mid"`
	Typeid       string        `json:"typeid"`
	Typename     string        `json:"typename"`
	Arcurl       string        `json:"arcurl"`
	Aid          int           `json:"aid"`
	Bvid         string        `json:"bvid"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Arcrank      string        `json:"arcrank"`
	Pic          string        `json:"pic"`
	Play         int           `json:"play"`
	VideoReview  int           `json:"video_review"`
	Favorites    int           `json:"favorites"`
	Tag          string        `json:"tag"`
	Review       int           `json:"review"`
	Pubdate      int           `json:"pubdate"`
	Senddate     int           `json:"senddate"`
	Duration     string        `json:"duration"`
	Badgepay     bool          `json:"badgepay"`
	HitColumns   []string      `json:"hit_columns"`
	ViewType     string        `json:"view_type"`
	IsPay        int           `json:"is_pay"`
	IsUnionVideo int           `json:"is_union_video"`
	RecTags      interface{}   `json:"rec_tags"`
	NewRecTags   []interface{} `json:"new_rec_tags"`
	RankScore    int           `json:"rank_score"`
	Like         int           `json:"like"`
	Upic         string        `json:"upic"`
	Corner       string        `json:"corner"`
	Cover        string        `json:"cover"`
	Desc         string        `json:"desc"`
	Url          string        `json:"url"`
	RecReason    string        `json:"rec_reason"`
}

type SearchResponse struct {
	Seid           string `json:"seid"`
	Page           int    `json:"page"`
	Pagesize       int    `json:"pagesize"`
	NumResults     int    `json:"numResults"`
	NumPages       int    `json:"numPages"`
	SuggestKeyword string `json:"suggest_keyword"`
	RqtType        string `json:"rqt_type"`
	CostTime       struct {
		ParamsCheck         string `json:"params_check"`
		IllegalHandler      string `json:"illegal_handler"`
		AsResponseFormat    string `json:"as_response_format"`
		AsRequest           string `json:"as_request"`
		SaveCache           string `json:"save_cache"`
		DeserializeResponse string `json:"deserialize_response"`
		AsRequestFormat     string `json:"as_request_format"`
		Total               string `json:"total"`
		MainHandler         string `json:"main_handler"`
	} `json:"cost_time"`
	ExpList    interface{}       `json:"exp_list"`
	EggHit     int               `json:"egg_hit"`
	Result     []VideoSearchInfo `json:"result"`
	ShowColumn int               `json:"show_column"`
}

type VideoInfoResponse struct {
	Bvid      string `json:"bvid"`
	Aid       int    `json:"aid"`
	Videos    int    `json:"videos"`
	Tid       int    `json:"tid"`
	Tname     string `json:"tname"`
	Copyright int    `json:"copyright"`
	Pic       string `json:"pic"`
	Title     string `json:"title"`
	Pubdate   int    `json:"pubdate"`
	Ctime     int    `json:"ctime"`
	Desc      string `json:"desc"`
	DescV2    []struct {
		RawText string `json:"raw_text"`
		Type    int    `json:"type"`
		BizId   int    `json:"biz_id"`
	} `json:"desc_v2"`
	State     int `json:"state"`
	Duration  int `json:"duration"`
	MissionId int `json:"mission_id"`
	Rights    struct {
		Bp            int `json:"bp"`
		Elec          int `json:"elec"`
		Download      int `json:"download"`
		Movie         int `json:"movie"`
		Pay           int `json:"pay"`
		Hd5           int `json:"hd5"`
		NoReprint     int `json:"no_reprint"`
		Autoplay      int `json:"autoplay"`
		UgcPay        int `json:"ugc_pay"`
		IsCooperation int `json:"is_cooperation"`
		UgcPayPreview int `json:"ugc_pay_preview"`
		NoBackground  int `json:"no_background"`
		CleanMode     int `json:"clean_mode"`
		IsSteinGate   int `json:"is_stein_gate"`
		Is360         int `json:"is_360"`
		NoShare       int `json:"no_share"`
	} `json:"rights"`
	Owner struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"owner"`
	Stat struct {
		Aid        int    `json:"aid"`
		View       int    `json:"view"`
		Danmaku    int    `json:"danmaku"`
		Reply      int    `json:"reply"`
		Favorite   int    `json:"favorite"`
		Coin       int    `json:"coin"`
		Share      int    `json:"share"`
		NowRank    int    `json:"now_rank"`
		HisRank    int    `json:"his_rank"`
		Like       int    `json:"like"`
		Dislike    int    `json:"dislike"`
		Evaluation string `json:"evaluation"`
		ArgueMsg   string `json:"argue_msg"`
	} `json:"stat"`
	Dynamic   string `json:"dynamic"`
	Cid       int    `json:"cid"`
	Dimension struct {
		Width  int `json:"width"`
		Height int `json:"height"`
		Rotate int `json:"rotate"`
	} `json:"dimension"`
	NoCache bool `json:"no_cache"`
	Pages   []struct {
		Cid       int    `json:"cid"`
		Page      int    `json:"page"`
		From      string `json:"from"`
		Part      string `json:"part"`
		Duration  int    `json:"duration"`
		Vid       string `json:"vid"`
		Weblink   string `json:"weblink"`
		Dimension struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			Rotate int `json:"rotate"`
		} `json:"dimension"`
		FirstFrame string `json:"first_frame"`
	} `json:"pages"`
	Subtitle struct {
		AllowSubmit bool          `json:"allow_submit"`
		List        []interface{} `json:"list"`
	} `json:"subtitle"`
	IsSeasonDisplay bool `json:"is_season_display"`
	UserGarb        struct {
		UrlImageAniCut string `json:"url_image_ani_cut"`
	} `json:"user_garb"`
	HonorReply struct {
	} `json:"honor_reply"`
}

type VideoTagInfo struct {
	TagId        int    `json:"tag_id"`
	TagName      string `json:"tag_name"`
	Cover        string `json:"cover"`
	HeadCover    string `json:"head_cover"`
	Content      string `json:"content"`
	ShortContent string `json:"short_content"`
	Type         int    `json:"type"`
	State        int    `json:"state"`
	Ctime        int    `json:"ctime"`
	Count        struct {
		View  int `json:"view"`
		Use   int `json:"use"`
		Atten int `json:"atten"`
	} `json:"count"`
	IsAtten         int    `json:"is_atten"`
	Likes           int    `json:"likes"`
	Hates           int    `json:"hates"`
	Attribute       int    `json:"attribute"`
	Liked           int    `json:"liked"`
	Hated           int    `json:"hated"`
	ExtraAttr       int    `json:"extra_attr"`
	TagType         string `json:"tag_type"`
	IsActivity      bool   `json:"is_activity"`
	Color           string `json:"color"`
	Alpha           int    `json:"alpha"`
	IsSeason        bool   `json:"is_season"`
	SubscribedCount int    `json:"subscribed_count"`
	ArchiveCount    string `json:"archive_count"`
	FeaturedCount   int    `json:"featured_count"`
	JumpUrl         string `json:"jump_url"`
}

type VideoTagResponse []VideoTagInfo

func (vtr VideoTagResponse) ToTagStringSlice() []string {
	list := make([]string, 0, 10)
	for _, info := range vtr {
		list = append(list, info.TagName)
	}
	return list
}

func (sdk *SDK) fastGet(url string, data interface{}) error {
	result := &ResponseBasic{Data: &data}

	client := resty.New()
	if resp, err := client.R().SetResult(result).Get(url); err != nil {
		return err
	} else {
		if resp.StatusCode() != http.StatusOK {
			sdk.logger.Error("request bilibili error", zap.Int("http_code", resp.StatusCode()), zap.String("request_url", url))
			return errors.New("request bilibili fail")
		}

		if result.Code != 0 {
			sdk.logger.Error("get bilbil error", zap.Int("result_code", result.Code), zap.String("result_msg", result.Message), zap.String("request_url", url))
			return errors.New(result.Message)
		}

		return nil
	}
}

func (sdk *SDK) VideoWebSearch(keyword string, page int) (data *SearchResponse, err error) {
	if err = sdk.fastGet(fmt.Sprintf(webVideoSearchURL, page, keyword), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (sdk *SDK) VideoWebSearchToInfoList(keyword string, page int) (list []VideoSearchInfo, totalPage int, err error) {
	if data, err := sdk.VideoWebSearch(keyword, page); err != nil {
		return nil, 0, err
	} else {
		return data.Result, data.NumPages, nil
	}
}

func (sdk *SDK) VideoWebInfo(bvid string) (data *VideoInfoResponse, err error) {
	if err = sdk.fastGet(fmt.Sprintf(webVideoInfoURL, bvid), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (sdk SDK) VideoWebTagInfo(aid string) (data *VideoTagResponse, err error) {
	if err = sdk.fastGet(fmt.Sprintf(webVideoTagInfoURL, aid), &data); err != nil {
		return nil, err
	}
	return data, nil
}
