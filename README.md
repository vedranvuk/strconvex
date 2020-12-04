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
func StringToInterface(in string, out interface{})
func StringToValue(in string, out reflect.Value)
func StringToBoolValue(in string, out reflect.Value)
func StringToIntValue(in string, out reflect.Value)
func StringToUintValue(in string, out reflect.Value)
func StringToFloat32Value(in string, out reflect.Value)
func StringToFloat64Value(in string, out reflect.Value)
func StringToComplex64Value(in string, out reflect.Value)
func StringToComplex128Value(in string, out reflect.Value)
func StringToStringValue(in string, out reflect.Value)
func StringToArrayValue(in string, out reflect.Value)
func StringToSliceValue(in string, out reflect.Value)
func StringToMapValue(in string, out reflect.Value)
func StringToStructValue(in string, out reflect.Value)
func StringToPointerValue(in string, out reflect.Value)
```
## License

MIT. See included LICENSE file.