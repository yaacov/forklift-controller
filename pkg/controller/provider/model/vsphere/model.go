package vsphere

import (
	"encoding/json"
	libmodel "github.com/konveyor/controller/pkg/inventory/model"
	"github.com/konveyor/virt-controller/pkg/controller/provider/model/base"
)

// Errors
var NotFound = libmodel.NotFound

const (
	Assign = "assign"
)

//
// Types
type Model = libmodel.Model
type Annotation = base.Annotation

//
// Base VMWare model.
type Base struct {
	// Managed object ID.
	ID string `sql:"pk"`
	// Name
	Name string `sql:""`
	// Parent
	Parent string `sql:"index(a)"`
	// Annotations
	Annotations string `sql:""`
}

//
// Get the PK.
func (m *Base) Pk() string {
	return m.ID
}

func (m *Base) String() string {
	return m.ID
}

func (m *Base) Labels() libmodel.Labels {
	return nil
}

func (m *Base) Equals(other libmodel.Model) bool {
	if vm, cast := other.(*VM); cast {
		return m.ID == vm.ID
	}

	return false
}

//
// An object reference.
type Ref struct {
	// The kind (type) of the referenced.
	Kind string
	// The ID of object referenced.
	ID string
}

//
// Encode the ref.
func (r *Ref) Encode() string {
	j, _ := json.Marshal(r)
	return string(j)
}

//
// Unmarshal the json `j` into self.
func (r *Ref) With(j string) *Ref {
	json.Unmarshal([]byte(j), r)
	return r
}

//
// Ref pointer.
func RefPtr() *Ref {
	r := Ref{}
	return &r
}

//
// List of `Ref`.
type RefList []Ref

//
// Encode the list.
func (r *RefList) Encode() string {
	j, _ := json.Marshal(r)
	return string(j)
}

//
// Unmarshal the json `j` into self.
func (r *RefList) With(j string) *RefList {
	json.Unmarshal([]byte(j), r)
	return r
}

//
// RefList pointer.
func RefListPtr() *RefList {
	r := RefList{}
	return &r
}

type Folder struct {
	Base
	Children string `sql:""`
}

type Datacenter struct {
	Base
	Clusters   string `sql:""`
	Networks   string `sql:""`
	Datastores string `sql:""`
	Vms        string `sql:""`
}

type Cluster struct {
	Base
	Hosts       string `sql:""`
	Networks    string `sql:""`
	Datastores  string `sql:""`
	DasEnabled  bool   `sql:""`
	DasVms      string `sql:""`
	DrsEnabled  bool   `sql:""`
	DrsBehavior string `sql:""`
	DrsVms      string `sql:""`
}

type Host struct {
	Base
	InMaintenanceMode bool   `sql:""`
	ProductName       string `sql:""`
	ProductVersion    string `sql:""`
	Networks          string `sql:""`
	Datastores        string `sql:""`
	Vms               string `sql:""`
}

type Network struct {
	Base
	Tag string `sql:""`
}

type Datastore struct {
	Base
	Type            string `sql:""`
	Capacity        int64  `sql:""`
	Free            int64  `sql:""`
	MaintenanceMode string `sql:""`
}

type VM struct {
	Base
	UUID                string `sql:""`
	Firmware            string `sql:""`
	CpuAffinity         string `sql:""`
	CpuHotAddEnabled    bool   `sql:""`
	CpuHotRemoveEnabled bool   `sql:""`
	MemoryHotAddEnabled bool   `sql:""`
	CpuCount            int32  `sql:""`
	MemorySizeMB        int32  `sql:""`
	GuestName           string `sql:""`
	BalloonedMemory     int32  `sql:""`
	IpAddress           string `sql:""`
}