package types

type Image struct {
	Thumb    string `json:"thumb,omitempty" mapstructure:"thumb"`
	Original string `json:"original,omitempty" mapstructure:"original"`
	Ext      string `json:"ext,omitempty" mapstructure:"ext"`
	Width    uint   `json:"width,omitempty" mapstructure:"width"`
	Height   uint   `json:"height,omitempty" mapstructure:"height"`
	Size     uint   `json:"size,omitempty" mapstructure:"size"`
}
