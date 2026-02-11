package metric

import "github.com/prometheus/client_golang/prometheus"

var (
	ActivePlayers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_players",
		Help: "Aktif oyuncu sayısı",
	})
	TotalRegisters = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "mc_total_registers_total",
		Help: "Toplam kayıt olan oyuncu sayısı",
	})
)

func init() {
	prometheus.MustRegister(ActivePlayers, TotalRegisters)
}
