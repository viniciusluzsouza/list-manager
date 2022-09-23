package models

type ListDTO struct {
	ListParam tinyList `json:"list"`
}

type tinyList struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type ItemDTO struct {
	ID          uint64  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	UserID      *uint64 `json:"user_id"`
}

type ItemsDTO struct {
	Items []ItemDTO `json:"items"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthRequestSSO struct {
	Login    string `json:"login"`
	APPToken string `json:"app_token"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func NewListDTO(list *List) *ListDTO {
	tinyList := tinyList{ID: list.ID, Title: list.Title}
	return &ListDTO{tinyList}
}

func NewItemDTO(item *Item) *ItemDTO {
	return &ItemDTO{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		UserID:      item.UserID,
	}
}

func NewItemsDTO(items *[]Item) *ItemsDTO {
	itemsDTO := make([]ItemDTO, len(*items))
	for i, item := range *items {
		itemsDTO[i] = *NewItemDTO(&item)
	}
	return &ItemsDTO{Items: itemsDTO}
}
