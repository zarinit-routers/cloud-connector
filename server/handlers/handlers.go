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

		if data.GroupId != organizationId && !user.IsAdmin() {
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

// TODO: refactor this function, increase its size
func AddTagsHandler() gin.HandlerFunc {
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

		var request struct {
			Id   string   `uri:"id" binding:"required"`
			Tags []string `json:"tags" binding:"required"`
		}

		if err := c.BindJSON(&request); err != nil {
			log.Error("Failed bind json", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		id, err := models.ParseUUID(request.Id)
		if err != nil {
			log.Error("Failed parse id", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		node, err := repository.GetQueries().GetNode(c.Request.Context(), id)
		if err != nil {
			log.Error("Failed get node from repository", "error", err, "nodeId", id.String())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if node.GroupId != organizationId && !user.IsAdmin() {
			log.Error("Try to access to node outside of own organization", "node.OrganizationId", node.GroupId.String(), "organizationId", organizationId.String())
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		for _, tag := range request.Tags {
			if tag == "" {
				log.Warn("Tag is empty")
				continue
			}
			err = repository.GetQueries().AddTag(c.Request.Context(), repository.AddTagParams{
				NodeId: id,
				Tag:    tag,
			})
			if err != nil {
				log.Error("Failed add tags to repository", "error", err, "nodeId", id.String())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		c.Status(http.StatusOK)
	}
}

// TODO: refactor this function, increase its size
func RemoveTagsHandler() gin.HandlerFunc {
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

		var request struct {
			Id   string   `uri:"id" binding:"required"`
			Tags []string `json:"tags" binding:"required"`
		}

		if err := c.BindJSON(&request); err != nil {
			log.Error("Failed bind json", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		id, err := models.ParseUUID(request.Id)
		if err != nil {
			log.Error("Failed parse id", "error", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		node, err := repository.GetQueries().GetNode(c.Request.Context(), id)
		if err != nil {
			log.Error("Failed get node from repository", "error", err, "nodeId", id.String())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if node.GroupId != organizationId && !user.IsAdmin() {
			log.Error("Try to access to node outside of own organization", "node.OrganizationId", node.GroupId.String(), "organizationId", organizationId.String())
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		for _, tag := range request.Tags {
			if tag == "" {
				log.Warn("Tag is empty")
				continue
			}
			err = repository.GetQueries().RemoveTag(c.Request.Context(), repository.RemoveTagParams{
				NodeId: id,
				Tag:    tag,
			})
			if err != nil {
				log.Error("Failed add tags to repository", "error", err, "nodeId", id.String())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		c.Status(http.StatusOK)
	}
}
