package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"stok-hadiah/models"
)

// parseStoreIDsFromString membaca store_id dalam format JSON array ataupun string dipisah koma.
func parseStoreIDsFromString(storeIDStr string) []int {
	storeIDStr = strings.TrimSpace(storeIDStr)
	if storeIDStr == "" {
		return []int{}
	}

	var ids []int
	if err := json.Unmarshal([]byte(storeIDStr), &ids); err == nil {
		return ids
	}

	parts := strings.Split(storeIDStr, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if id, err := strconv.Atoi(p); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// getAllowedStoreIDs mengambil daftar store id dari session user.
func getAllowedStoreIDs(c *gin.Context) []int {
	session := sessions.Default(c)
	userAny := session.Get("user")
	if userAny == nil {
		return []int{}
	}
	if su, ok := userAny.(models.SessionUser); ok {
		return parseStoreIDsFromString(su.StoreID)
	}
	return []int{}
}
