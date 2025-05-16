package metric

import (
	"encoding/json"
	"fmt"

	"github.com/result17/codeBeatCli/internal/summary"
)

type (
	// TODO generate metric
	MetricRatio[T string | uint32] struct {
		Value        T       `json:"value"`
		Duration     uint64  `json:"duration"`
		Ratio        float64 `json:"ratio"`
		DurationText string  `json:"durationText"`
	}

	MetricRatioData[T string | uint32] struct {
		GrandTotal summary.GrandTotal `json:"grandTotal"`
		Ratios     []MetricRatio[T]   `json:"ratios"`
		Metric     string             `json:"metric"`
	}
)

// custom MarshalJSON
func (m *MetricRatioData[T]) MarshalJSON() ([]byte, error) {
	type Alias MetricRatioData[T]
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// custom UnmarshalJSON
func (m *MetricRatioData[T]) UnmarshalJSON(data []byte) error {
	var raw struct {
		GrandTotal summary.GrandTotal `json:"grandTotal"`
		Ratios     []json.RawMessage  `json:"ratios"`
		Metric     string             `json:"metric"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	m.GrandTotal = raw.GrandTotal
	m.Metric = raw.Metric
	m.Ratios = make([]MetricRatio[T], len(raw.Ratios))

	for i, item := range raw.Ratios {
		var ratio MetricRatio[T]
		var value interface{}

		if err := json.Unmarshal(item, &value); err != nil {
			return err
		}

		// 直接解析到临时结构体获取value字段
		var tmp struct {
			Value json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal(item, &tmp); err != nil {
			return err
		}

		// 根据泛型类型T处理value字段
		var zero T
		switch any(zero).(type) {
		case string:
			var strVal string
			if err := json.Unmarshal(tmp.Value, &strVal); err != nil {
				return err
			}
			ratio.Value = any(strVal).(T)
		case uint32:
			var numVal uint32
			if err := json.Unmarshal(tmp.Value, &numVal); err != nil {
				return err
			}
			ratio.Value = any(numVal).(T)
		default:
			return fmt.Errorf("unsupported generic type for MetricRatio value")
		}

		if err := json.Unmarshal(item, &struct {
			Duration     *uint64
			Ratio        *float64
			DurationText *string
		}{
			Duration:     &ratio.Duration,
			Ratio:        &ratio.Ratio,
			DurationText: &ratio.DurationText,
		}); err != nil {
			return err
		}

		m.Ratios[i] = ratio
	}

	return nil
}
