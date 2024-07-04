package models

const ( //ClassEnum
	Bussines = "Bussines"
	Comfort  = "Comfort"
	Econom   = "Econom"
	Cargo    = "Cargo"
)

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
	DriverID  int     `json:"driver_id"` // вопрос, как-то можно не дублировать это поле? Просто если мы определим только Cord, как понять чьи это координаты?
	Longitude float32 `json:"longitude"` // ширина
	Latitude  float32 `json:"latitude"`  // долгота
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
