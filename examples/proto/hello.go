package proto

func (this *World) MarshalPacket() ([]byte, error) {
	return nil, nil
}

func (this *World) UnmarshalPacket([]byte) error {
	return nil
}

func (this *Hello) MarshalPacket() ([]byte, error) {
	return nil, nil
}

func (this *Hello) UnmarshalPacket([]byte) error {
	return nil
}
