package ipmi

import "context"

// 22.10 Get BT Interface Capabilities Command
type GetBTInterfaceCapabilitiesRequest struct {
	// empty
}

type GetBTInterfaceCapabilitiesResponse struct {
	OutstandingRequestsCountSupported uint8
	InputBufferMessageSizeBytes       uint8
	OutputBufferMessageSizeBytes      uint8
	BMCRequestToResponseTimeSec       uint8
	RecommendedRetries                uint8
}

func (req *GetBTInterfaceCapabilitiesRequest) Command() Command {
	return CommandGetBTInterfaceCapabilities
}

func (req *GetBTInterfaceCapabilitiesRequest) Pack() []byte {
	return []byte{}
}

func (res *GetBTInterfaceCapabilitiesResponse) Unpack(msg []byte) error {
	// at least 3 bytes
	if len(msg) < 5 {
		return ErrUnpackedDataTooShortWith(len(msg), 5)
	}

	res.OutstandingRequestsCountSupported, _, _ = unpackUint8(msg, 0)
	res.InputBufferMessageSizeBytes, _, _ = unpackUint8(msg, 1)
	res.OutputBufferMessageSizeBytes, _, _ = unpackUint8(msg, 2)
	res.BMCRequestToResponseTimeSec, _, _ = unpackUint8(msg, 3)
	res.RecommendedRetries, _, _ = unpackUint8(msg, 4)
	return nil
}

func (*GetBTInterfaceCapabilitiesResponse) CompletionCodes() map[uint8]string {
	// no command-specific cc
	return map[uint8]string{}
}

func (res *GetBTInterfaceCapabilitiesResponse) Format() string {
	return ""
}

func (c *Client) GetBTInterfaceCapabilities(ctx context.Context) (response *GetBTInterfaceCapabilitiesResponse, err error) {
	request := &GetBTInterfaceCapabilitiesRequest{}
	response = &GetBTInterfaceCapabilitiesResponse{}
	err = c.Exchange(ctx, request, response)
	return
}
