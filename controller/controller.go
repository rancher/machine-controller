package controller

import (
	"time"

	"os"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Register(management *config.ManagementContext) {
	machineLifecycle := &MachineLifecycle{}
	machineClient := management.Management.Machines("")
	machineLifecycle.machineClient = machineClient
	machineLifecycle.machineTemplateClient = management.Management.MachineTemplates("")

	machineClient.
		Controller().
		AddHandler(v3.NewMachineLifecycleAdapter("machine-controller", machineClient, machineLifecycle))
}

type MachineLifecycle struct {
	machineClient         v3.MachineInterface
	machineTemplateClient v3.MachineTemplateInterface
}

func (m *MachineLifecycle) Create(obj *v3.Machine) error {
	// No need to create a deepcopy of obj, obj is already a deepcopy
	if obj.Status.Provisioned && obj.Status.ExtractedConfig != "" {
		return nil
	}
	machineDir, err := buildBaseHostDir(obj.Name)
	if err != nil {
		return err
	}
	logrus.Debugf("Creating machine storage directory %s", machineDir)
	if !obj.Status.Provisioned {
		configRawMap := map[string]interface{}{}
		if obj.Spec.MachineTemplateName != "" {
			machineTemplate, err := m.machineTemplateClient.Get(obj.Spec.MachineTemplateName, metav1.GetOptions{})
			if err != nil {
				return err
			}
			for k, v := range machineTemplate.Spec.PublicValues {
				configRawMap[k] = v
			}
			for k, v := range machineTemplate.Spec.SecretValues {
				configRawMap[k] = v
			}
		} else {
			var err error
			switch obj.Spec.Driver {
			case "amazonec2":
				configRawMap, err = toMap(obj.Spec.AmazonEC2Config)
				if err != nil {
					return err
				}
			case "digitalocean":
				configRawMap, err = toMap(obj.Spec.DigitalOceanConfig)
				if err != nil {
					return err
				}
			case "azure":
				configRawMap, err = toMap(obj.Spec.AzureConfig)
				if err != nil {
					return err
				}
			}
		}

		createCommandsArgs := buildCreateCommand(obj, configRawMap)
		cmd := buildCommand(machineDir, createCommandsArgs)
		logrus.Infof("Provisioning machine %s", obj.Name)
		stdoutReader, stderrReader, err := startReturnOutput(cmd)
		if err != nil {
			return err
		}
		if err := reportStatus(stdoutReader, stderrReader, obj, m.machineClient); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
		obj, err = m.machineClient.Get(obj.Name, metav1.GetOptions{})
		if err != nil {
			logrus.Error(err)
			return err
		}
		obj.Status.Provisioned = true
		now := time.Now().Format(time.RFC3339)
		obj.Status.Conditions = append(obj.Status.Conditions, v3.MachineCondition{
			LastTransitionTime: now,
			LastUpdateTime:     now,
			Type:               ProvisionedState,
			Status:             v1.ConditionTrue,
		})
		if obj, err = m.machineClient.Update(obj); err != nil {
			logrus.Error(err)
			return err
		}
		logrus.Infof("Provisioning machine %s done", obj.Name)
	}
	if obj.Status.ExtractedConfig == "" {
		logrus.Infof("Generating and uploading machine config %s", obj.Name)
		sshkey, err := getSSHPrivateKey(machineDir, obj)
		if err != nil {
			return err
		}
		destFile, err := createExtractedConfig(machineDir, obj)
		if err != nil {
			return err
		}
		extractedConf, err := encodeFile(destFile)
		if err != nil {
			return err
		}
		obj, err = m.machineClient.Get(obj.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		obj.Status.ExtractedConfig = extractedConf
		obj.Status.SSHPrivateKey = sshkey
		if obj, err = m.machineClient.Update(obj); err != nil {
			return err
		}
		logrus.Infof("Generating and uploading machine config %s done", obj.Name)
	}
	obj, err = m.machineClient.Get(obj.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	os.RemoveAll(machineDir)
	return nil
}

func (m *MachineLifecycle) Updated(obj *v3.Machine) error {
	// YOU MUST CALL DEEPCOPY
	return nil
}

func (m *MachineLifecycle) Remove(obj *v3.Machine) error {
	// No need to create a deepcopy of obj, obj is already a deepcopy
	machineDir, err := buildBaseHostDir(obj.Name)
	if err != nil {
		return err
	}
	logrus.Debugf("Creating machine storage directory %s", machineDir)
	err = restoreMachineDir(obj, machineDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(machineDir)

	mExists, err := machineExists(machineDir, obj.Name)
	if err != nil {
		return err
	}

	if mExists {
		logrus.Infof("Removing machine %s", obj.Name)
		if err := deleteMachine(machineDir, obj); err != nil {
			return err
		}
		logrus.Infof("Removing machine %s done", obj.Name)
	}
	return nil
}
