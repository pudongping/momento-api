package validator

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/thedevsaddam/govalidator"
)

func init() {
	// number_min:0.01
	govalidator.AddCustomRule("number_min", func(field string, rule string, message string, value interface{}) error {
		// Parse the rule parameter
		minStr := strings.TrimPrefix(rule, "number_min:")
		minVal, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return errors.Wrap(err, "number_min 规则参数无效")
		}

		// Get the value
		var floatVal float64
		switch v := value.(type) {
		case float64:
			floatVal = v
		case float32:
			floatVal = float64(v)
		case int:
			floatVal = float64(v)
		case int64:
			floatVal = float64(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				// If not a valid float string, we skip this check and let other type checks (like "float") fail if present
				return nil
			}
			floatVal = f
		default:
			// For other types, we can't compare, so we assume valid or let other rules handle type checking
			return nil
		}

		if floatVal < minVal {
			if message != "" {
				return errors.New(message)
			}
			return errors.Errorf("%s 字段必须至少为 %v", field, minVal)
		}

		return nil
	})

	// number_max:100
	govalidator.AddCustomRule("number_max", func(field string, rule string, message string, value interface{}) error {
		// Parse the rule parameter
		maxStr := strings.TrimPrefix(rule, "number_max:")
		maxVal, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return errors.Wrap(err, "number_max 规则参数无效")
		}

		// Get the value
		var floatVal float64
		switch v := value.(type) {
		case float64:
			floatVal = v
		case float32:
			floatVal = float64(v)
		case int:
			floatVal = float64(v)
		case int64:
			floatVal = float64(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil
			}
			floatVal = f
		default:
			return nil
		}

		if floatVal > maxVal {
			if message != "" {
				return errors.New(message)
			}
			return errors.Errorf("%s 字段必须最多为 %v", field, maxVal)
		}

		return nil
	})
}
