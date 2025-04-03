package api

func (h *handler) getPublicIP() (ip string, err error) {
	body, err := h.installer.Fetch("https://ipv4.ident.me")
	if err != nil {
		return
	}

	ip = string(body)

	return
}
