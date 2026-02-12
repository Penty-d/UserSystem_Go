package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
	"user_system/database"
)

type TokenInfo struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	Username  string    `json:"username" binding:"required,max=50"`
	Role      string    `json:"role" binding:"required,oneof=admin user"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CreateTokenRequset struct {
	Username  string    `json:"username" binding:"required,max=50"`
	Role      string    `json:"role" binding:"required,oneof=admin user"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewAuthDBHandler() error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	//新建Tokens表
	_, err := database.DB.Exec(`
    CREATE TABLE IF NOT EXISTS tokens (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		token VARCHAR(64) NOT NULL UNIQUE,
        username VARCHAR(50) NOT NULL UNIQUE,
		role VARCHAR(5) NOT NULL DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		expired_at TIMESTAMP
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`) //FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	if err != nil {
		return fmt.Errorf("Failed to create tokens table: %w", err)
	}
	return nil
}

func GetToken(Info *CreateTokenRequset) (string, error) {
	if time.Now().After(Info.ExpiredAt) {
		return "", fmt.Errorf("CreateToken: ExpiredAt must be after now")
	}
	if database.DB == nil {
		return "", fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	var token string
	err := database.DB.QueryRow(`
		SELECT token, expired_at FROM tokens WHERE username = ?`,
		Info.Username,
	).Scan(&token, &Info.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			//token不存在，生成token
			token, err = GernerateToken()
			if err != nil {
				return "", err
			}
			_, err = database.DB.Exec(`
			INSERT INTO tokens (token, username, role, expired_at) VALUES(?, ?, ?, ?)`,
				token, Info.Username, Info.Role, Info.ExpiredAt,
			)
			if err != nil {
				return "", err
			}
			return token, nil
		}
		return "", err
	}
	//token存在，判断是否过期
	if time.Now().Before(Info.ExpiredAt) {
		return token, nil
	}
	//token过期，更新token
	token, err = GernerateToken()
	if err != nil {
		return "", err
	}
	Info.ExpiredAt = time.Now().Add(15 * time.Minute)
	err = UpdateToken(token, Info)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetInfobyToken(Token string) (*TokenInfo, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	var tokeninfo TokenInfo
	err := database.DB.QueryRow(`
	SELECT 
	id, username, role, created_at, expired_at
	FROM tokens
	WHERE token = ?`, Token,
	).Scan(&tokeninfo.ID, &tokeninfo.Username, &tokeninfo.Role, &tokeninfo.CreatedAt, &tokeninfo.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Token not found")
		}
		return nil, err
	}
	tokeninfo.Token = Token
	return &tokeninfo, nil
}

func GernerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	return token, nil
}

func DeleteToken(Token string) error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	_, err := database.DB.Exec(`
	DELETE FROM tokens WHERE token = ?`, Token,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetTokenCount() (int, error) {
	if database.DB == nil {
		return 0, fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	var count int
	err := database.DB.QueryRow(`
	SELECT COUNT(*) FROM tokens`,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetAllTokens() ([]TokenInfo, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	rows, err := database.DB.Query(`
	SELECT 
	id, token, username, role, created_at, expired_at
	FROM tokens
	ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tokens []TokenInfo
	for rows.Next() {
		var tokeninfo TokenInfo
		err := rows.Scan(&tokeninfo.ID, &tokeninfo.Token, &tokeninfo.Username, &tokeninfo.Role, &tokeninfo.CreatedAt, &tokeninfo.ExpiredAt)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tokeninfo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func UpdateToken(Token string, Info *CreateTokenRequset) error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	_, err := database.DB.Exec(`
	UPDATE tokens SET role = ?, expired_at = ?, token = ? WHERE username = ?`,
		Info.Role, Info.ExpiredAt, Token, Info.Username,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllTokens() error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	_, err := database.DB.Exec(`
	DELETE FROM tokens`,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteExpiredTokens() error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	_, err := database.DB.Exec(`
	DELETE FROM tokens WHERE expired_at < ?`, time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTokenByUsername(username string) error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	_, err := database.DB.Exec(`
	DELETE FROM tokens WHERE username = ?`, username,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTokenByID(id int) error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	_, err := database.DB.Exec(`
	DELETE FROM tokens WHERE id = ?`, id,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTokenByToken(token string) error {
	if database.DB == nil {
		return fmt.Errorf("NewAuthDBHandler: Database connection is not initialized")
	}
	_, err := database.DB.Exec(`
	DELETE FROM tokens WHERE token = ?`, token,
	)
	if err != nil {
		return err
	}
	return nil
}
