# strconvex

Package strconvex is an extension to strconv package that implements a common
interface for converting text to simple Go values using reflect to guess
type information and strconv for actual conversion.

As in strconv, only simple Go types are supported with a few minor helpful
additions that help with compound types but have the limitation that only the
first level is parsed and their elements or fields must be simple types.

As input, standard GoValue format from the fmt package is understood.

The api consists of the following:

```Go
func StringToInterface(in string, out interface{}) error
func StringToValue(in string, out reflect.Value) error
func StringToBoolValue(in string, out reflect.Value) error
func StringToIntValue(in string, out reflect.Value) error
func StringToUintValue(in string, out reflect.Value) error
func StringToFloat32Value(in string, out reflect.Value) error
func StringToFloat64Value(in string, out reflect.Value) error
func StringToComplex64Value(in string, out reflect.Value) error
func StringToComplex128Value(in string, out reflect.Value) error
func StringToStringValue(in string, out reflect.Value) error
func StringToArrayValue(in string, out reflect.Value) error
func StringToSliceValue(in string, out reflect.Value) error
func StringToMapValue(in string, out reflect.Value) error
func StringToStructValue(in string, out reflect.Value) error
func StringToPointerValue(in string, out reflect.Value) error
```
## License

MIT. See included LICENSE file.