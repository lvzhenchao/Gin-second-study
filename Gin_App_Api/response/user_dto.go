package response

import "Gin-second-study/Gin_App_Api/model"

type UserDto struct {
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

//DTO就是数据传输对象(Data Transfer Object)的缩写; 用于 展示层与服务层之间的数据之间的数据传输对象

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name: user.Name,
		Telephone: user.Telephone,
	}
}
