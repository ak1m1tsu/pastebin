package paste

// @description Тело запроса для создания пасты.
type CreatePasteBody struct {
	// Текст пасты
	Text string `example:"this is my paste"`
	// Формат текста
	Format string `example:"json" enums:"json,yaml,toml"`
	// Время, через которое паста становится не доступной
	Expires string `example:"30m"`
	// Пароль для получения доступа к пасте
	Password string `example:"hello"`
	// Название
	Name string `example:"thie is paste"`
} // @name Paste
