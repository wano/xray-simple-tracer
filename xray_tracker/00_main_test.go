package xray_tracker

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestMain(m *testing.M){
	code := m.Run()
	os.Exit(code)
}


func TestTrace(t *testing.T) {

	traceId := CreateNewTraceId()

	saas1 := CreateTracer(XRayTracerSetting{
		ServiceName : `SaaS-1`,
		TraceId : traceId,
	})

	err := saas1.Success()
	assert.NoError(t , err)

	parentId := saas1.GetId()

	metadata := map[string]interface{}{
		`inputEvent` : parentId,
	}

	saas2 := CreateTracer(XRayTracerSetting{
		ServiceName : `SaaS-2`,
		TraceId : traceId,
		ParentId: &parentId,
		Metadata: &metadata,
	})

	err = saas2.Success()
	assert.NoError(t , err)
	err = saas2.Fail(errors.New(`fail`))
	assert.NoError(t , err)

}
