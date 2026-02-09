package repositories

import (
	"fmt"
	"user_system/database"
	"user_system/models"
	"user_system/utils"
)

func NewDBHandler() error {
	if database.DB == nil {
		return fmt.Errorf("NewDBHandler: Database connection is not initialized")
	}
	//新建用户表
	_, err := database.DB.Exec(`
    CREATE TABLE IF NOT EXISTS users (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		return fmt.Errorf("Failed to create users table: %w", err)
	}
	return nil
}

func CreateUser(userInfo *models.UserInfo) *models.Response {
	if database.DB == nil {
		return &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//bcrypt加密密码
	hashedPassword, err := utils.HashPassword(userInfo.Password)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to hash password: %v", err), Type: 400}
	}
	//插入用户数据
	_, err = database.DB.Exec(`
		INSERT INTO users (username, password) VALUES (?, ?)`,
		userInfo.Username, hashedPassword,
	) //这里本来想查询一下是否存在同名用户，但mysql的唯一索引会自动帮我们处理这个问题，如果插入重复用户名会返回错误，我们直接捕获这个错误就行了
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to create user: %v", err), Type: 400}
	}
	/*
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return &models.Response{Message: fmt.Sprintf("Failed to get last insert ID: %v", err), Type: 400}
		}
		fmt.Printf("新用户ID: %d\n", lastInsertID)
	*/
	return &models.Response{Message: "User created successfully", Type: 200}
}

func UserLogin(userInfo *models.UserInfo) *models.Response {
	if database.DB == nil {
		return &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	var storedHashedPassword string
	err := database.DB.QueryRow(`
		SELECT password FROM users WHERE username = ?`,
		userInfo.Username,
	).Scan(&storedHashedPassword)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to query user: %v", err), Type: 404}
	}
	//检查密码
	if !utils.CheckPasswordHash(userInfo.Password, storedHashedPassword) {
		return &models.Response{Message: "Invalid username or password", Type: 400}
	}
	return &models.Response{Message: "Login successful", Type: 200}
}

func ChangePassword(userInfo *models.UserInfo) *models.Response {
	if database.DB == nil {
		return &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户存在
	var storedHashedPassword string
	database.DB.QueryRow(`
		SELECT username FROM users WHERE username = ?`,
		userInfo.Username,
	).Scan(&storedHashedPassword)
	if storedHashedPassword == "" {
		return &models.Response{Message: "User does not exist", Type: 400}
	}
	//bcrypt加密新密码
	hashedPassword, err := utils.HashPassword(userInfo.Password)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to hash password: %v", err), Type: 400}
	}
	//更新用户密码
	_, err = database.DB.Exec(`
		UPDATE users SET password = ? WHERE username = ?`,
		hashedPassword, userInfo.Username,
	)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to change password: %v", err), Type: 400}
	}
	return &models.Response{Message: "Password changed successfully", Type: 200}
}

func DeleteUser(userInfo *models.UserInfo) *models.Response {
	if database.DB == nil {
		return &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户存在
	var storedHashedPassword string
	database.DB.QueryRow(`
		SELECT username FROM users WHERE username = ?`,
		userInfo.Username,
	).Scan(&storedHashedPassword)
	if storedHashedPassword == "" {
		return &models.Response{Message: "User does not exist", Type: 400}
	}
	//删除用户数据
	_, err := database.DB.Exec(`
		DELETE FROM users WHERE username = ?`,
		userInfo.Username,
	)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to delete user: %v", err), Type: 400}
	}
	return &models.Response{Message: "User deleted successfully", Type: 200}
}

func GetAllUsers() ([]string, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询所有用户数据
	rows, err := database.DB.Query(`
		SELECT username FROM users`,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		usernames = append(usernames, username)
	}
	return usernames, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUserCount() (int, *models.Response) {
	if database.DB == nil {
		return 0, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数量
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM users`,
	).Scan(&count)
	if err != nil {
		return 0, &models.Response{Message: fmt.Sprintf("Failed to query user count: %v", err), Type: 400}
	}
	return count, &models.Response{Message: "User count retrieved successfully", Type: 200}
}

func GetUsernamesByPrefix(prefix string) ([]string, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户名以prefix开头的用户数据
	rows, err := database.DB.Query(`
		SELECT username FROM users WHERE username LIKE ?`,
		prefix+"%",
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		usernames = append(usernames, username)
	}
	return usernames, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsernamesBySuffix(suffix string) ([]string, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户名以suffix结尾的用户数据
	rows, err := database.DB.Query(`
		SELECT username FROM users WHERE username LIKE ?`,
		"%"+suffix,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		usernames = append(usernames, username)
	}
	return usernames, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsernamesBySubstring(substring string) ([]string, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户名包含substring的用户数据
	rows, err := database.DB.Query(`
		SELECT username FROM users WHERE username LIKE ?`,
		"%"+substring+"%",
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		usernames = append(usernames, username)
	}
	return usernames, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsernamesByLength(length int) ([]string, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户名长度为length的用户数据
	rows, err := database.DB.Query(`
		SELECT username FROM users WHERE LENGTH(username) = ?`,
		length,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		usernames = append(usernames, username)
	}
	return usernames, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsernamesByLengthRange(minLength, maxLength int) ([]string, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户名长度在minLength和maxLength之间的用户数据
	rows, err := database.DB.Query(`
		SELECT username FROM users WHERE LENGTH(username) BETWEEN ? AND ?`,
		minLength, maxLength,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		usernames = append(usernames, username)
	}
	return usernames, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

/*  因为用户信息表目前只有用户名和密码，所以这个接口就没什么意义了，暂时放着，以后如果要加更多用户信息再完善这个接口
func GetUserInfo(username string) (*models.UserInfo, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	var userInfo models.UserInfo
	err := database.DB.QueryRow(`
		SELECT username FROM users WHERE username = ?`,
		username,
	).Scan(&userInfo.Username)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query user: %v", err), Type: 400}
	}
	return &userInfo, &models.Response{Message: "User info retrieved successfully", Type: 200}
}
*/
