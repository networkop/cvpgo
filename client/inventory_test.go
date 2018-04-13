package cvpgo

import (
	"testing"
)

type CVPInfo struct {
	IPAddress string
	Username  string
	Password  string
	Container string
}

type InventoryTestData struct {
	CVP             *CVPInfo
	DeviceIP        string
	DeviceContainer string
	DeviceHostname  string
	DeviceMAC       string
}

func buildTestData() *InventoryTestData {
	result := InventoryTestData{}
	result.CVP = &CVPInfo{
		IPAddress: "192.168.133.1",
		Username:  "cvpadmin",
		Password:  "cvpadmin1",
		Container: "Test"}
	result.DeviceIP = "172.26.0.2"
	result.DeviceContainer = "Test"
	result.DeviceHostname = "localhost"
	result.DeviceMAC = "02:42:c0:be:b2:37"
	return &result
}

func TestAddDevices(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.AddDevice(testdata.DeviceIP, testdata.DeviceContainer)
	if err != nil {
		t.Errorf("%+v", err)
	}
	err = cvp.SaveInventory()
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestSaveInventory(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.SaveInventory()
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestGetDevices(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, err := cvp.GetDevice(testdata.DeviceHostname)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if dev == nil {
		t.Errorf("Did not retreive any devices\n")
	}
	t.Logf("Retrieved %+v", dev)
}

func TestGetDeviceConfig(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	config, err := cvp.GetInventoryConfig(testdata.DeviceMAC)
	if err != nil {
		t.Errorf("%+v", err)
	}
	t.Logf("Retrieved config %s", config)
}

func TestRemoveDevices(t *testing.T) {
	//cvpInfo := CVPInfo{IPAddress: "10.90.224.178", Username: "cvpadmin", Password: "arista123", Container: "CoreSite"}
	testdata := buildTestData()
	cvpInfo := *testdata.CVP
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.RemoveDevice(testdata.DeviceMAC)
	if err != nil {
		t.Errorf("%+v", err)
	}
}
