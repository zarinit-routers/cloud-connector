package repository

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/zarinit-routers/cloud-connector/storage/database"
	"gorm.io/gorm"
)

func mustConnect() *gorm.DB {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to database", "error", err.Error())
	}
	return db
}

type ModelBase struct {
	ID uuid.UUID `gorm:"primary_key" json:"id"`
}

type Node struct {
	*ModelBase
	OrganizationID  uuid.UUID  `json:"organizationId"`
	Name            string     `json:"name"`
	FirstConnection time.Time  `json:"firstConnection"`
	LastConnection  *time.Time `json:"lastConnection"`
	Tags            []*Tag     `gorm:"foreignKey:NodeID" json:"tags"`
}

type Tag struct {
	NodeID uuid.UUID `gorm:"primary_key" json:"nodeId"`
	Tag    string    `gorm:"primary_key" json:"tag"`
}

func GetNode(id uuid.UUID) (*Node, error) {
	db := mustConnect()
	var node Node
	err := db.Preload("Tags").Where("id = ?", id).First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}
func GetNodes(organizationID uuid.UUID) ([]*Node, error) {
	db := mustConnect()
	var nodes []*Node
	err := db.Preload("Tags").Where("organization_id = ?", organizationID).Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func NewTag(nodeId uuid.UUID, tag string) (*Tag, error) {
	model := &Tag{
		NodeID: nodeId,
		Tag:    tag,
	}
	db := mustConnect()
	err := db.Create(model).Error
	if err != nil {
		log.Error("failed to create tag", "error", err.Error())
		return nil, fmt.Errorf("failed to create tag: %s", err)
	}
	return model, nil
}

func RemoveTag(nodeID uuid.UUID, tag string) error {
	db := mustConnect()
	err := db.Unscoped().Model(&Tag{}).Delete("WHERE node_id = ? AND tag = ?", nodeID, tag).Error
	if err != nil {
		log.Error("failed to delete tag", "error", err.Error())
		return fmt.Errorf("failed to create tag: %s", err)
	}
	return nil
}

func NewNode(id uuid.UUID, organizationID uuid.UUID, name string) (*Node, error) {
	now := time.Now()
	model := &Node{
		ModelBase: &ModelBase{
			ID: id,
		},
		OrganizationID:  organizationID,
		Name:            name,
		LastConnection:  &now,
		FirstConnection: now,
	}
	db := mustConnect()
	err := db.Create(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func ReconnectNode(id uuid.UUID, organizationID uuid.UUID) (*Node, error) {
	db := mustConnect()
	var node Node
	err := db.Model(&node).Where("id = ?", id).Update("organization_id", organizationID).Update("last_connection", time.Now()).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func UpdateLastConnection(id uuid.UUID) (*Node, error) {
	db := mustConnect()
	var node Node
	err := db.Model(&node).Where("id = ?", id).Update("last_connection", time.Now()).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}
