package baseencode

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	t.Run("EncodeSmallNum", func(t *testing.T) {
		require.Equal(t, "", string(Base58Encode(0)))
		require.Equal(t, "2", string(Base58Encode(1)))
		require.Equal(t, "z", string(Base58Encode(57)))
		require.Equal(t, "21", string(Base58Encode(58)))
	})

	t.Run("EncodePositiveNumber", func(t *testing.T) {
		require.Equal(t, "BukQL", string(Base58Encode(123456789)))
		require.Equal(t, "24rgcX", string(Base58Encode(700000000)))
		require.Equal(t, "24rgcY", string(Base58Encode(700000001)))
	})

	t.Run("EncodeLargeNumber", func(t *testing.T) {
		require.Equal(t, "NQm6nKp8qFC", string(Base58Encode(9223372036854775807)))
	})
}

func TestBase58Decode(t *testing.T) {
	t.Run("DecodeSmallNum", func(t *testing.T) {
		num, err := Base58Decode([]byte(""))
		require.NoError(t, err)
		require.Equal(t, uint64(0), num)

		num, err = Base58Decode([]byte("2"))
		require.NoError(t, err)
		require.Equal(t, uint64(1), num)

		num, err = Base58Decode([]byte("z"))
		require.NoError(t, err)
		require.Equal(t, uint64(57), num)

		num, err = Base58Decode([]byte("21"))
		require.NoError(t, err)
		require.Equal(t, uint64(58), num)
	})

	t.Run("DecodePositiveNumber", func(t *testing.T) {
		num, err := Base58Decode([]byte("BukQL"))
		require.NoError(t, err)
		require.Equal(t, uint64(123456789), num)
	})

	t.Run("DecodeLargeNumber", func(t *testing.T) {
		num, err := Base58Decode([]byte("zzzzzzzz"))
		require.NoError(t, err)
		require.Equal(t, uint64(128063081718015), num)
	})

	t.Run("DecodeInvalidNumber", func(t *testing.T) {
		num, err := Base58Decode([]byte("zzzzzzzzz"))
		require.ErrorIs(t, err, ErrBase58Overflow)
		require.Equal(t, uint64(0), num)

		num, err = Base58Decode([]byte("000"))
		require.ErrorIs(t, err, ErrorInvalidCharacter)
		require.Equal(t, uint64(0), num)

		num, err = Base58Decode([]byte("12345l"))
		require.ErrorIs(t, err, ErrorInvalidCharacter)
		require.Equal(t, uint64(0), num)
	})
}
