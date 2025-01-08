package ixml

// GpmMetricsGetType includes interface types for GpmSample instead of nvmlGpmSample
type GpmMetricsGetType struct {
	Version    uint32
	NumMetrics uint32
	Sample1    GpmSample
	Sample2    GpmSample
	Metrics    [98]GpmMetric
}

func (g *nvmlGpmMetricsGetType) convert() *GpmMetricsGetType {
	out := &GpmMetricsGetType{
		Version:    g.Version,
		NumMetrics: g.NumMetrics,
		Sample1:    g.Sample1,
		Sample2:    g.Sample2,
	}
	for i := range g.Metrics {
		out.Metrics[i] = g.Metrics[i]
	}
	return out
}

func GpmMetricsGet(metricsGet *GpmMetricsGetType) Return {
	metricsGet.Version = GPM_METRICS_GET_VERSION
	return gpmMetricsGet(metricsGet)
}

func gpmMetricsGet(metricsGet *GpmMetricsGetType) Return {
	nvmlMetricsGet := (*nvmlGpmMetricsGetType)(metricsGet)
	ret := nvmlGpmMetricsGet(nvmlMetricsGet)
	*metricsGet = *nvmlMetricsGet.convert()
	return ret
}

func (gpmSample GpmSample) Free() Return {
	return nvmlGpmSampleFree(gpmSample)
}

func (gpmSample GpmSample) Get(device Device) Return {
	return nvmlGpmSampleGet(device, gpmSample)
}

func GpmSampleAlloc() (GpmSample, Return) {
	var gpmSample GpmSample
	ret := nvmlGpmSampleAlloc(&gpmSample)
	return gpmSample, ret
}

func (device Device) GpmSampleGet(gpmSample GpmSample) Return {
	return nvmlGpmSampleGet(device, gpmSample)
}

func (device Device) GpmQueryDeviceSupport() (GpmSupport, Return) {
	return gpmQueryDeviceSupport(device)
}

func gpmQueryDeviceSupport(device Device) (GpmSupport, Return) {
	var gpmSupport GpmSupport
	gpmSupport.Version = GPM_SUPPORT_VERSION
	ret := nvmlGpmQueryDeviceSupport(device, &gpmSupport)
	return gpmSupport, ret
}
