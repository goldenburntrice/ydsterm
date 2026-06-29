// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// YdstermKeys is the golang structure for table ydsterm_keys.
type YdstermKeys struct {
	Id            string      `json:"id"            orm:"id"              description:""`
	Name          string      `json:"name"          orm:"name"            description:""`
	PrivateKeyEnc string      `json:"privateKeyEnc" orm:"private_key_enc" description:""`
	PublicKey     string      `json:"publicKey"     orm:"public_key"      description:""`
	PassphraseEnc string      `json:"passphraseEnc" orm:"passphrase_enc"  description:""`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:""`
	SyncUser      string      `json:"syncUser"      orm:"sync_user"       description:""`
}
