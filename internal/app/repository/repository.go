package repository

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Repository struct {
	orders        []Order
	cartData      []CartData
	cartToService []CartToService
}


type Order struct {
	ID              int
	Title           string
	Description     string
	FullDescription string
	IMG             string
	OS              string
	Version         string
	Size            float32
	Functions       string
}

type CartData struct {
	ID    int
	Phone string
}


type CartToService struct {
	RequestID   int
	ServiceID   int
	InstallTime string
}



func NewRepository() (*Repository, error) {
	orders := []Order{
		{
			ID:              1,
			Title:           "GIT",
			Description:     "GIT - система контроля версий",
			FullDescription: "Git — система контроля версий для отслеживания изменений в файлах и совместной работы.",
			IMG:             "http://127.0.0.1:9000/lb1/Git_icon.png",
			OS:              "Linux",
			Version:         "2.45",
			Size:            1.5,
			Functions: `1. Создание репозитория
2. Добавление изменений в индекс
3. Фиксация изменений
4. Просмотр статуса и истории изменений
5. Ветвление
6. Слияние веток`,
		},
		{
			ID:              2,
			Title:           "Open Shift",
			Description:     "Open Shift - управление контейнерами",
			FullDescription: "Платформа контейнеризации и оркестрации приложений на базе Kubernetes.",
			IMG:             "http://127.0.0.1:9000/lb1/open_shift.png",
			OS:              "Linux",
			Version:         "4.15",
			Size:            18.2,
			Functions: `1. Оркестрация контейнеров
2. Автоматическое масштабирование
3. Непрерывное развертывание
4. Мониторинг и логирование
5. Управление сетями
6. Балансировка нагрузки`,
		},
		{
			ID:              3,
			Title:           "PostgreSQL",
			Description:     "PostgreSQL - система управления БД",
			FullDescription: "Надёжная объектно-реляционная СУБД с открытым исходным кодом.",
			IMG:             "http://127.0.0.1:9000/lb1/Postgre_SQL.png",
			OS:              "Linux",
			Version:         "15.1",
			Size:            22.3,
			Functions: `1. Подключение к различным СУБД
2. Редактор SQL запросов
3. Визуальное построение запросов
4. Управление схемами БД
5. Экспорт и импорт данных
6. Генерация ER-диаграмм`,
		},
		{
			ID:              4,
			Title:           "DBeaver",
			Description:     "DBeaver - инструмент для работы с БД",
			FullDescription: "Многоплатформенный инструмент для управления и проектирования СУБД.",
			IMG:             "http://127.0.0.1:9000/lb1/dbeaver.png",
			OS:              "Linux",
			Version:         "15.1",
			Size:            22.3,
			Functions: `1. Подключение к различным СУБД
2. Редактор SQL запросов
3. Визуальное построение запросов
4. Управление схемами БД
5. Экспорт и импорт данных
6. Генерация ER-диаграмм`,
		},
	}

	cartData := []CartData{
		{ID: 1, Phone: "8-800-555-35-35"},
	}

	cartToService := []CartToService{
		{RequestID: 1, ServiceID: 1, InstallTime: "12:00"},
		{RequestID: 1, ServiceID: 2, InstallTime: "15:00"},
	}

	return &Repository{
		orders:        orders,
		cartData:      cartData,
		cartToService: cartToService,
	}, nil
}



func (r *Repository) GetOrders() ([]Order, error) {
	if len(r.orders) == 0 {
		return nil, fmt.Errorf("массив заказов пуст")
	}
	return r.orders, nil
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	var result []Order
	for _, o := range r.orders {
		if strings.Contains(strings.ToLower(o.Title), strings.ToLower(title)) {
			result = append(result, o)
		}
	}
	return result, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	for _, o := range r.orders {
		if o.ID == id {
			return o, nil
		}
	}
	return Order{}, fmt.Errorf("товар не найден")
}

func (r *Repository) GetCartOrders() ([]Order, error) {
	var result []Order
	if len(r.cartData) == 0 {
		return result, nil
	}
	for _, rel := range r.cartToService {
		order, err := r.GetOrder(rel.ServiceID)
		if err == nil {
			result = append(result, order)
		}
	}
	return result, nil
}

func (r *Repository) GetCartPhone() string {
	if len(r.cartData) == 0 {
		return ""
	}
	return r.cartData[0].Phone
}

func (r *Repository) GetCartCount() int {
	return len(r.cartToService)
}

func (r *Repository) GetCartInstallTime(orderID int) string {
	for _, rel := range r.cartToService {
		if rel.ServiceID == orderID {
			return rel.InstallTime
		}
	}
	return ""
}


func CalculateEndTime(start string, size float32) string {
	layout := "15:04"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		return "ошибка времени"
	}
	speed := rand.Intn(41) + 10
	minutes := int(size / float32(speed) * 60)
	endTime := startTime.Add(time.Duration(minutes) * time.Minute)
	return endTime.Format(layout)
}


func (r *Repository) GetCartID() int {
	if len(r.cartData) == 0 {
		return 0
	}
	return r.cartData[0].ID
}
