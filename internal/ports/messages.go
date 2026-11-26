package ports

func PortMessage(available bool) string {
	if available {
		return "Port is free and ready"
	}

	return "Port is currently in use"
}
