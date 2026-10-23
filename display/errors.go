package display

type DeviceNotFoundError struct {
}

func (e *DeviceNotFoundError) Error() string {
	return "display interface not found"
}
