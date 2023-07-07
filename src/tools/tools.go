// пакет для служебных функций
package tools

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// MD5 - Превращает содержимое из переменной data в md5-хеш
func MD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// соединить две строки
func StrBuild(s1, s2 string) string {
	var out strings.Builder
	out.WriteString(s1)
	out.WriteString(s2)
	return out.String()
}

// генирирует пароль на алгоритме MD5
func PW5(login, passwd string) string {
	var out strings.Builder
	out.WriteString(login)
	out.WriteString(passwd)
	return fmt.Sprintf("%x", md5.Sum([]byte(out.String())))
}

// служебные функции
func PasswdEncrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return ""
	}
	return string(hash)
}

func PasswdVerify(passwd, hashPasswd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd))
	verified := err == nil
	if err == bcrypt.ErrMismatchedHashAndPassword {
		err = nil
	}
	return verified, err
}
