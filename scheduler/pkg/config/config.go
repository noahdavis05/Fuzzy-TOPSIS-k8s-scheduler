package config

import "scheduler/pkg/types"

// prometheus variables
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

// TOPSIS variables e.g. weights and ideals
var (
	// TOPSIS ideals
	PosCPUIdeal = types.FuzzyNumber{
		A: 75,
		B: 75,
		C: 75,
	}
	NegCPUIdeal = types.FuzzyNumber{
		A: 100,
		B: 100,
		C: 100,
	}

	PosRAMIdeal = types.FuzzyNumber{
		A: 75,
		B: 75,
		C: 75,
	}
	NegRAMIdeal = types.FuzzyNumber{
		A: 100,
		B: 100,
		C: 100,
	}

	PosCPURangeIdeal = types.FuzzyNumber{
		A: 0,
		B: 0,
		C: 0,
	}
	NegCPURangeIdeal = types.FuzzyNumber{
		A: 100,
		B: 100,
		C: 100,
	}

	PosRAMRangeIdeal = types.FuzzyNumber{
		A: 0,
		B: 0,
		C: 0,
	}
	NegRAMRangeIdeal = types.FuzzyNumber{
		A: 100,
		B: 100,
		C: 100,
	}

	// TOPSIS Weights
	CPUWeights = types.FuzzyNumber{
		A: 0.5,
		B: 0.5,
		C: 0.5,
	}

	RAMWeights = types.FuzzyNumber{
		A: 0.5,
		B: 0.5,
		C: 0.5,
	}

	CPURangeWeights = types.FuzzyNumber{
		A: 0.3,
		B: 0.3,
		C: 0.3,
	}

	RAMRangeWeights = types.FuzzyNumber{
		A: 0.3,
		B: 0.3,
		C: 0.3,
	}
)
