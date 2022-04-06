package bot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBot_isPaymentDocumentUrl(t *testing.T) {
	required := require.New(t)

	b := Bot{}

	txt := "https://lknpd.nalog.ru/api/v1/receipt/532128531749/2001rb3pdo/print"
	ok := b.isPaymentDocumentUrl(txt)
	required.True(ok)

	txt = "https://lknpd.naldg.ru/api/v1/receipt/532128531749/2001rb3pdo/prina"
	ok = b.isPaymentDocumentUrl(txt)
	required.False(ok)

	txt = "https://lnpd.nalog.ru/api/v1/receipt/532128531749/2001rb3pdo/print"
	ok = b.isPaymentDocumentUrl(txt)
	required.False(ok)

	txt = "https://lknpd.nalog.ru/api/v1/receipt/532128535749/2002ab3ddo/print"
	ok = b.isPaymentDocumentUrl(txt)
	required.True(ok)
}
