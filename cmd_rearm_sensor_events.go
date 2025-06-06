package ipmi

import (
	"context"
)

// 35.12 Re-arm Sensor Events Command
type RearmSensorEventsRequest struct {
	SensorNumber uint8

	RearmAllEventStatus bool

	DiscreteEvents bool

	// If the field of SensorEventFlag is true, it means to re-arm the specific event
	SensorEventFlag
}

type RearmSensorEventsResponse struct {
}

func (req *RearmSensorEventsRequest) Command() Command {
	return CommandRearmSensorEvents
}

func (req *RearmSensorEventsRequest) Pack() []byte {
	out := make([]byte, 2)

	out[0] = req.SensorNumber

	var b1 uint8
	// [7] - 0b = re-arm all event status from this sensor.
	b1 = setOrClearBit7(b1, !req.RearmAllEventStatus)
	out[1] = b1

	if req.RearmAllEventStatus {
		return out
	}

	var b2, b3, b4, b5 uint8
	if req.DiscreteEvents {
		b2 = setOrClearBit7(b2, req.SensorEventFlag.SensorEvent_State_7_Assert)
		b2 = setOrClearBit6(b2, req.SensorEventFlag.SensorEvent_State_6_Assert)
		b2 = setOrClearBit5(b2, req.SensorEventFlag.SensorEvent_State_5_Assert)
		b2 = setOrClearBit4(b2, req.SensorEventFlag.SensorEvent_State_4_Assert)
		b2 = setOrClearBit3(b2, req.SensorEventFlag.SensorEvent_State_3_Assert)
		b2 = setOrClearBit2(b2, req.SensorEventFlag.SensorEvent_State_2_Assert)
		b2 = setOrClearBit1(b2, req.SensorEventFlag.SensorEvent_State_1_Assert)
		b2 = setOrClearBit0(b2, req.SensorEventFlag.SensorEvent_State_0_Assert)

		b3 = setOrClearBit6(b3, req.SensorEventFlag.SensorEvent_State_14_Assert)
		b3 = setOrClearBit5(b3, req.SensorEventFlag.SensorEvent_State_13_Assert)
		b3 = setOrClearBit4(b3, req.SensorEventFlag.SensorEvent_State_12_Assert)
		b3 = setOrClearBit3(b3, req.SensorEventFlag.SensorEvent_State_11_Assert)
		b3 = setOrClearBit2(b3, req.SensorEventFlag.SensorEvent_State_10_Assert)
		b3 = setOrClearBit1(b3, req.SensorEventFlag.SensorEvent_State_9_Assert)
		b3 = setOrClearBit0(b3, req.SensorEventFlag.SensorEvent_State_8_Assert)

		b4 = setOrClearBit7(b4, req.SensorEventFlag.SensorEvent_State_7_Deassert)
		b4 = setOrClearBit6(b4, req.SensorEventFlag.SensorEvent_State_6_Deassert)
		b4 = setOrClearBit5(b4, req.SensorEventFlag.SensorEvent_State_5_Deassert)
		b4 = setOrClearBit4(b4, req.SensorEventFlag.SensorEvent_State_4_Deassert)
		b4 = setOrClearBit3(b4, req.SensorEventFlag.SensorEvent_State_3_Deassert)
		b4 = setOrClearBit2(b4, req.SensorEventFlag.SensorEvent_State_3_Deassert)
		b4 = setOrClearBit1(b4, req.SensorEventFlag.SensorEvent_State_1_Deassert)
		b4 = setOrClearBit0(b4, req.SensorEventFlag.SensorEvent_State_0_Deassert)

		b5 = setOrClearBit6(b5, req.SensorEventFlag.SensorEvent_State_14_Deassert)
		b5 = setOrClearBit5(b5, req.SensorEventFlag.SensorEvent_State_13_Deassert)
		b5 = setOrClearBit4(b5, req.SensorEventFlag.SensorEvent_State_12_Deassert)
		b5 = setOrClearBit3(b5, req.SensorEventFlag.SensorEvent_State_11_Deassert)
		b5 = setOrClearBit2(b5, req.SensorEventFlag.SensorEvent_State_10_Deassert)
		b5 = setOrClearBit1(b5, req.SensorEventFlag.SensorEvent_State_9_Deassert)
		b5 = setOrClearBit0(b5, req.SensorEventFlag.SensorEvent_State_8_Deassert)

	} else {
		b2 = setOrClearBit7(b2, req.SensorEventFlag.SensorEvent_UNC_High_Assert)
		b2 = setOrClearBit6(b2, req.SensorEventFlag.SensorEvent_UNC_Low_Assert)
		b2 = setOrClearBit5(b2, req.SensorEventFlag.SensorEvent_LNR_High_Assert)
		b2 = setOrClearBit4(b2, req.SensorEventFlag.SensorEvent_LNR_Low_Assert)
		b2 = setOrClearBit3(b2, req.SensorEventFlag.SensorEvent_LCR_High_Assert)
		b2 = setOrClearBit2(b2, req.SensorEventFlag.SensorEvent_LCR_Low_Assert)
		b2 = setOrClearBit1(b2, req.SensorEventFlag.SensorEvent_LNC_High_Assert)
		b2 = setOrClearBit0(b2, req.SensorEventFlag.SensorEvent_LNC_Low_Assert)

		b3 = setOrClearBit3(b3, req.SensorEventFlag.SensorEvent_UNR_High_Assert)
		b3 = setOrClearBit2(b3, req.SensorEventFlag.SensorEvent_UNR_Low_Assert)
		b3 = setOrClearBit1(b3, req.SensorEventFlag.SensorEvent_UCR_High_Assert)
		b3 = setOrClearBit0(b3, req.SensorEventFlag.SensorEvent_UCR_Low_Assert)

		b4 = setOrClearBit7(b4, req.SensorEventFlag.SensorEvent_UNC_High_Deassert)
		b4 = setOrClearBit6(b4, req.SensorEventFlag.SensorEvent_UNC_Low_Deassert)
		b4 = setOrClearBit5(b4, req.SensorEventFlag.SensorEvent_LNR_High_Deassert)
		b4 = setOrClearBit4(b4, req.SensorEventFlag.SensorEvent_LNR_Low_Deassert)
		b4 = setOrClearBit3(b4, req.SensorEventFlag.SensorEvent_LCR_High_Deassert)
		b4 = setOrClearBit2(b4, req.SensorEventFlag.SensorEvent_LCR_Low_Deassert)
		b4 = setOrClearBit1(b4, req.SensorEventFlag.SensorEvent_LNC_High_Deassert)
		b4 = setOrClearBit0(b4, req.SensorEventFlag.SensorEvent_LNC_Low_Deassert)

		b5 = setOrClearBit3(b5, req.SensorEventFlag.SensorEvent_UNR_High_Deassert)
		b5 = setOrClearBit2(b5, req.SensorEventFlag.SensorEvent_UNR_Low_Deassert)
		b5 = setOrClearBit1(b5, req.SensorEventFlag.SensorEvent_UCR_High_Deassert)
		b5 = setOrClearBit0(b5, req.SensorEventFlag.SensorEvent_UCR_Low_Deassert)
	}

	out = append(out, []byte{b2, b3, b4, b5}...)
	return out
}

func (res *RearmSensorEventsResponse) Unpack(msg []byte) error {
	return nil
}

func (r *RearmSensorEventsResponse) CompletionCodes() map[uint8]string {
	return map[uint8]string{}
}

func (res *RearmSensorEventsResponse) Format() string {
	return ""
}

func (c *Client) RearmSensorEvents(ctx context.Context, request *RearmSensorEventsRequest) (response *RearmSensorEventsResponse, err error) {
	response = &RearmSensorEventsResponse{}
	err = c.Exchange(ctx, request, response)
	return
}
