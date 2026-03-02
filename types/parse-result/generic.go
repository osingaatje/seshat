package parseresult

// Properties of values, e.g.: "field" : Visibility=private, Type=string
type ValueProperties struct {
	Visibility ValuePropVisibilityVar `json:"visibility,omitempty"`
	Type       string                 `json:"type,omitempty"` // type string,class,bool,etc.
}

// More detailed datatypes
type ValuePropVisibilityVar string

const (
	VisibilityPublic    ValuePropVisibilityVar = "public"
	VisibilityProtected ValuePropVisibilityVar = "protected"
	VisibilityPrivate   ValuePropVisibilityVar = "private"
	VisibilityUnknown   ValuePropVisibilityVar = "unknown"
)
