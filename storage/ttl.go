package storage

import (
	"fmt"
	"time"
)

//TTL
func TTL(ttl *int, ms MetricStore) {

	for {
		timer1 := time.NewTimer(time.Second * 5)
		<-timer1.C

		currentTime := time.Now()
		// fmt.Println(currentTime)
		// fmt.Println(ms)
		metricsList := ms.GetMetricFamilies()
		for i, v := range metricsList {
			fmt.Println(v)
			// name:"some_metric" type:UNTYPED metric:<label:<name:"instance" value:"" > label:<name:"job" value:"some_job" > untyped:<value:3.14 > >
			// name:"push_time_seconds" help:"Last Unix time when this group was changed in the Pushgateway." type:GAUGE metric:<label:<name:"instance" value:"" > label:<name:"job" value:"some_job" > gauge:<value:1.5163936722265959e+09 > >

			metricName := v.GetName()
			if metricName == "push_time_seconds" {
				// fmt.Printf("\nCurrent time: %d", currentTime.Unix())
				// fmt.Printf("\nJob's push_time_seconds time: %f", v.GetMetric()[0].GetGauge().GetValue())
				if int(currentTime.Unix())-int(v.GetMetric()[0].GetGauge().GetValue()) > *ttl {
					// When an entry of push_time_seconds exceeds TTL, grab value from the "job" label, submit request to remove all metrics related the job

					// fmt.Println(v.GetMetric()[0].GetLabel())
					// [name:"instance" value:""  name:"job" value:"some_job" ]
					jobName := v.GetMetric()[0].GetLabel()[1].GetValue()
					fmt.Println(jobName)
					fmt.Println(metricsList[i])

					labels := make(map[string]string)
					labels["job"] = jobName
					wr := WriteRequest{
						Labels:         labels,
						Timestamp:      time.Now(),
						MetricFamilies: nil,
					}

					ms.SubmitWriteRequest(wr)
				}
			}
		}
		timer1.Reset(time.Second * 5)
	}
}

// // DiskMetricStore is an implementation of MetricStore that persists metrics to
// // disk.
// type DiskMetricStore struct {
// 	lock            sync.RWMutex // Protects metricFamilies.
// 	writeQueue      chan WriteRequest
// 	drain           chan struct{}
// 	done            chan error
// 	metricGroups    GroupingKeyToMetricGroup
// 	persistenceFile string
// }

// // WriteRequest is a request to change the MetricStore, i.e. to process it, a
// // write lock has to be acquired. If MetricFamilies is nil, this is a request to
// // delete metrics that share the given Labels as a grouping key. Otherwise, this
// // is a request to update the MetricStore with the MetricFamilies. The key in
// // MetricFamilies is the name of the mapped metric family. All metrics in
// // MetricFamilies MUST have already set job and other labels that are consistent
// // with the Labels fields. The Timestamp field marks the time the request was
// // received from the network. It is not related to the timestamp_ms field in the
// // Metric proto message.
// type WriteRequest struct {
// 	Labels         map[string]string
// 	Timestamp      time.Time
// 	MetricFamilies map[string]*dto.MetricFamily
// }

// func (dms *DiskMetricStore) processWriteRequest(wr WriteRequest) {
// 	dms.lock.Lock()
// 	defer dms.lock.Unlock()

// 	key := model.LabelsToSignature(wr.Labels)

// 	if wr.MetricFamilies == nil {
// 		// Delete.
// 		delete(dms.metricGroups, key)
// 		return
// 	}
