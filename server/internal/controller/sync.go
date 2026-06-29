package controller

import (
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"ydsterm-server/internal/crypto"
	"ydsterm-server/internal/middleware"
)

type SyncController struct{}

var Sync = SyncController{}

var moduleTableMap = map[string]string{
	"hosts":           "sync_hosts",
	"host_groups":     "sync_host_groups",
	"keys":            "sync_keys",
	"snippets":        "sync_snippets",
	"snippet_folders": "sync_snippet_folders",
	"port_forwards":   "sync_port_forwards",
}

func (c *SyncController) Push(r *ghttp.Request) {
	ctx := r.Context()
	module := r.GetRouter("module").String()
	username := r.GetQuery("user").String()

	if module == "" || username == "" {
		r.Response.WriteJsonExit(g.Map{"error": "module and user required"})
		return
	}

	table, ok := moduleTableMap[module]
	if !ok {
		r.Response.WriteJsonExit(g.Map{"error": "unknown module"})
		return
	}

	dataMap, ok := r.GetParam(middleware.DecryptedDataKey).Val().(map[string]interface{})
	if !ok {
		r.Response.WriteJsonExit(g.Map{"error": "no decrypted data"})
		return
	}

	dataMap["sync_user"] = username
	dataMap["updated_at"] = time.Now().UTC().Format("2006-01-02 15:04:05")
	if _, ok := dataMap["created_at"]; !ok {
		dataMap["created_at"] = dataMap["updated_at"]
	}

	if _, err := g.DB().Model(table).Ctx(ctx).Data(dataMap).Save(); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "upsert failed"})
		return
	}

	r.Response.WriteJson(g.Map{"status": "ok"})
}

func (c *SyncController) Delete(r *ghttp.Request) {
	ctx := r.Context()
	module := r.GetRouter("module").String()
	username := r.GetQuery("user").String()

	if module == "" || username == "" {
		r.Response.WriteJsonExit(g.Map{"error": "module and user required"})
		return
	}

	table, ok := moduleTableMap[module]
	if !ok {
		r.Response.WriteJsonExit(g.Map{"error": "unknown module"})
		return
	}

	dataMap, ok := r.GetParam(middleware.DecryptedDataKey).Val().(map[string]interface{})
	if !ok {
		r.Response.WriteJsonExit(g.Map{"error": "no decrypted data"})
		return
	}

	id, _ := dataMap["id"].(string)
	if id == "" {
		r.Response.WriteJsonExit(g.Map{"error": "id required"})
		return
	}

	if _, err := g.DB().Model(table).Ctx(ctx).Where("id", id).Where("sync_user", username).Delete(); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "delete failed"})
		return
	}

	r.Response.WriteJson(g.Map{"status": "ok"})
}

func (c *SyncController) Pull(r *ghttp.Request) {
	ctx := r.Context()
	module := r.GetRouter("module").String()
	username := r.GetQuery("user").String()
	serverKey := r.GetQuery("server_key").String()

	if module == "" || username == "" {
		r.Response.WriteJsonExit(g.Map{"error": "module and user required"})
		return
	}

	table, ok := moduleTableMap[module]
	if !ok {
		r.Response.WriteJsonExit(g.Map{"error": "unknown module"})
		return
	}

	records, err := g.DB().Model(table).Ctx(ctx).Where("sync_user", username).All()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "query failed"})
		return
	}

	result := make([]g.Map, 0)
	for _, record := range records {
		m := record.Map()
		delete(m, "sync_user")
		result = append(result, m)
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "encode failed"})
		return
	}

	encrypted, err := crypto.EncryptWithKey(serverKey, string(jsonBytes))
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "encryption failed"})
		return
	}

	r.Response.WriteExit(encrypted)
}
