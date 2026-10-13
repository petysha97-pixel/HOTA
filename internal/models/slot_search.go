package models

// фильтры вкладки «Слоты» (приходят из адресной строки)
type SlotFilter struct {
	Tab       string // open — открытые роли, mine — мои отклики
	Query     string // текст: роль, название слота, техзадание, технологии
	Rolle     string // роль слота из таблицы roles
	StackID   []int  // слоту нужны ВСЕ эти технологии, максимум 6
	OnlyFree  bool   // только свободные: статус open (без review)
	FullMatch bool   // полное совпадение: у меня есть все технологии слота
}

// карточка слота во вкладке «Слоты»
type SlotCard struct {
	ID          int     `json:"id"`
	ProjectID   int     `json:"projectID"`
	ProjectName string  `json:"projectName"`
	Name        string  `json:"name"`
	Description string  `json:"description"` // техзадание
	Rolle       string  `json:"rolle"`
	Status      string  `json:"status"` // open или review (идёт знакомство)
	UserID      *int    `json:"userID"` // кандидат на review, null — никого
	Stack       []Stack `json:"stack"`

	// сравнение с моим стеком
	StackHits    int      `json:"stackHits"`    // сколько технологий слота у меня есть
	StackMissing []string `json:"stackMissing"` // каких технологий мне не хватает

	Applied bool `json:"applied"` // я уже откликнулся и жду ответа
}

// ответ вкладки «Слоты»
type SlotSearchResponse struct {
	OpenCount int        `json:"openCount"` // цифра на вкладке «Открытые роли»
	MineCount int        `json:"mineCount"` // цифра на вкладке «Мои отклики»
	Slots     []SlotCard `json:"slots"`
}
