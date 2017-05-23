package goscheduler

import (
	"fmt"
	"reflect"
)

// checkFunc will check that f is a function and that the the correct amount of parameters ...p are supplied and that they are of the correct type.
// Returns *reflect.Value, *[]reflect.Value and error.
func (job *Job) checkFunc(f *interface{}, p *[]interface{}) (*reflect.Value, *[]reflect.Value, error) {
	function := reflect.ValueOf(*f)

	// Check that f is a function.
	if function.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf("Function argument supplied needs to be of type %v but is of type %v. Job %v", reflect.Func, function.Kind(), job.name)
	}
	// Check to see that length of p is the same as the number of declared parameters in f.
	if len(*p) != function.Type().NumIn() {
		return nil, nil, fmt.Errorf("Invalid number of parameters supplied to function. Takes %v and %v was supplied. Job %v", function.Type().NumIn(), len(*p), job.name)
	}
	// Check to see that the parameters supplied in p are of the correct types as declared in f.
	parameters := []reflect.Value{}
	for num, param := range *p {
		if function.Type().In(num) != reflect.TypeOf(param) {
			return nil, nil, fmt.Errorf("Parameter %v is of type %v but should be of type %v. Job %v", num+1, reflect.TypeOf(param), function.Type().In(num), job.name)
		}
		parameters = append(parameters, reflect.ValueOf(param))
	}

	return &function, &parameters, nil
}
