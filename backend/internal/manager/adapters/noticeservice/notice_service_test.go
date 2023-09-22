package noticeservice

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	svc, err := NewNoticeService()
	require.Nil(t, err)
	data, err := svc.Create(&domain.Revision{
		To:        "Kaspi",
		Address:   "Almaty, Beibitwilik 4",
		No:        758,
		CreatedAt: time.Now(),
		Remarks:   "\u200b         Таблица Заявление:\n1) Возмещение затрат по экспорту - Производственная мощность, возможности увеличения (Клетка-B13), Замечания: Уахахахахаха\n\u200b         Таблица Затраты на доставку транспортом:\n2) Подтверждающий документ - БИН производителя (Клетка-L10), Замечания: Не правильный БИН\n3) Договор с производителем - сторона 2 (Клетка-E4), Замечания: Не правильно\n",
	})
	require.Nil(t, err)
	file, err := os.Create("./temp2.docx")
	require.Nil(t, err)
	defer file.Close()
	_, err = io.Copy(file, data)
	require.Nil(t, err)

}
