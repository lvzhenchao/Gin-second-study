package main

import "github.com/gin-gonic/gin"

//RSA实现Token
type RsaUser struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Telephone string `json:"telephone"`
	Password string `json:"password"`
}

func init()  {
	
}

func main()  {
	r := gin.Default()

	

	r.Run(":9090")
}
