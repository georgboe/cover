package formatters

import "github.com/georgboe/cover/models"

type Formatter interface {
	Render(pairs []models.FileAndTest)
}
