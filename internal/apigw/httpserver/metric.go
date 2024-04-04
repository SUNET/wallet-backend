package httpserver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// metrics is the metrics object for httpserver
type metrics struct {
	UploadCounter              prometheus.Counter
	NotificationCounter        prometheus.Counter
	DocumentDelCounter         prometheus.Counter
	DocumentGetCounter         prometheus.Counter
	DocumentAttestationCounter prometheus.Counter
	PortalCounter              prometheus.Counter
	IDMappingCounter           prometheus.Counter

	SignCounter     prometheus.Counter
	GetSignCounter  prometheus.Counter
	ValidateCounter prometheus.Counter
	RevokeCounter   prometheus.Counter

	HealthCounter prometheus.Counter
}

func (m *metrics) init() {
	m.DocumentAttestationCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_document_attestation_total",
		Help: "The total number of request to endpoint /api/v1/document/attestation",
	})

	m.DocumentDelCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_document_del_total",
		Help: "The total number of request to endpoint /api/v1/document/del",
	})

	m.DocumentGetCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_document_get_total",
		Help: "The total number of request to endpoint /api/v1/document",
	})

	m.NotificationCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_notification_total",
		Help: "The total number of request to endpoint /api/v1/notification",
	})

	m.PortalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_portal_total",
		Help: "The total number of request to endpoint /api/v1/portal",
	})

	m.UploadCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_upload_total",
		Help: "The total number of request to endpoint /api/v1/upload",
	})

	m.IDMappingCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_id_mapping_total",
		Help: "The total number of request to endpoint /api/v1/id/mapping",
	})

	m.SignCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_sign_total",
		Help: "The total number of request to endpoint /api/v1/eduseal/pdf/sign",
	})

	m.GetSignCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_get_sign_total",
		Help: "The total number of request to endpoint /api/v1/eduseal/pdf/:transaction_id",
	})

	m.ValidateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_validate_total",
		Help: "The total number of request to endpoint /api/v1/eduseal/pdf/validate",
	})

	m.RevokeCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_revoke_total",
		Help: "The total number of request to endpoint /api/v1/eduseal/pdf/revoke",
	})

	m.HealthCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_http_endpoint_health_total",
		Help: "The total number of request to endpoint /health",
	})
}
