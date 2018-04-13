package cvpgo

import (
	"os"
	"strings"
	"testing"
)

type TestData struct {
	CvpIP        string
	CvpUser      string
	CvpPwd       string
	CvpContainer string
	Device       string
	Ckls         []string
	Cnls         []string
	NewConfig    string
	NewCfglet    string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvList(key string, fallback []string) []string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.Split(value, ",")
	}
	return fallback
}

func updateFromEnvironment(data *TestData) {
	data.CvpIP = getEnv("CVP_IP", data.CvpIP)
	data.CvpUser = getEnv("CVP_USER", data.CvpUser)
	data.CvpPwd = getEnv("CVP_PWD", data.CvpPwd)
	data.CvpContainer = getEnv("CVP_CONT", data.CvpContainer)
	data.Device = getEnv("CVP_DEVICE", data.Device)
	data.Ckls = getEnvList("CVP_CKLS", data.Ckls)
	data.Cnls = getEnvList("CVP_CNLS", data.Cnls)
}

func buildConfigletTestData() *TestData {
	return &TestData{
		CvpIP:        "192.168.133.1",
		CvpUser:      "cvpadmin",
		CvpPwd:       "cvpadmin1",
		CvpContainer: "Test",
		Device:       "localhost",
		Cnls:         []string{"Test1"},
		NewConfig:    "username TEST nopassword",
		NewCfglet:    "Test1",
	}
}

func TestAddConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.AddConfiglet(data.NewCfglet, data.NewConfig)
	if err != nil {
		t.Errorf("Error adding configlet : %s", err)
	}
}

func TestGetConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	cfglet, err := cvp.GetConfigletByName(data.NewCfglet)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
	if cfglet.Key == "" {
		t.Logf("No Configlets were found")
	} else {
		t.Logf("Returned configlet with key %s", cfglet.Key)
	}
}

func TestSyncConfig(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	device, err := cvp.GetDevice(data.Device)
	if err != nil {
		t.Errorf("Error getting device : %s", data.Device)
	}
	err = cvp.ConfigSync(device.Fqdn)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
}

func TestValidateConfig(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	device, err := cvp.GetDevice(data.Device)
	if err != nil {
		t.Errorf("Error getting device : %s", data.Device)
	}
	err = cvp.ValidateReconcileAll(device.Key)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
}

func TestValidateReconcileAll(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	device, err := cvp.GetDevice(data.Device)
	if err != nil {
		t.Errorf("Error getting device : %s", data.Device)
	}
	err = cvp.ValidateReconcileAll(device.Key)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
}

func TestGetConfigletByDeviceID(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	device, err := cvp.GetDevice(data.Device)
	_, err = cvp.getConfigletByDeviceID(device.Key)
	if err != nil {
		t.Errorf("Error getting configlet : %s", err)
	}
}

func TestApplyConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, _ := cvp.GetDevice(data.Device)
	t.Logf("Retrieved %+v", dev)
	cnl := data.Cnls
	err := cvp.ApplyConfigletToDevice(dev.IPAddress, dev.Fqdn, dev.SystemMacAddress, cnl, true)
	if err != nil {
		t.Errorf("Error applying configlet : %s", err)
	}
}

func TestRemoveConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	dev, _ := cvp.GetDevice(data.Device)
	t.Logf("Retrieved %+v", dev)
	err := cvp.RemoveConfigletFromDevice(dev.IPAddress, dev.Fqdn, dev.Key, []string{data.NewCfglet}, true)
	if err != nil {
		t.Errorf("Error applying configlet : %s", err)
	}
}

func TestDeleteConfiglet(t *testing.T) {
	data := buildConfigletTestData()
	updateFromEnvironment(data)
	t.Logf("Test data: %+v", data)
	cvpInfo := CVPInfo{IPAddress: data.CvpIP, Username: data.CvpUser, Password: data.CvpPwd, Container: data.CvpContainer}
	cvp := New(cvpInfo.IPAddress, cvpInfo.Username, cvpInfo.Password)
	err := cvp.DeleteConfiglet(data.NewCfglet)
	if err != nil {
		t.Errorf("Erorr adding configlet : %s", err)
	}
}
