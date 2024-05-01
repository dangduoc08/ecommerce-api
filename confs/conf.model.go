package confs

type ConfModel struct {
	Port   int    `bind:"PORT" validate:"required"`
	AppENV string `bind:"APP_ENV" validate:"required,AppENV"`

	Username  string `bind:"USERNAME" validate:"required"`
	Password  string `bind:"PASSWORD" validate:"required"`
	Email     string `bind:"EMAIL" validate:"required,email"`
	FirstName string `bind:"FIRST_NAME" validate:"required"`
	LastName  string `bind:"LAST_NAME" validate:"required"`

	JWTAccessAPIKey      string `bind:"JWT_ACCESS_API_KEY" validate:"required"`
	JWTAccessAPIExpIn    int    `bind:"JWT_ACCESS_API_EXP_IN" validate:"required"`
	JWTRefreshTokenKey   string `bind:"JWT_REFRESH_TOKEN_KEY" validate:"required"`
	JWTRefreshTokenExpIn int    `bind:"JWT_REFRESH_TOKEN_EXP_IN" validate:"required"`

	PostgresHost     string `bind:"POSTGRES_HOST" validate:"required"`
	PostgresUser     string `bind:"POSTGRES_USER" validate:"required"`
	PostgresPassword string `bind:"POSTGRES_PASSWORD" validate:"required"`
	PostgresDB       string `bind:"POSTGRES_DB" validate:"required"`
	PostgresPort     int    `bind:"POSTGRES_PORT" validate:"required"`
	PostgresSSL      bool   `bind:"POSTGRES_SSL" validate:"boolean"`
}
