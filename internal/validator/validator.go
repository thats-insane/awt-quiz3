package validator

/* Validator struct
 */
type Validator struct {
	Errors map[string]string
}

/* Create a new validator
 */
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

/* Returns whether or not validator is empty
 */
func (v *Validator) IsEmpty() bool {
	return len(v.Errors) == 0
}

/* Adds a validator to the errors map
 */
func (v *Validator) AddError(key string, message string) {
	_, exists := v.Errors[key]

	if !exists {
		v.Errors[key] = message
	}
}

/* Checks the incoming data to see if its valid
 */
func (v *Validator) Check(acceptable bool, key string, message string) {
	if !acceptable {
		v.AddError(key, message)
	}
}
