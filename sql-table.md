資料表格設計
====

## 使用者 Users
| Column | Type | 
| ------ | ---- |
| ID     | INT  |
| Email  | VARCHAR(64) |
| Name   | VARCHAR |
| Type | ENUM('Admin', 'Teacher')
| Password | BLOB |

## 設備 Devices
| Column | Type | 
| -- | -- |
| ID | VARCHAR(8) |
| Type | ENUM('Student-iPad', 'Teacher-iPad', 'Chromebook', 'WAP', 'WirelessProjector') |

## 預約紀錄 Bookings
| Column | Type | 
| - | - |
| ID | INT |
| User | Users.ID |
| LendingTime | DATETIME |
| ReturnTime | DATETIME |
| Teacher | INT |
| Student | INT |
| Chromebook | INT |
| WAP | INT |
| Projector | INT |

## 借出紀錄 Records
| Column | Type | 
| - | - |
| Booking | Bookings.ID |
| Device | Devices.ID |
| LentFrom | DATETIME |
| LentUntil | DATETIME |

