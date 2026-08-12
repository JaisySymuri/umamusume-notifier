package metrics

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	once                  sync.Once
	registry              *prometheus.Registry
	commandTotal          *prometheus.CounterVec
	remindersScheduled    prometheus.Counter
	remindersSent         *prometheus.CounterVec
	reminderDeliveryDelay prometheus.Histogram
	telegramRequests      *prometheus.CounterVec
	telegramDuration      *prometheus.HistogramVec
	telegramErrors        *prometheus.CounterVec
	storageDuration       *prometheus.HistogramVec
	fullOverHourTotal     prometheus.Counter
)

func initMetrics() {
	registry = prometheus.NewRegistry()

	commandTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_commands_total",
			Help: "Total bot commands by command and outcome.",
		},
		[]string{"command", "outcome"},
	)

	remindersScheduled = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_reminders_scheduled_total",
			Help: "Total scheduled reminder events.",
		},
	)

	remindersSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_reminders_sent_total",
			Help: "Total reminder delivery attempts by outcome.",
		},
		[]string{"outcome"},
	)

	reminderDeliveryDelay = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "bot_reminder_delivery_delay_seconds",
			Help:    "Time between reminder scheduling and successful delivery.",
			Buckets: prometheus.DefBuckets,
		},
	)

	telegramRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_telegram_api_requests_total",
			Help: "Total Telegram API requests by method and outcome.",
		},
		[]string{"method", "outcome"},
	)

	telegramDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bot_telegram_api_duration_seconds",
			Help:    "Telegram API request duration by method.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	telegramErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_telegram_api_errors_total",
			Help: "Total Telegram API request errors by method and error type.",
		},
		[]string{"method", "error_type"},
	)

	storageDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bot_storage_op_duration_seconds",
			Help:    "SQLite operation duration by storage method.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"op"},
	)

	fullOverHourTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_full_over_hour_total",
			Help: "Total point systems that remained full for more than one hour.",
		},
	)

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			PidFn: func() (int, error) {
				return os.Getpid(), nil
			},
		}),
		commandTotal,
		remindersScheduled,
		remindersSent,
		reminderDeliveryDelay,
		telegramRequests,
		telegramDuration,
		telegramErrors,
		storageDuration,
		fullOverHourTotal,
	)
}

func ensureMetrics() {
	once.Do(initMetrics)
}

func ObserveCommand(command, outcome string) {
	ensureMetrics()
	commandTotal.WithLabelValues(command, outcome).Inc()
}

func ObserveReminderScheduled() {
	ensureMetrics()
	remindersScheduled.Inc()
}

func ObserveReminderSent(outcome string) {
	ensureMetrics()
	remindersSent.WithLabelValues(outcome).Inc()
}

func ObserveReminderDeliveryDelay(d time.Duration) {
	ensureMetrics()
	reminderDeliveryDelay.Observe(d.Seconds())
}

func ObserveTelegramAPIRequest(method, outcome string, d time.Duration) {
	ensureMetrics()
	telegramRequests.WithLabelValues(method, outcome).Inc()
	telegramDuration.WithLabelValues(method).Observe(d.Seconds())
}

func ObserveTelegramAPIError(method, errorType string) {
	ensureMetrics()
	telegramErrors.WithLabelValues(method, errorType).Inc()
}

func ObserveStorageOp(op string, d time.Duration) {
	ensureMetrics()
	storageDuration.WithLabelValues(op).Observe(d.Seconds())
}

func ObserveFullOverHour() {
	ensureMetrics()
	fullOverHourTotal.Inc()
}

func Handler() http.Handler {
	ensureMetrics()
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}

func ListenAddr() string {
	return "127.0.0.1:9091"
}
