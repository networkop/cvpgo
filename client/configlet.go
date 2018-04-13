package cvpgo

import (
	"encoding/json"
	"fmt"
	"log"
)

type ValidateConfig struct {
	NetElementId string `json:"netElementId"`
	Config       string `json:"config"`
}

type JsonData struct {
	Data         interface{} `json:"data,omitempty"`
	ErrorCode    string      `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
}

type AddConfigletRepsonse struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Config string `json:"config"`
	User   string `json:"user"`
}

type AddConfigletResp struct {
	Data []struct {
		Key                  string `json:"key"`
		Name                 string `json:"name"`
		Reconciled           bool   `json:"reconciled"`
		Config               string `json:"config"`
		User                 string `json:"user"`
		Note                 string `json:"note"`
		ContainerCount       int    `json:"containerCount"`
		NetElementCount      int    `json:"netElementCount"`
		DateTimeInLongFormat int    `json:"dateTimeInLongFormat"`
		IsDefault            string `json:"isDefault"`
		IsAutoBuilder        string `json:"isAutoBuilder"`
		Type                 string `json:"type"`
		FactoryID            int    `json:"factoryId"`
		ID                   int    `json:"id"`
	} `json:"data"`
}

type ApplyConfigletData struct {
	Data []ApplyConfiglet `json:"data"`
}

type ApplyConfiglet struct {
	Info                            string   `json:"info"`
	InfoPreview                     string   `json:"infoPreview"`
	Action                          string   `json:"action"`
	NodeType                        string   `json:"nodeType"`
	NodeID                          string   `json:"nodeId"`
	ToID                            string   `json:"toId"`
	ToIDType                        string   `json:"toIdType"`
	FromID                          string   `json:"fromId"`
	NodeName                        string   `json:"nodeName"`
	FromName                        string   `json:"fromName"`
	ToName                          string   `json:"toName"`
	NodeIPAddress                   string   `json:"nodeIpAddress"`
	NodeTargetIPAddress             string   `json:"nodeTargetIpAddress"`
	ConfigletList                   []string `json:"configletList"`
	ConfigletNamesList              []string `json:"configletNamesList"`
	IgnoreConfigletList             []string `json:"ignoreConfigletList"`
	IgnoreConfigletNamesList        []string `json:"ignoreConfigletNamesList"`
	ConfigletBuilderList            []string `json:"configletBuilderList"`
	ConfigletBuilderNamesList       []string `json:"configletBuilderNamesList"`
	IgnoreConfigletBuilderList      []string `json:"ignoreConfigletBuilderList"`
	IgnoreConfigletBuilderNamesList []string `json:"ignoreConfigletBuilderNamesList"`
}

type Configlet struct {
	Config string `json:"config,omitempty"`
	Name   string `json:"name"`
	Key    string `json:"key,omitempty"`
}

type ConfigletList struct {
	List []Configlet `json:"configletList"`
}

type ValidateRequest struct {
	NetElementID string   `json:"netElementId"`
	ConfigIDList []string `json:"configIdList"`
	PageType     string   `json:"pageType"`
}

type ReconcileBody struct {
	Name       string `json:"name"`
	Config     string `json:"config"`
	Reconciled bool   `json:"reconciled"`
}

func checkErrors(data JsonData) error {
	if data.ErrorCode != "" {
		log.Printf("Error from CVP: %s", data.ErrorMessage)
		return fmt.Errorf("CVP returned error code: %s, %s", data.ErrorCode, data.ErrorMessage)
	}
	return nil
}

// AddConfiglet adds configlet to CVP
func (c *CvpClient) AddConfiglet(name, config string) error {
	addConfigletURL := "/configlet/addConfiglet.do"
	configlet := Configlet{
		Name:   name,
		Config: config,
	}
	resp, err := c.Call(configlet, addConfigletURL)
	responseBody := JsonData{}
	if err = json.Unmarshal(resp, &responseBody); err != nil {
		log.Printf("Error adding configlet %+v", err)
	}
	return checkErrors(responseBody)
}

// DeleteConfiglet deletes configlet from CVP
func (c *CvpClient) DeleteConfiglet(cfgletName string) error {
	url := "/configlet/deleteConfiglet.do"
	cfglet, err := c.GetConfigletByName(cfgletName)
	// properties which are not allowed by the schema: ["config"]
	cfglet.Config = ""
	body := []Configlet{cfglet}
	resp, err := c.Call(body, url)
	responseBody := JsonData{}
	if err = json.Unmarshal(resp, &responseBody); err != nil {
		log.Printf("Error adding configlet %+v", err)
	}
	return checkErrors(responseBody)
}

// ValidateConfiglet takes the netElementId (MAC Address) and a Configlet Name
// given as Key from adding a configlet, and validates it
func (c *CvpClient) ValidateConfiglet(netElementID, cfgletName string) error {
	url := "/provisioning/v2/validateAndCompareConfiglets.do"
	cfglet, err := c.GetConfigletByName(cfgletName)
	req := ValidateRequest{
		NetElementID: netElementID,
		ConfigIDList: []string{cfglet.Key},
	}
	resp, err := c.Call(req, url)
	responseBody := JsonData{}
	if err = json.Unmarshal(resp, &responseBody); err != nil {
		log.Printf("Error validating configlet %+v", err)
	}
	return checkErrors(responseBody)
}

// ConfigSync udpates reconcile configlet with the current running configuration of the device
func (c *CvpClient) ConfigSync(hostname string) error {
	device, err := c.GetDevice(hostname)
	url := "/provisioning/updateReconcileConfiglet.do?netElementId=" + device.Key
	config, err := c.GetInventoryConfig(device.Key)
	log.Printf("Got running config from device %+v", device.Key)
	req := ReconcileBody{
		Name:       "RECONCILE_" + device.Fqdn,
		Config:     config,
		Reconciled: false,
	}
	_, err = c.Call(req, url)
	if err != nil {
		log.Printf("Error creating Reconcile configlet config %+v", err)
	}
	err = c.ApplyConfigletToDevice(device.IPAddress, device.Fqdn, device.SystemMacAddress, []string{"RECONCILE_" + device.Fqdn}, true)
	if err != nil {
		log.Printf("Error applying Reconcile configlet to device %+v", err)
	}
	return nil
}

// ValidateReconcileAll takes netElementID and validates and reconciles all configlets assigned to a device
func (c *CvpClient) ValidateReconcileAll(netElementID string) error {
	url := "/provisioning/v2/validateAndCompareConfiglets.do"
	cfglets, err := c.getConfigletByDeviceID(netElementID)
	log.Printf("Got configlets from device %+v", cfglets)
	req := ValidateRequest{
		NetElementID: netElementID,
		ConfigIDList: getKeys(cfglets),
	}
	resp, err := c.Call(req, url)
	responseBody := JsonData{}
	if err = json.Unmarshal(resp, &responseBody); err != nil {
		log.Printf("Error validating configlet %+v", err)
	}
	return checkErrors(responseBody)
}

/* ApplyConfigletToDevice applies configlets to a device
   deviceIpAddress -- Ip address of the device (type: string)
   deviceFqdn -- Fully qualified domain name for device (type: string)
   deviceKey -- mac address of the device (type: string)
   cnl -- List of name of configlets to be applied
   (type: List of Strings)
   ckl -- Keys of configlets to be applied (type: List of Strings)
*/
func (c *CvpClient) ApplyConfigletToDevice(deviceIP, deviceName, deviceMac string, cnl []string, save bool) error {
	// func (c *CvpClient) ApplyConfigletToDevice(deviceName, deviceMac string, cnl, ckl, []string) error {
	// func (c *CvpClient) ApplyConfigletToDevice(deviceIP, deviceName, deviceMac string, cnl, ckl, cbnl, cbkl []string) error {
	cfgletCurrent, err := c.getConfigletByDeviceID(deviceMac)
	if err != nil {
		log.Printf("Error retrieving configlets from a device")
		return err
	}
	cfgletNew, err := c.getConfigletsByName(cnl)
	if err != nil {
		log.Printf("Error retrieving configlets by its name")
		return err
	}
	cfgletAll, _ := c.mergeCfglet(cfgletCurrent, cfgletNew)
	applyCfglet := ApplyConfiglet{
		Info:                            "Configlet Assign to device: " + deviceName,
		InfoPreview:                     "<b>Configlet assign</b> to Device " + deviceName,
		Action:                          "associate",
		NodeIPAddress:                   deviceIP,
		NodeTargetIPAddress:             deviceIP,
		NodeType:                        "configlet",
		ToID:                            deviceMac,
		ToIDType:                        "netelement",
		ToName:                          deviceName,
		ConfigletList:                   getKeys(cfgletAll),
		ConfigletNamesList:              getNames(cfgletAll),
		ConfigletBuilderList:            []string{},
		ConfigletBuilderNamesList:       []string{},
		IgnoreConfigletList:             []string{},
		IgnoreConfigletNamesList:        []string{},
		IgnoreConfigletBuilderList:      []string{},
		IgnoreConfigletBuilderNamesList: []string{},
	}
	log.Printf("Applying configlet : %+v", applyCfglet)
	err = c.addTempAction(applyCfglet)
	if save {
		return c.saveTopologyV2(applyCfglet)
	}
	return err
}

func (c *CvpClient) addTempAction(action ApplyConfiglet) error {
	url := "/provisioning/addTempAction.do?format=topology&queryParam=&nodeId=root"
	dataArray := []ApplyConfiglet{action}
	data := ApplyConfigletData{
		Data: dataArray,
	}
	resp, err := c.Call(data, url)
	log.Printf("Response from temp action %+s", resp)
	return err
}

func (c *CvpClient) saveTopologyV2(action ApplyConfiglet) error {
	url := "/provisioning/v2/saveTopology.do"
	dataArray := []ApplyConfiglet{action}
	data := ApplyConfigletData{
		Data: dataArray,
	}
	resp, err := c.Call(data, url)
	log.Printf("Response from temp action %+s", resp)
	return err
}

func (c *CvpClient) ValidateConfig(deviceMac, config string) error {
	url := "/configlet/validateConfig.do"
	data := ValidateConfig{
		NetElementId: deviceMac,
		Config:       config,
	}
	body := JsonData{
		Data: data,
	}
	resp, err := c.Call(body, url)
	log.Printf("Response from temp action %+s", resp)
	return err
}

func (c *CvpClient) getConfigletByDeviceID(deviceMac string) ([]Configlet, error) {
	url := "/provisioning/getConfigletsByNetElementId.do?netElementId=" + deviceMac + "&queryParam=&startIndex=0&endIndex=15"
	respbody, err := c.Get(url)
	respConfiglet := ConfigletList{}
	err = json.Unmarshal(respbody, &respConfiglet)
	if err != nil {
		log.Printf("Error decoding getConfigletByDeviceID :%s\n", err)
		return nil, err
	}
	//if len(respConfiglet.List) == 0 {
	//	return nil, fmt.Errorf("No configlets returned")
	//}
	return respConfiglet.List, err
}

func (c *CvpClient) GetConfigletByName(cfglet string) (Configlet, error) {
	url := "/configlet/getConfigletByName.do?name=" + cfglet
	respbody, err := c.Get(url)
	respConfiglet := Configlet{}
	err = json.Unmarshal(respbody, &respConfiglet)
	if err != nil {
		log.Printf("Error decoding getConfigletByName :%s\n", err)
		return respConfiglet, err
	}
	return respConfiglet, err
}

func (c *CvpClient) getConfigletsByName(cfglets []string) (result []Configlet, err error) {
	for _, cfgletName := range cfglets {
		cfglet, err := c.GetConfigletByName(cfgletName)
		if err != nil {
			return result, err
		}
		result = append(result, cfglet)
	}
	return result, nil
}

func contains(list []Configlet, elem Configlet) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

func (c *CvpClient) filterCfglet(all, remove []Configlet) (stay []Configlet, err error) {
	for _, cfglet := range all {
		if !contains(remove, cfglet) {
			stay = append(stay, cfglet)
		}
	}
	return stay, nil
}

func (c *CvpClient) mergeCfglet(current, new []Configlet) (all []Configlet, err error) {
	all = append(current, all...)
	for _, cfglet := range new {
		if !contains(current, cfglet) {
			all = append(all, cfglet)
		}
	}
	return all, nil
}

func getNames(cfglets []Configlet) []string {
	result := make([]string, 0)
	for _, cfglet := range cfglets {
		result = append(result, cfglet.Name)
	}
	return result
}

func getKeys(cfglets []Configlet) []string {
	result := make([]string, 0)
	for _, cfglet := range cfglets {
		result = append(result, cfglet.Key)
	}
	return result
}

// RemoveConfigletFromDevice removes configlets (list of strings) from device
func (c *CvpClient) RemoveConfigletFromDevice(deviceIP, deviceName, deviceMac string, cfgletRemoveNames []string, save bool) error {
	cfgletAll, err := c.getConfigletByDeviceID(deviceMac)
	if err != nil {
		return err
	}
	cfgletRemove, err := c.getConfigletsByName(cfgletRemoveNames)
	if err != nil {
		return err
	}
	cfgletRemain, err := c.filterCfglet(cfgletAll, cfgletRemove)
	removeCfglet := ApplyConfiglet{
		Info:                            "Configlet Remove from device: " + deviceName,
		InfoPreview:                     "<b>Configlet remove</b> from Device " + deviceName,
		Action:                          "associate",
		NodeIPAddress:                   deviceIP,
		NodeTargetIPAddress:             deviceIP,
		NodeType:                        "configlet",
		ToID:                            deviceMac,
		ToIDType:                        "netelement",
		ToName:                          deviceName,
		ConfigletList:                   getKeys(cfgletRemain),
		ConfigletNamesList:              getNames(cfgletRemain),
		ConfigletBuilderList:            []string{},
		ConfigletBuilderNamesList:       []string{},
		IgnoreConfigletList:             getKeys(cfgletRemove),
		IgnoreConfigletNamesList:        getNames(cfgletRemove),
		IgnoreConfigletBuilderList:      []string{},
		IgnoreConfigletBuilderNamesList: []string{},
	}
	log.Printf("Removing configlet : %+v", removeCfglet)
	err = c.addTempAction(removeCfglet)
	if save {
		return c.saveTopologyV2(removeCfglet)
	}
	return err
}
