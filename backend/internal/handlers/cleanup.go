package handlers

func (cfg *App) Cleanup() {
	cfg.deleteExpiredJTI()
	cfg.deleteExpiredRefresh()
}
