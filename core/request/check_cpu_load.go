package request

import "context"

// CheckCPULoadRequest
//
// CheckCPULoadRequest is a the request struct for the check cpu load request.
//
// swagger:model
type CheckCPULoadRequest struct {
	CheckDeviceRequest
	CPULoadThresholds CheckThresholds `json:"cpuLoadThresholds" xml:"cpuLoadThresholds"`
}

func (r *CheckCPULoadRequest) validate(ctx context.Context) error {
	if err := r.CPULoadThresholds.validate(); err != nil {
		return err
	}
	return r.CheckDeviceRequest.validate(ctx)
}
