package powersoftau

import "testing"
import "github.com/FiloSottile/powersoftau/internal/chacha20"
import "encoding/hex"

func TestChaChaRng(t *testing.T) {
	r := chacha20.NewRng(&[32]byte{})
	if r.ReadUint32() != 0xade0b876 {
		t.Fail()
	}
}

func TestHashToG2(t *testing.T) {
	// From an instrumented Rust implementation:
	// hash_to_g2(0) = G2(x=Fq2(Fq(0x13f6c72ded114c2f55c291abdf68b032c7adb95c91dd3411606b7870703fd9d08c0e5c711d850611860c07522ec6cb00) + Fq(0x01db4a3b72b7e09ae15918061d2e02110926be25716c1b1614f3ef88be59c57ce58308bf3606159e33d845144350d924) * u), y=Fq2(Fq(0x12cca5c9e7e975092de2ceab7b7c82c0ab3c8875ac9f667525393a3786b282de9b7d84fb51ea0e235a4bb30367ccffe7) + Fq(0x05633798cffebe08ca83b1ac89851a96acdb4e34cd1ebc902d5a62fbcca8abfafd3581a17a6dd9f3026578eef578cc1c) * u)) = 81db4a3b72b7e09ae15918061d2e02110926be25716c1b1614f3ef88be59c57ce58308bf3606159e33d845144350d92413f6c72ded114c2f55c291abdf68b032c7adb95c91dd3411606b7870703fd9d08c0e5c711d850611860c07522ec6cb00

	res := HashToG2(make([]byte, 32)).EncodeCompressed()
	if hex.EncodeToString(res) != "81db4a3b72b7e09ae15918061d2e02110926be25716c1b1614f3ef88be59c57ce58308bf3606159e33d845144350d92413f6c72ded114c2f55c291abdf68b032c7adb95c91dd3411606b7870703fd9d08c0e5c711d850611860c07522ec6cb00" {
		t.Fail()
	}
}
