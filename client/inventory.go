package cvpgo

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

// AddDevice adds a device into CVP's inventory
func (c *CvpClient) AddDevice(ipAddr string) error {
	element := []AddInventoryElement{
		AddInventoryElement{
			ContainerName: "Tenant",
			ContainerId:   "root",
			ContainerType: "Existing",
			IpAddress:     ipAddr,
			ContainerList: []ContainerListElement{},
		},
	}
	addInventory := AddInventory{
		Data: element,
	}
	addInventoryURL := "/inventory/add/addToInventory.do?startIndex=0&endIndex=15"
	_, err := c.Call(addInventory, addInventoryURL)
	return err
}