package config

var (
	PrometheusServiceName = "kube-prom-stack-kube-prome-prometheus"
	PrometheusNamespace   = "monitoring"
	PrometheusPort        = "9090"

	/*
		PrometheusURL = fmt.Sprintf(
			"http://%s.%s.svc.cluster.local:%s",
			PrometheusServiceName,
			PrometheusNamespace,
			PrometheusPort,
		)
	*/
	PrometheusURL = "http://localhost:9090"
)
