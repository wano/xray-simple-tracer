# xray-simple-tracer

simple wrapper for AWS x-ray by golang.
you need add trace_id , parent_id manually. 

```sh

glide install github.com/wano/xray-simple-tracer
# or
dep enduore -add github.com/wano/xray-simple-tracer

```

```go
traceId := xray_tracer.CreateNewTraceId()

sass1Task := xray_tracer.CreateXrayTraceInstance(xray_tracer.XRayTracerSetting{
	ServiceName : `SaSS-1`,
	TraceId : traceId,
})

err := sass1Task.CallOnSuccess()

parentId := sass1Task.GetId()

metadata := map[string]interface{}{
	`inputEvent` : parentId,
}

sass3Task := xray_tracer.CreateXrayTraceInstance(xray_tracer.XRayTracerSetting{
	ServiceName : `SaSS-2`,
	TraceId : traceId,
	ParentId: &parentId,
	Metadata: &metadata,
})

err := sass2Task.CallOnSuccess()

```

![image](https://user-images.githubusercontent.com/1452731/36366079-d405ea20-158f-11e8-84bc-a0aa12c08197.png)

