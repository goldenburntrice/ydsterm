package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"ydsterm-server/internal/crypto"
)

const DecryptedDataKey = "sync_decrypted_data"

func SyncAuthMiddleware(serverKey string) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		timestamp := r.Header.Get("X-Timestamp")
		md5Sig := r.Header.Get("X-MD5")

		if timestamp == "" || md5Sig == "" {
			r.Response.WriteStatusExit(401, g.Map{"error": "missing X-Timestamp or X-MD5 header"})
			return
		}

		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			r.Response.WriteStatusExit(401, g.Map{"error": "invalid timestamp"})
			return
		}

		diff := math.Abs(float64(time.Now().Unix() - ts))
		if diff > 300 {
			r.Response.WriteStatusExit(401, g.Map{"error": "timestamp expired"})
			return
		}

		encBody, err := io.ReadAll(r.Body)
		if err != nil {
			r.Response.WriteStatusExit(400, g.Map{"error": "read body failed"})
			return
		}

		decrypted, err := crypto.DecryptWithKey(serverKey, string(encBody))
		if err != nil {
			r.Response.WriteStatusExit(400, g.Map{"error": "decryption failed"})
			return
		}

		var dataMap map[string]interface{}
		if err := json.Unmarshal([]byte(decrypted), &dataMap); err != nil {
			r.Response.WriteStatusExit(400, g.Map{"error": "invalid json"})
			return
		}

		computed := computeMD5Sign(dataMap, timestamp, serverKey)
		if computed != md5Sig {
			r.Response.WriteStatusExit(401, g.Map{"error": "signature mismatch"})
			return
		}

		r.SetParam(DecryptedDataKey, dataMap)
		r.Middleware.Next()
	}
}

func computeMD5Sign(dataMap map[string]interface{}, timestamp, serverKey string) string {
	keys := make([]string, 0, len(dataMap))
	for k := range dataMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys))
	for _, k := range keys {
		values = append(values, fmt.Sprintf("%v", dataMap[k]))
	}

	suffix := serverKey
	if len(suffix) > 8 {
		suffix = suffix[len(suffix)-8:]
	}

	raw := strings.Join(values, "|") + timestamp + suffix
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}
