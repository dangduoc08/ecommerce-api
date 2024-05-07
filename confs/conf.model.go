package confs

type ConfModel struct {
	Port int `bind:"PORT" validate:"required"`

	JWTAccessAPIKey       string `bind:"JWT_ACCESS_API_KEY" validate:"required"`
	JWTAccessAPIExpIn     int    `bind:"JWT_ACCESS_API_EXP_IN" validate:"required"`
	JWTRefreshTokenKey    string `bind:"JWT_REFRESH_TOKEN_KEY" validate:"required"`
	JWTRefreshTokenExpIn  int    `bind:"JWT_REFRESH_TOKEN_EXP_IN" validate:"required"`
	JWTResetPasswordKey   string `bind:"JWT_RESET_PASSWORD_KEY" validate:"required"`
	JWTResetPasswordExpIn int    `bind:"JWT_RESET_PASSWORD_EXP_IN" validate:"required"`

	PostgresHost     string `bind:"POSTGRES_HOST" validate:"required"`
	PostgresUser     string `bind:"POSTGRES_USER" validate:"required"`
	PostgresPassword string `bind:"POSTGRES_PASSWORD" validate:"required"`
	PostgresDB       string `bind:"POSTGRES_DB" validate:"required"`
	PostgresPort     int    `bind:"POSTGRES_PORT" validate:"required"`
	PostgresSSL      bool   `bind:"POSTGRES_SSL" validate:"boolean"`
}
