package xray_tracer

import (
	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"time"
	"fmt"
	"crypto/rand"
	"encoding/json"
	"errors"
)

type TraceId string

type XRayTracerSetting struct {
	ServiceName string
	TraceId string

	//optional
	ParentId *string
	AwsConfig *aws.Config

	Annotations *map[string]interface{} //検索用index
	Metadata *map[string]interface{} //segmentの追加データ
}

func CreateXrayTraceInstance(setting XRayTracerSetting) XRayTrace {

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	awsConfig := setting.AwsConfig
	if awsConfig == nil {
		awsConfig = &aws.Config{
			Region: aws.String(`ap-northeast-1`),
		}
	}

	svc := xray.New(sess  , awsConfig)
	svc.ServiceName = setting.ServiceName

	startTime := time.Now().Unix()

	return &implXrayTrace {
		xRaySession : svc,
		XRayTracerSetting: setting,
		StartTime: startTime,
		Id : getRandom(8), //16桁
	}
}

func CreateNewTraceId() string {
	epoch := fmt.Sprintf("%x", time.Now().Unix()) //8桁の16進数
	return `1-` + epoch + `-` + getRandom(12) // 24桁の16進数
}

// getRandom generates a random hex encoded string
func getRandom(i int) string {
	b := make([]byte, i)
	for {
		// keep trying till we get it
		if _, err := rand.Read(b); err != nil {
			continue
		}
		return fmt.Sprintf("%x", b)
	}
}

type XRayTrace interface {
	CallOnSuccess() error
	CallOnFail(error) error
	GetId() string
}

type implXrayTrace struct {
	xRaySession *xray.XRay
	XRayTracerSetting XRayTracerSetting
	StartTime int64
	Id string
}

type XRayTraceBody struct {
	Id string `json:"id"`
	TraceId string `json:"trace_id"`
	StartTime int64 `json:"start_time"`
	EndTime int64 	`json:"end_time"`
	Name string `json:"name"`

	ParentId *string `json:"parent_id,omitempty"`
	Fault *bool `json:"fault,omitempty"` //5XX Client Error
	Error *bool `json:"error,omitempty"`//4XX Client Error
	Cause *XRayCause `json:"cause,omitempty"`

	Annotations *map[string]interface{} `json:"annotations,omitempty"`
	Metadata *map[string]interface{} `json:"metadata,omitempty"`
}

type XRayCause struct {
	Exceptions []XRayException `json:"exceptions"`
}

type XRayException struct {
	Message string `json:"message"`
}

func (instance *implXrayTrace) GetId() string {
	return instance.Id
}

func (instance *implXrayTrace) CallOnSuccess() error {

	xrayCliBody := XRayTraceBody{
		Id : instance.Id,
		TraceId: instance.XRayTracerSetting.TraceId,
		Name: instance.XRayTracerSetting.ServiceName,
		StartTime: instance.StartTime,
		EndTime: time.Now().Unix(),
		ParentId: instance.XRayTracerSetting.ParentId,
		Metadata:instance.XRayTracerSetting.Metadata,
		Annotations: instance.XRayTracerSetting.Annotations,
	}

	marshalBody  , err := json.Marshal(xrayCliBody)
	if err != nil {
		return err
	}

	s , err := instance.xRaySession.PutTraceSegments(&xray.PutTraceSegmentsInput{
		TraceSegmentDocuments: []*string{
			aws.String(string(marshalBody)),
		},
	})
	if err != nil {
		return err
	}
	if len(s.UnprocessedTraceSegments) > 0 {
		return errors.New(*s.UnprocessedTraceSegments[0].Message)
	}
	return nil
}

func (instance *implXrayTrace) CallOnFail(err error) error {

	cause := XRayCause{
		Exceptions : []XRayException{
			XRayException{
				Message: err.Error(),
			},
		},
	}

	xrayCliBody := XRayTraceBody{
		Id : instance.Id,
		TraceId: instance.XRayTracerSetting.TraceId,
		Name: instance.XRayTracerSetting.ServiceName,
		StartTime: instance.StartTime,
		EndTime: time.Now().Unix(),
		ParentId: instance.XRayTracerSetting.ParentId,
		Fault: aws.Bool(true),
		Cause: &cause,
		Metadata:instance.XRayTracerSetting.Metadata,
		Annotations: instance.XRayTracerSetting.Annotations,
	}

	marshalBody  , err := json.Marshal(xrayCliBody)
	if err != nil {
		return err
	}

	s , err := instance.xRaySession.PutTraceSegments(&xray.PutTraceSegmentsInput{
		TraceSegmentDocuments: []*string{
			aws.String(string(marshalBody)),
		},
	})
	if err != nil {
		return err
	}
	if len(s.UnprocessedTraceSegments) > 0 {
		return errors.New(*s.UnprocessedTraceSegments[0].Message)
	}
	return nil
}
