package main

ет платежи
type PaymentProcessor struct{}

// NewPaymentProcessor создает новый процессор платежей
func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{}
}

// ProcessPayment обрабатывает платеж
func (pp *PaymentProcessor) ProcessPayment(method string, amount float64) error {
	if method != "крипта" {
		return ErrUnsupportedPaymentMethod
	}

	// Имитация проверки средств
	maxAmount := 1000.0
	if amount > maxAmount {
		return &InsufficientFundsError{
			RequestedAmount: amount,
			MaxAmount:       maxAmount,
		}
	}

	return nil
}

gfjfgjfgjfgjfgjfghfgnfgn

func HandlePayments() error {
	var method string
	var amountStr string

	// Спросить метод перевода в консоли
	fmt.Print("Введите метод перевода (карта/СБП): ")
	fmt.Scanln(&method)

	// Спросить сумму перевода в консоли
	fmt.Print("Введите сумму перевода: ")
	fmt.Scanln(&amountStr)

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	pp := NewPaymentProcessor()
	if err := pp.ProcessPayment(method, amount); err != nil {
		return fmt.Errorf("process payment: %w", err)
	}

	return nil
}

func main() {
	err := HandlePayments()
	var InsufficientFundsError *InsufficientFundsError

	if err != nil {
		switch {
		case errors.Is(err, ErrUnsupportedPaymentMethod):
			{
				fmt.Println("Ошибка: неподдерживаемый метод платежа.")
				os.Exit(1)
			}
		case errors.As(err, &InsufficientFundsError):
			{
				fmt.Printf("Ошибка: недостаточно средств. %s\n", InsufficientFundsError.Error())
				os.Exit(1)
			}
		default:
			{
				fmt.Printf("Неизвестная ошибка: %s.\n", err.Error())
				os.Exit(1)
			}
		}

	}
	fmt.Println("Платеж успешно обработан!")
}



// QUIC
// ​IT-KAMASUTRA > плейлист back-end путь самурая 03 урок
// package main

// import (
// 	"fmt"
// 	"math/rand"

// 	"strings"
// )

// func main() {
// 	village := Village{}

// 	// Создаем жителей деревни
// 	resident1 := &Resident{
// 		Name:    "Алиса",
// 		Age:     30,
// 		Married: false,
// 		Alive:   true,
// 		Events:  []string{}} //(список событий за год жизни)

// 	resident2 := &Resident{
// 		Name:    "Борис",
// 		Age:     40,
// 		Married: true,
// 		Alive:   true,
// 		Events:  []string{}}

// 	// Создаем животных
// 	animal1 := &Animal{
// 		Name:   "Бобик",
// 		Age:    5,
// 		Type:   "собака",
// 		Alive:  true,
// 		Events: []string{}}

// 	animal2 := &Animal{
// 		Name:   "Мурка",
// 		Age:    3,
// 		Type:   "кошка",
// 		Alive:  true,
// 		Events: []string{}}

// 	// Добавляем элементы в деревню
// 	village.AddElement(resident1)
// 	village.AddElement(resident2)
// 	village.AddElement(animal1)
// 	village.AddElement(animal2)

// 	// Симуляция обновления деревни на несколько лет
// 	for i := 0; i < 5; i++ {
// 		fmt.Printf("Год %d:\n", i+1)
// 		village.UpdateAll()
// 		fmt.Println(village.ShowAllInfo())
// 	}

// }

// type VillageElement interface {
// 	Update()
// 	FlushInfo() string
// }

// type Resident struct {
// 	Name    string
// 	Age     int
// 	Married bool
// 	Alive   bool
// 	Events  []string
// }

// // +1 год жизни для жителя
// func (r Resident) upYear() {
// 	r.Age++
// }

// // изменения статуса брака у жителя
// func (r Resident) Marriagestatus() {

// 	rand := rand.Intn(10) + 1
// 	if rand%2 == 0 {
// 		r.Alive = false
// 	}
// }

// // статуса умирания жителя
// func (r Resident) Alivestatus() {

// 	rand := rand.Intn(10) + 1

// 	if rand%2 == 0 {
// 		r.Alive = false
// 	}

// }

// func (r Resident) Update() {

// 	r.upYear()
// 	r.Alivestatus()

// 	rand := rand.Intn(10) + 1

// 	if rand%2 == 0 {
// 		r.Alivestatus()
// 	}

// }

// func (r Resident) FlushInfo() string {

// 	var result string

// 	for _, event := range r.Events {
// 		result += event + "\n"
// 	}

// 	result = strings.Join(r.Events, "\n")

// 	r.Events = make([]string, 0)

// 	return result

// }

// // животное
// type Animal struct {
// 	Name   string
// 	Age    int
// 	Type   string
// 	Alive  bool
// 	Events []string
// }

// // +1 год жизни для животного
// func (a Animal) upYear() {
// 	a.Age++
// }

// // метод умипания животного
// func (a Animal) Alivestatus() {

// 	if a.Age > 1 {
// 		a.Alive = false
// 		a.Events = append(a.Events, "умер")
// 	}
// }

// func (a Animal) Update() {
// 	rand := rand.Intn(10) + 1

// 	if rand%2 == 0 {
// 		a.Alivestatus()
// 	}
// }

// func (a Animal) FlushInfo() string {

// 	var result string

// 	for _, event := range a.Events {
// 		result += event + "\n"
// 	}

// 	result = strings.Join(a.Events, "\n")

// 	a.Events = make([]string, 0)

// 	return result

// }

// type Village struct {
// 	elements []VillageElement
// }

// func (v *Village) AddElement(el VillageElement) {
// 	v.elements = append(v.elements, el)
// }

// func (v Village) UpdateAll() {

// 	for _, element := range v.elements {
// 		element.Update()

// 	}
// }

// func (v *Village) ShowAllInfo() string {

// 	var info string

// 	for _, element := range v.elements {
// 		info += element.FlushInfo()

// 	}

// 	return info

// }

// https://t.me/GO_Blc
// https://t.me/GO_Blc
// https://t.me/GO_Blc
// goto
// faltruni

//Debian ubunta - на них крутиться сервер
// Obsidian
// func printInfo(user User) {

// 	fmt.Printf("Покупатель %s. Телефон: %s. Адрес: г. %s, %s.\n", user.Name, user.Phone, user.Address.City, user.Address.Street)

// Python + Django 2 - 3 тыс запросов, а Go - 100 тыс запросов
//Tailwind CSS и будстрап
