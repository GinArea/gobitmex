package bitmexapi

import "time"

type GetUser struct{}

type User struct {
	Id                    int         `json:"id"`
	Firstname             string      `json:"firstname"`
	Lastname              string      `json:"lastname"`
	Username              string      `json:"username"`
	AccountName           string      `json:"accountName"`
	IsUser                bool        `json:"isUser"`
	Email                 string      `json:"email"`
	DateOfBirth           string      `json:"dateOfBirth"` // format "1980-04-12"
	Phone                 string      `json:"phone"`
	Created               time.Time   `json:"created"`
	LastUpdated           time.Time   `json:"lastUpdated"`
	Preferences           Preferences `json:"preferences"`
	TFAEnabled            string      `json:"TFAEnabled"`
	AffiliateID           string      `json:"affiliateID"`
	Country               string      `json:"country"`
	GeoipCountry          string      `json:"geoipCountry"`
	GeoipRegion           string      `json:"geoipRegion"`
	FirstTradeTimestamp   time.Time   `json:"firstTradeTimestamp"`
	FirstDepositTimestamp time.Time   `json:"firstDepositTimestamp"`
	IsElite               bool        `json:"isElite"`
	LastKnownAuthority    string      `json:"lastKnownAuthority"`
	Typ                   string      `json:"typ"`

	// absent in docs
	IsRestricted             bool `json:"isRestricted"`
	HasPendingPeriodicReview bool `json:"hasPendingPeriodicReview"`
}

type Preferences struct {
	AlertOnLiquidations        bool           `json:"alertOnLiquidations"`
	AnimationsEnabled          bool           `json:"animationsEnabled"`
	AnnouncementsLastSeen      time.Time      `json:"announcementsLastSeen"`
	BotsAdvancedMode           bool           `json:"botsAdvancedMode"`
	BotVideosHidden            bool           `json:"botVideosHidden"`
	ChatChannelID              int            `json:"chatChannelID"`
	ColorTheme                 string         `json:"colorTheme"`
	Currency                   string         `json:"currency"`
	Debug                      bool           `json:"debug"`
	DisableChartQuotes         bool           `json:"disableChartQuotes"`
	DisableEmails              []string       `json:"disableEmails"`
	DisablePush                []string       `json:"disablePush"`
	DisplayCorpEnrollUpsell    bool           `json:"displayCorpEnrollUpsell"`
	EquivalentCurrency         string         `json:"equivalentCurrency"`
	Features                   []string       `json:"features"`
	Favourites                 []string       `json:"favourites"`
	FavouritesAssets           []string       `json:"favouritesAssets"`
	FavouritesOrdered          []string       `json:"favouritesOrdered"`
	FavouriteBots              []string       `json:"favouriteBots"`
	FavouriteContracts         []string       `json:"favouriteContracts"`
	HasSetTradingCurrencies    bool           `json:"hasSetTradingCurrencies"`
	HideConfirmDialogs         []string       `json:"hideConfirmDialogs"`
	HideConnectionModal        bool           `json:"hideConnectionModal"`
	HideFromLeaderboard        bool           `json:"hideFromLeaderboard"`
	HideNameFromLeaderboard    bool           `json:"hideNameFromLeaderboard"`
	HidePnlInGuilds            bool           `json:"hidePnlInGuilds"`
	HideRoiInGuilds            bool           `json:"hideRoiInGuilds"`
	HideNotifications          []string       `json:"hideNotifications"`
	HidePhoneConfirm           bool           `json:"hidePhoneConfirm"`
	GuidesShownVersion         int            `json:"guidesShownVersion"`
	IsSensitiveInfoVisible     bool           `json:"isSensitiveInfoVisible"`
	IsWalletZeroBalanceHidden  bool           `json:"isWalletZeroBalanceHidden"`
	Locale                     string         `json:"locale"`
	LocaleSetTime              int64          `json:"localeSetTime"`
	MarginPnlRow               string         `json:"marginPnlRow"`
	MarginPnlRowKind           string         `json:"marginPnlRowKind"`
	MobileLocale               string         `json:"mobileLocale"`
	MsgsSeen                   []string       `json:"msgsSeen"`
	Notifications              map[string]any `json:"notifications"`
	OptionsBeta                bool           `json:"optionsBeta"`
	OrderBookBinning           map[string]any `json:"orderBookBinning"`
	OrderBookType              string         `json:"orderBookType"`
	OrderClearImmediate        bool           `json:"orderClearImmediate"`
	OrderControlsPlusMinus     bool           `json:"orderControlsPlusMinus"`
	OrderControlsOpenCloseTabs bool           `json:"orderControlsOpenCloseTabs"`
	OrderInputType             string         `json:"orderInputType"`
	PlatformLayout             string         `json:"platformLayout"`
	SelectedFiatCurrency       string         `json:"selectedFiatCurrency"`
	ShowChartBottomToolbar     bool           `json:"showChartBottomToolbar"`
	ShowLocaleNumbers          bool           `json:"showLocaleNumbers"`
	Sounds                     []string       `json:"sounds"`
	SpacingPreference          string         `json:"spacingPreference"`
	StrictIPCheck              bool           `json:"strictIPCheck"`
	StrictTimeout              bool           `json:"strictTimeout"`
	TickerGroup                string         `json:"tickerGroup"`
	TickerPinned               bool           `json:"tickerPinned"`
	TradeLayout                string         `json:"tradeLayout"`
	UserColor                  string         `json:"userColor"`
	VideosSeen                 []string       `json:"videosSeen"`
}

func (o *Client) GetUser() Response[User] {
	return GetUser{}.Do(o)
}

func (o GetUser) Do(c *Client) Response[User] {
	return Get(c, "v1/user", o, identity[User])
}

// https://www.bitmex.com/api/explorer/#!/APIKey/APIKey_get

type ApiKey struct {
	Id              string
	Secret          string
	Name            string
	Nonce           int
	Cidr            string
	Cidrs           []string
	TargetAccountId int
	Permissions     []string
	Enabled         bool
	UserId          int
	Created         time.Time
}

type GetApiKey struct{}

func (o *Client) GetApiKey() Response[[]ApiKey] {
	return GetApiKey{}.Do(o)
}

func (o GetApiKey) Do(c *Client) Response[[]ApiKey] {
	return Get(c, "v1/apiKey", o, identity[[]ApiKey])
}

// https://docs.bitmex.com/api-explorer/user-affiliates-get

type UserAffiliates struct {
	UserID    int       `json:"userId"`
	Lifecycle Lifecycle `json:"lifecycle"`
}

type Lifecycle struct {
	Country       string    `json:"COUNTRY"`
	Registered    time.Time `json:"REGISTERED"`
	EmailVerified time.Time `json:"EMAIL_VERIFIED"`
	KYCCompleted  time.Time `json:"KYC_COMPLETED"`
	FirstDeposit  time.Time `json:"FIRST_DEPOSIT"`
	FirstTrade    time.Time `json:"FIRST_TRADE"`
}

type GetAffiliates struct {
	Depth           string `url:"depth,omitempty"`
	SelectUserId    string `url:"selectUserId,omitempty"`
	TargetAccountId string `url:"targetAccountId,omitempty"`
}

func (o *Client) GetUserAffiliates(v GetAffiliates) Response[[]UserAffiliates] {
	return v.Do(o)
}

func (o GetAffiliates) Do(c *Client) Response[[]UserAffiliates] {
	return Get(c, "v1/userAffiliates", o, identity[[]UserAffiliates])
}
