package webmoney_test

import (
	wm "github.com/C0nstantin/go-webmoney" // Adjust the import path as necessary
	"os"
	"testing"
)

func TestXInterface_GetUrl(t *testing.T) {
	// Setup environment variables for testing
	os.Setenv("MERCHANT_DOMAIN", "test.merchant.web.money")
	os.Setenv("CLASSIC_DOMAIN", "test.w3s.webmoney.com")
	os.Setenv("LIGHT_DOMAIN", "test.light.webmoney.com")

	// Define test cases
	tests := []struct {
		name       string
		xInterface wm.XInterface
		isClassic  bool
		want       string
	}{
		{
			name:       "Classic w3s",
			xInterface: wm.XInterface{Name: "InterfaceName", Type: "w3s"},
			isClassic:  true,
			want:       "https://test.w3s.webmoney.com/asp/XMLInterfaceName.asp",
		},
		{
			name:       "Light FindWMPurseNew",
			xInterface: wm.XInterface{Name: "FindWMPurseNew", Type: "w3s"},
			isClassic:  false,
			want:       "https://test.light.webmoney.com/asp/XMLFindWMPurseNew.asp",
		},
		{
			name:       "Light other",
			xInterface: wm.XInterface{Name: "OtherInterface", Type: "w3s"},
			isClassic:  false,
			want:       "https://test.light.webmoney.com/asp/XMLOtherInterfaceCert.asp",
		},
		{
			name:       "Merchant type",
			xInterface: wm.XInterface{Name: "MerchantInterface", Type: "merchant"},
			isClassic:  false,
			want:       "https://test.merchant.web.money/conf/xml/XMLMerchantInterface.asp",
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.xInterface.GetUrl(tt.isClassic); got != tt.want {
				t.Errorf("XInterface.GetUrl() = %v, want %v", got, tt.want)
			}
		})
	}

	// Cleanup environment variables
	os.Unsetenv("MERCHANT_DOMAIN")
	os.Unsetenv("CLASSIC_DOMAIN")
	os.Unsetenv("LIGHT_DOMAIN")
}
