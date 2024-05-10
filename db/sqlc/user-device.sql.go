// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user-device.sql

package db

import (
	"context"
)

const addUserDevice = `-- name: AddUserDevice :one
INSERT INTO user_device_lists (
    user_id,
    device_id,
    ip_address,
    client_details
) VALUES (
    $1, $2, $3, $4
)RETURNING id, user_id, device_id, ip_address, client_details, is_blocked, created_at
`

type AddUserDeviceParams struct {
	UserID        int64  `json:"user_id"`
	DeviceID      string `json:"device_id"`
	IpAddress     string `json:"ip_address"`
	ClientDetails string `json:"client_details"`
}

func (q *Queries) AddUserDevice(ctx context.Context, arg AddUserDeviceParams) (UserDeviceList, error) {
	row := q.db.QueryRowContext(ctx, addUserDevice,
		arg.UserID,
		arg.DeviceID,
		arg.IpAddress,
		arg.ClientDetails,
	)
	var i UserDeviceList
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.DeviceID,
		&i.IpAddress,
		&i.ClientDetails,
		&i.IsBlocked,
		&i.CreatedAt,
	)
	return i, err
}
