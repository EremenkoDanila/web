package auth

type FixedUser struct {
    UserID uint
}

var instance *FixedUser

func GetFixedUser() *FixedUser {
    if instance == nil {
        instance = &FixedUser{
            UserID: 2, // ID пользователя из вашей БД
        }
    }
    return instance
}

func (u *FixedUser) GetUserID() uint {
    return u.UserID
}