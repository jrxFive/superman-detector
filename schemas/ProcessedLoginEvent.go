package schemas

type ProcessedLoginEvent struct {
	CurrentGeo                     Geo       `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool      `json:"travelToCurrentGeoSuspicious,omitempty"`
	TravelFromCurrentGeoSuspicious bool      `json:"travelFromCurrentGeoSuspicious,omitempty"`
	PrecedingIpAccess              *IPAccess `json:"precedingIPAccess,omitempty"`
	SubsequentIpAccess             *IPAccess `json:"subsequentIPAccess,omitempty"`
}
