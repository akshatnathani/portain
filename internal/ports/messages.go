package ports

func PortMessage(available bool) string {
	if available {
		return "Port is available"
	}

	return "Port is occupied"
}