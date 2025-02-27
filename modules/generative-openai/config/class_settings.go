//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package config

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/moduletools"
)

const (
	modelProperty            = "model"
	temperatureProperty      = "temperature"
	maxTokensProperty        = "maxTokens"
	frequencyPenaltyProperty = "frequencyPenalty"
	presencePenaltyProperty  = "presencePenalty"
	topPProperty             = "topP"
)

var (
	DefaultOpenAIModel                    = "text-davinci-003"
	DefaultOpenAITemperature      float64 = 0.0
	DefaultOpenAIMaxTokens        float64 = 1200
	DefaultOpenAIFrequencyPenalty float64 = 0.0
	DefaultOpenAIPresencePenalty  float64 = 0.0
	DefaultOpenAITopP             float64 = 1.0
)

var maxTokensForModel = map[string]float64{
	"text-davinci-003": 1200,
}

type classSettings struct {
	cfg moduletools.ClassConfig
}

func NewClassSettings(cfg moduletools.ClassConfig) *classSettings {
	return &classSettings{cfg: cfg}
}

func (ic *classSettings) Validate(class *models.Class) error {
	if ic.cfg == nil {
		// we would receive a nil-config on cross-class requests, such as Explore{}
		return errors.New("empty config")
	}

	temperature := ic.getFloatProperty(temperatureProperty, &DefaultOpenAITemperature)
	if temperature == nil || (*temperature < 0 || *temperature > 1) {
		return errors.Errorf("Wrong temperature configuration, values are between 0.0 and 1.0")
	}

	maxTokens := ic.getFloatProperty(maxTokensProperty, &DefaultOpenAIMaxTokens)
	if maxTokens == nil || (*maxTokens < 0 || *maxTokens > getMaxTokensForModel(DefaultOpenAIModel)) {
		return errors.Errorf("Wrong maxTokens configuration, values are should have a minimal value of 1 and max is dependant on the model used")
	}

	frequencyPenalty := ic.getFloatProperty(frequencyPenaltyProperty, &DefaultOpenAIFrequencyPenalty)
	if frequencyPenalty == nil || (*frequencyPenalty < 0 || *frequencyPenalty > 1) {
		return errors.Errorf("Wrong frequencyPenalty configuration, values are between 0.0 and 1.0")
	}

	presencePenalty := ic.getFloatProperty(presencePenaltyProperty, &DefaultOpenAIPresencePenalty)
	if presencePenalty == nil || (*presencePenalty < 0 || *presencePenalty > 1) {
		return errors.Errorf("Wrong presencePenalty configuration, values are between 0.0 and 1.0")
	}

	topP := ic.getFloatProperty(topPProperty, &DefaultOpenAITopP)
	if topP == nil || (*topP < 0 || *topP > 5) {
		return errors.Errorf("Wrong topP configuration, values are should have a minimal value of 1 and max of 5")
	}

	return nil
}

func (ic *classSettings) getStringProperty(name, defaultValue string) *string {
	if ic.cfg == nil {
		// we would receive a nil-config on cross-class requests, such as Explore{}
		return &defaultValue
	}

	model, ok := ic.cfg.ClassByModuleName("generative-openai")[name]
	if ok {
		asString, ok := model.(string)
		if ok {
			return &asString
		}
		var empty string
		return &empty
	}
	return &defaultValue
}

func (ic *classSettings) getFloatProperty(name string, defaultValue *float64) *float64 {
	if ic.cfg == nil {
		// we would receive a nil-config on cross-class requests, such as Explore{}
		return defaultValue
	}

	val, ok := ic.cfg.ClassByModuleName("generative-openai")[name]
	if ok {
		asFloat, ok := val.(float64)
		if ok {
			return &asFloat
		}
		asNumber, ok := val.(json.Number)
		if ok {
			asFloat, _ := asNumber.Float64()
			return &asFloat
		}
		asInt, ok := val.(int)
		if ok {
			asFloat := float64(asInt)
			return &asFloat
		}
		var wrongVal float64 = -1.0
		return &wrongVal
	}

	if defaultValue != nil {
		return defaultValue
	}
	return nil
}

func getMaxTokensForModel(model string) float64 {
	return maxTokensForModel[model]
}

func (ic *classSettings) Model() string {
	return *ic.getStringProperty(modelProperty, DefaultOpenAIModel)
}

func (ic *classSettings) MaxTokens() float64 {
	return *ic.getFloatProperty(maxTokensProperty, &DefaultOpenAIMaxTokens)
}

func (ic *classSettings) Temperature() float64 {
	return *ic.getFloatProperty(temperatureProperty, &DefaultOpenAITemperature)
}

func (ic *classSettings) FrequencyPenalty() float64 {
	return *ic.getFloatProperty(frequencyPenaltyProperty, &DefaultOpenAIFrequencyPenalty)
}

func (ic *classSettings) PresencePenalty() float64 {
	return *ic.getFloatProperty(presencePenaltyProperty, &DefaultOpenAIPresencePenalty)
}

func (ic *classSettings) TopP() float64 {
	return *ic.getFloatProperty(topPProperty, &DefaultOpenAITopP)
}
