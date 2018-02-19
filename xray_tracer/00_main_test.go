package xray_tracer

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M){
	code := m.Run()
	os.Exit(code)
}


func TestTrace(t *testing.T) {

	traceId := CreateNewTraceId()

	xr := CreateXrayTraceInstance(XRayTracerSetting{
		ServiceName : `SaSS-1`,
		TraceId : traceId,
	})

	err := xr.CallOnSuccess()
	assert.NoError(t , err)

	parentId := xr.GetId()

	metadata := map[string]interface{}{
		`inputEvent` : parentId,
	}

	xr2 := CreateXrayTraceInstance(XRayTracerSetting{
		ServiceName : `SaSS-2`,
		TraceId : traceId,
		ParentId: &parentId,
		Metadata: &metadata,
	})

	err = xr2.CallOnSuccess()
	assert.NoError(t , err)

}
