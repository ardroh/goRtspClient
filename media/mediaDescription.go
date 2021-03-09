package media

type MediaDescription struct {
	Media          string
	Port           MediaPort
	Proto          string
	Fmt            int
	Attributes     []Attribute
	ConnectionData *ConnectionData
	Bandwidth      *Bandwidth
}
