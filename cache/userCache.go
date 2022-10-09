package cache

import "go_ecommerce/models"

type UserCache interface {
	Set(key string, value *models.User)
	Get(key string) *models.User
}
