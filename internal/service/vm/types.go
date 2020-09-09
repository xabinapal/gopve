package vm

import (
	"fmt"
	"net/http"

	"github.com/xabinapal/gopve/internal/client"
	"github.com/xabinapal/gopve/pkg/types/node"
	"github.com/xabinapal/gopve/pkg/types/vm"
)

type Service struct {
	client client.Client
	api    client.API
}

func NewService(cli client.Client, api client.API) *Service {
	return &Service{
		client: cli,
		api:    api,
	}
}

type VirtualMachine struct {
	svc  *Service
	full bool

	node       string
	kind       vm.Kind
	vmid       uint
	name       string
	isTemplate bool
}

func (obj *VirtualMachine) load() (vm.VirtualMachine, error) {
	vm, err := obj.svc.Get(obj.vmid)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

type QEMUVirtualMachine struct {
	VirtualMachine

	cpu    vm.QEMUCPUProperties
	memory vm.QEMUMemoryProperties
}

func (obj *QEMUVirtualMachine) Load() error {
	if obj.full {
		return nil
	}

	vm, err := obj.VirtualMachine.load()
	if err != nil {
		return err
	}

	switch x := vm.(type) {
	case *QEMUVirtualMachine:
		*obj = *x
	default:
		panic(fmt.Sprintf("This should never happen: %s", err.Error()))
	}

	return nil
}

type LXCVirtualMachine struct {
	VirtualMachine

	cpu    vm.LXCCPUProperties
	memory vm.LXCMemoryProperties
}

func (obj *LXCVirtualMachine) Load() error {
	if obj.full {
		return nil
	}

	vm, err := obj.VirtualMachine.load()
	if err != nil {
		return err
	}

	switch x := vm.(type) {
	case *LXCVirtualMachine:
		*obj = *x
	default:
		panic(fmt.Sprintf("This should never happen: %s", err.Error()))
	}

	return nil
}

func (vm *VirtualMachine) GetNode() (node.Node, error) {
	return vm.svc.api.Node().Get(vm.node)
}

func (vm *VirtualMachine) Node() string {
	return vm.node
}

func (vm *VirtualMachine) Kind() vm.Kind {
	return vm.kind
}

func (vm *VirtualMachine) VMID() uint {
	return vm.vmid
}

func (vm *VirtualMachine) Name() string {
	return vm.name
}

func (vm *VirtualMachine) IsTemplate() bool {
	return vm.isTemplate
}

func (obj *VirtualMachine) Status() (vm.Status, error) {
	var res struct {
		Status vm.Status `json:"status"`
	}

	if err := obj.svc.client.Request(http.MethodPost, fmt.Sprintf("node/%s/qemu/%d/status/current", obj.node, obj.vmid), nil, &res); err != nil {
		return vm.StatusStopped, err
	}

	if err := res.Status.IsValid(); err != nil {
		return vm.StatusStopped, err
	}

	return res.Status, nil
}

func (obj *VirtualMachine) Delete(purge bool, force bool) error {
	return fmt.Errorf("not implemented")
}

func (obj *QEMUVirtualMachine) CPU() (vm.QEMUCPUProperties, error) {
	if err := obj.Load(); err != nil {
		return vm.QEMUCPUProperties{}, err
	}

	return obj.cpu, nil
}

func (obj *QEMUVirtualMachine) Memory() (vm.QEMUMemoryProperties, error) {
	if err := obj.Load(); err != nil {
		return vm.QEMUMemoryProperties{}, err
	}

	return obj.memory, nil
}

func (obj *QEMUVirtualMachine) SetProperties(props vm.QEMUProperties) error {
	form, err := props.MapToValues()
	if err != nil {
		return err
	}

	if err := obj.svc.client.Request(http.MethodPost, fmt.Sprintf("node/%s/qemu/%d/config", obj.node, obj.vmid), form, nil); err != nil {
		return err
	}

	return nil
}

func (obj *LXCVirtualMachine) CPU() (vm.LXCCPUProperties, error) {
	if err := obj.Load(); err != nil {
		return vm.LXCCPUProperties{}, err
	}

	return obj.cpu, nil
}

func (obj *LXCVirtualMachine) Memory() (vm.LXCMemoryProperties, error) {
	if err := obj.Load(); err != nil {
		return vm.LXCMemoryProperties{}, err
	}

	return obj.memory, nil
}

func (obj *LXCVirtualMachine) SetProperties(props vm.LXCProperties) error {
	form, err := props.MapToValues()
	if err != nil {
		return err
	}

	if err := obj.svc.client.Request(http.MethodPost, fmt.Sprintf("node/%s/lxc/%d/config", obj.node, obj.vmid), form, nil); err != nil {
		return err
	}

	return nil
}
