package db

import "github.com/longzekun/specter/server/db/models"

func SelectOperatorByToken(token string) *models.Operator {
	var operator models.Operator
	err := Session().Where(&models.Operator{
		Token: token,
	}).First(&operator).Error
	if err != nil {
		return nil
	}
	return &operator
}
