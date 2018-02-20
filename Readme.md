# xray-simple-tracer

simple wrapper for AWS x-ray by golang.  
you need to add trace_id , parent_id manually. 

```sh

glide install github.com/wano/xray-simple-tracer

# or

dep ensure -add github.com/wano/xray-simple-tracer

```

```go

traceId := xray_tracer.CreateNewTraceId()

saas1 := CreateTracer(xray_tracer.XRayTracerSetting{
	ServiceName : `SaaS-1`,
	TraceId : traceId,
})

err := saas1.Success()

parentId := saas1.GetId()

metadata := map[string]interface{}{
	`inputEvent` : parentId,
}

saas2 := CreateTracer(xray_tracer.XRayTracerSetting{
	ServiceName : `SaaS-2`,
	TraceId : traceId,
	ParentId: &parentId,
	Metadata: &metadata,
})

err = saas2.Success()

// on Fail
err = saas2.Fault(errors.New(`fail`))
err = saas2.Warn(errors.New(`fail`))


```

![image](https://user-images.githubusercontent.com/1452731/36366079-d405ea20-158f-11e8-84bc-a0aa12c08197.png)

