package ecsmeta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/itchyny/gojq"
	"github.com/lestrrat-go/backoff"
	"github.com/pkg/errors"
)

//Object is ECS metadata object
type Object struct {
	Value interface{}
}

//New gets the ECS metadata. The acquisition method changes depending on the option.
// if return ECSMeta(nil), failed get ECS metadata
func New(opts ...Option) *Object {
	s := newSetting()
	for _, opt := range opts {
		opt(s)
	}
	return &Object{Value: getMetadata(s)}
}

func getMetadata(s *setting) interface{} {
	if s.endpoint == "" {
		s.Logf("endpoint is not set")
		return nil
	}

	b, cancel := s.Start(context.Background())
	defer cancel()
	for i := 1; backoff.Continue(b); i++ {
		metadata, err := getMetadataOnce(s)
		if err == nil {
			return metadata
		}
		s.Logf("[%d]: unable to get ecs metadata response: %v", i, err)
	}
	s.Logf("max retries count reached")
	return nil
}

func getMetadataOnce(s *setting) (interface{}, error) {
	resp, err := s.client.Get(s.endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var metadata interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&metadata); err != nil {
		return nil, errors.Wrap(err, "unable to decode response body")
	}
	return metadata, nil
}

// Query queries object by go-jq
func (o *Object) Query(query string) (*Object, error) {
	v, err := o.query(query)
	return &Object{Value: v}, err
}

func (o *Object) query(query string) (interface{}, error) {
	if o.Value == nil {
		return nil, errors.New("object value is nil")
	}
	jq, err := gojq.Parse(query)
	if err != nil {
		return nil, err
	}
	iter := jq.Run(o.Value)
	vs := make([]interface{}, 0, 2)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		vs = append(vs, v)
	}
	switch len(vs) {
	case 0:
		return nil, errors.Errorf("%s is not found in the state", query)
	case 1:
		return vs[0], nil
	}
	return vs, nil
}

//Bytes returns a Byte array representing the metadata object
func (o Object) Bytes() []byte {
	switch v := (o.Value).(type) {
	case string:
		return []byte(v)
	default:
		b, _ := json.Marshal(v)
		return b
	}
}

//String returns a json string
func (o Object) String() string {
	return string(o.Bytes())
}
