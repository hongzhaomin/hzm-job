package sdk

type Param[T any] interface {
	GetUrl() string
	GetAccessToken() string
}

func NewBaseParam[T any](url, accessToken string) *BaseParam[T] {
	return &BaseParam[T]{
		Url:         url,
		AccessToken: accessToken,
	}
}

type BaseParam[T any] struct {
	Url         string `json:"url"`
	AccessToken string `json:"accessToken"`
}

func (p *BaseParam[T]) GetUrl() string {
	return p.Url
}

func (p *BaseParam[T]) GetAccessToken() string {
	return p.AccessToken
}
