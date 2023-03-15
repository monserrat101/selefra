package analytics

import (
	"context"
	"github.com/rudderlabs/analytics-go/v4"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-utils/pkg/id_util"
	"github.com/selefra/selefra/pkg/selefra_workspace"
	"sync"
	"time"
)

// ------------------------------------------------ ---------------------------------------------------------------------

// Analytics Represents an interface for analysis
type Analytics interface {

	// Init Initialization analyzer
	Init(ctx context.Context) *schema.Diagnostics

	// Submit the information to be collected
	Submit(ctx context.Context, event *Event) *schema.Diagnostics

	// Close Turn off analyzer
	Close(ctx context.Context) *schema.Diagnostics
}

// ------------------------------------------------ ---------------------------------------------------------------------

type Event struct {
	Name    string
	Payload string
}

func NewEvent(name, payload string) *Event {
	return &Event{
		Name:    name,
		Payload: payload,
	}
}

func (x *Event) SetName(name string) *Event {
	x.Name = name
	return x
}

func (x *Event) SetPayload(payload string) *Event {
	x.Payload = payload
	return x
}

// ------------------------------------------------ ---------------------------------------------------------------------

type RudderstackAnalytics struct {
	client   analytics.Client
	deviceId string
}

var _ Analytics = &RudderstackAnalytics{}

func (x *RudderstackAnalytics) Init(ctx context.Context) *schema.Diagnostics {
	// Instantiates a client to use send messages to the Rudder API.
	// TODO Wait to confirm whether the write key can be exposed
	client := analytics.New("xxxxx", "https://selefralefsm.dataplane.rudderstack.com")

	x.client = client

	deviceId, diagnostics := selefra_workspace.GetDeviceID()
	x.deviceId = deviceId

	return diagnostics
}

func (x *RudderstackAnalytics) Submit(ctx context.Context, event *Event) *schema.Diagnostics {
	err := x.client.Enqueue(analytics.Track{
		AnonymousId: x.deviceId,
		// Don't track specific users, just collect some basic usage statistics
		UserId:    "",
		MessageId: id_util.RandomId(),
		Event:     event.Name,
		Timestamp: time.Now(),
	})
	return schema.NewDiagnostics().AddError(err)
}

func (x *RudderstackAnalytics) Close(ctx context.Context) *schema.Diagnostics {
	if x.client != nil {
		err := x.client.Close()
		if err != nil {
			return schema.NewDiagnostics().AddErrorMsg("close Rudderstack client failed: %s", err.Error())
		}
	}
	return nil
}

// ------------------------------------------------ ---------------------------------------------------------------------

var DefaultAnalytics Analytics
var InitOnce sync.Once

func Init(ctx context.Context) *schema.Diagnostics {
	DefaultAnalytics = &RudderstackAnalytics{}
	return DefaultAnalytics.Init(ctx)
}

func Submit(ctx context.Context, event *Event) *schema.Diagnostics {

	InitOnce.Do(func() {
		_ = Init(context.Background())
	})

	return DefaultAnalytics.Submit(ctx, event)
}

func Close(ctx context.Context) *schema.Diagnostics {
	if DefaultAnalytics != nil {
		return DefaultAnalytics.Close(ctx)
	} else {
		return nil
	}
}
