/*
Copyright 2020 The Knative Authors

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

package v1

import (
	"context"

	"github.com/google/go-cmp/cmp/cmpopts"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/kmp"

	"knative.dev/eventing/pkg/apis/config"
)

const (
	BrokerClassAnnotationKey = "eventing.knative.dev/broker.class"
)

func (b *Broker) Validate(ctx context.Context) *apis.FieldError {
	ctx = apis.WithinParent(ctx, b.ObjectMeta)

	cfg := config.FromContextOrDefaults(ctx)
	var brConfig *config.ClassAndBrokerConfig
	if cfg.Defaults != nil {
		if c, ok := cfg.Defaults.NamespaceDefaultsConfig[b.GetNamespace()]; ok {
			brConfig = c
		} else {
			brConfig = cfg.Defaults.ClusterDefault
		}
	}

	withNS := ctx
	if brConfig == nil || brConfig.DisallowDifferentNamespaceConfig == nil || !*brConfig.DisallowDifferentNamespaceConfig {
		withNS = apis.AllowDifferentNamespace(ctx)
	}

	// Make sure a BrokerClassAnnotation exists
	var errs *apis.FieldError
	if bc, ok := b.GetAnnotations()[BrokerClassAnnotationKey]; !ok || bc == "" {
		errs = errs.Also(apis.ErrMissingField(BrokerClassAnnotationKey))
	}

	errs = errs.Also(b.Spec.Validate(withNS).ViaField("spec"))
	if apis.IsInUpdate(ctx) {
		original := apis.GetBaseline(ctx).(*Broker)
		errs = errs.Also(b.CheckImmutableFields(ctx, original))
	}
	return errs
}

func (bs *BrokerSpec) Validate(ctx context.Context) *apis.FieldError {
	var errs *apis.FieldError

	// Validate the Config
	if bs.Config != nil {
		if ce := bs.Config.Validate(ctx); ce != nil {
			errs = errs.Also(ce.ViaField("config"))
		}
	}

	if bs.Delivery != nil {
		if de := bs.Delivery.Validate(ctx); de != nil {
			errs = errs.Also(de.ViaField("delivery"))
		}
	}
	return errs
}

func (b *Broker) CheckImmutableFields(ctx context.Context, original *Broker) *apis.FieldError {
	if original == nil {
		return nil
	}

	// Only Delivery options are mutable.
	ignoreArguments := cmpopts.IgnoreFields(BrokerSpec{}, "Delivery")
	if diff, err := kmp.ShortDiff(original.Spec, b.Spec, ignoreArguments); err != nil {
		return &apis.FieldError{
			Message: "Failed to diff Broker",
			Paths:   []string{"spec"},
			Details: err.Error(),
		}
	} else if diff != "" {
		return &apis.FieldError{
			Message: "Immutable fields changed (-old +new)",
			Paths:   []string{"spec"},
			Details: diff,
		}
	}

	// Make sure you can't change the class annotation.
	if diff, _ := kmp.ShortDiff(original.GetAnnotations()[BrokerClassAnnotationKey], b.GetAnnotations()[BrokerClassAnnotationKey]); diff != "" {
		return &apis.FieldError{
			Message: "Immutable annotations changed (-old +new)",
			Paths:   []string{"annotations"},
			Details: diff,
		}
	}

	return nil
}
