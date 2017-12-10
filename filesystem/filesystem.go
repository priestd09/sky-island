package filesystem

import (
	"github.com/Sirupsen/logrus"
	"github.com/briandowns/sky-island/config"
	"github.com/briandowns/sky-island/utils"
	"gopkg.in/alexcesaro/statsd.v2"
)

// FSServicer defines the behavior of the filesystem service
type FSServicer interface {
	CreateBaseJailDataset() error
	CloneBaseToJail(string) error
	CreateDataset() error
	CreateSnapshot() error
	RemoveDataset(string) error
}

// fsService
type fsService struct {
	logger  *logrus.Logger
	conf    *config.Config
	wrapper utils.Wrapper
	metrics *statsd.Client
}

// NewFilesystemService creates a new value of type FileSystemService which provides the dependencies
// to the service methods
func NewFilesystemService(conf *config.Config, l *logrus.Logger, metrics *statsd.Client, w utils.Wrapper) FSServicer {
	return &fsService{
		logger:  l,
		conf:    conf,
		wrapper: w,
		metrics: metrics,
	}
}

// CreateBaseJailDataset creates a Dataset and mounts it for the base jail
func (f *fsService) CreateBaseJailDataset() error {
	t := f.metrics.NewTiming()
	defer t.Send("dataset_create")
	_, err := f.wrapper.Output("zfs", "create", "-o", "mountpoint="+f.conf.Jails.BaseJailDir, f.conf.Filesystem.ZFSDataset, "")
	return err
}

// CloneBaseToJail does a ZFS clone from the base jail to the new jail
func (f *fsService) CloneBaseToJail(jname string) error {
	t := f.metrics.NewTiming()
	defer t.Send("dataset_create")
	_, err := f.wrapper.Output("zfs", "clone", f.conf.Filesystem.ZFSDataset+"/jails/releases/"+f.conf.Release+"@p1", f.conf.Filesystem.ZFSDataset+"/jails/"+jname)
	return err
}

// CreateDataset creates a new ZFS Dataset
func (f *fsService) CreateDataset() error {
	_, err := f.wrapper.Output("zfs", "create", "-p", f.conf.Filesystem.ZFSDataset+"/jails/releases/"+f.conf.Release)
	return err
}

// CreateSnapshot creates a ZFS snapshot
func (f *fsService) CreateSnapshot() error {
	t := f.metrics.NewTiming()
	defer t.Send("snapshot_create")
	_, err := f.wrapper.Output("zfs", "snapshot", f.conf.Filesystem.ZFSDataset+"/jails/releases/"+f.conf.Release+"@p1")
	return err
}

// RemoveDataset removes the Dataset associated with the given id
func (f *fsService) RemoveDataset(id string) error {
	t := f.metrics.NewTiming()
	defer t.Send("dataset_remove")
	_, err := f.wrapper.Output("zfs", "destroy", "-rf", f.conf.Filesystem.ZFSDataset+"/jails/"+id)
	return err
}
