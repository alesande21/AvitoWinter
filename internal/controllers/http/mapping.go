package http

import "AvitoWinter/internal/entity"

func MapUserInfoToInfoResponse(userInfo *entity.UserInfo) InfoResponse {
	inventory := make([]struct {
		Quantity *int    `json:"quantity,omitempty"`
		Type     *string `json:"type,omitempty"`
	}, 0, len(userInfo.GetItems()))

	for _, item := range userInfo.GetItems() {
		quantity := item.GetQuantity()
		itemType := item.GetItem()
		inventory = append(inventory, struct {
			Quantity *int    `json:"quantity,omitempty"`
			Type     *string `json:"type,omitempty"`
		}{
			Quantity: &quantity,
			Type:     &itemType,
		})
	}

	received := make([]struct {
		Amount   *int    `json:"amount,omitempty"`
		FromUser *string `json:"fromUser,omitempty"`
	}, 0, len(userInfo.GetUserTransfers().GetReceived()))

	for _, transfer := range userInfo.GetUserTransfers().GetReceived() {
		amount := transfer.GetAmount()
		fromUser := transfer.GetUsername()
		received = append(received, struct {
			Amount   *int    `json:"amount,omitempty"`
			FromUser *string `json:"fromUser,omitempty"`
		}{
			Amount:   &amount,
			FromUser: &fromUser,
		})
	}

	sent := make([]struct {
		Amount *int    `json:"amount,omitempty"`
		ToUser *string `json:"toUser,omitempty"`
	}, 0, len(userInfo.GetUserTransfers().GetSent()))

	for _, transfer := range userInfo.GetUserTransfers().GetSent() {
		amount := transfer.GetAmount()
		toUser := transfer.GetUsername()
		sent = append(sent, struct {
			Amount *int    `json:"amount,omitempty"`
			ToUser *string `json:"toUser,omitempty"`
		}{
			Amount: &amount,
			ToUser: &toUser,
		})
	}

	return InfoResponse{
		Coins:     userInfo.GetCoins(),
		Inventory: &inventory,
		CoinHistory: &struct {
			Received *[]struct {
				Amount   *int    `json:"amount,omitempty"`
				FromUser *string `json:"fromUser,omitempty"`
			} `json:"received,omitempty"`
			Sent *[]struct {
				Amount *int    `json:"amount,omitempty"`
				ToUser *string `json:"toUser,omitempty"`
			} `json:"sent,omitempty"`
		}{
			Received: &received,
			Sent:     &sent,
		},
	}
}
