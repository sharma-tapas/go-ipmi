package ipmi

import (
	"context"
	"fmt"
)

// 28.2 Get Chassis Status Command
type GetChassisStatusRequest struct {
	// empty
}

type GetChassisStatusResponse struct {
	// Current Power State
	PowerRestorePolicy PowerRestorePolicy
	PowerControlFault  bool // Controller attempted to turn system power on or off, but system did not enter desired state.
	PowerFault         bool // fault detected in main power subsystem
	InterLock          bool // chassis is presently shut down because a chassis	panel interlock switch is active
	PowerOverload      bool // system shutdown because of power overload condition.
	PowerIsOn          bool // 系统电源：上电

	// Last Power Event
	LastPowerOnByCommand                   bool
	LastPowerDownByPowerFault              bool
	LastPowerDownByPowerInterlockActivated bool
	LastPowerDownByPowerOverload           bool
	ACFailed                               bool

	// Last Power Event

	// Misc. Chassis State
	ChassisIdentifySupported bool
	ChassisIdentifyState     ChassisIdentifyState
	CollingFanFault          bool
	DriveFault               bool
	FrontPanelLockoutActive  bool // (power off and reset via chassis push-buttons disabled. 前面板锁定)
	ChassisIntrusionActive   bool // 机箱入侵:（机箱盖被打开）

	// Front Panel Button Capabilities and disable/enable status (Optional)
	SleepButtonDisableAllowed      bool
	DiagnosticButtonDisableAllowed bool
	ResetButtonDisableAllowed      bool
	PoweroffButtonDisableAllowed   bool
	SleepButtonDisabled            bool
	DiagnosticButtonDisabled       bool
	ResetButtonDisabled            bool
	PoweroffButtonDisabled         bool
}

type ChassisIdentifyState uint8

const (
	ChassisIdentifyStateOff          ChassisIdentifyState = 0
	ChassisIdentifyStateTemporaryOn  ChassisIdentifyState = 1
	ChassisIdentifyStateIndefiniteOn ChassisIdentifyState = 2
)

func (c ChassisIdentifyState) String() string {
	m := map[ChassisIdentifyState]string{
		0: "Off",
		1: "Temporary (timed) On",
		2: "Indefinite On",
	}
	s, ok := m[c]
	if ok {
		return s
	}
	return "reserved"
}

// PowerRestorePolicy
// 通电开机策略
type PowerRestorePolicy uint8

const (
	PowerRestorePolicyAlwaysOff PowerRestorePolicy = 0 // 保持下电（关机）
	PowerRestorePolicyPrevious  PowerRestorePolicy = 1 // 与之前保持一致（恢复断电前状态）
	PowerRestorePolicyAlwaysOn  PowerRestorePolicy = 2 // 保持上电（开机）
)

var SupportedPowerRestorePolicies = []string{
	"always-off", "always-on", "previous",
}

func (p PowerRestorePolicy) String() string {
	m := map[PowerRestorePolicy]string{
		0: "always-off", // chassis stays powered off after AC/mains returns
		1: "previous",   // after AC returns, power is restored to the state that was in effect when AC/mains was lost
		2: "always-on",  // chassis always powers up after AC/mains returns
	}
	s, ok := m[p]
	if ok {
		return s
	}
	return "unknown"
}

func (req *GetChassisStatusRequest) Pack() []byte {
	return []byte{}
}

func (req *GetChassisStatusRequest) Command() Command {
	return CommandGetChassisStatus
}

func (res *GetChassisStatusResponse) CompletionCodes() map[uint8]string {
	return map[uint8]string{}
}

func (res *GetChassisStatusResponse) Unpack(msg []byte) error {
	if len(msg) < 3 {
		return ErrUnpackedDataTooShortWith(len(msg), 3)
	}

	b1, _, _ := unpackUint8(msg, 0)
	// first clear bit 7, then shift right 5 bits
	b := (b1 & 0x7f) >> 5
	res.PowerRestorePolicy = PowerRestorePolicy(b)
	res.PowerControlFault = isBit4Set(b1)
	res.PowerFault = isBit3Set(b1)
	res.InterLock = isBit2Set(b1)
	res.PowerOverload = isBit1Set(b1)
	res.PowerIsOn = isBit0Set(b1)

	b2, _, _ := unpackUint8(msg, 1)
	res.LastPowerOnByCommand = isBit4Set(b2)
	res.LastPowerDownByPowerFault = isBit3Set(b2)
	res.LastPowerDownByPowerInterlockActivated = isBit2Set(b2)
	res.LastPowerDownByPowerOverload = isBit1Set(b2)
	res.ACFailed = isBit0Set(b2)

	b3, _, _ := unpackUint8(msg, 2)
	res.ChassisIdentifySupported = isBit6Set(b3)
	res.ChassisIdentifyState = ChassisIdentifyState((b3 & 0x30) >> 4)
	res.CollingFanFault = isBit3Set(b3)
	res.DriveFault = isBit2Set(b3)
	res.FrontPanelLockoutActive = isBit1Set(b3)
	res.ChassisIntrusionActive = isBit0Set(b3)

	if len(msg) == 4 {
		b4, _, _ := unpackUint8(msg, 3)
		res.SleepButtonDisableAllowed = isBit7Set(b4)
		res.DiagnosticButtonDisableAllowed = isBit6Set(b4)
		res.ResetButtonDisableAllowed = isBit5Set(b4)
		res.PoweroffButtonDisableAllowed = isBit4Set(b4)
		res.SleepButtonDisabled = isBit3Set(b4)
		res.DiagnosticButtonDisabled = isBit2Set(b4)
		res.ResetButtonDisabled = isBit1Set(b4)
		res.PoweroffButtonDisabled = isBit0Set(b4)
	}
	return nil
}

func (res *GetChassisStatusResponse) Format() string {
	return "" +
		fmt.Sprintf("System Power         : %s\n", formatBool(res.PowerIsOn, "on", "off")) +
		fmt.Sprintf("Power Overload       : %v\n", res.PowerOverload) +
		fmt.Sprintf("Power Interlock      : %s\n", formatBool(res.InterLock, "active", "inactive")) +
		fmt.Sprintf("Main Power Fault     : %v\n", res.PowerFault) +
		fmt.Sprintf("Power Control Fault  : %v\n", res.PowerControlFault) +
		fmt.Sprintf("Power Restore Policy : %s\n", res.PowerRestorePolicy.String()) +
		fmt.Sprintf("Last Power Event     : %s\n", formatBool(res.ChassisIntrusionActive, "active", "inactive")) +
		fmt.Sprintf("Chassis Intrusion    : %s\n", formatBool(res.ChassisIntrusionActive, "active", "inactive")) +
		fmt.Sprintf("Front-Panel Lockout  : %s\n", formatBool(res.FrontPanelLockoutActive, "active", "inactive")) +
		fmt.Sprintf("Drive Fault          : %v\n", res.DriveFault) +
		fmt.Sprintf("Cooling/Fan Fault    : %v\n", res.CollingFanFault) +
		fmt.Sprintf("Sleep Button Disable : %s\n", formatBool(res.SleepButtonDisableAllowed, "allowed", "disallowed")) +
		fmt.Sprintf("Diag Button Disable  : %s\n", formatBool(res.DiagnosticButtonDisableAllowed, "allowed", "disallowed")) +
		fmt.Sprintf("Reset Button Disable : %s\n", formatBool(res.ResetButtonDisableAllowed, "allowed", "disallowed")) +
		fmt.Sprintf("Power Button Disable : %s\n", formatBool(res.PoweroffButtonDisableAllowed, "allowed", "disallowed")) +
		fmt.Sprintf("Sleep Button Disabled: %v\n", res.SleepButtonDisabled) +
		fmt.Sprintf("Diag Button Disabled : %v\n", res.DiagnosticButtonDisabled) +
		fmt.Sprintf("Reset Button Disabled: %v\n", res.ResetButtonDisabled) +
		fmt.Sprintf("Power Button Disabled: %v\n", res.PoweroffButtonDisabled)
}

func (c *Client) GetChassisStatus(ctx context.Context) (response *GetChassisStatusResponse, err error) {
	request := &GetChassisStatusRequest{}
	response = &GetChassisStatusResponse{}
	err = c.Exchange(ctx, request, response)
	return
}
