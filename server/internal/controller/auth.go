package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/crypto/bcrypt"

	"ydsterm-server/internal/dao"
)

type AuthController struct{}

var Auth = AuthController{}

type VerifyRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (c *AuthController) Verify(r *ghttp.Request) {
	ctx := r.Context()
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "invalid request body"})
		return
	}
	if req.Username == "" || req.Password == "" {
		r.Response.WriteJsonExit(g.Map{"error": "username and password required"})
		return
	}

	record, err := dao.SyncUsers.Ctx(ctx).Where("username", req.Username).One()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": "database error"})
		return
	}

	if record.IsEmpty() {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{"error": "failed to hash password"})
			return
		}
		b := make([]byte, 16)
		rand.Read(b)
		id := hex.EncodeToString(b)
		_, err = dao.SyncUsers.Ctx(ctx).Data(g.Map{
			"id":            id,
			"username":      req.Username,
			"password_hash": string(hash),
		}).Insert()
		if err != nil {
			r.Response.WriteJsonExit(g.Map{"error": "failed to create user"})
			return
		}
		r.Response.WriteJson(VerifyResponse{ID: id, Username: req.Username})
		return
	}

	hash := record["password_hash"].String()
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		r.Response.WriteStatusExit(401, g.Map{"error": "invalid password"})
		return
	}

	r.Response.WriteJson(VerifyResponse{
		ID:       record["id"].String(),
		Username: record["username"].String(),
	})
}
