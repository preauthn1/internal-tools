package model

// DenyIPAddressSettingViewModel ?
type DenyIPAddressSettingViewModel struct {
	IPAddress string `json:"ipAddress"`
}

// LoadBlockIPListBody ?
type LoadBlockIPListBody struct {
	BodyBase
	Data []DenyIPAddressSettingViewModel `json:"data"`
}
