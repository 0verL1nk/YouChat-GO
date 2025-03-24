package utils

import "golang.org/x/crypto/bcrypt"

// 加密密码
func HashAndSalt(pwd string) (hashed string, err error) {
	res, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	hashed = string(res)
	return
}

// 验证密码
func ComparePasswords(hashedPwd string, plainPwd string) (err error) {
	byteHash := []byte(hashedPwd)
	err = bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		return err
	}
	return
}
