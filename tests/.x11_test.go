//

package tests

import (
	"fmt"
	"testing"
	"webmoney"
)

func TestPassport(t *testing.T) {
	result, err := webmoney.GetInfoWmid(`128756507061`)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(result.CertInfo)
	}
	fmt.Println(result.CertInfo)

}
