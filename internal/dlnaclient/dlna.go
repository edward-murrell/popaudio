package dlnaclient

import (
	"context"
	"encoding/xml"
	"errors"
	"github.com/huin/goupnp/dcps/av1"
	"popaudio/pkg/model"
)

func GetServerByName(ctx context.Context, name string) (*DlnaClient, error) {
	c, errs, err := av1.NewContentDirectory1ClientsCtx(ctx)
	if err != nil {
		return nil, err
	}
	if len(errs) > 0 {
		for i, err := range errs {
			println("error %d: %s", i, err)
		}
		return nil, err
	}

	for _, server := range c {
		if server.RootDevice.Device.FriendlyName != name {
			continue
		}
		return NewDlnaClient(server), nil
	}
	return nil, errors.New("server not found")
}

type DlnaClient struct {
	server *av1.ContentDirectory1
}

func NewDlnaClient(server *av1.ContentDirectory1) *DlnaClient {
	return &DlnaClient{server: server}
}

func (d *DlnaClient) GetRoot() (*model.DIDLLite, error) {
	result, _, _, _, err := d.server.Browse("0", "BrowseDirectChildren", "*", 0, 16384, "")
	var didl model.DIDLLite
	err = xml.Unmarshal([]byte(result), &didl)
	if err != nil {
		return nil, err
	}
	return &didl, nil
}

func (d *DlnaClient) Browse(dir *model.Container) (*model.DIDLLite, error) {
	result, _, _, _, err := d.server.Browse(dir.ID, "BrowseDirectChildren", "*", 0, 16384, "")
	var didl model.DIDLLite
	err = xml.Unmarshal([]byte(result), &didl)
	if err != nil {
		return nil, err
	}
	return &didl, nil
}
