package collector

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

// If the extension labels have to added in this exporter, you only have
// to add the code to the following two parts.
// 1. Added the new labels name to defaultExtensionLabels
// 2. Added the new labels values to getAllResource
var defaultExtensionLabels = map[string][]string{
	"sys_elb":                 {"name", "provider", "vip_address"},
	"sys_elb_listener":        {"name", "port"},
	"sys_nat":                 {"name"},
	"sys_rds":                 {"port", "name", "role"},
	"sys_dcs":                 {"ip", "port", "name", "engine"},
	"sys_dms":                 {"name"},
	"sys_dms_instance":        {"name", "engine_version", "resource_spec_code", "connect_address", "port"},
	"sys_dms_instance_broker": {"name", "engine_version", "resource_spec_code", "connect_address", "port"},
	"sys_dms_instance_topics": {"name", "engine_version", "resource_spec_code", "connect_address", "port"},
	"sys_vpc_bandwidth":       {"name", "size", "share_type", "bandwidth_type", "charge_mode"},
	"sys_vpc_eip":             {"name", "public_ip_address", "type"},
	"sys_evs":                 {"name", "server_id", "device"},
	"sys_ecs":                 {"hostname"},
	"sys_as":                  {"name", "status"},
	"sys_functiongraph":       {"func_urn"},
}

// TTL represents the time to life
const TTL = time.Hour * 3

var (
	elbInfo serversInfo
	natInfo serversInfo
	rdsInfo serversInfo
	dmsInfo serversInfo
	dcsInfo serversInfo
	vpcInfo serversInfo
	evsInfo serversInfo
	ecsInfo serversInfo
	asInfo  serversInfo
	fgsInfo serversInfo
)

type serversInfo struct {
	TTL       int64
	LenMetric int
	Info      map[string][]string
	sync.Mutex
}

func (exporter *BaseHuaweiCloudExporter) getElbResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	elbInfo.Lock()
	defer elbInfo.Unlock()
	if elbInfo.Info == nil || time.Now().Unix() > elbInfo.TTL || elbInfo.LenMetric != exporter.MetricLen {
		allELBs, err := getAllLoadBalancer(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all LoadBalancer error:", err.Error())
		}
		if allELBs != nil {
			for _, elb := range *allELBs {
				resourceInfos[elb.ID] = []string{elb.Name, elb.Provider, elb.VipAddress}
			}
		}

		allListeners, err := getAllListener(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Listener error:", err.Error())
		}
		if allListeners != nil {
			for _, listener := range *allListeners {
				resourceInfos[listener.ID] = []string{listener.Name, fmt.Sprintf("%d", listener.ProtocolPort)}
			}
		}

		elbInfo.Info = resourceInfos
		elbInfo.TTL = time.Now().Add(TTL).Unix()
		elbInfo.LenMetric = exporter.MetricLen
	}
	return elbInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getNatResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	natInfo.Lock()
	defer natInfo.Unlock()
	if natInfo.Info == nil || time.Now().Unix() > natInfo.TTL || natInfo.LenMetric != exporter.MetricLen {
		allnat, err := getAllNat(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Nat error:", err.Error())
		}
		if allnat != nil {
			for _, nat := range *allnat {
				resourceInfos[nat.ID] = []string{nat.Name}
			}
		}

		natInfo.Info = resourceInfos
		natInfo.TTL = time.Now().Add(TTL).Unix()
		natInfo.LenMetric = exporter.MetricLen
	}
	return natInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getRdsResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	rdsInfo.Lock()
	defer rdsInfo.Unlock()
	if rdsInfo.Info == nil || time.Now().Unix() > rdsInfo.TTL || rdsInfo.LenMetric != exporter.MetricLen {
		allrds, err := getAllRds(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Rds error:", err.Error())
		}
		if allrds != nil {
			for _, rds := range allrds.Instances {
				for _, node := range rds.Nodes {
					resourceInfos[node.Id] = []string{fmt.Sprintf("%d", rds.Port), node.Name, node.Role}
				}
			}
		}

		rdsInfo.Info = resourceInfos
		rdsInfo.TTL = time.Now().Add(TTL).Unix()
		rdsInfo.LenMetric = exporter.MetricLen
	}
	return rdsInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getDmsResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	dmsInfo.Lock()
	defer dmsInfo.Unlock()
	if dmsInfo.Info == nil || time.Now().Unix() > dmsInfo.TTL || dmsInfo.LenMetric != exporter.MetricLen {
		allDmsInstance, err := getAllDms(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Dms error:", err.Error())
		}
		if allDmsInstance != nil {
			for _, dms := range allDmsInstance.Instances {
				resourceInfos[dms.InstanceID] = []string{dms.Name, dms.EngineVersion, dms.ResourceSpecCode, dms.ConnectAddress,
					fmt.Sprintf("%d", dms.Port)}
			}
		}

		allQueues, err := getAllDmsQueue(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Dms Queue error:", err.Error())
		}
		if allQueues != nil {
			for _, queue := range *allQueues {
				resourceInfos[queue.ID] = []string{queue.Name}
			}
		}

		dmsInfo.Info = resourceInfos
		dmsInfo.TTL = time.Now().Add(TTL).Unix()
		dmsInfo.LenMetric = exporter.MetricLen
	}
	return dmsInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getDcsResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	dcsInfo.Lock()
	defer dcsInfo.Unlock()
	if dcsInfo.Info == nil || time.Now().Unix() > dcsInfo.TTL || dcsInfo.LenMetric != exporter.MetricLen {
		allDcs, err := getAllDcs(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Dcs error:", err.Error())
		}
		if allDcs != nil {
			for _, dcs := range allDcs.Instances {
				resourceInfos[dcs.InstanceID] = []string{dcs.IP, fmt.Sprintf("%d", dcs.Port), dcs.Name, dcs.Engine}
			}
		}

		dcsInfo.Info = resourceInfos
		dcsInfo.TTL = time.Now().Add(TTL).Unix()
		dcsInfo.LenMetric = exporter.MetricLen
	}
	return dcsInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getVpcResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	vpcInfo.Lock()
	defer vpcInfo.Unlock()
	if vpcInfo.Info == nil || time.Now().Unix() > vpcInfo.TTL || vpcInfo.LenMetric != exporter.MetricLen {
		allPublicIps, err := getAllPublicIP(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all PublicIp error:", err.Error())
		}
		if allPublicIps != nil {
			for _, publicIP := range *allPublicIps {
				resourceInfos[publicIP.ID] = []string{publicIP.BandwidthName, publicIP.PublicIpAddress, publicIP.Type}
			}
		}

		allBandwidth, err := getAllBandwidth(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Bandwidth error:", err.Error())
			return resourceInfos
		}
		if allBandwidth != nil {
			for _, bandwidth := range *allBandwidth {
				resourceInfos[bandwidth.ID] = []string{bandwidth.Name, fmt.Sprintf("%d", bandwidth.Size), bandwidth.ShareType, bandwidth.BandwidthType, bandwidth.ChargeMode}
			}
		}

		vpcInfo.Info = resourceInfos
		vpcInfo.TTL = time.Now().Add(TTL).Unix()
		vpcInfo.LenMetric = exporter.MetricLen
	}
	return vpcInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getEvsResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	evsInfo.Lock()
	defer evsInfo.Unlock()
	if evsInfo.Info == nil || time.Now().Unix() > evsInfo.TTL || evsInfo.LenMetric != exporter.MetricLen {
		allVolumes, err := getAllVolume(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Volume error:", err.Error())
		}
		if allVolumes != nil {
			for _, volume := range *allVolumes {
				if len(volume.Attachments) > 0 {
					device := strings.Split(volume.Attachments[0].Device, "/")
					resourceInfos[fmt.Sprintf("%s-%s", volume.Attachments[0].ServerID, device[len(device)-1])] = []string{volume.Name, volume.Attachments[0].ServerID, volume.Attachments[0].Device}
				}
			}
		}

		evsInfo.Info = resourceInfos
		evsInfo.TTL = time.Now().Add(TTL).Unix()
		evsInfo.LenMetric = exporter.MetricLen
	}
	return evsInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getEcsResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	ecsInfo.Lock()
	defer ecsInfo.Unlock()
	if ecsInfo.Info == nil || time.Now().Unix() > ecsInfo.TTL || ecsInfo.LenMetric != exporter.MetricLen {
		allServers, err := getAllServer(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Server error:", err.Error())
		}
		if allServers != nil {
			for _, server := range *allServers {
				resourceInfos[server.ID] = []string{server.Name}
			}
		}

		ecsInfo.Info = resourceInfos
		ecsInfo.TTL = time.Now().Add(TTL).Unix()
		ecsInfo.LenMetric = exporter.MetricLen
	}
	return ecsInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getAsResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	asInfo.Lock()
	defer asInfo.Unlock()
	if asInfo.Info == nil || time.Now().Unix() > asInfo.TTL || asInfo.LenMetric != exporter.MetricLen {
		allGroups, err := getAllGroup(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Group error:", err.Error())
		}
		if allGroups != nil {
			for _, group := range *allGroups {
				resourceInfos[group.ID] = []string{group.Name, group.Status}
			}
		}

		asInfo.Info = resourceInfos
		asInfo.TTL = time.Now().Add(TTL).Unix()
		asInfo.LenMetric = exporter.MetricLen
	}
	return asInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getFunctionGraphResourceInfo() map[string][]string {
	resourceInfos := map[string][]string{}
	fgsInfo.Lock()
	defer fgsInfo.Unlock()
	if fgsInfo.Info == nil || time.Now().Unix() > fgsInfo.TTL || fgsInfo.LenMetric != exporter.MetricLen {
		functionList, err := getAllFunction(exporter.ClientConfig)
		if err != nil {
			logs.Logger.Errorln("Get all Function error:", err.Error())
		}
		if functionList != nil {
			for _, function := range functionList.Functions {
				resourceInfos[fmt.Sprintf("%s-%s", function.Package, function.FuncName)] = []string{function.FuncUrn}
			}
		}

		fgsInfo.Info = resourceInfos
		fgsInfo.TTL = time.Now().Add(TTL).Unix()
		fgsInfo.LenMetric = exporter.MetricLen
	}
	return fgsInfo.Info
}

func (exporter *BaseHuaweiCloudExporter) getAllResource(namespace string) map[string][]string {
	switch namespace {
	case "SYS.ELB":
		return exporter.getElbResourceInfo()
	case "SYS.NAT":
		return exporter.getNatResourceInfo()
	case "SYS.RDS":
		return exporter.getRdsResourceInfo()
	case "SYS.DMS":
		return exporter.getDmsResourceInfo()
	case "SYS.DCS":
		return exporter.getDcsResourceInfo()
	case "SYS.VPC":
		return exporter.getVpcResourceInfo()
	case "SYS.EVS":
		return exporter.getEvsResourceInfo()
	case "SYS.ECS":
		return exporter.getEcsResourceInfo()
	case "SYS.AS":
		return exporter.getAsResourceInfo()
	case "SYS.FunctionGraph":
		return exporter.getFunctionGraphResourceInfo()
	default:
		return map[string][]string{}
	}
}
