package xmlrpc

const (
	ValueEmpty   = 1
	ValueInt     = 2
	ValueDouble  = 3
	ValueBoolean = 4
	ValueString  = 5
	ValueDate    = 6
	ValueBase64  = 7
	ValueArray   = 8
)

type Value struct {
	Type    int     //internal
	Number  int     //i4 i8 int
	Boolean bool    //boolean
	Double  float64 //double
	String  string  //string
	Date    string  //dateTime.iso8601
	Base64  string  //base64
	Array   []Value //array
}

func (v *Value) setType(rpcName string) {
	switch rpcName {
	case "int", "i4", "i8":
		v.Type = ValueInt
	case "double":
		v.Type = ValueDouble
	case "boolean":
		v.Type = ValueBoolean
	case "string":
		v.Type = ValueString
	case "dateTime.iso8601":
		v.Type = ValueDate
	case "base64":
		v.Type = ValueBase64
	case "array":
		v.Type = ValueArray
	default:
		v.Type = 0
	}
}
