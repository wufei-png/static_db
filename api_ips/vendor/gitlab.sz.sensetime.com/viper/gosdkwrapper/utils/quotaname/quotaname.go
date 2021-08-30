package quotaname

import (
	"fmt"
	"runtime"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/monitor"
)

var InternalNvidiaDeviceName = map[string]string{
	monitor.P4:                 "p4",
	monitor.GraphicsDevice:     "p4",
	monitor.T4:                 "t4",
	monitor.NvidiaGTX1080:      "p4",
	monitor.NvidiaGTX1080Ti:    "p4",
	monitor.NvidiaRTX2070Super: "rtx2070s",
	monitor.RTX4000:            "rtx4000",
	monitor.ATLAS:              "ascend310",
	monitor.HI3559A:            "hi3559a",
	monitor.StpuDevice:         "stpu",
	monitor.StpuDeviceSV114A:   "sv114a",
	monitor.StpuDeviceSV116E:   "sv116e",
	monitor.AscendDevice310:    "ascend310",
	monitor.MLU220Edge:         "mlu220edge",
	monitor.MLU270:             "mlu270",
}

func GetQuotaPostfix(n string) (string, error) {
	switch n {
	case monitor.GTX1060, monitor.GTX1660, monitor.GTX1070, monitor.GTX2080:
		// 开发机不需要做特殊处理
		return "", nil
	case monitor.P4, monitor.GraphicsDevice:
		// 兼容老系统，quota名字不变
		return "", nil
		// 兼容cpu模式, cpu模式下quota后缀使用arch名
	case runtime.GOARCH:
		return runtime.GOARCH, nil
	default:
		// 如果以后有新增加会混布的卡，在这里添加
		if quotaName, ok := InternalNvidiaDeviceName[n]; ok {
			// 兼容老系统，quota名字不变
			if quotaName == "p4" {
				return "", nil
			}
			return quotaName, nil
		}
	}

	return "", fmt.Errorf("unknown device %v", n)
}
