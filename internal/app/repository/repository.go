package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
	orders        []Order
	versions      []Version
	orderVersion  []OrderVersion
	cartData      []CartData
	orderRequests []OrderRequest
}

// Структуры корзины
type CartData struct {
	Phone  string
	Orders []CartOrder
}

type CartOrder struct {
	ID      int
	OS      string
	Version string
}

type OrderVersion struct {
	OrderID   int
	VersionID int
}

// Новая структура для заказов
type OrderRequest struct {
	ID          int
	PhoneNumber string
	Services    []ServiceDetail
}

// Структура для деталей услуги в заказе
type ServiceDetail struct {
	ServiceID int
	OS        string
	Version   string
}

func NewRepository() (*Repository, error) {
	orders := []Order{
		{
			ID:              1,
			Title:           "GIT",
			Description:     "GIT - система контроля версии",
			FullDescription: "Git — Система контроля версии, которая позволяет отслеживать изменения файлов и хранить их версии и оперативно возвращать в любое сохраненное состояние",
			IMG:             "http://127.0.0.1:9000/lb1/Git_icon.png",
			OS:              "Linux",
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
			FullDescription: "Open Shift - платформа контейнеризации, построенная на основе Kubernetes, которая предоставляет комплексное решение для разработки, развертывания и управления контейнерными приложениями",
			IMG:             "http://127.0.0.1:9000/lb1/open_shift.png",
			OS:              "Linux",
			Functions: `1. Оркестрация контейнеров
2. Автоматическое масштабирование
3. Непрерывное развертывание
4. Мониторинг и логирование
5. Управление сетями
6. Балансировка нагрузки`,
		},
		{
			ID:              3,
			Title:           "DBeaver",
			Description:     "DBeaver - инструмент для работы с БД",
			FullDescription: "DBeaver - многоплатформенный инструмент с открытым исходным кодом для работы с базами данных, используемый разработчиками, администраторами баз данных и аналитиками для управления и проектирования СУБД",
			IMG:             "http://127.0.0.1:9000/lb1/dbeaver.png",
			OS:              "Linux",
			Functions: `1. Подключение к различным СУБД
2. Редактор SQL запросов
3. Визуальное построение запросов
4. Управление схемами БД
5. Экспорт и импорт данных
6. Генерация ER-диаграмм`,
		},
		{
			ID:              4,
			Title:           "PostgreSQL",
			Description:     "PostgreSQL - система управления БД",
			FullDescription: "PostgreSQL - мощная и надёжная система управления объектно-реляционными базами данных (СУБД) с открытым исходным кодом",
			IMG:             "http://127.0.0.1:9000/lb1/Postgre_SQL.png",
			OS:              "Linux",
			Functions: `1. Транзакции ACID
2. Репликация данных
3. Расширяемость через расширения
4. Полнотекстовый поиск
5. JSON поддержка
6. Геопространственные данные`,
		},
	}

	versions := []Version{
		{ID: 1, Number: "1.2", Size: 1},
		{ID: 2, Number: "1.3", Size: 1.5},
		{ID: 3, Number: "13.0", Size: 20},
		{ID: 4, Number: "14.0", Size: 20.2},
		{ID: 5, Number: "1.0", Size: 10},
		{ID: 6, Number: "1.1", Size: 10},
		{ID: 7, Number: "11.0", Size: 12},
		{ID: 8, Number: "12.0", Size: 13},
	}

orderVersion := []OrderVersion{
		{OrderID: 1, VersionID: 1},
		{OrderID: 1, VersionID: 2},
		{OrderID: 2, VersionID: 3},
		{OrderID: 2, VersionID: 4},
		{OrderID: 3, VersionID: 5},
		{OrderID: 3, VersionID: 6},
		{OrderID: 4, VersionID: 7},
		{OrderID: 4, VersionID: 8},
	}

	cartData := []CartData{
		{
			Phone: "8-800-555-35-35",
			Orders: []CartOrder{
				{ID: 1, OS: "Linux", Version: "1.3"},
				{ID: 2, OS: "Lunyx", Version: "14.0"},
			},
		},
	}

	orderRequests := []OrderRequest{}

	return &Repository{
		orders:        orders,
		versions:      versions,
		orderVersion:  orderVersion,
		cartData:      cartData,
		orderRequests: orderRequests,
	}, nil
}

type Version struct {
	ID     int
	Number string
	Size   float32
}

type Order struct {
	ID              int
	Title           string
	Description     string
	FullDescription string
	IMG             string
	OS              string
	Functions       string
}

func (r *Repository) getVersionsByOrderID(orderID int) []Version {
	var result []Version
	for _, ov := range r.orderVersion {
		if ov.OrderID == orderID {
			for _, version := range r.versions {
				if version.ID == ov.VersionID {
					result = append(result, version)
					break
				}
			}
		}
	}
	return result
}

func (o *Order) GetLatestVersion(versions []Version) Version {
	if len(versions) == 0 {
		return Version{}
	}
	latest := versions[0]
	for _, version := range versions {
		if version.Number > latest.Number {
			latest = version
		}
	}
	return latest
}

func (r *Repository) GetOrders() ([]Order, error) {
	if len(r.orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return r.orders, nil
}

func (r *Repository) GetCartOrders() ([]Order, error) {
	var result []Order
	if len(r.cartData) == 0 {
		return result, nil
	}

	for _, orderInfo := range r.cartData[0].Orders {
		order, err := r.GetOrder(orderInfo.ID)
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

func (r *Repository) GetCartDataVersion(orderID int) string {
	if len(r.cartData) == 0 {
		return ""
	}
	for _, order := range r.cartData[0].Orders {
		if order.ID == orderID {
			return order.Version
		}
	}
	return ""
}

func (r *Repository) GetOrder(id int) (Order, error) {
	for _, order := range r.orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("заказ не найден")
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	var result []Order
	for _, order := range r.orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}
	return result, nil
}

func (r *Repository) FindVersionByNumber(orderID int, versionQuery string) (Version, error) {
	versions := r.getVersionsByOrderID(orderID)
	if len(versions) == 0 {
		return Version{}, fmt.Errorf("версии не найдены")
	}
	if versionQuery == "" {
		return versions[len(versions)-1], nil
	}
	for _, version := range versions {
		if strings.Contains(strings.ToLower(version.Number), strings.ToLower(versionQuery)) {
			return version, nil
		}
	}
	return Version{}, fmt.Errorf("версия не найдена")
}

func (r *Repository) GetVersionsByOrderID(orderID int) ([]Version, error) {
	return r.getVersionsByOrderID(orderID), nil
}

func (r *Repository) GetCartCount() int {
	if len(r.cartData) == 0 {
		return 0
	}
	return len(r.cartData[0].Orders)
}
