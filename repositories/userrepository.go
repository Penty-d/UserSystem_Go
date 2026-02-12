package repositories

import (
	"fmt"
	"time"
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
        password VARCHAR(64) NOT NULL,
		fullname VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL,
		role VARCHAR(5) NOT NULL DEFAULT 'user',
		status VARCHAR(20) NOT NULL DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		return fmt.Errorf("Failed to create users table: %w", err)
	}
	return nil
}

func CreateUser(userInfo *models.CreateUserRequest) *models.Response {
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
		INSERT INTO users (username, password, fullname, email, role) VALUES (?, ?, ?, ?, ?)`,
		userInfo.Username, hashedPassword, userInfo.FullName, userInfo.Email, userInfo.Role,
	) //这里本来想查询一下是否存在同名用户，但mysql的唯一索引会自动帮我们处理这个问题，如果插入重复用户名会返回错误，我们直接捕获这个错误就行了
	if err != nil {
		return &models.Response{Message: "Failed to create user", Type: 400}
	}
	return &models.Response{Message: "User created successfully", Type: 200}
}

func UserLogin(userInfo *models.LoginRequest) (*models.Response, string) {
	if database.DB == nil {
		return &models.Response{Message: "Database connection is not initialized", Type: 400}, ""
	}
	//查询用户数据
	var storedHashedPassword, status, role string
	err := database.DB.QueryRow(`
		SELECT password, status, role FROM users WHERE username = ?`,
		userInfo.Username,
	).Scan(&storedHashedPassword, &status, &role)
	if status == "deleted" {
		return &models.Response{Message: "User account is deleted", Type: 400}, ""
	}
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to query user: %v", err), Type: 404}, ""
	}
	//检查密码
	if !utils.CheckPasswordHash(userInfo.Password, storedHashedPassword) {
		return &models.Response{Message: "Invalid username or password", Type: 400}, ""
	}
	Request := utils.CreateTokenRequset{
		ExpiredAt: time.Now().Add(time.Minute * 15),
		Role:      role,
		Username:  userInfo.Username,
	}
	token, err := utils.GetToken(&Request)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to create token: %v", err), Type: 400}, ""
	}
	return &models.Response{Message: "Login successful", Type: 200}, token
}

func UpdateUser(userInfo *models.UpdateUserRequest) *models.Response {
	if database.DB == nil {
		return &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	query := "UPDATE users SET "
	args := []interface{}{}
	if userInfo.Password != nil {
		//bcrypt加密新密码
		hashedPassword, err := utils.HashPassword(*userInfo.Password)
		if err != nil {
			return &models.Response{Message: fmt.Sprintf("Failed to hash password: %v", err), Type: 400}
		}
		query += "password = ?, "
		args = append(args, hashedPassword)
	}
	if userInfo.Role != nil {
		query += "role = ?, "
		args = append(args, *userInfo.Role)
	}
	if userInfo.Email != nil {
		query += "email = ?, "
		args = append(args, *userInfo.Email)
	}
	if userInfo.FullName != nil {
		query += "fullname = ?, "
		args = append(args, *userInfo.FullName)
	}
	if userInfo.Status != nil {
		query += "status = ?, "
		args = append(args, *userInfo.Status)
	}
	if len(args) == 0 {
		return &models.Response{Message: "No fields to update", Type: 400}
	}
	//去掉最后一个逗号和空格
	query = query[:len(query)-2]
	query += " WHERE username = ?"
	args = append(args, userInfo.Username)
	//更新用户数据
	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to update user: %v", err), Type: 400}
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to get affected rows: %v", err), Type: 400}
	}
	if rows == 0 {
		return &models.Response{Message: "User does not exist", Type: 400}
	}
	return &models.Response{Message: "User updated successfully", Type: 200}
}

func RemoveUser(ID int) *models.Response { //硬删除用户数据，慎用
	if database.DB == nil {
		return &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//删除用户数据
	result, err := database.DB.Exec(`
		DELETE FROM users WHERE id = ?`,
		ID,
	)
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to delete user: %v", err), Type: 400}
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return &models.Response{Message: fmt.Sprintf("Failed to get affected rows: %v", err), Type: 400}
	}
	if rows == 0 {
		return &models.Response{Message: "User does not exist", Type: 400}
	}
	return &models.Response{Message: "User deleted successfully", Type: 200}
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

func GetUserInfoByID(ID uint) (*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	var userInfo models.User
	err := database.DB.QueryRow(`
        SELECT 
            id, username, password, fullname, email, 
            role, status, created_at, updated_at
        FROM users
        WHERE id = ?`, ID,
	).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query user: %v", err), Type: 400}
	}
	return &userInfo, &models.Response{Message: "User info retrieved successfully", Type: 200}
}

func GetAllUsers() ([]*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}

	// 查询所有用户数据
	rows, err := database.DB.Query(`
		SELECT 
			id, username, password, fullname, email, 
			role, status, created_at, updated_at
		FROM users`)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
		if err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		users = append(users, &userInfo)
	}

	// 检查是否在遍历过程中发生错误
	if err := rows.Err(); err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Error during iteration: %v", err), Type: 400}
	}

	return users, &models.Response{Message: "All users retrieved successfully", Type: 200}
}

func GetUsersByStatus(status string) ([]*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	rows, err := database.DB.Query(`
		SELECT 
			id, username, password, fullname, email,
			role, status, created_at, updated_at
		FROM users
		WHERE status = ?`, status,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
		if err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		users = append(users, &userInfo)
	}
	return users, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsersByRole(role string) ([]*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	rows, err := database.DB.Query(`
		SELECT
			id, username, password, fullname, email,
			role, status, created_at, updated_at
		FROM users
		WHERE role = ?`, role,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
		if err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		users = append(users, &userInfo)
	}
	return users, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUserByUsername(username string) (*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	var userInfo models.User
	err := database.DB.QueryRow(`
		SELECT
			id, username, password, fullname, email,
			role, status, created_at, updated_at
		FROM users
		WHERE username = ?`, username,
	).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query user: %v", err), Type: 400}
	}
	return &userInfo, &models.Response{Message: "User retrieved successfully", Type: 200}
}

func GetUsersByFullname(fullname string) ([]*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	rows, err := database.DB.Query(`
		SELECT
			id, username, password, fullname, email,
			role, status, created_at, updated_at
		FROM users
		WHERE fullname = ?`, fullname,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
		if err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		users = append(users, &userInfo)
	}
	return users, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsersByEmail(email string) ([]*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	rows, err := database.DB.Query(`
		SELECT
			id, username, password, fullname, email,
			role, status, created_at, updated_at
		FROM users
		WHERE email = ?`, email,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
		if err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		users = append(users, &userInfo)
	}
	return users, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsersByCreatedAt(createdAt time.Time) ([]*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	rows, err := database.DB.Query(`
		SELECT
			id, username, password, fullname, email,
			role, status, created_at, updated_at
		FROM users
		WHERE created_at = ?`, createdAt,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
		if err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		users = append(users, &userInfo)
	}
	return users, &models.Response{Message: "Users retrieved successfully", Type: 200}
}

func GetUsersByUpdateAt(updateAt time.Time) ([]*models.User, *models.Response) {
	if database.DB == nil {
		return nil, &models.Response{Message: "Database connection is not initialized", Type: 400}
	}
	//查询用户数据
	rows, err := database.DB.Query(`
		SELECT
			id, username, password, fullname, email,
			role, status, created_at, updated_at
		FROM users
		WHERE updated_at = ?`, updateAt,
	)
	if err != nil {
		return nil, &models.Response{Message: fmt.Sprintf("Failed to query users: %v", err), Type: 400}
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		var userInfo models.User
		err := rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.FullName, &userInfo.Email, &userInfo.Role, &userInfo.Status, &userInfo.CreatedAt, &userInfo.UpdatedAt)
		if err != nil {
			return nil, &models.Response{Message: fmt.Sprintf("Failed to scan user: %v", err), Type: 400}
		}
		users = append(users, &userInfo)
	}
	return users, &models.Response{Message: "Users retrieved successfully", Type: 200}
}
