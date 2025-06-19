package shared_type

type ProviderInterface struct {
	URL                 string             `json:"url" valdate:"required,url"`
	AuthStrategy        string             `json:"authStrategy" validate:"required"`
	RequestMethod       string             `json:"requestMethod" validate:"required,oneof=GET POST PUT DELETE"`
	RequestContentType  string             `json:"requestContentType" validate:"required"`
	ResponseContentType string             `json:"responseContentType" validate:"required"`
	RequestInterface    []InterfaceElement `json:"requestInterface" validate:"required,min=1,dive"`
	ResponseInterface   []InterfaceElement `json:"responseInterface" validate:"required,min=1,dive"`
}

type InterfaceElement struct {
	ID                string            `json:"id" validate:"required"`
	Type              string            `json:"type" validate:"required,oneof=body query header"`
	Required          bool              `json:"required"`
	Key               string            `json:"key" validate:"required"`
	ValueType         string            `json:"valueType" validate:"required,oneof=string number boolean"`
	BindedElementType BindedElementType `json:"bindedElementType" validate:"required"`
}

type BindedElementType struct {
	Label           string `json:"label" validate:"required"`
	HTMLElementType string `json:"htmlElementType" validate:"required"`
	ValueType       string `json:"valueType" validate:"required,oneof=string number boolean"`
}
