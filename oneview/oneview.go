package oneview

import (
	"errors"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/HewlettPackard/oneview-golang/icsp"
	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

type Attrib struct {
	Name string  `json:"name,omitempty"`
	Value string  `json:"value,omitempty"`
}

// Driver OneView driver structure
type Driver struct {
	*drivers.BaseDriver
	ClientICSP           *icsp.ICSPClient
	ClientOV             *ov.OVClient
	IloUser              string
	IloPassword          string
	IloPort              int
	OSBuildPlans         []string
	OSDeployPlan         string
	OSDeployAttributes   []Attrib
	SSHUser              string
	SSHPort              int
	SSHPublicKey         string
	ServerTemplate       string
	HWName               string
	PublicSlotID         int
	PublicConnectionName string
	Profile              ov.ServerProfile
	Hardware             ov.ServerHardware
	Server               icsp.Server
}

const (
	driverName     = "oneview"
	defaultTimeout = 1 * time.Second
)

// Error messages
var (
	ErrDriverMissingEndPointOptionOV   = errors.New("Missing option --oneview-ov-endpoint or environment ONEVIEW_OV_ENDPOINT")
	ErrDriverMissingEndPointOptionICSP = errors.New("Missing option --oneview-icsp-endpoint or environment ONEVIEW_ICSP_ENDPOINT")
	ErrDriverMissingTemplateOption     = errors.New("Missing option --oneview-server-template or environment ONEVIEW_SERVER_TEMPLATE")
	ErrDriverMissingBuildPlanOption    = errors.New("Missing option --oneview-os-plans or ONEVIEW_OS_PLANS")
)

// NewDriver - create a OneView object driver
func NewDriver(machineName string, storePath string) drivers.Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: machineName,
			StorePath:   storePath,
		},
	}
}

// GetCreateFlags registers the flags this driver adds to
// "docker hosts create"
//
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			Name:   "oneview-ov-user",
			Usage:  "User Name to OneView Server",
			Value:  "",
			EnvVar: "ONEVIEW_OV_USER",
		},
		mcnflag.StringFlag{
			Name:   "oneview-ov-password",
			Usage:  "Password to OneView Server",
			Value:  "",
			EnvVar: "ONEVIEW_OV_PASSWORD",
		},
		mcnflag.StringFlag{
			Name:   "oneview-ov-domain",
			Usage:  "Domain to OneView Server",
			Value:  "LOCAL",
			EnvVar: "ONEVIEW_OV_DOMAIN",
		},
		mcnflag.StringFlag{
			Name:   "oneview-ov-endpoint",
			Usage:  "OneView Server URL Endpoint",
			Value:  "",
			EnvVar: "ONEVIEW_OV_ENDPOINT",
		},
		mcnflag.IntFlag{
			Name:   "oneview-ov-apiversion",
			Usage:  "Force api version to an older release, ie; 201",
			Value:  1,
			EnvVar: "ONEVIEW_OV_APIVERSION",
		},
		mcnflag.StringFlag{
			Name:   "oneview-icsp-user",
			Usage:  "User Name to OneView Insight Controller",
			Value:  "",
			EnvVar: "ONEVIEW_ICSP_USER",
		},
		mcnflag.StringFlag{
			Name:   "oneview-icsp-password",
			Usage:  "Password to OneView Insight Controller",
			Value:  "",
			EnvVar: "ONEVIEW_ICSP_PASSWORD",
		},
		mcnflag.StringFlag{
			Name:   "oneview-icsp-domain",
			Usage:  "Domain to OneView Insight Controller",
			Value:  "LOCAL",
			EnvVar: "ONEVIEW_ICSP_DOMAIN",
		},
		mcnflag.StringFlag{
			Name:   "oneview-icsp-endpoint",
			Usage:  "OneView Insight Controller URL Endpoint",
			Value:  "",
			EnvVar: "ONEVIEW_ICSP_ENDPOINT",
		},
		mcnflag.IntFlag{
			Name:   "oneview-icsp-apiversion",
			Usage:  "Force api version to an older release, ie; 200",
			Value:  1,
			EnvVar: "ONEVIEW_ICSP_APIVERSION",
		},
		mcnflag.BoolFlag{
			Name:   "oneview-sslverify",
			Usage:  "SSH private key path",
			EnvVar: "ONEVIEW_SSLVERIFY",
		},
		mcnflag.StringFlag{
			Name:   "oneview-ssh-user",
			Usage:  "OneView build plan ssh user account",
			Value:  "docker",
			EnvVar: "ONEVIEW_SSH_USER",
		},
		mcnflag.IntFlag{
			Name:   "oneview-ssh-port",
			Usage:  "OneView build plan ssh host port",
			Value:  22,
			EnvVar: "ONEVIEW_SSH_PORT",
		},
		mcnflag.StringFlag{
			Name:   "oneview-server-template",
			Usage:  "OneView server template to use for blade provisioning, see OneView Server Template for setup.",
			Value:  "DOCKER_1.8_OVTEMP",
			EnvVar: "ONEVIEW_SERVER_TEMPLATE",
		},
		mcnflag.StringFlag{
			Name:   "oneview-os-deploy-plan",
			Usage:  "HPE Synergy Image Streamer OS Deploy Plan",
			Value:  "RHEL72_Docker",
			EnvVar: "ONEVIEW_OS_DEPLOY_PLAN",
		},
		mcnflag.StringFlag{
			Name:   "oneview-os-attribs",
			Usage:  "HPE Synergy Image Streamer OS Deploy Attributes",
			Value:  "[{\"name\":\"fqdn\",\"value\":\"system.company.com\"}]",
			EnvVar: "ONEVIEW_OS_ATTRIBS",
		},
		mcnflag.StringFlag{
			Name:   "oneview-server-hw-name",
			Usage:  "Server HW name to select for deployment",
			Value:  ",CN999999999F Bay 99",
			EnvVar: "ONEVIEW_SERVER_HW_NAME",
		},
		mcnflag.StringFlag{
			Name:   "oneview-os-plans",
			Usage:  "Comma separated list of OneView ICsp OS Build plans to use for OS provisioning, see ICsp OS Plan for setup.",
			Value:  "RHEL71_DOCKER_1.8",
			EnvVar: "ONEVIEW_OS_PLANS",
		},
		mcnflag.StringFlag{
			Name:   "oneview-ilo-user",
			Usage:  "ILO User id that is used during ICsp server creation.",
			Value:  "docker",
			EnvVar: "ONEVIEW_ILO_USER",
		},
		mcnflag.StringFlag{
			Name:   "oneview-ilo-password",
			Usage:  "ILO password that is used during ICsp server creation.",
			Value:  "",
			EnvVar: "ONEVIEW_ILO_PASSWORD",
		},
		mcnflag.IntFlag{
			Name:   "oneview-ilo-port",
			Usage:  "optional ILO port to use.",
			Value:  443,
			EnvVar: "ONEVIEW_ILO_PORT",
		},
		mcnflag.IntFlag{
			Name:   "oneview-public-slotid",
			Usage:  "Optional slot id of the public interface to use for connecting with docker.",
			Value:  1,
			EnvVar: "ONEVIEW_PUBLIC_SLOTID",
		},
		mcnflag.StringFlag{
			Name:   "oneview-public-connection-name",
			Usage:  "Optional Connection name to use for public interface to connect with docker, the name can be defined in server template.  Overrides slotid option.",
			Value:  "",
			EnvVar: "ONEVIEW_PUBLIC_CONNECTION_NAME",
		},
	}
}

// DriverName - get the name of the driver
func (d *Driver) DriverName() string {
	log.Debug("DriverName...%s", driverName)
	return driverName
}

// GetSSHHostname - gets the hostname that docker-machine connects to
func (d *Driver) GetSSHHostname() (string, error) {
	log.Debug("GetSSHHostname...")
	return d.GetIP()
}

// GetSSHUsername - gets the ssh user that will be connected to
func (d *Driver) GetSSHUsername() string {
	log.Debug("GetSSHUsername...")
	return d.SSHUser
}

// SetConfigFromFlags - gets the mcnflag configuration flags
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	log.Debug("SetConfigFromFlags...")

	if flags.String("oneview-icsp-endpoint") == "" {
		d.ClientICSP = nil
	} else {
		d.ClientICSP = d.ClientICSP.NewICSPClient(flags.String("oneview-icsp-user"),
			flags.String("oneview-icsp-password"),
			flags.String("oneview-icsp-domain"),
			flags.String("oneview-icsp-endpoint"),
			flags.Bool("oneview-sslverify"),
			flags.Int("oneview-icsp-apiversion"))

		if flags.Int("oneview-icsp-apiversion") == 1 {
			d.ClientICSP.RefreshVersion()
		}
	}

	d.ClientOV = d.ClientOV.NewOVClient(flags.String("oneview-ov-user"),
		flags.String("oneview-ov-password"),
		flags.String("oneview-ov-domain"),
		flags.String("oneview-ov-endpoint"),
		flags.Bool("oneview-sslverify"),
		flags.Int("oneview-ov-apiversion"))

	// we only get the version from /version if it's not setup becuse 1 is not a real version
	if flags.Int("oneview-ov-apiversion") == 1 {
		d.ClientOV.RefreshVersion()
	}

	d.IloUser = flags.String("oneview-ilo-user")
	d.IloPassword = flags.String("oneview-ilo-password")
	d.IloPort = flags.Int("oneview-ilo-port")

	d.PublicSlotID = flags.Int("oneview-public-slotid")
	d.PublicConnectionName = flags.String("oneview-public-connection-name")

	d.SSHUser = flags.String("oneview-ssh-user")
	d.SSHPort = flags.Int("oneview-ssh-port")

	d.ServerTemplate = flags.String("oneview-server-template")
	d.OSBuildPlans = strings.Split(flags.String("oneview-os-plans"), ",")
	d.OSDeployPlan = flags.String("oneview-os-deploy-plan")
	d.HWName = flags.String("oneview-server-hw-name")

	var attribs []Attrib
	if err :=json.Unmarshal([]byte(flags.String("oneview-os-attribs")), &attribs); err != nil { 
	log.Errorf("Error with un-marshalling OS Attribute data: %s", err) 
		os.Exit(1) 
	} 
	d.OSDeployAttributes = attribs

	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")

	// TODO : we should verify settings for each client

	// check for the ov endpoint
	if d.ClientOV.Endpoint == "" {
		return ErrDriverMissingEndPointOptionOV
	}
	// check for the icsp endpoint
	if d.ClientICSP != nil {
		if  d.ClientICSP.Endpoint == "" {
			return ErrDriverMissingEndPointOptionICSP
		}

		// error if one of the plans is empty string
		for _, osplan := range d.OSBuildPlans {
				if osplan == "" {
					return ErrDriverMissingBuildPlanOption
				}
		}
	}
	// check for the template name
	if d.ServerTemplate == "" {
			return ErrDriverMissingTemplateOption
	}

	return nil
}

// PreCreateCheck - pre create check
func (d *Driver) PreCreateCheck() (err error) {
	log.Debug("PreCreateCheck...")
	// verify you can connect to ov
	ovVersion, err := d.ClientOV.GetAPIVersion()
	if err != nil {
		return err
	}
	if ovVersion.CurrentVersion <= 0 {
		return fmt.Errorf("Unable to get a valid version from OneView,  %+v\n", ovVersion)
	}
	if d.ClientICSP != nil {
		// verify you can connect to icsp
		icspVersion, err := d.ClientICSP.GetAPIVersion()
		if err != nil {
			return err
		}
		if icspVersion.CurrentVersion <= 0 {
			return fmt.Errorf("Unable to get a valid version from ICsp,  %+v\n", icspVersion)
		}
	}
	return nil
}

// Create - create server for docker
func (d *Driver) Create() error {
	log.Infof("Generating SSH keys...")
	if err := d.createKeyPair(); err != nil {
		return fmt.Errorf("unable to create key pair: %s", err)
	}

	log.Debugf("OV Endpoint is: %s", d.ClientOV.Endpoint)
	// create the server profile in oneview, we need a hostname and a template name

	log.Debugf("***> GetProfileTemplateByName(%s)", d.ServerTemplate)
	serverProfileTemplate, error := d.ClientOV.GetProfileTemplateByName(d.ServerTemplate)
	if error != nil || serverProfileTemplate.URI.IsNil() {
	  return fmt.Errorf("Could not find Server Profile Template\n%+v", d.ServerTemplate)
	}
	var serverHardware ov.ServerHardware
	// If it is the bogus "CN999999999F Bay 99" example treat as unset
	if d.HWName != "" && d.HWName != "CN999999999F Bay 99"{
	  log.Debugf("***> GetServerHardwareByName(%s)", d.HWName)
	  serverHardware, error = d.ClientOV.GetServerHardwareByName(d.HWName)
	  if(error != nil){
		return error
	  }
	} else {
	  return fmt.Errorf("no machine name specified")
	}

	if d.OSDeployPlan != "" {
		deploymentAttributes := make(map[string]string)
		log.Debugf("***> GetOsDeploymentPlanByName(%s)", d.OSDeployPlan)
		osDeploymentPlan, error := d.ClientOV.GetOSDeploymentPlanByName(d.OSDeployPlan)
		if error != nil || osDeploymentPlan.URI.IsNil() {
			return fmt.Errorf("Could not find osDeploymentPlan: %s", d.OSDeployPlan)
		}

		if d.OSDeployAttributes != nil {
			for i, v := range d.OSDeployAttributes {
				log.Debugf("Attribute %d is %+v", int(i), v)
				deploymentAttributes[v.Name] = v.Value
			} 
		}
 		if d.SSHPublicKey != "" {
                        // Split ssh keys into smaller parts to workaround issues in I3S
			parts := splitStringIntoParts(d.SSHPublicKey, 250)
			for i := 0; i < 2; i++ {
				attribPrefix := fmt.Sprintf("sshkey%d", i+1)
                                v := ""
                                if i < len(parts) {
					v = parts[i]
                                }
				log.Debugf("Attribute %s is %+v", attribPrefix, v)
				deploymentAttributes[attribPrefix] = v
                        }
		}
		log.Debugf("***> CreateProfileFromTemplateWithI3S")
		SPerror := d.ClientOV.CreateProfileFromTemplateWithI3S(d.MachineName, serverProfileTemplate, serverHardware, osDeploymentPlan, deploymentAttributes)
		if SPerror != nil {
			return SPerror
		}
		log.Debugf("***> get profile and server hardware info")
		if err := d.getBlade(); err != nil {
			return err
		}
		// Sleep for 1s to let power lock clear
		time.Sleep(1 * time.Second)

		log.Debugf("***> power on server")
		// power on the server, and leave it in that state
		if err := d.Hardware.PowerOn(); err != nil {
			return err
		}
		ip, err := d.GetIP()
		if err != nil {
			return err
		}
		d.IPAddress = ip

		// use ssh wait for the system to boot
                sshAvailable := 0
		for i := 0; i < 200 && sshAvailable == 0; i++ {
		    sshClient, err := d.getLocalSSHClient()
		    if err != nil {
		        time.Sleep(10 * time.Second)
			} else {
				sshAvailable = 1 
                                sshClient.Output("hostname")
			}
		}
	} else {
		// create d.Hardware and d.Profile
		if err := d.ClientOV.CreateMachine(d.MachineName, d.ServerTemplate); err != nil {
			return err
		}

		if err := d.getBlade(); err != nil {
			return err
		}

		// power off let customization bring the server online
		if err := d.Hardware.PowerOff(); err != nil {
			return err
		}

		// add the server to icsp, TestCreateServer
		// apply a build plan, TestApplyDeploymentJobs
		var sp *icsp.CustomServerAttributes
		sp = sp.New()
		sp.Set("docker_user", d.SSHUser)
		sp.Set("public_key", d.SSHPublicKey)
		// TODO: make a util for this
		if len(os.Getenv("proxy_enable")) > 0 {
			sp.Set("proxy_enable", os.Getenv("proxy_enable"))
		} else {
			sp.Set("proxy_enable", "false")
		}

		strProxy := os.Getenv("proxy_config")
		sp.Set("proxy_config", strProxy)

		sp.Set("docker_hostname", d.MachineName+"-@server_name@")

		sp.Set("interface", "@interface@") // this is populated later

		// Get the mac address for public Connection on server profile
		var publicmac string
		if d.PublicConnectionName != "" {
			conn, err := d.Profile.GetConnectionByName(d.PublicConnectionName)
			if err != nil {
				return err
			}
			publicmac = conn.MAC.String()
		} else {
			publicmac = ""
		}

		// arguments for customize server
		cs := icsp.CustomizeServer{
			HostName:         d.MachineName,                   // machine-rack-enclosure-bay
			SerialNumber:     d.Profile.SerialNumber.String(), // get it
			ILoUser:          d.IloUser,
			IloPassword:      d.IloPassword,
			IloIPAddress:     d.Hardware.GetIloIPAddress(), // MpIpAddress for v1
			IloPort:          d.IloPort,
			OSBuildPlans:     d.OSBuildPlans, // array of OS Build Plans to apply
			PublicSlotID:     d.PublicSlotID, // this is the slot id of the public interface
			PublicMAC:        publicmac,      // Server profile mac address, overrides slotid
			ServerProperties: sp,
		}
		// create d.Server and apply a build plan and configure the custom attributes
		if err := d.ClientICSP.CustomizeServer(cs); err != nil {
			return err
		}

		ip, err := d.GetIP()
		if err != nil {
			return err
		}
		d.IPAddress = ip

		// use ssh to set keys, and test ssh
		sshClient, err := d.getLocalSSHClient()
		if err != nil {
			return err
		}

		pubKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
		if err != nil {
			return err
		}

		if out, err := sshClient.Output(fmt.Sprintf(
			"printf '%%s' '%s' | tee /home/%s/.ssh/authorized_keys",
			string(pubKey),
			d.GetSSHUsername(),
		)); err != nil {
			log.Error(out)
			return err
		}
	}
	log.Infof("%s, Completed all create steps, provisioning will complete at first boot.", d.DriverName())

	defer closeAll(d)
	return nil
}

// closeAll - cleanup sessions on the OV and ICSP appliances
func closeAll(d *Driver) {
	err := d.ClientOV.SessionLogout()
	if err != nil {
		log.Warnf("OV Session Logout : %s", err)
	}
	
	if d.ClientICSP != nil {
	err = d.ClientICSP.SessionLogout()
	  if err != nil {
		log.Warnf("ICsp Session Logout : %s", err)
	  }
	}
}

// GetURL - get docker url
func (d *Driver) GetURL() (string, error) {
	log.Debug("GetURL...")
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

// GetIP - get server host or ip address
// TODO: we need to get ip of server from icsp or ov??
// currently the only way i can see to get this is with sudo ifconfig|grep inet
func (d *Driver) GetIP() (string, error) {
	log.Debug("GetIP...")

	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return "", err
	}

	//REVISIT: Image streamer has no agent - use custom attrib to guess
	if d.ClientICSP == nil {
		if d.Profile.ServerProfilev300 != nil {
			var SP3 ov.ServerProfilev300 = *(d.Profile.ServerProfilev300)	
			for _, v := range SP3.OSDeploymentSettings.OSCustomAttributes {
				// if it starts with "ip" assume that is our IP
				if strings.Index(strings.ToLower(v.Name), "ip") == 0 {
					log.Debugf("Got IP %s from attribute %+v", v.Value, v)
					return v.Value, nil
				}
			}  
		} else {
			// maybe don't treat no IP as an error for image streamer?
			return "", nil
		}
	}  
	
	sPublicIPv4, err := d.Server.GetPublicIPV4()
	if err != nil {
		return "", err
	}
	if sPublicIPv4 == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return sPublicIPv4, nil
}

// GetState - get the running state of the target machine
func (d *Driver) GetState() (state.State, error) {
	log.Debug("GetState...")

	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return state.Error, err
	}
	if d.ClientICSP != nil {
		if icsp.Provisioning.Equal(d.Server.OpswLifecycle) {
			return state.Starting, nil
		}
		if icsp.Unprovisioned.Equal(d.Server.OpswLifecycle) ||
			icsp.PreUnProvisioned.Equal(d.Server.OpswLifecycle) {
			return state.Stopping, nil
		}
		if icsp.Deactivated.Equal(d.Server.OpswLifecycle) {
			return state.Stopped, nil
		}
		if icsp.ProvisionedFailed.Equal(d.Server.OpswLifecycle) {
			return state.Error, nil
		}
	}
	// use power state to determine status
	ps, err := d.Hardware.GetPowerState()
	if err != nil {
		return state.Error, err
	}
	switch ps {
	case ov.P_ON:
		return state.Running, nil
	case ov.P_OFF:
		return state.Stopped, nil
	case ov.P_UKNOWN:
		return state.Error, nil
	default:
		return state.None, nil
	}

}

// Start - start the docker machine target
func (d *Driver) Start() error {
	log.Infof("Starting ... %s", d.MachineName)

	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return err
	}

	// power on the server, and leave it in that state
	if err := d.Hardware.PowerOn(); err != nil {
		return err
	}
	if d.ClientICSP != nil {
		// implement icsp check for is in maintenance mode or started
		isManaged, err := d.ClientICSP.IsServerManaged(d.Hardware.SerialNumber.String())
		if err != nil {
			return err
		}
		if !isManaged {
			return errors.New("Server was started but not ready, check icsp status")
		}
	}
	return nil
}

// Stop - stop the docker machine target
func (d *Driver) Stop() error {
	log.Debug("Stop...")
	log.Infof("Stop ... %s", d.MachineName)
	// gracefully attempt to stop the os

	if _, err := drivers.RunSSHCommandFromDriver(d, "sudo shutdown -P now"); err != nil {
		log.Warnf("Problem shutting down gracefully : %s", err)
	}

	// get the blade for this driver
	if err := d.getBlade(); err != nil {
		return err
	}

	// power on the server, and leave it in that state
	if err := d.Hardware.PowerOff(); err != nil {
		return err
	}
	// cleanup
	defer closeAll(d)
	return nil
}

// Remove - remove the docker machine target
//    Should remove the ICSP provisioned plan and the Server Profile from OV
func (d *Driver) Remove() error {
	log.Debug("Remove...")
	// remove the ssh keys
	if err := d.deleteKeyPair(); err != nil {
		return err
	}
	if err := d.Stop(); err != nil {
		return err
	}
	if err := d.getBlade(); err != nil {
		return err
	}
	if d.ClientICSP != nil {
		// destroy the server in icsp
		isDeleted, err := d.ClientICSP.DeleteServer(d.Server.MID)
		if err != nil {
			return err
		}
		if !isDeleted {
			return fmt.Errorf("Unable to delete the server from icsp : %s, %s", d.MachineName, d.Server.MID)
		}
	}
	// delete the server profile in ov : TestDeleteProfile
	t, err := d.ClientOV.SubmitDeleteProfile(d.Profile)
	err = t.Wait()
	if err != nil {
		return err
	}
	// cleanup
	defer closeAll(d)
	return nil
}

// Restart - restart the target machine
func (d *Driver) Restart() error {
	log.Debug("Restarting...")
	if err := d.Stop(); err != nil {
		return err
	}
	return d.Start()
}

// Kill - kill the docker machine
func (d *Driver) Kill() error {
	log.Debug("Killing...")
	//TODO: implement power off , is there a force?
	return nil
}

// publicSSHKeyPath - get the path to public ssh key
func (d *Driver) publicSSHKeyPath() string {
	log.Debug("publicSSHKeyPath...")
	return d.GetSSHKeyPath() + ".pub"
}

// /////////  HELPLERS /////////////

func (d *Driver) getBlade() (err error) {
	log.Debug("In getBlade()")

	d.Profile, err = d.ClientOV.GetProfileByName(d.MachineName)
	if err != nil {
		return err
	}

	log.Debugf("***> check if we got a profile")
	if d.Profile.URI.IsNil() {
		err = fmt.Errorf("Attempting to get machine profile information, unable to find machine in oneview: %s", d.MachineName)
		return err
	}

	// get the server hardware associated with that test profile
	log.Debugf("***> GetServerHardware")
	d.Hardware, err = d.ClientOV.GetServerHardware(d.Profile.ServerHardwareURI)
	if d.Hardware.URI.IsNil() {
		err = fmt.Errorf("Attempting to get machine blade information, unable to find machine: %s", d.MachineName)
		return err
	}

	if d.ClientICSP != nil {
		// get server entry in ICsp
		if d.Hardware.VirtualSerialNumber.IsNil() {
			// get the server profile with SerialNumber
			d.Server, err = d.ClientICSP.GetServerBySerialNumber(d.Hardware.SerialNumber.String())
		} else {
			// get the server profile with the VirtualSerialNumber
			d.Server, err = d.ClientICSP.GetServerBySerialNumber(d.Hardware.VirtualSerialNumber.String())
		}
		if err != nil {
			return err
		}
	}
	return err
}

// createKeyPair - generate key files needed
func (d *Driver) createKeyPair() error {

	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	publicKey, err := ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return err
	}

	log.Debugf("created keys => %s", string(publicKey))
	d.SSHPublicKey = string(publicKey)
	return nil
}

// deleteKeyPair
func (d *Driver) deleteKeyPair() error {
	if err := os.Remove(d.GetSSHKeyPath()); err != nil {
		return err
	}
	if err := os.Remove(d.GetSSHKeyPath() + ".pub"); err != nil {
		return err
	}
	return nil
}

func (d *Driver) getLocalSSHClient() (ssh.Client, error) {
	sshAuth := &ssh.Auth{
		Passwords: []string{"docker"},
		Keys:      []string{d.GetSSHKeyPath()},
	}
	sshClient, err := ssh.NewNativeClient(d.GetSSHUsername(), d.IPAddress, d.SSHPort, sshAuth)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}

func splitStringIntoParts(orig string, partSize int) []string {
	a := []rune(orig )
	var parts []string
	var res string = ""
	for i, r := range a {
		res = res + string(r)
		if i > 0 && (i+1)%partSize == 0 {
		        parts = append(parts, res)
			res = ""
		}
	}
	if res != "" {
		parts = append(parts, res)     
        }
	return parts
}
