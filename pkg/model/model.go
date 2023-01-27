package model

type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	EndpointsNum int    `json:"-"`
}

type Endpoint struct {
	Id        int    `json:"-"`
	Url       string `json:"url"`
	Threshold int    `json:"threshold"`
}

type EndpointResponse struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}
type Alert struct {
	Id          int    `json:"id"`
	EndpointID  string `json:"-"`
	Description string `json:"description"`
}
