/*
Copyright 2020 The Knative Authors. All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/utils/pointer"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

func TestDeliverySpecValidation(t *testing.T) {
	invalidString := "invalid time"
	bop := BackoffPolicyExponential
	validDuration := "PT2S"
	invalidDuration := "1985-04-12T23:20:50.52Z"
	tests := []struct {
		name string
		spec *DeliverySpec
		want *apis.FieldError
	}{{
		name: "nil is valid",
		spec: nil,
		want: nil,
	}, {
		name: "invalid time format",
		spec: &DeliverySpec{BackoffDelay: &invalidString},
		want: func() *apis.FieldError {
			return apis.ErrInvalidValue(invalidString, "backoffDelay")
		}(),
	}, {
		name: "invalid deadLetterSink",
		spec: &DeliverySpec{DeadLetterSink: &duckv1.Destination{}},
		want: func() *apis.FieldError {
			return apis.ErrGeneric("expected at least one, got none", "ref", "uri").ViaField("deadLetterSink")
		}(),
	}, {
		name: "valid timeout",
		spec: &DeliverySpec{Timeout: &validDuration},
		want: nil,
	}, {
		name: "invalid timeout",
		spec: &DeliverySpec{Timeout: &invalidDuration},
		want: func() *apis.FieldError {
			return apis.ErrInvalidValue(invalidDuration, "timeout")
		}(),
	}, {
		name: "valid backoffPolicy",
		spec: &DeliverySpec{BackoffPolicy: &bop},
		want: nil,
	}, {
		name: "valid backoffDelay",
		spec: &DeliverySpec{BackoffDelay: &validDuration},
		want: nil,
	}, {
		name: "invalid backoffDelay",
		spec: &DeliverySpec{BackoffDelay: &invalidDuration},
		want: func() *apis.FieldError {
			return apis.ErrInvalidValue(invalidDuration, "backoffDelay")
		}(),
	}, {
		name: "negative retry",
		spec: &DeliverySpec{Retry: pointer.Int32(-1)},
		want: func() *apis.FieldError {
			return apis.ErrInvalidValue("-1", "retry")
		}(),
	}, {
		name: "valid retry 0",
		spec: &DeliverySpec{Retry: pointer.Int32(0)},
	}, {
		name: "valid retry 1",
		spec: &DeliverySpec{Retry: pointer.Int32(1)},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.spec.Validate(context.TODO())
			if diff := cmp.Diff(test.want.Error(), got.Error()); diff != "" {
				t.Error("DeliverySpec.Validate (-want, +got) =", diff)
			}
		})
	}
}
