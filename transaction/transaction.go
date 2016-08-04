package transaction

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Transaction represents a Dredd transaction object.
// http://dredd.readthedocs.io/en/latest/data-structures/#transaction
type Transaction struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	FullPath string `json:"fullPath,omitempty"`
	Request  *struct {
		Body    string  `json:"body,omitempty"`
		Headers Headers `json:"headers,omitempty"`
		URI     string  `json:"uri,omitempty"`
		Method  string  `json:"method,omitempty"`
	} `json:"request,omitempty"`
	Expected *struct {
		StatusCode string                 `json:"statusCode,omitempty"`
		Body       string                 `json:"body,omitempty"`
		Headers    map[string]interface{} `json:"headers,omitempty"`
		Schema     *json.RawMessage       `json:"bodySchema,omitempty"`
	} `json:"expected,omitempty"`
	Real     *struct {
		Body       string  `json:"body"`
		Headers    Headers `json:"headers"`
		StatusCode int     `json:"statusCode"`
	} `json:"real,omitempty"`
	Origin  *json.RawMessage `json:"origin,omitempty"`
	Test    *json.RawMessage `json:"test,omitempty"`
	Results *json.RawMessage `json:"results,omitempty"`
	Skip    bool             `json:"skip"`
	Fail    interface{}      `json:"fail,omitempty"`

	TestOrder []string `json:"hooks_modifications,omitempty"`
}

type Headers map[string][]string

func (hdrs *Headers) UnmarshalJSON(data []byte) error {
	var headers = make(map[string]interface{})
	err := json.Unmarshal(data, &headers)

	if err != nil {
		return err
	}

	*hdrs = make(map[string][]string)
	for k := range headers {
		header := headers[k]
		value := reflect.ValueOf(header)
		kind := value.Kind()

		switch kind {
		case reflect.Float32, reflect.Float64, reflect.Int, reflect.String:
			(*hdrs)[k] = []string{
				stringFromValue(value),
			}
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				(*hdrs)[k] = append((*hdrs)[k], stringFromValueAtIndex(value, i))
			}

		default:
			// TODO: Figure out how to use this
			panic(fmt.Sprintf("unsupported Kind %s", kind.String()))
		}

	}

	return nil
}

func stringFromValueAtIndex(val reflect.Value, index int) string {
	return fmt.Sprint(val.Index(index).Interface())
}

func stringFromValue(val reflect.Value) string {
	return fmt.Sprint(val.Interface())
}

// func (t *Transaction) UnmarshalJSON(data []byte) error {
// 	fmt.Printf("%v", string(data))
// 	return json.Unmarshal(data, Transaction{})
// }

// AddTestOrderPoint adds a value to the hooks_modification key used when
// running dredd with TEST_DREDD_HOOKS_HANDLER_ORDER enabled.
func (t *Transaction) AddTestOrderPoint(value string) {
	if t.TestOrder == nil {
		t.TestOrder = []string{}
	}
	t.TestOrder = append(t.TestOrder, value)
}
