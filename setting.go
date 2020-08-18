package ecsmeta

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lestrrat-go/backoff"
)

//Option provides a way to specify how to get the ECS metadata
type Option func(*setting)

type setting struct {
	client   *http.Client
	endpoint string
	policy   backoff.Policy
	logger   *log.Logger
}

func newSetting() *setting {
	var endpoint string
	if e := os.Getenv("ECS_CONTAINER_METADATA_URI"); e != "" {
		endpoint = e + "/task"
	}
	if e := os.Getenv("ECS_CONTAINER_METADATA_URI_V4"); e != "" {
		endpoint = e + "/task"
	}
	return &setting{
		endpoint: endpoint,
		client:   http.DefaultClient,
	}
}

func (s *setting) Logf(format string, args ...interface{}) {
	if s.logger == nil {
		return
	}
	s.logger.Printf(format, args...)
}

func (s *setting) Start(ctx context.Context) (backoff.Backoff, backoff.CancelFunc) {
	if s.policy != nil {
		return s.policy.Start(ctx)
	}
	defaultPolicy := backoff.NewExponential(
		backoff.WithInterval(500*time.Millisecond),
		backoff.WithJitterFactor(0.5),
		backoff.WithMaxRetries(5),
	)
	return defaultPolicy.Start(ctx)
}

//WithEndpoint specifies the Endpoint to get ECS metadata.
func WithEndpoint(endpoint string) Option {
	return func(s *setting) {
		s.endpoint = endpoint
	}
}

//WithEnableV2 enables automatic detection of V2 endpoints.
func WithEnableV2() Option {
	return func(s *setting) {
		if s.endpoint == "" {
			s.endpoint = "http://169.254.170.2/v2/metadata" //for v2
		}
	}
}

//WithHTTPClient specifies the HTTP client used to download ECS Metadata.
func WithHTTPClient(client *http.Client) Option {
	return func(s *setting) {
		s.client = client
	}
}

//WithRetryPolicy specifies the retry policy when ECS metadata download fails.
func WithRetryPolicy(policy backoff.Policy) Option {
	return func(s *setting) {
		s.policy = policy
	}
}

//WithLogger specifies the output method of Log
func WithLogger(l *log.Logger) Option {
	return func(s *setting) {
		s.logger = l
	}
}
