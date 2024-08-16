package types

type MDBDeviceOnline struct {
	DeviceId string `json:"device_id" bson:"device_id"`
	AddTime  int64  `json:"add_time" bson:"add_time"`
}

type MDBDeviceInfo struct {
	DeviceId string `json:"device_id" bson:"device_id"`
	WsMachineInfoRequest
	AddTime    int64 `json:"add_time" bson:"add_time"`
	UpdateTime int64 `json:"update_time" bson:"update_time"`
}
