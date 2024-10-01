package apophis

import "github.com/ninesbr/sheeps.toolkit.go/apophis/pb"

type ConfirmDelivery struct {
	OK     func()
	Retry  func()
	Ignore func()
}

type MessageResponse struct {
	ConfirmDelivery
	*pb.SubscribeMessage
}

type QueueCreateRequest struct {
	*pb.PubRequest
}

type MessageRequest struct {
	*pb.PubMessageRequest
}
