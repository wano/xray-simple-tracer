# xray-simple-tracer

simple wrapper for AWS x-ray by golang.
you need add trace_id , parent_id manually. 

```go

glide install github.com/wano/xray-simple-tracer
# or
dep enduore -add github.com/wano/xray-simple-tracer

```

```go
	traceId := CreateNewTraceId()

	xr := CreateXrayTraceInstance(XRayTracerSetting{
		ServiceName : `SaSS-1`,
		TraceId : traceId,
	})

	err := xr.CallOnSuccess()

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
```


