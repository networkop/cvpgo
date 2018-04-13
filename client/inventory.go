package cvpgo

import (
	"encoding/json"
	"fmt"
	"log"
)

type GetInventory struct {
	Total          int          `json:"total"`
	NetElementList []NetElement `json:"netElementList"`
}

type GetContainer struct {
	Total         int         `json:"total"`
	ContainerList []Container `json:"containerList"`
}

type Container struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type NetElement struct {
	ModeName         string `json:"modelName"`
	InternalVersion  string `json:"internalVersion"`
	SystemMacAddress string `json:"systemMacAddress"`
	SerialNumber     string `json:"serialNumber"`
	Version          string `json:"version"`
	Fqdn             string `json:"fqdn"`
	Key              string `json:"key"`
	IPAddress        string `json:"ipAddress"`
}

type InventoryData struct {
	Total      int `json:"total"`
	Containers struct {
		HierarchyNetElementCount int    `json:"hierarchyNetElementCount"`
		ChildNetElementCount     int    `json:"childNetElementCount"`
		Key                      string `json:"key"`
		Name                     string `json:"name"`
		Type                     string `json:"type"`
		ChildContainerCount      int    `json:"childContainerCount"`
		ParentContainerID        string `json:"parentContainerId"`
		Mode                     string `json:"mode"`
		DeviceStatus             string `json:"deviceStatus"`
		ChildTaskCount           int    `json:"childTaskCount"`
		ChildNetElementList      string `json:"childNetElementList"`
		TempAction               string `json:"tempAction"`
		TempEvent                string `json:"tempEvent"`
		ChildContainerList       []struct {
		} `json:"childContainerList"`
	} `json:"containers"`
	TempContainer []struct {
		Name             string `json:"name"`
		ParentID         string `json:"parentId"`
		Key              string `json:"key"`
		ChildContainerID string `json:"childContainerId"`
		Type             string `json:"type"`
		UserID           string `json:"userId"`
		FactoryID        int    `json:"factoryId"`
		ID               int    `json:"id"`
	} `json:"tempContainer"`
	Dashboard []struct {
		Count     int    `json:"count"`
		Status    string `json:"status"`
		FactoryID int    `json:"factoryId"`
		ID        int    `json:"id"`
	} `json:"dashboard"`
	TempNetElement []struct {
		ClassID          string   `json:"classId"`
		ModelName        string   `json:"modelName"`
		InternalVersion  string   `json:"internalVersion"`
		SystemMacAddress string   `json:"systemMacAddress"`
		MemTotal         string   `json:"memTotal"`
		BootupTimeStamp  string   `json:"bootupTimeStamp"`
		MemFree          string   `json:"memFree"`
		Architecture     string   `json:"architecture"`
		InternalBuildID  string   `json:"internalBuildId"`
		HardwareRevision string   `json:"hardwareRevision"`
		Fqdn             string   `json:"fqdn"`
		IPAddress        string   `json:"ipAddress"`
		TaskIDList       []string `json:"taskIdList"`
		FactoryID        string   `json:"factoryId"`
		ZtpMode          string   `json:"ztpMode"`
		IsDANZEnabled    string   `json:"isDANZEnabled"`
		ContainerID      string   `json:"containerId"`
		StatusMessage    string   `json:"statusMessage"`
		ContainerName    string   `json:"containerName"`
		UserID           string   `json:"userId"`
		Version          string   `json:"version"`
		Key              string   `json:"key"`
		ID               string   `json:"id"`
		Type             string   `json:"type"`
		SerialNumber     string   `json:"serialNumber"`
		Status           string   `json:"status"`
		IsMLAGEnabled    string   `json:"isMLAGEnabled"`
		UserName         string   `json:"userName"`
		Passowrd         string   `json:"passowrd"`
	} `json:"tempNetElement"`
}

type AddInventory struct {
	Data []AddInventoryElement `json:"data"`
}

type AddInventoryElement struct {
	ContainerName string                 `json:"containerName"`
	ContainerId   string                 `json:"containerId"`
	ContainerType string                 `json:"containerType"`
	IpAddress     string                 `json:"ipAddress"`
	ContainerList []ContainerListElement `json:"containerList"`
}

type ContainerListElement struct {
	Name             string `json:"name"`
	Key              string `json:"key"`
	ChildContainerId string `json:"childContainerId"`
	Type             string `json:"type"`
}

type GetInventoryConfigResponse struct {
	Output   string   `json:"output"`
	Warnings []string `json:"warnings"`
}

// AddDevice adds a device into CVP's inventory
func (c *CvpClient) AddDevice(ipAddr string, cn string) error {
	container, err := c.GetContainerByName(cn)
	element := []AddInventoryElement{
		AddInventoryElement{
			ContainerName: cn,
			ContainerId:   container.Key,
			ContainerType: "Existing",
			IpAddress:     ipAddr,
			ContainerList: []ContainerListElement{},
		},
	}
	addInventory := AddInventory{
		Data: element,
	}
	addInventoryURL := "/inventory/add/addToInventory.do?startIndex=0&endIndex=15"
	_, err = c.Call(addInventory, addInventoryURL)
	return err
}

// AddDevice adds a device into CVP's inventory
func (c *CvpClient) SaveInventory() error {
	saveInventoryURL := "/inventory/v2/saveInventory.do"
	_, err := c.Call("", saveInventoryURL)
	return err
}

// Removes device from CVP's inventory
func (c *CvpClient) RemoveDevice(mac string) error {
	data := struct {
		Data []string `json:"data"`
	}{[]string{mac}}
	RemoveInventoryURL := "/inventory/deleteDevices.do?"
	_, err := c.Call(data, RemoveInventoryURL)
	return err
}

// GetDevice uses the hostname of a device to lookup the full entry in CVP
// and returns the full NetElement entry
func (c *CvpClient) GetDevice(hostname string) (*NetElement, error) {
	getDeviceURL := "/inventory/getInventory.do?queryparam=" + hostname + "&startIndex=0&endIndex=0"
	respbody, err := c.Get(getDeviceURL)
	respDevice := GetInventory{}
	err = json.Unmarshal(respbody, &respDevice)
	if err != nil {
		log.Printf("Error decoding getdevice :%s\n", err)
		return nil, err
	}
	if len(respDevice.NetElementList) == 0 {
		return nil, fmt.Errorf("No devices returned")
	}
	return &respDevice.NetElementList[0], err
}

// GetInventoryConfig returns the configuration retrieved from a netelement.
func (c *CvpClient) GetInventoryConfig(netElemendID string) (string, error) {
	url := "/inventory/getInventoryConfiguration.do?netElementId=" + netElemendID
	respbody, err := c.Get(url)
	respData := GetInventoryConfigResponse{}
	err = json.Unmarshal(respbody, &respData)
	if err != nil {
		log.Printf("Error decoding getInventoryConfiguration :%s\n", err)
		return "", err
	}
	return respData.Output, nil
}

// Returns a container that exactly matches the name.
func (c *CvpClient) GetContainerByName(query string) (*Container, error) {
	getContainerURL := "/provisioning/searchTopology.do?queryParam=" + query + "&startIndex=0&endIndex=0"
	respbody, err := c.Get(getContainerURL)
	respContainer := GetContainer{}
	err = json.Unmarshal(respbody, &respContainer)
	if err != nil {
		log.Printf("Error decoding getcontainer :%s\n", err)
		return nil, err
	}
	if len(respContainer.ContainerList) == 0 {
		return nil, fmt.Errorf("No devices returned")
	}
	return &respContainer.ContainerList[0], err
}

// GetInventory will return all the devices in CVP
func (c *CvpClient) GetInventory(query string) (*[]NetElement, error) {
	getDeviceURL := "/inventory/getInventory.do?queryparam=" + query + "&startIndex=0&endIndex=0"
	respbody, err := c.Get(getDeviceURL)
	respDevice := GetInventory{}
	err = json.Unmarshal(respbody, &respDevice)
	if err != nil {
		log.Printf("Error decoding getdevice :%s\n", err)
		return nil, err
	}
	if len(respDevice.NetElementList) == 0 {
		return nil, fmt.Errorf("No devices returned")
	}
	return &respDevice.NetElementList, err
}
