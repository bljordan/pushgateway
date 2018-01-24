package storage

import (
	"time"
)

//TTL
func TTL(ttl *int, ms MetricStore) {

	for {
		// Timer for how frequently to check
		timer1 := time.NewTimer(time.Second * 5)
		<-timer1.C

		currentTime := time.Now()
		metricsList := ms.GetMetricFamilies()
		for _, v := range metricsList {
			// fmt.Println(v)

			metricName := v.GetName()
			if metricName == "push_time_seconds" {
				jobPushTimeSeconds := v.GetMetric()[0].GetGauge().GetValue()
				if int(currentTime.Unix())-int(jobPushTimeSeconds) > *ttl {
					// When an entry of push_time_seconds exceeds TTL, grab value from the "job" label, submit request to remove all metrics related the job

					jobName := v.GetMetric()[0].GetLabel()[1].GetValue()
					// fmt.Println(jobName)
					// fmt.Println(metricsList[i])

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
