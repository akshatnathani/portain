package ports

func PortMessage(available bool) string {
	if available {
		return "Port is available and ready"
	}

	return "Port is occupied"
}
