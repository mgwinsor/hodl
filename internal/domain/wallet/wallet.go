package wallet

import "github.com/google/uuid"

type WalletType int

const (
	WalletTypeExchange WalletType = iota
	WalletTypeHardware
	WalletTypeSoftware
	WalletTypeCustodial
)

type Wallet struct {
	ID   uuid.UUID
	Name string
	Type WalletType
}
