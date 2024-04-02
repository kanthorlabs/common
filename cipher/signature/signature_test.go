package signature

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var key = uuid.NewString()
var data = fmt.Sprintf(
	"%s.%s.%d",
	testdata.Fake.Lorem().Sentence(1),
	testdata.Fake.Lorem().Sentence(1),
	time.Now().UTC().UnixMilli(),
)

func TestSignature(t *testing.T) {
	t.Run("English", func(st *testing.T) {
		engkey := uuid.NewString()
		engdata := "The quick brown fox jumps over the lazy dog"

		sign := Sign(engkey, engdata)
		require.NoError(st, Verify(engkey, engdata, sign))
	})

	t.Run("Chinese", func(st *testing.T) {
		chinesekey := "相应中文可简译为“快狐跨懒狗”，完整翻译则是“敏捷的棕色狐狸跨过懒狗"
		chinesedata := "相应中文可简译为“快狐跨懒狗”，完整翻译则是“敏捷的棕色狐狸跨过懒狗"

		sign := Sign(chinesekey, chinesedata)
		require.NoError(st, Verify(chinesekey, chinesedata, sign))
	})
}
