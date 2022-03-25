package account

import (
	"context"
	"errors"
	"log"
)

type ACL interface {
	Guard() error
	Create() (*ACLManager, error)
	CreateNewGroup(name string) *Group
	AddItemsToGroup(gId string, gItem ...*GroupItem) error
	AddPermission(name string, gId string, r ResourceID, pt PermissionType) error
}

type ResourceID uint32

const (
	// R_NONE ...
	R_NONE ResourceID = iota
	R_PROXY
	R_FILE
	R_URL
	R_CONTRACT
	R_SYSTEM
	R_FULLCHAIN
	R_BLOCK
)

type PermissionType uint32

const (
	PT_NONE PermissionType = iota
	PT_READ
	PT_DELETE
	PT_UPDATE
	PT_WRITE
	PT_SYSTEM_UPGRADE
	PT_VALIDATE
	PT_VERIFY
	PT_LEAD
	PT_MINTER
	PT_BURNER
	PT_ADMIN
	PT_SUPERADMIN
	PT_BAN
)

var Groups map[string]*Group

type GroupItem struct {
	ID           string //keypair id
	Permission   string
	PermissionID PermissionType
	Resource     ResourceID
	IsKeyPair    bool
	Weight       int
}

type Group struct {
	Name  string
	Items []*GroupItem
}

type Permission struct {
	Name      string
	Groups    []string
	Items     []*GroupItem
	Threshold int
}

// ACLManager implements ACL interface
type ACLManager struct {
	IsInitialized  bool
	AclContext     context.Context
	currentAccount *Account
	accounts       map[string]*Account
	blacklist      map[string]bool
}

func (A ACLManager) Guard() error {
	if A.IsInitialized {
		return errors.New("acl manager already initialized")
	}
	return nil
}

func (A ACLManager) Create() (*ACLManager, error) {

	am := A.NewACLManager()
	if am == nil {
		return nil, errors.New("failed initializing acl manager")
	}
	return am, nil
}

func (A ACLManager) CreateNewGroup(name string) *Group {
	if len(Groups) == 0 {
		Groups = map[string]*Group{}
	}
	if Groups[name] != nil {
		return Groups[name]
	}
	return createGroup(name)
}

func (A ACLManager) AddItemsToGroup(gName string, gItems ...*GroupItem) error {
	group := Groups[gName]
	if group != nil {
		addGroupItems(group, gItems...)
	} else {
		return errors.New("group cannot be found")
	}
	return nil

}

func (A ACLManager) AddPermission(name string, gId string, r ResourceID, pt PermissionType) error {
	//TODO implement me
	panic("implement me")
}

func (A ACLManager) NewACLManager() *ACLManager {
	err := A.Guard()
	if err != nil {
		log.Println(err)
		return nil
	}

	a := &ACLManager{}
	a.AclContext = context.Background()
	return a
}

func createGroup(name string) *Group {
	g := &Group{Name: name}
	Groups[name] = g
	return g
}

func addGroupItems(group *Group, gi ...*GroupItem) *Group {

	group.Items = make([]*GroupItem, len(gi))
	for i := 0; i < len(gi); i++ {
		group.Items = append(group.Items, gi[i])
	}
	return group
}
