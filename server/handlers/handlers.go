package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/cloud-connector/models"
	"github.com/zarinit-routers/cloud-connector/storage/repository"
	"github.com/zarinit-routers/middleware/auth"
)

func GetClientsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *auth.AuthData
		if u, err := auth.GetUser(c); err != nil {
			log.Error("Failed get user", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			user = u
		}
		log.Info("User", "user", user)
		organizationId, err := models.ParseUUID(user.OrganizationId)
		if err != nil {
			log.Error("Failed parse organization id", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		data, err := repository.GetQueries().GetNodes(c.Request.Context(), organizationId)
		if err != nil {
			log.Error("Failed get nodes from repository", "error", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Info("Nodes", "nodes", data)
		nodes := []Node{}
		for _, d := range data {

			tags, err := repository.GetQueries().GetTags(c.Request.Context(), d.Id)
			if err != nil {
				log.Error("Failed get tags from repository", "error", err, "nodeId", d.Id.String())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			nodes = append(nodes, Node{
				Id:   d.Id.String(),
				Name: d.Name.String,
				Tags: tags,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"nodes": nodes,
		})
	}
}

type Node struct {
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func GetSingleClientHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *auth.AuthData
		if u, err := auth.GetUser(c); err != nil {
			log.Error("Failed get user", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			user = u
		}
		log.Info("User", "user", user)
		organizationId, err := models.ParseUUID(user.OrganizationId)
		if err != nil {
			log.Error("Failed parse organization id", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var uri struct {
			Id string `uri:"id" binding:"required"`
		}

		if err := c.BindUri(&uri); err != nil {
			log.Error("Failed bind uri", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		id, err := models.ParseUUID(uri.Id)
		if err != nil {
			log.Error("Failed parse id", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		data, err := repository.GetQueries().GetNode(c.Request.Context(), id)
		if err != nil {
			log.Error("Failed get nodes from repository", "error", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if data.GroupId != organizationId {
			log.Error("Try to access to node outside of own organization", "node.OrganizationId", data.GroupId.String(), "organizationId", organizationId.String())
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		tags, err := repository.GetQueries().GetTags(c.Request.Context(), data.Id)
		if err != nil {
			log.Error("Failed get tags from repository", "error", err, "nodeId", data.Id.String())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		node := Node{
			Id:   data.Id.String(),
			Name: data.Name.String,
			Tags: tags,
		}
		c.JSON(http.StatusOK, gin.H{
			"node": node,
		})
	}
}
