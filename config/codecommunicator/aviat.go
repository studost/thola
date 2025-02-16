package codecommunicator

import (
	"context"
	"github.com/inexio/thola/internal/device"
	"github.com/inexio/thola/internal/deviceclass/groupproperty"
	"github.com/inexio/thola/internal/network"
	"github.com/pkg/errors"
	"strings"
)

type aviatCommunicator struct {
	codeCommunicator
}

func (c *aviatCommunicator) GetInterfaces(ctx context.Context, filter ...groupproperty.Filter) ([]device.Interface, error) {
	interfaces, err := c.deviceClass.GetInterfaces(ctx, filter...)
	if err != nil {
		return nil, err
	}

	con, ok := network.DeviceConnectionFromContext(ctx)
	if !ok || con.SNMP == nil {
		return nil, errors.New("snmp client is empty")
	}

	// aviatModemStatusMaxCapacity
	res, err := con.SNMP.SnmpClient.SNMPWalk(ctx, "1.3.6.1.4.1.2509.9.3.2.4.1.1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aviatModemStatusMaxCapacity")
	}

	var maxCapacity uint64
	for _, r := range res {
		capacityVal, err := r.GetValue()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get aviatModemStatusMaxCapacity value")
		}
		capacity, err := capacityVal.UInt64()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse aviatModemStatusMaxCapacity value")
		}
		maxCapacity += capacity
	}

	// aviatModemCurCapacityTx
	res, err = con.SNMP.SnmpClient.SNMPWalk(ctx, "1.3.6.1.4.1.2509.9.3.2.1.1.11")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aviatModemCurCapacityTx")
	}

	var maxBitRateTx uint64
	for _, r := range res {
		bitRateVal, err := r.GetValue()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get aviatModemCurCapacityTx value")
		}
		bitRate, err := bitRateVal.UInt64()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse aviatModemCurCapacityTx value")
		}
		maxBitRateTx += bitRate * 1000
	}

	// aviatModemCurCapacityRx
	res, err = con.SNMP.SnmpClient.SNMPWalk(ctx, "1.3.6.1.4.1.2509.9.3.2.1.1.12")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aviatModemCurCapacityRx")
	}

	var maxBitRateRx uint64
	for _, r := range res {
		bitRateVal, err := r.GetValue()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get aviatModemCurCapacityRx value")
		}
		bitRate, err := bitRateVal.UInt64()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse aviatModemCurCapacityRx value")
		}
		maxBitRateRx += bitRate * 1000
	}

	for i, interf := range interfaces {
		if interf.IfName != nil && strings.HasPrefix(*interf.IfName, "Radio") {
			interfaces[i].MaxSpeedIn = &maxCapacity
			interfaces[i].MaxSpeedOut = &maxCapacity
			interfaces[i].Radio = &device.RadioInterface{
				MaxbitrateOut: &maxBitRateTx,
				MaxbitrateIn:  &maxBitRateRx,
			}
		}
	}

	return interfaces, nil
}
