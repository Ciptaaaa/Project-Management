package utils

import (
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/google/uuid"
)

func SortListByPosition(lists []models.List, order []uuid.UUID) []models.List{
	ordered := make ([]models.List,0,len(order))

	listMap := make(map[uuid.UUID]models.List)
	for _,l := range lists{
		listMap[l.PublicID]= l
	}

	//grouping by order

	for _,id := range order{
		if list,ok := listMap[id];ok {
			ordered = append(ordered, list)
		}
	}

	return ordered
}