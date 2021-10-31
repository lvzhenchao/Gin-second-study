package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string  {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")

	result := make([]byte, n)
	rand.Seed(time.Now().Unix())//初始化随机数的种子
	for i := range result{
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
