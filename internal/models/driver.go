package models

type Driver struct {
	ID       int
	FullName string
	Age      uint8
	Exp      uint8 // опыт вождения
	Rating   uint8 // делить на 10, пример 45 == 4.5 звезд
	Cord     Cord
	Car      Car
}

type Cord struct {
	Longitude string // ширина
	Latitude  string // долгота
}

type Car struct {
	ID            int
	Model         string
	Name          string
	Year          uint16
	Number        string // номер машины
	Color         string // цвет машины
	Class         string // бизнес, комфорт, эконом, грузовой
	MaxPassangers uint8  // максимум пассажиров
	BabyChair     bool   // наличие детского сидения
	WithAnimals   bool   // разрешенно с животными
}
