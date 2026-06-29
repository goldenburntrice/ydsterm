// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SyncKeys is the golang structure for table sync_keys.
type SyncKeys struct {
	Id            string      `json:"id"            orm:"id"              ` //
	Name          string      `json:"name"          orm:"name"            ` //
	PrivateKeyEnc string      `json:"privateKeyEnc" orm:"private_key_enc" ` //
	PublicKey     string      `json:"publicKey"     orm:"public_key"      ` //
	PassphraseEnc string      `json:"passphraseEnc" orm:"passphrase_enc"  ` //
	SyncUser      string      `json:"syncUser"      orm:"sync_user"       ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      ` //
}
