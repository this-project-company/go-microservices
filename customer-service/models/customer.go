package models

type Customer struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string
    Email string `gorm:"uniqueIndex"`
}

func (Customer) TableName() string {
    return "customer"
}