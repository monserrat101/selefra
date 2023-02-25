package log

import (
	"github.com/songzhibin97/gkit/coding"
	_ "github.com/songzhibin97/gkit/coding/proto"
)

var protoCoding = coding.GetCode("proto")

func TransportWsMsg(rec *UploadLogStream_Request) (ret []byte, err error) {
	return protoCoding.Marshal(rec)
}
